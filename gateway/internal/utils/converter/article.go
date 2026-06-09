package converter

import (
	core_client "core-rpc/core_client"
	"gateway/internal/types"
)

func ToArticleInfo(a *core_client.ArticleInfo) *types.ArticleInfo {
	if a == nil {
		return &types.ArticleInfo{}
	}
	return &types.ArticleInfo{
		Id:           a.GetId(),
		UserId:       a.GetUserId(),
		CategoryId:   a.GetCategoryId(),
		Title:        a.GetTitle(),
		Summary:      a.GetSummary(),
		Content:      a.GetContent(),
		Cover:        a.GetCover(),
		ViewCount:    a.GetViewCount(),
		LikeCount:    a.GetLikeCount(),
		FavorCount:   a.GetFavorCount(),
		CommentCount: a.GetCommentCount(),
		CreatedAt:    a.GetCreatedAt(),
		UpdatedAt:    a.GetUpdatedAt(),
		AuthorName:   a.GetAuthorName(),
		AuthorAvatar: a.GetAuthorAvatar(),
	}
}

func ToArticleList(articles []*core_client.ArticleInfo) []types.ArticleInfo {
	out := make([]types.ArticleInfo, 0, len(articles))
	for _, a := range articles {
		out = append(out, *ToArticleInfo(a))
	}
	return out
}
