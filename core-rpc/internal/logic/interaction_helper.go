package logic

import (
	"core-rpc/internal/model/entity"
	"core-rpc/internal/utils"
	"core-rpc/pb/core"
)

// collectUserIDsFromComments extracts unique user ids from comments
func collectUserIDsFromComments(comments []entity.Comment) []uint64 {
	seen := make(map[uint64]struct{})
	var ids []uint64
	for _, c := range comments {
		if _, ok := seen[c.UserID]; !ok {
			seen[c.UserID] = struct{}{}
			ids = append(ids, c.UserID)
		}
	}
	return ids
}

// commentsToProtoList converts comments to proto objects using provided user map and preview replies
func commentsToProtoList(comments []entity.Comment, userMap map[uint64]entity.User, previewReplies map[uint64][]*core.CommentInfo) []*core.CommentInfo {
	out := make([]*core.CommentInfo, 0, len(comments))
	for i := range comments {
		name, avatar := utils.UserDisplay(userMap, comments[i].UserID)
		var replies []*core.CommentInfo
		if previewReplies != nil {
			replies = previewReplies[comments[i].ID]
		}
		out = append(out, utils.CommentToProto(&comments[i], name, avatar, replies))
	}
	return out
}
