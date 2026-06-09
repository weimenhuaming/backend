package converter

import (
	"core-rpc/internal/model/entity"
	"core-rpc/internal/utils"
	"core-rpc/pb/core"
)

func ArticleToProto(a *entity.Article, authorName, authorAvatar string) *core.ArticleInfo {
	if a == nil {
		return nil
	}
	return &core.ArticleInfo{
		Id:           a.ID,
		UserId:       a.UserID,
		CategoryId:   a.CategoryID,
		Title:        a.Title,
		Summary:      a.Summary,
		Content:      a.Content,
		Cover:        a.Cover,
		ViewCount:    a.ViewCount,
		LikeCount:    a.LikeCount,
		FavorCount:   a.FavorCount,
		CommentCount: a.CommentCount,
		CreatedAt:    utils.FormatTime(a.CreatedAt),
		UpdatedAt:    utils.FormatTime(a.UpdatedAt),
		AuthorName:   authorName,
		AuthorAvatar: authorAvatar,
	}
}
