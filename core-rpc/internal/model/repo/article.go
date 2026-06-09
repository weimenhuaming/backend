package repo

import (
	"core-rpc/internal/model/converter"
	"core-rpc/internal/model/entity"
	"core-rpc/internal/utils"
	"core-rpc/pb/core"

	"gorm.io/gorm"
)

type ArtModel struct {
	DB *gorm.DB
}

func NewArtModel(db *gorm.DB) *ArtModel {
	return &ArtModel{
		DB: db,
	}
}

// Create inserts an article
func (m *ArtModel) Create(a *entity.Article) error {
	return m.DB.Create(a).Error
}

// FindByID returns article by id
func (m *ArtModel) FindByID(id uint64) (*entity.Article, error) {
	var a entity.Article
	if err := m.DB.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

// List returns articles with total count
func (m *ArtModel) List(offset, limit int) ([]entity.Article, int64, error) {
	var total int64
	q := m.DB.Model(&entity.Article{}).Order("created_at DESC")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []entity.Article
	if err := q.Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ListByCategory returns articles filtered by category
func (m *ArtModel) ListByCategory(categoryId uint64, offset, limit int) ([]entity.Article, int64, error) {
	var total int64
	q := m.DB.Model(&entity.Article{}).Where("category_id = ?", categoryId).Order("created_at DESC")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []entity.Article
	if err := q.Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// IncView increments view count
func (m *ArtModel) IncView(id uint64) error {
	return m.DB.Model(&entity.Article{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// LoadArticlesWithAuthors converts articles to proto with author info
func (m *ArtModel) LoadArticlesWithAuthors(articles []entity.Article) ([]*core.ArticleInfo, error) {
	if len(articles) == 0 {
		return nil, nil
	}
	userIDs := make([]uint64, 0, len(articles))
	seen := make(map[uint64]struct{})
	for _, a := range articles {
		if _, ok := seen[a.UserID]; !ok {
			seen[a.UserID] = struct{}{}
			userIDs = append(userIDs, a.UserID)
		}
	}
	var users []entity.User
	if len(userIDs) > 0 {
		if err := m.DB.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
			return nil, err
		}
	}
	userMap := utils.LoadUsersMap(users)
	out := make([]*core.ArticleInfo, 0, len(articles))
	for i := range articles {
		name, avatar := utils.UserDisplay(userMap, articles[i].UserID)
		out = append(out, converter.ArticleToProto(&articles[i], name, avatar))
	}
	return out, nil
}

// Search searches articles by keyword and optional category
func (m *ArtModel) Search(keyword string, categoryID uint64, offset, limit int) ([]entity.Article, int64, error) {
	q := m.DB.Model(&entity.Article{})
	if categoryID > 0 {
		q = q.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("title LIKE ? OR content LIKE ?", like, like)
	}
	q = q.Order("created_at DESC")

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []entity.Article
	if err := q.Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ListByUserID returns articles authored by a user
func (m *ArtModel) ListByUserID(userID uint64, offset, limit int) ([]entity.Article, int64, error) {
	var total int64
	q := m.DB.Model(&entity.Article{}).Where("user_id = ?", userID).Order("created_at DESC")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []entity.Article
	if err := q.Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// UpdateByID updates fields for an article by id
func (m *ArtModel) UpdateByID(id uint64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	return m.DB.Model(&entity.Article{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID deletes an article by id
func (m *ArtModel) DeleteByID(id uint64) error {
	return m.DB.Delete(&entity.Article{}, id).Error
}
