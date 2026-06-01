package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	ActionLike    = "like"
	ActionUnknown = "unknown"

	ObjectTypeArticle = "article"
	ObjectTypeComment = "comment"
)

type InteractionLike struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	UserID     uint64         `gorm:"not null;uniqueIndex:uk_user_object_like,priority:1"`
	ObjectType string         `gorm:"size:16;not null;uniqueIndex:uk_user_object_like,priority:2"`
	ObjectID   uint64         `gorm:"not null;uniqueIndex:uk_user_object_like,priority:3"`
	ActionType string         `gorm:"size:16;not null"`
}

func (InteractionLike) TableName() string {
	return "interaction_like"
}
