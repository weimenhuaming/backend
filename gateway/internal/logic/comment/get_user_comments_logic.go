// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package comment

import (
	"context"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCommentsLogic {
	return &GetUserCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserCommentsLogic) GetUserComments(req *types.GetUserCommentsReq) (resp *types.GetUserCommentsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
