package article

import (
	"context"
	"core-rpc/core_client"
	"gateway/internal/utils/vaild"

	"gateway/internal/response"
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

func (l *DeleteArticleLogic) DeleteArticle(req *types.DeleteArticleReq) error {
	uid, ok, msg := vaild.GetAdminUserID(l.ctx)
	if !ok {
		return response.ErrorAdminAuth(msg)
	}

	// call core rpc and pass user id explicitly
	_, err := l.svcCtx.Core.DeleteArticle(l.ctx, &core_client.DeleteArticleReq{Id: req.Id, UserId: uid})
	if err != nil {
		return response.ErrorInternalServer(err.Error())
	}

	return nil
}
