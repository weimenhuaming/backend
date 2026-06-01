package repo

import "gorm.io/gorm"

type InteractionModel struct {
	DB *gorm.DB
}

func NewInteractionModel(db *gorm.DB) *InteractionModel {
	return &InteractionModel{
		DB: db,
	}
}
