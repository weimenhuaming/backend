package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	ActionLike  = "like"
	ActionFavor = "favor"
)

type InteractionLike struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	UserID     uint64         `gorm:"not null;uniqueIndex:uk_user_article_like,priority:1"`
	ArticleID  uint64         `gorm:"not null;uniqueIndex:uk_user_article_like,priority:2"`
	ActionType string         `gorm:"size:16;not null;uniqueIndex:uk_user_article_like,priority:3"`
}

func (InteractionLike) TableName() string {
	return "interaction_like"
}

type InteractionFavor struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	UserID     uint64         `gorm:"not null;uniqueIndex:uk_user_article_favor,priority:1"`
	ArticleID  uint64         `gorm:"not null;uniqueIndex:uk_user_article_favor,priority:2"`
	ActionType string         `gorm:"size:16;not null;uniqueIndex:uk_user_article_favor,priority:3"`
}

func (InteractionFavor) TableName() string {
	return "interaction_favor"
}

type InteractionCommentLike struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	UserID     uint64         `gorm:"not null;uniqueIndex:uk_user_comment_like,priority:1"`
	CommentID  uint64         `gorm:"not null;uniqueIndex:uk_user_comment_like,priority:2"`
	ActionType string         `gorm:"size:16;not null;uniqueIndex:uk_user_comment_like,priority:3"`
}

func (InteractionCommentLike) TableName() string {
	return "interaction_comment_like"
}

type TokenBlacklist struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	RefreshToken string    `gorm:"size:255;not null;uniqueIndex:uniq_refresh_token"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (TokenBlacklist) TableName() string {
	return "token_blacklist"
}
