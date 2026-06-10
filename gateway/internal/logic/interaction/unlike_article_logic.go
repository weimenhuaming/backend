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

type UnlikeArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnlikeArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikeArticleLogic {
	return &UnlikeArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnlikeArticleLogic) UnlikeArticle(req *types.UnlikeArticleReq) (resp *types.LikeArticleData, err error) {
	userID, ok, msg := vaild.GetActionUserID(l.ctx)
	if !ok {
		return nil, response.ErrorActionAuth(msg)
	}
	if req.ArticleId == 0 {
		return nil, response.ErrorBadRequest("文章ID无效")
	}

	r, err := l.svcCtx.Core.UnlikeArticle(l.ctx, &core_client.UnlikeArticleReq{
		ArticleId: req.ArticleId,
		UserId:    userID,
	})
	if err != nil {
		return nil, response.ErrorLikeOperation(err)
	}

	return &types.LikeArticleData{LikeCount: r.GetLikeCount()}, nil
}
