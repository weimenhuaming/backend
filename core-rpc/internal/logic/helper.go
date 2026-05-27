package logic

import (
	"fmt"
	"time"

	"core-rpc/internal/model/entity"
	"core-rpc/pb/core"
)

const timeLayout = "2006-01-02 15:04:05"

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(timeLayout)
}

func normalizePage(page int32) int {
	if page <= 0 {
		return 1
	}
	return int(page)
}

func normalizeSize(size int32, defaultSize int) int {
	if size <= 0 {
		return defaultSize
	}
	if size > 100 {
		return 100
	}
	return int(size)
}

func normalizePageUint32(page uint32) int {
	if page == 0 {
		return 1
	}
	return int(page)
}

func normalizePageSizeUint32(pageSize uint32, defaultSize int) int {
	if pageSize == 0 {
		return defaultSize
	}
	if pageSize > 100 {
		return 100
	}
	return int(pageSize)
}

func offsetLimit(page, size int) (offset, limit int) {
	return (page - 1) * size, size
}

func articleToProto(a *entity.Article, authorName, authorAvatar string) *core.ArticleInfo {
	if a == nil {
		return nil
	}
	return &core.ArticleInfo{
		Id:           a.ID,
		UserId:       a.UserID,
		CategoryId:   a.CategoryID,
		Title:        a.Title,
		Summary:      a.Summary,
		Content:      a.Content,
		Cover:        a.Cover,
		ViewCount:    a.ViewCount,
		LikeCount:    a.LikeCount,
		FavorCount:   a.FavorCount,
		CommentCount: a.CommentCount,
		CreatedAt:    formatTime(a.CreatedAt),
		UpdatedAt:    formatTime(a.UpdatedAt),
		AuthorName:   authorName,
		AuthorAvatar: authorAvatar,
	}
}

func commentToProto(c *entity.Comment, userName, userAvatar string, replies []*core.CommentInfo) *core.CommentInfo {
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
		CreatedAt:   formatTime(c.CreatedAt),
		UserName:    userName,
		UserAvatar:  userAvatar,
		Replies:     replies,
	}
}

// loadUsersMap 批量加载用户展示信息
func loadUsersMap(users []entity.User) map[uint64]entity.User {
	m := make(map[uint64]entity.User, len(users))
	for _, u := range users {
		m[u.ID] = u
	}
	return m
}

func userDisplay(m map[uint64]entity.User, userID uint64) (name, avatar string) {
	if u, ok := m[userID]; ok {
		return u.Name, u.Avatar
	}
	return fmt.Sprintf("用户%d", userID), ""
}
