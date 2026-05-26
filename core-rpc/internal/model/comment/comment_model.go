package comment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CommentModel = (*customCommentModel)(nil)

type (
	CommentModel interface {
		commentModel
		withSession(session sqlx.Session) CommentModel
		SoftDelete(ctx context.Context, id, userId uint64) error
		FindOneActive(ctx context.Context, id uint64) (*Comment, error)
		UpdateRootId(ctx context.Context, id uint64) error
		IncChildCount(ctx context.Context, rootId uint64, delta int64) error
		IncLikeCount(ctx context.Context, id uint64, delta int64) (uint64, error)
		ListTopLevelByArticle(ctx context.Context, articleId uint64, page, size int32, orderBy string) ([]*Comment, error)
		CountTopLevelByArticle(ctx context.Context, articleId uint64) (int64, error)
		ListRepliesByRoot(ctx context.Context, rootId uint64, page, size int32) ([]*Comment, error)
		CountRepliesByRoot(ctx context.Context, rootId uint64) (int64, error)
		ListPreviewReplies(ctx context.Context, rootId uint64, limit int32) ([]*Comment, error)
		ListByUser(ctx context.Context, userId uint64, page, size int32) ([]*Comment, error)
		CountByUser(ctx context.Context, userId uint64) (int64, error)
		FindCommentLikeActive(ctx context.Context, userId, commentId uint64) (bool, error)
		InsertCommentLike(ctx context.Context, userId, commentId uint64) error
		RestoreCommentLike(ctx context.Context, userId, commentId uint64) error
		SoftDeleteCommentLike(ctx context.Context, userId, commentId uint64) error
	}

	customCommentModel struct {
		*defaultCommentModel
	}
)

func NewCommentModel(conn sqlx.SqlConn) CommentModel {
	return &customCommentModel{
		defaultCommentModel: newCommentModel(conn),
	}
}

func (m *customCommentModel) withSession(session sqlx.Session) CommentModel {
	return NewCommentModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultCommentModel) SoftDelete(ctx context.Context, id, userId uint64) error {
	query := fmt.Sprintf("UPDATE %s SET deleted_at = ? WHERE id = ? AND user_id = ? AND deleted_at IS NULL", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id, userId)
	return err
}

func (m *defaultCommentModel) FindOneActive(ctx context.Context, id uint64) (*Comment, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = ? AND deleted_at IS NULL LIMIT 1", commentRows, m.table)
	var resp Comment
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

func (m *defaultCommentModel) UpdateRootId(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("UPDATE %s SET root_id = ? WHERE id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id, id)
	return err
}

func (m *defaultCommentModel) IncChildCount(ctx context.Context, rootId uint64, delta int64) error {
	query := fmt.Sprintf("UPDATE %s SET child_count = child_count + ? WHERE id = ? AND deleted_at IS NULL", m.table)
	_, err := m.conn.ExecCtx(ctx, query, delta, rootId)
	return err
}

func (m *defaultCommentModel) IncLikeCount(ctx context.Context, id uint64, delta int64) (uint64, error) {
	query := fmt.Sprintf("UPDATE %s SET like_count = GREATEST(CAST(like_count AS SIGNED) + ?, 0) WHERE id = ? AND deleted_at IS NULL", m.table)
	_, err := m.conn.ExecCtx(ctx, query, delta, id)
	if err != nil {
		return 0, err
	}
	c, err := m.FindOneActive(ctx, id)
	if err != nil {
		return 0, err
	}
	return c.LikeCount, nil
}

func (m *defaultCommentModel) ListTopLevelByArticle(ctx context.Context, articleId uint64, page, size int32, orderBy string) ([]*Comment, error) {
	page, size = normalizePage(page, size)
	order := topLevelOrderClause(orderBy)
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE article_id = ? AND parent_id = 0 AND deleted_at IS NULL %s LIMIT ? OFFSET ?",
		commentRows, m.table, order,
	)
	var resp []*Comment
	err := m.conn.QueryRowsCtx(ctx, &resp, query, articleId, size, (page-1)*size)
	return resp, err
}

func (m *defaultCommentModel) CountTopLevelByArticle(ctx context.Context, articleId uint64) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE article_id = ? AND parent_id = 0 AND deleted_at IS NULL", m.table)
	var cnt int64
	err := m.conn.QueryRowCtx(ctx, &cnt, query, articleId)
	return cnt, err
}

