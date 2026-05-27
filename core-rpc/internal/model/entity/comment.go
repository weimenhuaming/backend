package entity

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement;comment:评论ID"`
	ArticleID   uint64         `gorm:"not null;index:idx_article_parent,priority:1;index:idx_article_created,priority:1;comment:文章ID"`
	UserID      uint64         `gorm:"not null;index:idx_user_id;comment:用户ID"`
	ParentID    uint64         `gorm:"not null;default:0;index:idx_article_parent,priority:2;comment:父评论ID"`
	RootID      uint64         `gorm:"not null;default:0;index:idx_root_id;comment:根评论ID"`
	ReplyToID   uint64         `gorm:"default:0;comment:回复目标用户ID"`
	ReplyToName string         `gorm:"size:100;default:'';comment:回复目标用户名"`
	Content     string         `gorm:"type:text;not null;comment:评论内容"`
	LikeCount   uint32         `gorm:"not null;default:0;comment:点赞数"`
	ChildCount  uint32         `gorm:"not null;default:0;comment:子评论数"`
	CreatedAt   time.Time      `gorm:"index:idx_article_created,priority:2;comment:创建时间"`
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_deleted_at;comment:删除时间"`
}

func (Comment) TableName() string {
	return "comment"
}
