package logic

import (
	"context"
	"errors"
	"strings"

	"core-rpc/internal/model/entity"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReplyLogic {
	return &CreateReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateReplyLogic) CreateReply(in *core.CreateReplyReq) (*core.CreateReplyResp, error) {
	if in.UserId == 0 || in.RootId == 0 || in.ParentId == 0 {
		return nil, errors.New("参数无效")
	}
	if strings.TrimSpace(in.Content) == "" {
		return nil, errors.New("回复内容不能为空")
	}

	var replyID uint64
	err := l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
		var root entity.Comment
		if err := tx.Where("id = ? AND parent_id = 0", in.RootId).First(&root).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("根评论不存在")
			}
			return err
		}

		var parent entity.Comment
		if err := tx.First(&parent, in.ParentId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("父评论不存在")
			}
			return err
		}
		if parent.ArticleID != root.ArticleID {
			return errors.New("评论不属于同一篇文章")
		}

		reply := &entity.Comment{
			ArticleID:   root.ArticleID,
			UserID:      in.UserId,
			ParentID:    in.ParentId,
			RootID:      in.RootId,
			ReplyToID:   in.ReplyToId,
			ReplyToName: in.ReplyToName,
			Content:     strings.TrimSpace(in.Content),
		}
		if err := tx.Create(reply).Error; err != nil {
			return err
		}
		if err := tx.Model(&root).Update("child_count", gorm.Expr("child_count + ?", 1)).Error; err != nil {
			return err
		}
		if err := tx.Model(&entity.Article{}).Where("id = ?", root.ArticleID).
			Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			return err
		}
		replyID = reply.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &core.CreateReplyResp{ReplyId: replyID}, nil
}
