package article

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticlesLogic {
	return &ListArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListArticlesLogic) ListArticles(req *types.ListArticlesReq) (resp *types.ListArticlesResp, err error) {
	// call core rpc
	rpcReq := &core_client.ListArticlesReq{
		Page:       req.Page,
		PageSize:   req.PageSize,
		CategoryId: req.CategoryId,
		UserId:     req.UserId,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
	}

	r, err := l.svcCtx.Core.ListArticles(l.ctx, rpcReq)
	if err != nil {
		return &types.ListArticlesResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	var arts []types.ArticleInfo
	for _, a := range r.GetArticles() {
		arts = append(arts, types.ArticleInfo{
			Id:           uint64(a.GetId()),
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

	return &types.ListArticlesResp{
		Code: 200,
		Msg:  "ok",
		Data: types.ListArticlesData{
			Articles: arts,
			Total:    r.GetTotal(),
			Page:     r.GetPage(),
			PageSize: r.GetPageSize(),
		},
	}, nil
}
