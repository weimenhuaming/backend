package article

import (
	"context"
	"core-rpc/core_client"
	"gateway/internal/response"
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

func (l *CreateArticleLogic) CreateArticle(req *types.CreateArticleReq) error {
	// basic validation
	role := l.ctx.Value("X-user-Role")
	if role != "admin" {
		return response.NewError(403, "非管理员，没有权限执行")
	}

	if req.Title == "" && req.Content == "" {
		return response.NewError(400, "title is required")
	}

	// call core rpc
	UserId := l.ctx.Value("X-user-Id").(uint64)
	_, err := l.svcCtx.Core.CreateArticle(l.ctx, &core_client.CreateArticleReq{
		CategoryId: req.CategoryId,
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		Cover:      req.Cover,
		UserId:     UserId,
	})
	if err != nil {
		return response.NewError(500, err.Error())
	}

	return nil
}
