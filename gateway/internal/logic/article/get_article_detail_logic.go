package article

import (
	"context"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleDetailLogic {
	return &GetArticleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleDetailLogic) GetArticleDetail(req *types.GetArticleDetailReq) (resp *types.GetArticleDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
