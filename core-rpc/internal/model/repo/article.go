package repo

import "gorm.io/gorm"

type ArtModel struct {
	DB *gorm.DB
}

func NewArtModel(db *gorm.DB) *ArtModel {
	return &ArtModel{
		DB: db,
	}
}
