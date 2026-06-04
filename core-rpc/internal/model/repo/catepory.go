package repo

import (
	"core-rpc/internal/model/entity"
	"gorm.io/gorm"
)

type CateModel struct {
	DB *gorm.DB
}

func NewCateModel(db *gorm.DB) *CateModel {
	return &CateModel{
		DB: db,
	}
}

// Create inserts a new category with the given name. Returns the created entity or an error.
func (m *CateModel) Create(name string) (*entity.Category, error) {
	c := &entity.Category{Name: name}
	if err := m.DB.Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

// FindByName returns category by name or gorm.ErrRecordNotFound
func (m *CateModel) FindByName(name string) (*entity.Category, error) {
	var c entity.Category
	if err := m.DB.Where("name = ?", name).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// List returns all categories ordered by id asc
func (m *CateModel) List() ([]entity.Category, error) {
	var categories []entity.Category
	if err := m.DB.Order("id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// DeleteByID deletes a category by id and returns rows affected
func (m *CateModel) DeleteByID(id uint64) (int64, error) {
	res := m.DB.Delete(&entity.Category{}, id)
	return res.RowsAffected, res.Error
}

// FindByID returns category by id
func (m *CateModel) FindByID(id uint64) (*entity.Category, error) {
	var c entity.Category
	if err := m.DB.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
