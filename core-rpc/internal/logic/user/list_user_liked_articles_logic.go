package user

import (
	"context"
	"errors"

	"core-rpc/internal/svc"
	"core-rpc/internal/utils"

	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserLikedArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserLikedArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLikedArticlesLogic {
	return &ListUserLikedArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUserLikedArticlesLogic) ListUserLikedArticles(in *core.ListUserLikedArticlesReq) (*core.ListArticlesResp, error) {
	if in.GetUserId() == 0 {
		return nil, errors.New("用户 ID 无效")
	}

	page := utils.NormalizePageUint32(in.GetPage())
	size := utils.NormalizePageSizeUint32(in.GetPageSize(), 10)
	off, limit := utils.OffsetLimit(page, size)

	articles, total, err := l.svcCtx.InteractionRepo.ListLikedArticles(in.GetUserId(), off, limit)
	if err != nil {
		return nil, err
	}

	protoList, err := l.svcCtx.ArtRepo.LoadArticlesWithAuthors(articles)
	if err != nil {
		return nil, err
	}

	return &core.ListArticlesResp{
		Articles: protoList,
		Total:    uint32(total),
		Page:     uint32(page),
		PageSize: uint32(size),
	}, nil
}
