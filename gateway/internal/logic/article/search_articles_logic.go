package article

import (
	"context"
	"strings"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils/converter"
	"gateway/internal/utils/vaild"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchArticlesLogic {
	return &SearchArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchArticlesLogic) SearchArticles(req *types.SearchArticlesReq) (resp *types.SearchArticlesData, err error) {
	keyword := strings.TrimSpace(req.Keyword)
	if keyword == "" {
		return nil, response.ErrorBadRequest("搜索关键词不能为空")
	}

	page, pageSize := vaild.NormalizePageSize(req.Page, req.PageSize)

	r, err := l.svcCtx.Core.SearchArticles(l.ctx, &core_client.SearchArticlesReq{
		Keyword:    keyword,
		Page:       page,
		PageSize:   pageSize,
		CategoryId: req.CategoryId,
	})
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}

	return &types.SearchArticlesData{
		Articles: converter.ToArticleList(r.GetArticles()),
		Total:    r.GetTotal(),
		Page:     r.GetPage(),
		PageSize: r.GetPageSize(),
	}, nil
}