func (m *defaultCommentModel) ListRepliesByRoot(ctx context.Context, rootId uint64, page, size int32) ([]*Comment, error) {
	page, size = normalizePage(page, size)
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE root_id = ? AND parent_id > 0 AND deleted_at IS NULL ORDER BY created_at ASC LIMIT ? OFFSET ?",
		commentRows, m.table,
	)
	var resp []*Comment
	err := m.conn.QueryRowsCtx(ctx, &resp, query, rootId, size, (page-1)*size)
	return resp, err
}

func (m *defaultCommentModel) CountRepliesByRoot(ctx context.Context, rootId uint64) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE root_id = ? AND parent_id > 0 AND deleted_at IS NULL", m.table)
	var cnt int64
	err := m.conn.QueryRowCtx(ctx, &cnt, query, rootId)
	return cnt, err
}

func (m *defaultCommentModel) ListPreviewReplies(ctx context.Context, rootId uint64, limit int32) ([]*Comment, error) {
	if limit <= 0 {
		limit = 3
	}
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE root_id = ? AND parent_id > 0 AND deleted_at IS NULL ORDER BY created_at ASC LIMIT ?",
		commentRows, m.table,
	)
	var resp []*Comment
	err := m.conn.QueryRowsCtx(ctx, &resp, query, rootId, limit)
	return resp, err
}

func (m *defaultCommentModel) ListByUser(ctx context.Context, userId uint64, page, size int32) ([]*Comment, error) {
	page, size = normalizePage(page, size)
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE user_id = ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?",
		commentRows, m.table,
	)
	var resp []*Comment
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, size, (page-1)*size)
	return resp, err
}

func (m *defaultCommentModel) CountByUser(ctx context.Context, userId uint64) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE user_id = ? AND deleted_at IS NULL", m.table)
	var cnt int64
	err := m.conn.QueryRowCtx(ctx, &cnt, query, userId)
	return cnt, err
}

const commentLikeActionType = "comment"

func (m *defaultCommentModel) FindCommentLikeActive(ctx context.Context, userId, commentId uint64) (bool, error) {
	query := `SELECT 1 FROM interaction_like WHERE user_id = ? AND article_id = ? AND action_type = ? AND deleted_at IS NULL LIMIT 1`
	var n int
	err := m.conn.QueryRowCtx(ctx, &n, query, userId, commentId, commentLikeActionType)
	if err == sqlx.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *defaultCommentModel) RestoreCommentLike(ctx context.Context, userId, commentId uint64) error {
	query := `UPDATE interaction_like SET deleted_at = NULL WHERE user_id = ? AND article_id = ? AND action_type = ?`
	_, err := m.conn.ExecCtx(ctx, query, userId, commentId, commentLikeActionType)
	return err
}

func (m *defaultCommentModel) InsertCommentLike(ctx context.Context, userId, commentId uint64) error {
	query := `INSERT INTO interaction_like (user_id, article_id, action_type) VALUES (?, ?, ?)`
	_, err := m.conn.ExecCtx(ctx, query, userId, commentId, commentLikeActionType)
	return err
}

func (m *defaultCommentModel) SoftDeleteCommentLike(ctx context.Context, userId, commentId uint64) error {
	query := `UPDATE interaction_like SET deleted_at = ? WHERE user_id = ? AND article_id = ? AND action_type = ? AND deleted_at IS NULL`
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), userId, commentId, commentLikeActionType)
	return err
}

func normalizePage(page, size int32) (int32, int32) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	return page, size
}

func topLevelOrderClause(orderBy string) string {
	switch strings.ToLower(strings.TrimSpace(orderBy)) {
	case "hot":
		return "ORDER BY like_count DESC, created_at DESC"
	default:
		return "ORDER BY created_at DESC"
	}
}
