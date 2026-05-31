package logic

import (
	"core-rpc/internal/utils"
	"errors"

	"core-rpc/internal/model/entity"
	"core-rpc/pb/core"

	"gorm.io/gorm"
)

var (
	errAlreadyLiked = errors.New("已经点过赞了")
	errNotLiked     = errors.New("尚未点赞")
)

func interactionQuery(db *gorm.DB, userID uint64, objectType string, objectID uint64) *gorm.DB {
	return db.Where("user_id = ? AND object_type = ? AND object_id = ?",
		userID, objectType, objectID)
}

func likeQuery(db *gorm.DB, userID uint64, objectType string, objectID uint64) *gorm.DB {
	return interactionQuery(db, userID, objectType, objectID).
		Where("action_type = ?", entity.ActionLike)
}

func addLike(db *gorm.DB, userID uint64, objectType string, objectID uint64) (int32, error) {
	var row entity.InteractionLike
	err := interactionQuery(db, userID, objectType, objectID).First(&row).Error
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

func removeLike(db *gorm.DB, userID uint64, objectType string, objectID uint64) (int32, error) {
	var row entity.InteractionLike
	err := interactionQuery(db, userID, objectType, objectID).First(&row).Error
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

func addArticleLike(db *gorm.DB, userID, articleID uint64) (int32, error) {
	return addLike(db, userID, entity.ObjectTypeArticle, articleID)
}

func removeArticleLike(db *gorm.DB, userID, articleID uint64) (int32, error) {
	return removeLike(db, userID, entity.ObjectTypeArticle, articleID)
}

func adjustArticleCounter(db *gorm.DB, articleID uint64, field string, delta int32) (uint32, error) {
	if delta != 0 {
		if err := db.Model(&entity.Article{}).Where("id = ?", articleID).
			Update(field, gorm.Expr(field+" + ?", delta)).Error; err != nil {
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

func isObjectLiked(db *gorm.DB, userID uint64, objectType string, objectID uint64) (bool, error) {
	if userID == 0 || objectID == 0 {
		return false, nil
	}
	var count int64
	err := likeQuery(db, userID, objectType, objectID).Count(&count).Error
	return count > 0, err
}

func batchCommentLiked(db *gorm.DB, userID uint64, commentIDs []uint64) (map[uint64]bool, error) {
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
	if err := db.Where("user_id = ? AND object_type = ? AND action_type = ? AND object_id IN ?",
		userID, entity.ObjectTypeComment, entity.ActionLike, commentIDs).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.ObjectID] = true
	}
	return out, nil
}

func addCommentLike(db *gorm.DB, userID, commentID uint64) (int32, error) {
	return addLike(db, userID, entity.ObjectTypeComment, commentID)
}

func removeCommentLike(db *gorm.DB, userID, commentID uint64) (int32, error) {
	return removeLike(db, userID, entity.ObjectTypeComment, commentID)
}

func adjustCommentLikeCount(db *gorm.DB, commentID uint64, delta int32) (uint32, error) {
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

func loadArticlesWithAuthors(db *gorm.DB, articles []entity.Article) ([]*core.ArticleInfo, error) {
	if len(articles) == 0 {
		return nil, nil
	}
	userIDs := make([]uint64, 0, len(articles))
	seen := make(map[uint64]struct{})
	for _, a := range articles {
		if _, ok := seen[a.UserID]; !ok {
			seen[a.UserID] = struct{}{}
			userIDs = append(userIDs, a.UserID)
		}
	}
	var users []entity.User
	if len(userIDs) > 0 {
		if err := db.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
			return nil, err
		}
	}
	userMap := utils.LoadUsersMap(users)
	out := make([]*core.ArticleInfo, 0, len(articles))
	for i := range articles {
		name, avatar := utils.UserDisplay(userMap, articles[i].UserID)
		out = append(out, utils.ArticleToProto(&articles[i], name, avatar))
	}
	return out, nil
}

func fetchUserMap(db *gorm.DB, userIDs []uint64) (map[uint64]entity.User, error) {
	if len(userIDs) == 0 {
		return map[uint64]entity.User{}, nil
	}
	var users []entity.User
	if err := db.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	return utils.LoadUsersMap(users), nil
}

func collectUserIDsFromComments(comments []entity.Comment) []uint64 {
	seen := make(map[uint64]struct{})
	var ids []uint64
	for _, c := range comments {
		if _, ok := seen[c.UserID]; !ok {
			seen[c.UserID] = struct{}{}
			ids = append(ids, c.UserID)
		}
	}
	return ids
}

func commentsToProtoList(comments []entity.Comment, userMap map[uint64]entity.User, previewReplies map[uint64][]*core.CommentInfo) []*core.CommentInfo {
	out := make([]*core.CommentInfo, 0, len(comments))
	for i := range comments {
		name, avatar := utils.UserDisplay(userMap, comments[i].UserID)
		var replies []*core.CommentInfo
		if previewReplies != nil {
			replies = previewReplies[comments[i].ID]
		}
		out = append(out, utils.CommentToProto(&comments[i], name, avatar, replies))
	}
	return out
}
