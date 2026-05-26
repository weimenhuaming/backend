package interaction

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ InteractionCommentModel = (*customInteractionCommentModel)(nil)

type (
	// InteractionCommentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInteractionCommentModel.
	InteractionCommentModel interface {
		interactionCommentModel
		withSession(session sqlx.Session) InteractionCommentModel
	}

	customInteractionCommentModel struct {
		*defaultInteractionCommentModel
	}
)

// NewInteractionCommentModel returns a model for the database table.
func NewInteractionCommentModel(conn sqlx.SqlConn) InteractionCommentModel {
	return &customInteractionCommentModel{
		defaultInteractionCommentModel: newInteractionCommentModel(conn),
	}
}

func (m *customInteractionCommentModel) withSession(session sqlx.Session) InteractionCommentModel {
	return NewInteractionCommentModel(sqlx.NewSqlConnFromSession(session))
}
