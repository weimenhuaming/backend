package repo

import (
	"core-rpc/internal/model/entity"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserModel struct {
	DB *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		DB: db,
	}
}

//======================================================================================================================
// login section

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

func (u *UserModel) NameLogin(name, plainPassword string) (*entity.User, error) {
	user := &entity.User{}
	err := u.DB.Where("name = ?", name).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}
	if user.Password != plainPassword {
		return nil, errors.New("用户名或者密码错误")
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
		// 判断 MySQL 唯一键冲突
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, errors.New("该邮箱已注册")
		}
		return nil, err
	}

	return user, nil
}

func (u *UserModel) Logout() {
	// todo Logout暂时不做修改 后续在管理，当前是存入缓存
}

func (u *UserModel) ResetPasswordByEmail(email, newPassword string) error {
	user := &entity.User{}
	// 1. 查用户是否存在
	err := u.DB.Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("not found")
		}
		return err
	}

	// 2. 更新密码
	err = u.DB.Model(user).Update("password", newPassword).Error
	if err != nil {
		return err
	}
	return nil
}

// =====================================================================================================================
// user info section

func (u *UserModel) FindByID(id uint64) (*entity.User, error) {
	user := &entity.User{}
	if err := u.DB.Select("id", "name", "avatar").First(user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserModel) FindProfileByID(id uint64) (*entity.User, error) {
	user := &entity.User{}
	if err := u.DB.Select("id", "name", "phone", "email", "avatar", "role", "sex", "age").
		First(user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserModel) UpdateProfile(id uint64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	return u.DB.Model(&entity.User{}).Where("id = ?", id).Updates(updates).Error
}

func (u *UserModel) FindByIDs(ids []uint64) (map[uint64]entity.User, error) {
	out := make(map[uint64]entity.User)
	if len(ids) == 0 {
		return out, nil
	}
	var users []entity.User
	if err := u.DB.Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	for _, us := range users {
		out[us.ID] = us
	}
	return out, nil
}
