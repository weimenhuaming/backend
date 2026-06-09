package utils

import (
	"fmt"
	"time"

	"core-rpc/internal/model/entity"
	"core-rpc/pb/core"
)

const timeLayout = "2006-01-02 15:04:05"

func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(timeLayout)
}

func NormalizePage(page int32) int {
	if page <= 0 {
		return 1
	}
	return int(page)
}

func NormalizeSize(size int32, defaultSize int) int {
	if size <= 0 {
		return defaultSize
	}
	if size > 100 {
		return 100
	}
	return int(size)
}

func NormalizePageUint32(page uint32) int {
	if page == 0 {
		return 1
	}
	return int(page)
}

func NormalizePageSizeUint32(pageSize uint32, defaultSize int) int {
	if pageSize == 0 {
		return defaultSize
	}
	if pageSize > 100 {
		return 100
	}
	return int(pageSize)
}

func OffsetLimit(page, size int) (offset, limit int) {
	return (page - 1) * size, size
}

func CommentToProto(c *entity.Comment, userName, userAvatar string, replies []*core.CommentInfo) *core.CommentInfo {
	if c == nil {
		return nil
	}
	return &core.CommentInfo{
		Id:          c.ID,
		ArticleId:   c.ArticleID,
		UserId:      c.UserID,
		ParentId:    c.ParentID,
		RootId:      c.RootID,
		ReplyToId:   c.ReplyToID,
		ReplyToName: c.ReplyToName,
		Content:     c.Content,
		LikeCount:   c.LikeCount,
		ChildCount:  c.ChildCount,
		CreatedAt:   FormatTime(c.CreatedAt),
		UserName:    userName,
		UserAvatar:  userAvatar,
		Replies:     replies,
	}
}

// LoadUsersMap 批量加载用户展示信息
func LoadUsersMap(users []entity.User) map[uint64]entity.User {
	m := make(map[uint64]entity.User, len(users))
	for _, u := range users {
		m[u.ID] = u
	}
	return m
}

func UserDisplay(m map[uint64]entity.User, userID uint64) (name, avatar string) {
	if u, ok := m[userID]; ok {
		return u.Name, u.Avatar
	}
	return fmt.Sprintf("用户%d", userID), ""
}
