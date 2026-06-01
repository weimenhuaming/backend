package entity

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Name      string         `gorm:"size:128;not null;uniqueIndex:uk_category_name;comment:分类名称"`
	CreatedAt time.Time      `gorm:"comment:创建时间"`
	UpdatedAt time.Time      `gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:软删除时间"`
}

func (Category) TableName() string {
	return "category"
}
