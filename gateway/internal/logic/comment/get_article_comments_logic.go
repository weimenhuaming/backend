// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package comment

import (
	"context"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArticleCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleCommentsLogic {
	return &GetArticleCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleCommentsLogic) GetArticleComments(req *types.GetArticleCommentsReq) (resp *types.GetArticleCommentsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
