package repo

import "gorm.io/gorm"

type CommentModel struct {
	DB *gorm.DB
}

func NewCommentModel(db *gorm.DB) *CommentModel {
	return &CommentModel{
		DB: db,
	}
}
