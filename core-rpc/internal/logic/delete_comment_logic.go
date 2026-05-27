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

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCommentLogic) DeleteComment(in *core.DeleteCommentReq) (*core.DeleteCommentResp, error) {
	err := l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
		var c entity.Comment
		if err := tx.First(&c, in.Id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("评论不存在")
			}
			return err
		}
		if in.UserId > 0 && c.UserID != in.UserId {
			return errors.New("无权删除该评论")
		}

		if err := tx.Delete(&c).Error; err != nil {
			return err
		}

		// 一级评论：同时软删除其下回复
		if c.ParentID == 0 {
			var replies []entity.Comment
			if err := tx.Where("root_id = ? AND parent_id > 0", c.ID).Find(&replies).Error; err != nil {
				return err
			}
			if len(replies) > 0 {
				if err := tx.Delete(&replies).Error; err != nil {
					return err
				}
			}
			dec := 1 + int(c.ChildCount)
			return tx.Model(&entity.Article{}).Where("id = ?", c.ArticleID).
				Update("comment_count", gorm.Expr("GREATEST(comment_count - ?, 0)", dec)).Error
		}

		// 回复：减少根评论子评论数
		if err := tx.Model(&entity.Comment{}).Where("id = ?", c.RootID).
			Update("child_count", gorm.Expr("GREATEST(child_count - 1, 0)")).Error; err != nil {
			return err
		}
		return tx.Model(&entity.Article{}).Where("id = ?", c.ArticleID).
			Update("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)")).Error
	})
	if err != nil {
		return nil, err
	}
	return &core.DeleteCommentResp{}, nil
}
