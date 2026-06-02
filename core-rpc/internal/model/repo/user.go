package repo

import (
	"core-rpc/internal/model/entity"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

//======================================================================================================================

type UserModel struct {
	DB *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		DB: db,
	}
}

func (u *UserModel) EmailLogin(email string) (*entity.User, error) {
	user := &entity.User{}
	err := u.DB.Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("该邮箱尚未注册")
		}
		return nil, err
	}
	return user, nil
}

func (u *UserModel) EmailRegister(name, email, password string) (*entity.User, error) {
	user := &entity.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	// 直接插入，依靠数据库唯一索引防重复
	err := u.DB.Create(user).Error
	if err != nil {
		if isDuplicateError(err) {
			return nil, errors.New("该邮箱已注册")
		}
		return nil, err
	}

	return user, nil
}

// 判断 MySQL 唯一键冲突
func isDuplicateError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func (u *UserModel) Logout() {
	// todo Logout暂时不做修改 后续在管理，当前是存入缓存
}

func (u *UserModel) ResetPasswordByEmail(email, newPassword string) error {
	user := &entity.User{}
	// 第一步：查用户是否存在
	err := u.DB.Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("not found")
		}
		return err
	}

	err = u.DB.Model(user).Update("password", newPassword).Error
	if err != nil {
		return err
	}
	return nil
}
