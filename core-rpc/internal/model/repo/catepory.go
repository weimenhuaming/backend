package repo

import "gorm.io/gorm"

type CateModel struct {
	DB *gorm.DB
}

func NewCateModel(db *gorm.DB) *CateModel {
	return &CateModel{
		DB: db,
	}
}
