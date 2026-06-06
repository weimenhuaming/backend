package user

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/types"
)

func toTypesUserProfile(p *core_client.UserProfile) types.UserProfile {
	if p == nil {
		return types.UserProfile{}
	}
	return types.UserProfile{
		Id:     p.GetId(),
		Name:   p.GetName(),
		Phone:  p.GetPhone(),
		Email:  p.GetEmail(),
		Avatar: p.GetAvatar(),
		Role:   p.GetRole(),
		Sex:    p.GetSex(),
		Age:    p.GetAge(),
	}
}

func toTypesArticleList(articles []*core_client.ArticleInfo) []types.ArticleInfo {
	out := make([]types.ArticleInfo, 0, len(articles))
	for _, a := range articles {
		out = append(out, types.ArticleInfo{
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
		})
	}
	return out
}

func currentUserID(ctx context.Context) (uint64, bool) {
	userId, ok := ctx.Value("X-user-Id").(uint64)
	return userId, ok && userId > 0
}

func currentUserRole(ctx context.Context) string {
	role, _ := ctx.Value("X-user-Role").(string)
	return role
}
