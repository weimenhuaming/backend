package category

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// CategoryModel is the exported interface for category operations
type CategoryModel interface {
	Insert(ctx context.Context, data *Category) (sql.Result, error)
	FindAll(ctx context.Context) ([]*Category, error)
	FindOneByName(ctx context.Context, name string) (*Category, error)
	Delete(ctx context.Context, id uint64) error
}

type defaultCategoryModel struct {
	conn  sqlx.SqlConn
	table string
}

// Category 文章分类表结构
type Category struct {
	Id          uint64         `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

// NewCategoryModel returns a model for the database table.
func NewCategoryModel(conn sqlx.SqlConn) CategoryModel {
	return &defaultCategoryModel{
		conn:  conn,
		table: "`category`",
	}
}

func (m *defaultCategoryModel) Insert(ctx context.Context, data *Category) (sql.Result, error) {
	query := "insert into `category` (`name`, `description`, `created_at`, `updated_at`, `deleted_at`) values (?, ?, ?, ?, ?)"
	now := data.CreatedAt
	if now.IsZero() {
		now = time.Now()
	}
	ret, err := m.conn.ExecCtx(ctx, query, data.Name, data.Description, now, now, data.DeletedAt)
	return ret, err
}

func (m *defaultCategoryModel) FindAll(ctx context.Context) ([]*Category, error) {
	query := "select `id`, `name`, `description`, `created_at`, `updated_at`, `deleted_at` from `category` where deleted_at IS NULL order by id desc"
	var resp []*Category
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	return resp, err
}

func (m *defaultCategoryModel) FindOneByName(ctx context.Context, name string) (*Category, error) {
	var resp Category
	query := "select `id`, `name`, `description`, `created_at`, `updated_at`, `deleted_at` from `category` where `name` = ? limit 1"
	err := m.conn.QueryRowCtx(ctx, &resp, query, name)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultCategoryModel) Delete(ctx context.Context, id uint64) error {
	// soft delete: set deleted_at = now()
	query := "update `category` set `deleted_at` = ? where `id` = ?"
	now := time.Now()
	_, err := m.conn.ExecCtx(ctx, query, now, id)
	return err
}
