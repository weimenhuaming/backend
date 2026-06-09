package interaction

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleLikeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArticleLikeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleLikeStatusLogic {
	return &GetArticleLikeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleLikeStatusLogic) GetArticleLikeStatus(req *types.GetArticleLikeStatusReq) (resp *types.GetArticleLikeStatusData, err error) {
	if req.ArticleId == 0 {
		return nil, response.NewError(400, "文章ID无效")
	}

	r, err := l.svcCtx.Core.GetArticleLikeStatus(l.ctx, &core_client.GetArticleLikeStatusReq{
		ArticleId: req.ArticleId,
		UserId:    userIDFromCtx(l.ctx),
	})
	if err != nil {
		return nil, response.NewError(500, err.Error())
	}

	return &types.GetArticleLikeStatusData{Liked: r.GetLiked()}, nil
}
