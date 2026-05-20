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
		SoftDelete(ctx context.Context, id uint64, userId uint64) error
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

// SoftDelete 设置 deleted_at 字段（软删除），仅当 article.id 和 article.user_id 同时匹配时执行
func (m *defaultArticleModel) SoftDelete(ctx context.Context, id uint64, userId uint64) error {
	query := fmt.Sprintf("UPDATE %s SET deleted_at = ? WHERE id = ? AND user_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id, userId)
	return err
}
