package entity

import "time"

type TokenBlacklist struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	RefreshToken string    `gorm:"size:255;not null;uniqueIndex:uniq_refresh_token"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (TokenBlacklist) TableName() string {
	return "token_blacklist"
}
