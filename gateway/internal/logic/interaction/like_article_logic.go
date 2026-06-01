package interaction

import (
	"context"

	core_client "core-rpc/core_client"
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

func (l *LikeArticleLogic) LikeArticle(req *types.LikeArticleReq) (resp *types.LikeArticleResp, err error) {
	userID, ok, msg := likeUserFromCtx(l.ctx)
	if !ok {
		return &types.LikeArticleResp{Code: authFailCode(msg), Msg: msg}, nil
	}
	if req.ArticleId == 0 {
		return &types.LikeArticleResp{Code: 400, Msg: "文章ID无效"}, nil
	}

	r, err := l.svcCtx.Core.LikeArticle(l.ctx, &core_client.LikeArticleReq{
		ArticleId: req.ArticleId,
		UserId:    userID,
	})
	if err != nil {
		return &types.LikeArticleResp{Code: likeErrCode(err), Msg: err.Error()}, nil
	}

	return &types.LikeArticleResp{
		Code: 200,
		Msg:  "点赞成功",
		Data: types.LikeArticleData{LikeCount: r.GetLikeCount()},
	}, nil
}
