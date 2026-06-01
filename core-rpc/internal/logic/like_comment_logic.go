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

type LikeCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeCommentLogic {
	return &LikeCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikeCommentLogic) LikeComment(in *core.LikeCommentReq) (*core.LikeCommentResp, error) {
	if in.UserId == 0 || in.CommentId == 0 {
		return nil, errors.New("参数无效")
	}

	var count int64
	if err := l.svcCtx.Db.Model(&entity.Comment{}).Where("id = ?", in.CommentId).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("评论不存在")
	}

	var likeCount uint32
	err := l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
		delta, err := addCommentLike(tx, in.UserId, in.CommentId)
		if err != nil {
			return err
		}
		likeCount, err = adjustCommentLikeCount(tx, in.CommentId, delta)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &core.LikeCommentResp{LikeCount: likeCount}, nil
}
