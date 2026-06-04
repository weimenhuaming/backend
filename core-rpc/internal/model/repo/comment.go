package repo

import (
	"core-rpc/internal/model/entity"
	"errors"
	"gorm.io/gorm"
)

type CommentModel struct {
	DB *gorm.DB
}

func NewCommentModel(db *gorm.DB) *CommentModel {
	return &CommentModel{
		DB: db,
	}
}

// CreateComment creates a root comment and increments article comment_count
func (m *CommentModel) CreateComment(userID, articleID uint64, content string) (uint64, error) {
	var commentID uint64
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		var article entity.Article
		if err := tx.First(&article, articleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("文章不存在")
			}
			return err
		}
		c := &entity.Comment{
			ArticleID: articleID,
			UserID:    userID,
			ParentID:  0,
			Content:   content,
		}
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		c.RootID = c.ID
		if err := tx.Model(c).Update("root_id", c.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&article).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			return err
		}
		commentID = c.ID
		return nil
	})
	if err != nil {
		return 0, err
	}
	return commentID, nil
}

// CreateReply creates a reply under a root comment, updates counters atomically
func (m *CommentModel) CreateReply(userID, rootID, parentID, replyToID uint64, replyToName, content string) (uint64, error) {
	var replyID uint64
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		var root entity.Comment
		if err := tx.Where("id = ? AND parent_id = 0", rootID).First(&root).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("根评论不存在")
			}
			return err
		}
		var parent entity.Comment
		if err := tx.First(&parent, parentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("父评论不存在")
			}
			return err
		}
		if parent.ArticleID != root.ArticleID {
			return errors.New("评论不属于同一篇文章")
		}
		reply := &entity.Comment{
			ArticleID:   root.ArticleID,
			UserID:      userID,
			ParentID:    parentID,
			RootID:      rootID,
			ReplyToID:   replyToID,
			ReplyToName: replyToName,
			Content:     content,
		}
		if err := tx.Create(reply).Error; err != nil {
			return err
		}
		if err := tx.Model(&root).Update("child_count", gorm.Expr("child_count + ?", 1)).Error; err != nil {
			return err
		}
		if err := tx.Model(&entity.Article{}).Where("id = ?", root.ArticleID).
			Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			return err
		}
		replyID = reply.ID
		return nil
	})
	if err != nil {
		return 0, err
	}
	return replyID, nil
}

// Delete deletes a comment (root or reply) and adjusts counters accordingly
func (m *CommentModel) Delete(commentID, userID uint64) error {
	return m.DB.Transaction(func(tx *gorm.DB) error {
		var c entity.Comment
		if err := tx.First(&c, commentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("评论不存在")
			}
			return err
		}
		if userID > 0 && c.UserID != userID {
			return errors.New("无权删除该评论")
		}
		if err := tx.Delete(&c).Error; err != nil {
			return err
		}
		if c.ParentID == 0 {
			var replies []entity.Comment
			if err := tx.Where("root_id = ? AND parent_id > 0", c.ID).Find(&replies).Error; err != nil {
				return err
			}
			if len(replies) > 0 {
				if err := tx.Delete(&replies).Error; err != nil {
					return err
				}
			}
			dec := 1 + int(c.ChildCount)
			return tx.Model(&entity.Article{}).Where("id = ?", c.ArticleID).
				Update("comment_count", gorm.Expr("GREATEST(comment_count - ?, 0)", dec)).Error
		}
		if err := tx.Model(&entity.Comment{}).Where("id = ?", c.RootID).
			Update("child_count", gorm.Expr("GREATEST(child_count - 1, 0)")).Error; err != nil {
			return err
		}
		return tx.Model(&entity.Article{}).Where("id = ?", c.ArticleID).
			Update("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)")).Error
	})
}

// ListTopComments lists root comments with pagination and ordering
func (m *CommentModel) ListTopComments(articleID uint64, orderBy string, offset, limit int) (int64, []entity.Comment, error) {
	var total int64
	q := m.DB.Model(&entity.Comment{}).Where("article_id = ? AND parent_id = 0", articleID)
	if orderBy == "hot" {
		q = q.Order("like_count DESC, created_at DESC")
	} else {
		q = q.Order("created_at DESC")
	}
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var comments []entity.Comment
	if err := q.Offset(offset).Limit(limit).Find(&comments).Error; err != nil {
		return 0, nil, err
	}
	return total, comments, nil
}

// ListReplies returns replies under a root comment (ordered asc)
func (m *CommentModel) ListReplies(rootID uint64, offset, limit int) (int64, []entity.Comment, error) {
	q := m.DB.Model(&entity.Comment{}).Where("root_id = ? AND parent_id > 0", rootID).Order("created_at ASC")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var replies []entity.Comment
	if err := q.Offset(offset).Limit(limit).Find(&replies).Error; err != nil {
		return 0, nil, err
	}
	return total, replies, nil
}

// ListByUser lists comments created by a user (both root and replies), paginated by created_at desc
func (m *CommentModel) ListByUser(userID uint64, offset, limit int) (int64, []entity.Comment, error) {
	q := m.DB.Model(&entity.Comment{}).Where("user_id = ?", userID).Order("created_at DESC")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var comments []entity.Comment
	if err := q.Offset(offset).Limit(limit).Find(&comments).Error; err != nil {
		return 0, nil, err
	}
	return total, comments, nil
}

// LoadPreviewReplies loads up to previewLimit replies for each root id in a single query
func (m *CommentModel) LoadPreviewReplies(rootIDs []uint64, previewLimit int) (map[uint64][]entity.Comment, error) {
	out := make(map[uint64][]entity.Comment)
	if len(rootIDs) == 0 || previewLimit <= 0 {
		return out, nil
	}
	var replies []entity.Comment
	if err := m.DB.Where("root_id IN ? AND parent_id > 0", rootIDs).Order("created_at ASC").Find(&replies).Error; err != nil {
		return nil, err
	}
	grouped := make(map[uint64][]entity.Comment)
	for _, r := range replies {
		list := grouped[r.RootID]
		if len(list) < previewLimit {
			grouped[r.RootID] = append(list, r)
		}
	}
	return grouped, nil
}

// FindByID returns a comment by id
func (m *CommentModel) FindByID(id uint64) (*entity.Comment, error) {
	var c entity.Comment
	if err := m.DB.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// Exists checks whether a comment exists
func (m *CommentModel) Exists(id uint64) (bool, error) {
	var cnt int64
	if err := m.DB.Model(&entity.Comment{}).Where("id = ?", id).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}
