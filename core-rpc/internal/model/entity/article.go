package entity

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	CreatedAt    time.Time      `gorm:"comment:创建时间"`
	UpdatedAt    time.Time      `gorm:"comment:更新时间"`
	DeletedAt    gorm.DeletedAt `gorm:"index;comment:软删除时间"`
	UserID       uint64         `gorm:"not null;index;comment:作者用户ID"`
	CategoryID   uint64         `gorm:"not null;default:0;index;comment:分类ID"`
	Title        string         `gorm:"size:255;not null;default:'';comment:文章标题"`
	Summary      string         `gorm:"size:512;not null;default:'';comment:文章摘要"`
	Content      string         `gorm:"type:longtext;comment:文章内容"`
	Cover        string         `gorm:"size:255;not null;default:'';comment:封面图URL"`
	ViewCount    uint32         `gorm:"not null;default:0;comment:浏览量"`
	LikeCount    uint32         `gorm:"not null;default:0;comment:点赞数"`
	FavorCount   uint32         `gorm:"not null;default:0;comment:收藏数"`
	CommentCount uint32         `gorm:"not null;default:0;comment:评论数"`
}

func (Article) TableName() string {
	return "article"
}
