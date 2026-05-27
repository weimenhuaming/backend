package article

import (
	"context"
	"core-rpc/core_client"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleLogic {
	return &DeleteArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteArticleLogic) DeleteArticle(req *types.DeleteArticleReq) (resp *types.DeleteArticleResp, err error) {
	role := l.ctx.Value("X-user-Role")
	if role != "admin" {
		return &types.DeleteArticleResp{
			Code: 403,
			Msg:  "非管理员，没有权限执行",
		}, nil
	}

	uid := l.ctx.Value("X-user-Id").(uint64)
	if uid == 0 {
		return &types.DeleteArticleResp{
			Code: 401,
			Msg:  "该用户不存在",
		}, nil
	}

	// call core rpc and pass user id explicitly
	_, err = l.svcCtx.Core.DeleteArticle(l.ctx, &core_client.DeleteArticleReq{Id: req.Id, UserId: uid})
	if err != nil {
		return &types.DeleteArticleResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &types.DeleteArticleResp{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}
