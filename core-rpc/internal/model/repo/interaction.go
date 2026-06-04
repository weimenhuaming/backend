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

// AddLike records a like (or flips an existing record to like). db can be a transaction or nil to use model DB
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

// RemoveLike marks a like as removed
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

// AddArticleLike convenience
func (m *InteractionModel) AddArticleLike(db *gorm.DB, userID, articleID uint64) (int32, error) {
	return m.AddLike(db, userID, entity.ObjectTypeArticle, articleID)
}

// RemoveArticleLike convenience
func (m *InteractionModel) RemoveArticleLike(db *gorm.DB, userID, articleID uint64) (int32, error) {
	return m.RemoveLike(db, userID, entity.ObjectTypeArticle, articleID)
}

// AdjustArticleCounter adjusts a numeric field on article and returns the updated value
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

// IsObjectLiked checks if the user liked an object
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

// BatchCommentLiked returns a map of commentID->liked
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

// AddCommentLike / RemoveCommentLike / AdjustCommentLikeCount
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

// LikeArticle performs like operation within a transaction and returns updated like count
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

// UnlikeArticle performs unlike operation within a transaction and returns updated like count
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

// LikeComment performs like operation for a comment within a transaction and returns updated like count
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

// UnlikeComment performs unlike operation for a comment within a transaction and returns updated like count
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
