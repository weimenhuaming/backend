package article

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ArticleModel = (*customArticleModel)(nil)

type (
	// ArticleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleModel.
	ArticleModel interface {
		articleModel
		withSession(session sqlx.Session) ArticleModel
		SoftDelete(ctx context.Context, id uint64) error
	}

	customArticleModel struct {
		*defaultArticleModel
	}
)

// NewArticleModel returns a model for the database table.
func NewArticleModel(conn sqlx.SqlConn) ArticleModel {
	return &customArticleModel{
		defaultArticleModel: newArticleModel(conn),
	}
}

func (m *customArticleModel) withSession(session sqlx.Session) ArticleModel {
	return NewArticleModel(sqlx.NewSqlConnFromSession(session))
}

// SoftDelete 设置 deleted_at 字段（软删除）
func (m *defaultArticleModel) SoftDelete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("UPDATE %s SET deleted_at = ? WHERE id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}
