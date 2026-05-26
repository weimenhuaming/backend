package article

import (
	"context"
	"fmt"
	"strings"
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
		// FindOneActive returns the article by id only when it is not soft-deleted
		FindOneActive(ctx context.Context, id uint64) (*Article, error)
		// List returns a page of articles (filters by categoryId/userId when non-zero), excludes soft-deleted rows
		List(ctx context.Context, page, pageSize uint32, categoryId, userId uint64, sortBy, sortOrder string) ([]*Article, error)
		// Count returns total number of articles matching filters (excludes soft-deleted rows)
		Count(ctx context.Context, categoryId, userId uint64) (int64, error)
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

// FindOneActive 查询单篇文章，但会排除已软删除的文章（deleted_at IS NULL）
func (m *defaultArticleModel) FindOneActive(ctx context.Context, id uint64) (*Article, error) {
	// reuse articleRows constant from generated file
	query := fmt.Sprintf("select %s from %s where `id` = ? and deleted_at IS NULL limit 1", articleRows, m.table)
	var resp Article
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// List 返回文章列表，自动过滤 deleted_at IS NULL，并支持分页、按 category/user 过滤与排序
func (m *defaultArticleModel) List(ctx context.Context, page, pageSize uint32, categoryId, userId uint64, sortBy, sortOrder string) ([]*Article, error) {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	// 白名单校验排序字段，避免 SQL 注入
	allowedSortFields := map[string]bool{
		"created_at": true,
		"view_count": true,
		"like_count": true,
	}
	sb := strings.ToLower(sortBy)
	if !allowedSortFields[sb] {
		sb = "created_at"
	}
	so := strings.ToLower(sortOrder)
	if so != "asc" {
		so = "desc"
	}

	where := "WHERE deleted_at IS NULL"
	var args []interface{}
	if categoryId != 0 {
		where += " AND category_id = ?"
		args = append(args, categoryId)
	}
	if userId != 0 {
		where += " AND user_id = ?"
		args = append(args, userId)
	}

	order := fmt.Sprintf("ORDER BY %s %s", sb, so)
	limitOffset := "LIMIT ? OFFSET ?"
	args = append(args, pageSize, (page-1)*pageSize)

	query := fmt.Sprintf("SELECT %s FROM %s %s %s %s", articleRows, m.table, where, order, limitOffset)

	var resp []*Article
	if err := m.conn.QueryRowsCtx(ctx, &resp, query, args...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Count 返回匹配条件的总数（用于分页 total）
func (m *defaultArticleModel) Count(ctx context.Context, categoryId, userId uint64) (int64, error) {
	where := "WHERE deleted_at IS NULL"
	var args []interface{}
	if categoryId != 0 {
		where += " AND category_id = ?"
		args = append(args, categoryId)
	}
	if userId != 0 {
		where += " AND user_id = ?"
		args = append(args, userId)
	}

	query := fmt.Sprintf("SELECT COUNT(1) FROM %s %s", m.table, where)
	var cnt int64
	if err := m.conn.QueryRowCtx(ctx, &cnt, query, args...); err != nil {
		return 0, err
	}
	return cnt, nil
}
