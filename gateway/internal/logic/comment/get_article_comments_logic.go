// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package comment

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
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

func (l *GetArticleCommentsLogic) GetArticleComments(req *types.GetArticleCommentsReq) (resp *types.GetArticleCommentsData, err error) {
	if req.ArticleId == 0 {
		return nil, response.NewError(400, "文章ID不存在")
	}

	page := int32(req.Page)
	size := int32(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	r, err := l.svcCtx.Core.GetArticleComments(l.ctx, &core_client.GetArticleCommentsReq{
		ArticleId: req.ArticleId,
		Page:      page,
		Size:      size,
		OrderBy:   req.OrderBy,
	})
	if err != nil {
		return nil, response.NewError(500, err.Error())
	}

	return &types.GetArticleCommentsData{
		Comments: toTypesCommentList(r.GetComments()),
		Total:    uint32(r.GetTotal()),
		Page:     uint32(r.GetPage()),
		PageSize: uint32(r.GetSize()),
	}, nil
}
