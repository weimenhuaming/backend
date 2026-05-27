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

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCommentLogic) CreateComment(in *core.CreateCommentReq) (*core.CreateCommentResp, error) {
	if in.UserId == 0 || in.ArticleId == 0 {
		return nil, errors.New("参数无效")
	}
	if strings.TrimSpace(in.Content) == "" {
		return nil, errors.New("评论内容不能为空")
	}

	var commentID uint64
	err := l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
		var article entity.Article
		if err := tx.First(&article, in.ArticleId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("文章不存在")
			}
			return err
		}

		c := &entity.Comment{
			ArticleID: in.ArticleId,
			UserID:    in.UserId,
			ParentID:  0,
			Content:   strings.TrimSpace(in.Content),
		}
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		c.RootID = c.ID
		if err := tx.Model(c).Update("root_id", c.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&article).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			return err
		}
		commentID = c.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &core.CreateCommentResp{CommentId: commentID}, nil
}
