package logic

import (
	"strings"

	"core-rpc/internal/model/entity"

	"gorm.io/gorm"
)

func applyArticleFilters(q *gorm.DB, categoryID, userID uint64) *gorm.DB {
	if categoryID > 0 {
		q = q.Where("category_id = ?", categoryID)
	}
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	return q
}

func applyArticleSort(q *gorm.DB, sortBy, sortOrder string) *gorm.DB {
	field := "created_at"
	switch strings.ToLower(sortBy) {
	case "view_count", "like_count", "created_at":
		field = strings.ToLower(sortBy)
	}
	order := "DESC"
	if strings.EqualFold(sortOrder, "asc") {
		order = "ASC"
	}
	return q.Order(field + " " + order)
}

func listArticlesQuery(db *gorm.DB, categoryID, userID uint64, sortBy, sortOrder string) *gorm.DB {
	return applyArticleSort(applyArticleFilters(db.Model(&entity.Article{}), categoryID, userID), sortBy, sortOrder)
}
