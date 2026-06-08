package article

import (
	"context"
	"errors"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArticleLogic {
	return &UpdateArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateArticleLogic) UpdateArticle(in *core.UpdateArticleReq) (*core.UpdateArticleResp, error) {
	// ensure article exists
	if _, err := l.svcCtx.ArtRepo.FindByID(in.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if in.CategoryId > 0 {
		if _, err := l.svcCtx.CateRepo.FindByID(in.CategoryId); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("分类不存在")
			}
			return nil, err
		}
		updates["category_id"] = in.CategoryId
	}
	if in.Title != "" {
		updates["title"] = in.Title
	}
	if in.Summary != "" {
		updates["summary"] = in.Summary
	}
	if in.Content != "" {
		updates["content"] = in.Content
	}
	if in.Cover != "" {
		updates["cover"] = in.Cover
	}
	if len(updates) == 0 {
		return &core.UpdateArticleResp{}, nil
	}
	if err := l.svcCtx.ArtRepo.UpdateByID(in.Id, updates); err != nil {
		return nil, err
	}
	return &core.UpdateArticleResp{}, nil
}
