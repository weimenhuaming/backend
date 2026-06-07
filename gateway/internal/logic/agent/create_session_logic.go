package agent

import (
	"context"
	"strings"
	"time"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSessionLogic) CreateSession(req *types.CreateSessionReq) (resp *types.CreateSessionResp, err error) {
	sessionId := strings.TrimSpace(req.SessionId)
	if sessionId == "" {
		sessionId = uuid.NewString()
	}

	userId := strings.TrimSpace(req.UserId)
	if userId == "" {
		userId = "guest"
	}

	now := time.Now().Unix()

	return &types.CreateSessionResp{
		Code: 200,
		Msg:  "会话创建成功",
		Data: types.CreateSessionData{
			SessionId: sessionId,
			UserId:    userId,
			CreatedAt: now,
			Message:   "ready",
		},
	}, nil
}
