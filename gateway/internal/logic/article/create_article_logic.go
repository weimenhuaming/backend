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
	role := l.ctx.Value("X-user-Role")
	if role != "admin" {
		return &types.CreateArticleResp{
			Code: 403,
			Msg:  "非管理员，没有权限执行",
		}, nil
	}

	if req.Title == "" && req.Content == "" {
		return &types.CreateArticleResp{
			Code: 400,
			Msg:  "title is required",
		}, nil
	}

	// call core rpc
	UserId := l.ctx.Value("user_id").(uint64)
	_, err = l.svcCtx.Core.CreateArticle(l.ctx, &core_client.CreateArticleReq{
		CategoryId: req.CategoryId,
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		Cover:      req.Cover,
		UserId:     UserId,
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
