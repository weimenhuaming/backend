package logic

import (
	"context"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleCommentsLogic {
	return &GetArticleCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetArticleCommentsLogic) GetArticleComments(in *core.GetArticleCommentsReq) (*core.GetArticleCommentsResp, error) {
	page, size := normalizePageSize(in.GetPage(), in.GetSize())

	rows, err := l.svcCtx.CommentModel.ListTopLevelByArticle(l.ctx, in.GetArticleId(), page, size, in.GetOrderBy())
	if err != nil {
		return nil, err
	}
	total, err := l.svcCtx.CommentModel.CountTopLevelByArticle(l.ctx, in.GetArticleId())
	if err != nil {
		return nil, err
	}

	comments := buildCommentInfoList(l.ctx, l.svcCtx, rows)
	for i, row := range rows {
		preview, err := l.svcCtx.CommentModel.ListPreviewReplies(l.ctx, row.Id, 3)
		if err != nil {
			return nil, err
		}
		comments[i].Replies = buildCommentInfoList(l.ctx, l.svcCtx, preview)
	}

	return &core.GetArticleCommentsResp{
		Comments: comments,
		Page:     page,
		Size:     size,
		Total:    int32(total),
	}, nil
}
