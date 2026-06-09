package article

import (
	"context"
	"core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils/vaild"

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
	UserId, ok := vaild.GetUserID(l.ctx)
	if !ok {
		return response.ErrorUnauthorized("用户身份验证失败，请重新登录")
	}

	if ok = vaild.IsAdmin(l.ctx); !ok {
		return response.ErrorForbidden("非管理员，无权限操作")
	}

	if req.Title == "" && req.Content == "" {
		return response.ErrorBadRequest("title is required")
	}

	_, err := l.svcCtx.Core.CreateArticle(l.ctx, &core_client.CreateArticleReq{
		CategoryId: req.CategoryId,
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		Cover:      req.Cover,
		UserId:     UserId,
	})
	if err != nil {
		return response.ErrorInternalServer(err.Error())
	}

	return nil
}
