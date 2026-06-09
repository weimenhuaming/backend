package article

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticlesByCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArticlesByCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticlesByCategoryLogic {
	return &GetArticlesByCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticlesByCategoryLogic) GetArticlesByCategory(req *types.GetArticlesByCategoryReq) (resp *types.ListArticlesData, err error) {
	rpcReq := &core_client.GetArticlesByCategoryReq{
		CategoryId: req.CategoryId,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	r, err := l.svcCtx.Core.GetArticlesByCategory(l.ctx, rpcReq)
	if err != nil {
		return nil, response.NewError(500, err.Error())
	}

	var arts []types.ArticleInfo
	for _, a := range r.GetArticles() {
		arts = append(arts, types.ArticleInfo{
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

	return &types.ListArticlesData{
		Articles: arts,
		Total:    r.GetTotal(),
		Page:     r.GetPage(),
		PageSize: r.GetPageSize(),
	}, nil
}
