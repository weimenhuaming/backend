package logic

import (
	"context"

	"core-rpc/internal/model/entity"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentRepliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentRepliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentRepliesLogic {
	return &GetCommentRepliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentRepliesLogic) GetCommentReplies(in *core.GetCommentRepliesReq) (*core.GetCommentRepliesResp, error) {
	page := normalizePage(in.Page)
	size := normalizeSize(in.Size, 10)
	off, limit := offsetLimit(page, size)

	q := l.svcCtx.Db.Model(&entity.Comment{}).
		Where("root_id = ? AND parent_id > 0", in.RootId).
		Order("created_at ASC")

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	var replies []entity.Comment
	if err := q.Offset(off).Limit(limit).Find(&replies).Error; err != nil {
		return nil, err
	}

	userMap, err := fetchUserMap(l.svcCtx.Db, collectUserIDsFromComments(replies))
	if err != nil {
		return nil, err
	}

	return &core.GetCommentRepliesResp{
		Replies: commentsToProtoList(replies, userMap, nil),
		Page:    int32(page),
		Size:    int32(size),
		Total:   int32(total),
	}, nil
}
