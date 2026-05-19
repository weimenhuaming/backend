package article

import (
	"context"
	"core-rpc/core_client"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateArticleLogic {
	return &CreateArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateArticleLogic) CreateArticle(req *types.CreateArticleReq) (resp *types.CreateArticleResp, err error) {
	// basic validation
	if req.Title == "" || req.Content == "" {
		return &types.CreateArticleResp{
			Code: 400,
			Msg:  "title is required",
		}, nil
	}

	// call core rpc
	_, err = l.svcCtx.Core.CreateArticle(l.ctx, &core_client.CreateArticleReq{
		CategoryId: req.CategoryId,
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		Cover:      req.Cover,
		UserId:     req.UserId,
	})
	if err != nil {
		return &types.CreateArticleResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &types.CreateArticleResp{
		Code: 200,
		Msg:  "创建成功",
	}, nil
}
