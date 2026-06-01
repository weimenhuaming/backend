package logic

import (
	"context"
	"errors"

	"core-rpc/internal/model/entity"
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
	var article entity.Article
	if err := l.svcCtx.Db.First(&article, in.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if in.CategoryId > 0 {
		var cat entity.Category
		if err := l.svcCtx.Db.First(&cat, in.CategoryId).Error; err != nil {
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
	if err := l.svcCtx.Db.Model(&article).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &core.UpdateArticleResp{}, nil
}
