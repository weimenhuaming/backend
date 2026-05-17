package logic

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
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
	md, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}
	var uidStr string
	if v := md.Get("user-id"); len(v) > 0 {
		uidStr = v[0]
	} else if v := md.Get("user_id"); len(v) > 0 {
		uidStr = v[0]
	} else {
		return nil, errors.New("missing user id in metadata")
	}
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		return nil, err
	}

	// 检查文章存在与作者
	a, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	if a.UserId != uid {
		return nil, errors.New("not article owner")
	}

	// 更新可选字段
	if in.GetCategoryId() != 0 {
		a.CategoryId = in.GetCategoryId()
	}
	if in.GetTitle() != "" {
		a.Title = in.GetTitle()
	}
	if in.GetSummary() != "" {
		a.Summary = in.GetSummary()
	}
	if in.GetContent() != "" {
		a.Content = sql.NullString{String: in.GetContent(), Valid: true}
	}
	if in.GetCover() != "" {
		a.Cover = in.GetCover()
	}

	err = l.svcCtx.ArticleModel.Update(l.ctx, a)
	if err != nil {
		return nil, err
	}

	return &core.UpdateArticleResp{}, nil
}
