package logic

import (
	"context"
	"core-rpc/internal/utils"

	"core-rpc/internal/model/entity"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCommentsLogic {
	return &GetUserCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserCommentsLogic) GetUserComments(in *core.GetUserCommentsReq) (*core.GetUserCommentsResp, error) {
	page := utils.NormalizePage(in.Page)
	size := utils.NormalizeSize(in.Size, 10)
	off, limit := utils.OffsetLimit(page, size)

	q := l.svcCtx.Db.Model(&entity.Comment{}).
		Where("user_id = ?", in.UserId).
		Order("created_at DESC")

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	var comments []entity.Comment
	if err := q.Offset(off).Limit(limit).Find(&comments).Error; err != nil {
		return nil, err
	}

	userMap, err := fetchUserMap(l.svcCtx.Db, collectUserIDsFromComments(comments))
	if err != nil {
		return nil, err
	}

	return &core.GetUserCommentsResp{
		Comments: commentsToProtoList(comments, userMap, nil),
		Page:     int32(page),
		Size:     int32(size),
		Total:    int32(total),
	}, nil
}
