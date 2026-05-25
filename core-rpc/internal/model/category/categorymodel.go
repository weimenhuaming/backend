package category

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CategoryModel = (*customCategoryModel)(nil)

type (
	// CategoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCategoryModel.
	CategoryModel interface {
		categoryModel
		withSession(session sqlx.Session) CategoryModel
	}

	customCategoryModel struct {
		*defaultCategoryModel
	}
)

// NewCategoryModel returns a model for the database table.
func NewCategoryModel(conn sqlx.SqlConn) CategoryModel {
	return &customCategoryModel{
		defaultCategoryModel: newCategoryModel(conn),
	}
}

func (m *customCategoryModel) withSession(session sqlx.Session) CategoryModel {
	return NewCategoryModel(sqlx.NewSqlConnFromSession(session))
}
