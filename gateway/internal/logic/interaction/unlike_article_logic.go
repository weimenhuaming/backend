package interaction

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/svc"
	"gateway/internal/types"

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

func (l *UnlikeArticleLogic) UnlikeArticle(req *types.UnlikeArticleReq) (resp *types.UnlikeArticleResp, err error) {
	userID, ok, msg := likeUserFromCtx(l.ctx)
	if !ok {
		return &types.UnlikeArticleResp{Code: authFailCode(msg), Msg: msg}, nil
	}
	if req.ArticleId == 0 {
		return &types.UnlikeArticleResp{Code: 400, Msg: "文章ID无效"}, nil
	}

	r, err := l.svcCtx.Core.UnlikeArticle(l.ctx, &core_client.UnlikeArticleReq{
		ArticleId: req.ArticleId,
		UserId:    userID,
	})
	if err != nil {
		return &types.UnlikeArticleResp{Code: likeErrCode(err), Msg: err.Error()}, nil
	}

	return &types.UnlikeArticleResp{
		Code: 200,
		Msg:  "已取消点赞",
		Data: types.LikeArticleData{LikeCount: r.GetLikeCount()},
	}, nil
}
