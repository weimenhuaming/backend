package logic

import (
	"context"
	"time"

	"core-rpc/internal/model/comment"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"
)

func normalizePageSize(page, size int32) (int32, int32) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	return page, size
}

func buildCommentInfo(ctx context.Context, svcCtx *svc.ServiceContext, c *comment.Comment) *core.CommentInfo {
	if c == nil {
		return nil
	}
	info := &core.CommentInfo{
		Id:          c.Id,
		ArticleId:   c.ArticleId,
		UserId:      c.UserId,
		ParentId:    c.ParentId,
		RootId:      c.RootId,
		ReplyToId:   c.ReplyToId,
		ReplyToName: c.ReplyToName,
		Content:     c.Content,
		LikeCount:   uint32(c.LikeCount),
		ChildCount:  uint32(c.ChildCount),
		CreatedAt:   c.CreatedAt.Format(time.DateTime),
	}
	if u, err := svcCtx.UserModel.FindOne(ctx, c.UserId); err == nil && u != nil {
		info.UserName = u.Name
		info.UserAvatar = u.Avatar
	}
	return info
}

func buildCommentInfoList(ctx context.Context, svcCtx *svc.ServiceContext, rows []*comment.Comment) []*core.CommentInfo {
	list := make([]*core.CommentInfo, 0, len(rows))
	for _, row := range rows {
		list = append(list, buildCommentInfo(ctx, svcCtx, row))
	}
	return list
}
