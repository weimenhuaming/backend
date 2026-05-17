package logic

import (
	"context"
	"core-rpc/internal/model/article"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"
	"database/sql"

	"github.com/zeromicro/go-zero/core/logx"
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

	newArticle := &article.Article{
		UserId:     in.UserId,
		CategoryId: in.GetCategoryId(),
		Title:      in.GetTitle(),
		Summary:    in.GetSummary(),
		Content:    sql.NullString{String: in.GetContent(), Valid: in.GetContent() != ""},
		Cover:      in.GetCover(),
	}

	_, err := l.svcCtx.ArticleModel.Insert(l.ctx, newArticle)
	if err != nil {
		return nil, err
	}

	return &core.CreateArticleResp{}, nil
}
