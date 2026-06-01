package entity

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	CreatedAt time.Time      `gorm:"comment:创建时间"`
	UpdatedAt time.Time      `gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:软删除时间"`
	Name      string         `gorm:"size:64;not null;default:'';uniqueIndex:uk_name;comment:用户名"`
	Password  string         `gorm:"size:255;not null;default:'';comment:密码"`
	Phone     string         `gorm:"size:20;not null;default:'';index:idx_phone;comment:手机号"`
	Avatar    string         `gorm:"size:255;not null;default:'';comment:头像URL"`
	Email     string         `gorm:"size:128;not null;default:'';index:idx_email;comment:邮箱"`
	Role      string         `gorm:"size:32;not null;default:user;comment:权限"`
	Sex       string         `gorm:"size:8;not null;default:未知;comment:性别"`
	Age       uint64         `gorm:"not null;default:0;comment:年龄"`
}

func (User) TableName() string {
	return "user"
}

type UserModel struct {
	DB *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		DB: db,
	}
}

func (u *UserModel) EmailLogin(email string) (*User, error) {
	user := &User{}
	err := u.DB.Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("该邮箱尚未注册")
		}
		return nil, err
	}
	return user, nil
}
