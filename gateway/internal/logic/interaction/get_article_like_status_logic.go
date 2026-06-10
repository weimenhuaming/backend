package interaction

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils/vaild"

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
		return nil, response.ErrorBadRequest("文章ID无效")
	}

	// 1. 不存在，则为游客
	userId, ok := vaild.GetUserID(l.ctx)
	if !ok {
		return &types.GetArticleLikeStatusData{Liked: false}, nil
	}

	// 2. 存在传id
	r, err := l.svcCtx.Core.GetArticleLikeStatus(l.ctx, &core_client.GetArticleLikeStatusReq{
		ArticleId: req.ArticleId,
		UserId:    userId,
	})
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}

	return &types.GetArticleLikeStatusData{Liked: r.GetLiked()}, nil
}
