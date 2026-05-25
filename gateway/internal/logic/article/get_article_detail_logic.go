package article

import (
	"context"
	"core-rpc/core_client"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleDetailLogic {
	return &GetArticleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleDetailLogic) GetArticleDetail(req *types.GetArticleDetailReq) (resp *types.GetArticleDetailResp, err error) {
	// call core rpc to get article detail
	r, err := l.svcCtx.Core.GetArticleDetail(l.ctx, &core_client.GetArticleDetailReq{Id: req.Id})
	if err != nil {
		return &types.GetArticleDetailResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	if r == nil || r.Article == nil {
		return &types.GetArticleDetailResp{
			Code: 404,
			Msg:  "article not found",
		}, nil
	}

	a := r.Article

	return &types.GetArticleDetailResp{
		Code: 200,
		Msg:  "success",
		Data: types.ArticleInfo{
			Id:           a.Id,
			UserId:       a.UserId,
			CategoryId:   a.CategoryId,
			Title:        a.Title,
			Summary:      a.Summary,
			Content:      a.Content,
			Cover:        a.Cover,
			ViewCount:    a.ViewCount,
			LikeCount:    a.LikeCount,
			FavorCount:   a.FavorCount,
			CommentCount: a.CommentCount,
			CreatedAt:    a.CreatedAt,
			UpdatedAt:    a.UpdatedAt,
			AuthorName:   a.AuthorName,
			AuthorAvatar: a.AuthorAvatar,
		},
	}, nil
}
