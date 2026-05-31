package logic

import (
	"core-rpc/internal/utils"
	"errors"

	"core-rpc/internal/model/entity"
	"core-rpc/pb/core"

	"gorm.io/gorm"
)

func toggleArticleLike(db *gorm.DB, userID, articleID uint64, active bool) (int32, error) {
	var row entity.InteractionLike
	err := db.Unscoped().Where("user_id = ? AND article_id = ? AND action_type = ?",
		userID, articleID, entity.ActionLike).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if !active {
			return 0, nil
		}
		row = entity.InteractionLike{
			UserID:     userID,
			ArticleID:  articleID,
			ActionType: entity.ActionLike,
		}
		return 1, db.Create(&row).Error
	}
	if err != nil {
		return 0, err
	}
	if active {
		if row.DeletedAt.Valid {
			return 1, db.Unscoped().Model(&row).Update("deleted_at", nil).Error
		}
		return 0, nil
	}
	if !row.DeletedAt.Valid {
		return -1, db.Delete(&row).Error
	}
	return 0, nil
}

func toggleArticleFavor(db *gorm.DB, userID, articleID uint64, active bool) (int32, error) {
	var row entity.InteractionFavor
	err := db.Unscoped().Where("user_id = ? AND article_id = ? AND action_type = ?",
		userID, articleID, entity.ActionFavor).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if !active {
			return 0, nil
		}
		row = entity.InteractionFavor{
			UserID:     userID,
			ArticleID:  articleID,
			ActionType: entity.ActionFavor,
		}
		return 1, db.Create(&row).Error
	}
	if err != nil {
		return 0, err
	}
	if active {
		if row.DeletedAt.Valid {
			return 1, db.Unscoped().Model(&row).Update("deleted_at", nil).Error
		}
		return 0, nil
	}
	if !row.DeletedAt.Valid {
		return -1, db.Delete(&row).Error
	}
	return 0, nil
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

func toggleCommentLike(db *gorm.DB, userID, commentID uint64, active bool) (int32, error) {
	var row entity.InteractionCommentLike
	err := db.Unscoped().Where("user_id = ? AND comment_id = ? AND action_type = ?",
		userID, commentID, entity.ActionLike).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if !active {
			return 0, nil
		}
		row = entity.InteractionCommentLike{
			UserID:     userID,
			CommentID:  commentID,
			ActionType: entity.ActionLike,
		}
		return 1, db.Create(&row).Error
	}
	if err != nil {
		return 0, err
	}
	if active {
		if row.DeletedAt.Valid {
			return 1, db.Unscoped().Model(&row).Update("deleted_at", nil).Error
		}
		return 0, nil
	}
	if !row.DeletedAt.Valid {
		return -1, db.Delete(&row).Error
	}
	return 0, nil
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
