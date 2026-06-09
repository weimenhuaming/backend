package interaction

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeArticleLogic {
	return &LikeArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeArticleLogic) LikeArticle(req *types.LikeArticleReq) (resp *types.LikeArticleData, err error) {
	userID, ok, msg := likeUserFromCtx(l.ctx)
	if !ok {
		return nil, authFailError(msg)
	}
	if req.ArticleId == 0 {
		return nil, response.ErrorBadRequest("文章ID无效")
	}

	r, err := l.svcCtx.Core.LikeArticle(l.ctx, &core_client.LikeArticleReq{
		ArticleId: req.ArticleId,
		UserId:    userID,
	})
	if err != nil {
		return nil, likeRPCError(err)
	}

	return &types.LikeArticleData{LikeCount: r.GetLikeCount()}, nil
}
