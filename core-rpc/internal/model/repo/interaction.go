package repo

import (
	"core-rpc/internal/model/entity"
	"errors"
	"gorm.io/gorm"
)

var (
	errAlreadyLiked = errors.New("已经点过赞了")
	errNotLiked     = errors.New("尚未点赞")
)

type InteractionModel struct {
	DB *gorm.DB
}

func NewInteractionModel(db *gorm.DB) *InteractionModel {
	return &InteractionModel{
		DB: db,
	}
}

// AddLike 记录点赞（或将已有记录切换为点赞）。
// 参数 db 可以是事务（*gorm.DB）或 nil，若为 nil 则使用模型内置的 DB
func (m *InteractionModel) AddLike(db *gorm.DB, userID uint64, objectType string, objectID uint64) (int32, error) {
	if db == nil {
		db = m.DB
	}
	var row entity.InteractionLike
	err := db.Where("user_id = ? AND object_type = ? AND object_id = ?", userID, objectType, objectID).First(&row).Error
	if err == nil {
		if row.ActionType == entity.ActionLike {
			return 0, errAlreadyLiked
		}
		return 1, db.Model(&row).Update("action_type", entity.ActionLike).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	row = entity.InteractionLike{
		UserID:     userID,
		ObjectType: objectType,
		ObjectID:   objectID,
		ActionType: entity.ActionLike,
	}
	return 1, db.Create(&row).Error
}

// RemoveLike 将点赞标记为已取消
func (m *InteractionModel) RemoveLike(db *gorm.DB, userID uint64, objectType string, objectID uint64) (int32, error) {
	if db == nil {
		db = m.DB
	}
	var row entity.InteractionLike
	err := db.Where("user_id = ? AND object_type = ? AND object_id = ?", userID, objectType, objectID).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errNotLiked
	}
	if err != nil {
		return 0, err
	}
	if row.ActionType != entity.ActionLike {
		return 0, errNotLiked
	}
	return -1, db.Model(&row).Update("action_type", entity.ActionUnknown).Error
}

// AddArticleLike 文章点赞的便捷包装
func (m *InteractionModel) AddArticleLike(db *gorm.DB, userID, articleID uint64) (int32, error) {
	return m.AddLike(db, userID, entity.ObjectTypeArticle, articleID)
}

// RemoveArticleLike 文章取消点赞的便捷包装
func (m *InteractionModel) RemoveArticleLike(db *gorm.DB, userID, articleID uint64) (int32, error) {
	return m.RemoveLike(db, userID, entity.ObjectTypeArticle, articleID)
}

// AdjustArticleCounter 调整文章的数值字段（例如 like_count、favor_count）并返回更新后的值
func (m *InteractionModel) AdjustArticleCounter(db *gorm.DB, articleID uint64, field string, delta int32) (uint32, error) {
	if db == nil {
		db = m.DB
	}
	if delta != 0 {
		if err := db.Model(&entity.Article{}).Where("id = ?", articleID).Update(field, gorm.Expr(field+" + ?", delta)).Error; err != nil {
			return 0, err
		}
	}
	var article entity.Article
	if err := db.Select(field).First(&article, articleID).Error; err != nil {
		return 0, err
	}
	switch field {
	case "like_count":
		return article.LikeCount, nil
	case "favor_count":
		return article.FavorCount, nil
	default:
		return 0, errors.New("unknown counter field")
	}
}

// IsObjectLiked 检查指定用户是否对某个对象（文章/评论等）进行了点赞
func (m *InteractionModel) IsObjectLiked(db *gorm.DB, userID uint64, objectType string, objectID uint64) (bool, error) {
	if db == nil {
		db = m.DB
	}
	if userID == 0 || objectID == 0 {
		return false, nil
	}
	var count int64
	if err := db.Where("user_id = ? AND object_type = ? AND action_type = ? AND object_id = ?", userID, objectType, entity.ActionLike, objectID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// BatchCommentLiked 批量查询评论的点赞状态，返回 map[commentID]bool
func (m *InteractionModel) BatchCommentLiked(db *gorm.DB, userID uint64, commentIDs []uint64) (map[uint64]bool, error) {
	if db == nil {
		db = m.DB
	}
	out := make(map[uint64]bool, len(commentIDs))
	for _, id := range commentIDs {
		if id > 0 {
			out[id] = false
		}
	}
	if userID == 0 || len(commentIDs) == 0 {
		return out, nil
	}
	var rows []entity.InteractionLike
	if err := db.Where("user_id = ? AND object_type = ? AND action_type = ? AND object_id IN ?", userID, entity.ObjectTypeComment, entity.ActionLike, commentIDs).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.ObjectID] = true
	}
	return out, nil
}

// AddCommentLike / RemoveCommentLike / AdjustCommentLikeCount: 评论点赞/取消与计数调整的封装方法
func (m *InteractionModel) AddCommentLike(db *gorm.DB, userID, commentID uint64) (int32, error) {
	return m.AddLike(db, userID, entity.ObjectTypeComment, commentID)
}

func (m *InteractionModel) RemoveCommentLike(db *gorm.DB, userID, commentID uint64) (int32, error) {
	return m.RemoveLike(db, userID, entity.ObjectTypeComment, commentID)
}

func (m *InteractionModel) AdjustCommentLikeCount(db *gorm.DB, commentID uint64, delta int32) (uint32, error) {
	if db == nil {
		db = m.DB
	}
	if delta != 0 {
		if err := db.Model(&entity.Comment{}).Where("id = ?", commentID).
			Update("like_count", gorm.Expr("like_count + ?", delta)).Error; err != nil {
			return 0, err
		}
	}
	var c entity.Comment
	if err := db.Select("like_count").First(&c, commentID).Error; err != nil {
		return 0, err
	}
	return c.LikeCount, nil
}

// LikeArticle 在事务中执行文章点赞操作，并返回更新后的点赞数量
func (m *InteractionModel) LikeArticle(userID, articleID uint64) (uint32, error) {
	var likeCount uint32
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		delta, err := m.AddArticleLike(tx, userID, articleID)
		if err != nil {
			return err
		}
		likeCount, err = m.AdjustArticleCounter(tx, articleID, "like_count", delta)
		return err
	})
	return likeCount, err
}

// UnlikeArticle 在事务中执行文章取消点赞操作，并返回更新后的点赞数量
func (m *InteractionModel) UnlikeArticle(userID, articleID uint64) (uint32, error) {
	var likeCount uint32
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		delta, err := m.RemoveArticleLike(tx, userID, articleID)
		if err != nil {
			return err
		}
		likeCount, err = m.AdjustArticleCounter(tx, articleID, "like_count", delta)
		return err
	})
	return likeCount, err
}

// LikeComment 在事务中执行评论点赞操作，并返回更新后的点赞数量
func (m *InteractionModel) LikeComment(userID, commentID uint64) (uint32, error) {
	var likeCount uint32
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		delta, err := m.AddCommentLike(tx, userID, commentID)
		if err != nil {
			return err
		}
		likeCount, err = m.AdjustCommentLikeCount(tx, commentID, delta)
		return err
	})
	return likeCount, err
}

// UnlikeComment 在事务中执行评论取消点赞操作，并返回更新后的点赞数量
func (m *InteractionModel) UnlikeComment(userID, commentID uint64) (uint32, error) {
	var likeCount uint32
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		delta, err := m.RemoveCommentLike(tx, userID, commentID)
		if err != nil {
			return err
		}
		likeCount, err = m.AdjustCommentLikeCount(tx, commentID, delta)
		return err
	})
	return likeCount, err
}
