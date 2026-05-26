package logic

import (
	"context"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticlesLogic {
	return &ListArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListArticlesLogic) ListArticles(in *core.ListArticlesReq) (*core.ListArticlesResp, error) {
	// default paging
	page := in.GetPage()
	pageSize := in.GetPageSize()
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	// fetch list and total from model (model will exclude soft-deleted rows)
	arts, err := l.svcCtx.ArticleModel.List(l.ctx, page, pageSize, in.GetCategoryId(), in.GetUserId(), in.GetSortBy(), in.GetSortOrder())
	if err != nil {
		return nil, err
	}

	total, err := l.svcCtx.ArticleModel.Count(l.ctx, in.GetCategoryId(), in.GetUserId())
	if err != nil {
		return nil, err
	}

	var respArts []*core.ArticleInfo
	for _, a := range arts {
		content := ""
		if a.Content.Valid {
			content = a.Content.String
		}
		createdAt := a.CreatedAt.Format("2006-01-02 15:04:05")
		updatedAt := a.UpdatedAt.Format("2006-01-02 15:04:05")

		var authorName, authorAvatar string
		if au, e := l.svcCtx.UserModel.FindOne(l.ctx, a.UserId); e == nil && au != nil {
			authorName = au.Name
			authorAvatar = au.Avatar
		}

		respArts = append(respArts, &core.ArticleInfo{
			Id:           a.Id,
			UserId:       a.UserId,
			CategoryId:   a.CategoryId,
			Title:        a.Title,
			Summary:      a.Summary,
			Content:      content,
			Cover:        a.Cover,
			ViewCount:    uint32(a.ViewCount),
			LikeCount:    uint32(a.LikeCount),
			FavorCount:   uint32(a.FavorCount),
			CommentCount: uint32(a.CommentCount),
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
			AuthorName:   authorName,
			AuthorAvatar: authorAvatar,
		})
	}

	return &core.ListArticlesResp{
		Articles: respArts,
		Total:    uint32(total),
		Page:     page,
		PageSize: pageSize,
	}, nil
}
