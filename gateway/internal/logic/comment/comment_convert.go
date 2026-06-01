package comment

import (
	core_client "core-rpc/core_client"
	"gateway/internal/types"
)

func toTypesCommentInfo(c *core_client.CommentInfo) types.CommentInfo {
	if c == nil {
		return types.CommentInfo{}
	}
	info := types.CommentInfo{
		Id:          c.GetId(),
		ArticleId:   c.GetArticleId(),
		UserId:      c.GetUserId(),
		ParentId:    c.GetParentId(),
		RootId:      c.GetRootId(),
		ReplyToId:   c.GetReplyToId(),
		ReplyToName: c.GetReplyToName(),
		Content:     c.GetContent(),
		LikeCount:   c.GetLikeCount(),
		ChildCount:  c.GetChildCount(),
		CreatedAt:   c.GetCreatedAt(),
		UserName:    c.GetUserName(),
		UserAvatar:  c.GetUserAvatar(),
	}
	for _, r := range c.GetReplies() {
		reply := toTypesCommentInfo(r)
		info.Replies = append(info.Replies, reply)
	}
	return info
}

func toTypesCommentList(list []*core_client.CommentInfo) []types.CommentInfo {
	out := make([]types.CommentInfo, 0, len(list))
	for _, item := range list {
		out = append(out, toTypesCommentInfo(item))
	}
	return out
}
