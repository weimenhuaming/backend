package logic

import (
	"context"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleDetailLogic {
	return &GetArticleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetArticleDetailLogic) GetArticleDetail(in *core.GetArticleDetailReq) (*core.GetArticleDetailResp, error) {
	// 读取文章
	a, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	// 尝试读取作者信息（非阻塞，如果失败不影响文章内容返回）
	var authorName, authorAvatar string
	au, err := l.svcCtx.UserModel.FindOne(l.ctx, a.UserId)
	if err == nil && au != nil {
		authorName = au.Name
		authorAvatar = au.Avatar
	}

	createdAt := a.CreatedAt.Format("2006-01-02 15:04:05")
	updatedAt := a.UpdatedAt.Format("2006-01-02 15:04:05")
	content := ""
	if a.Content.Valid {
		content = a.Content.String
	}

	articleInfo := &core.ArticleInfo{
		Id:           a.Id,
		UserId:       a.UserId,
		CategoryId:   a.CategoryId,
		Title:        a.Title,
		Summary:      a.Summary,
		Content:      content,
		Cover:        a.Cover,
		ViewCount:    uint32(a.ViewCount),
		LikeCount:    uint32(a.LikeCount),
		FavorCount:   uint32(a.FavorCount),
		CommentCount: uint32(a.CommentCount),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		AuthorName:   authorName,
		AuthorAvatar: authorAvatar,
	}

	return &core.GetArticleDetailResp{
		Article: articleInfo,
	}, nil
}
