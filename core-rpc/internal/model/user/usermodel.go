package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		withSession(session sqlx.Session) UserModel
		FindOneByEmail(ctx context.Context, email string) (*User, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

func (m *customUserModel) withSession(session sqlx.Session) UserModel {
	return NewUserModel(sqlx.NewSqlConnFromSession(session))
}

// =====================================================================================================================
// 自定义查询

func (m *defaultUserModel) FindOneByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := fmt.Sprintf("SELECT %s FROM %s WHERE email = ? LIMIT 1", userRows, m.table)

	err := m.conn.QueryRowCtx(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
