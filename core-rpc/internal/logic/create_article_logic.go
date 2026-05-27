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

type CreateArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateArticleLogic {
	return &CreateArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateArticleLogic) CreateArticle(in *core.CreateArticleReq) (*core.CreateArticleResp, error) {
	if in.UserId == 0 {
		return nil, errors.New("用户未登录")
	}
	if in.Title == "" {
		return nil, errors.New("标题不能为空")
	}

	if in.CategoryId > 0 {
		var cat entity.Category
		if err := l.svcCtx.Db.First(&cat, in.CategoryId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("分类不存在")
			}
			return nil, err
		}
	}

	article := &entity.Article{
		UserID:     in.UserId,
		CategoryID: in.CategoryId,
		Title:      in.Title,
		Summary:    in.Summary,
		Content:    in.Content,
		Cover:      in.Cover,
	}
	if err := l.svcCtx.Db.Create(article).Error; err != nil {
		return nil, err
	}
	return &core.CreateArticleResp{}, nil
}
