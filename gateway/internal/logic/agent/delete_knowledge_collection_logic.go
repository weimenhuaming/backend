package agent

import (
	"context"
	"gateway/internal/utils/vaild"
	"strings"

	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	agent_client "other-rpc/agent_client"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteKnowledgeCollectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteKnowledgeCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteKnowledgeCollectionLogic {
	return &DeleteKnowledgeCollectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteKnowledgeCollectionLogic) DeleteKnowledgeCollection(req *types.DeleteKnowledgeCollectionReq) error {
	if ok := vaild.IsAdmin(l.ctx); !ok {
		return response.ErrorForbidden("非管理员，无权限操作")
	}

	collection := strings.TrimSpace(req.Collection)
	if collection == "" {
		return response.ErrorBadRequest("collection 名称不能为空")
	}

	_, err := l.svcCtx.Agent.DeleteCollection(l.ctx, &agent_client.DeleteCollectionRequest{
		Collection: collection,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return response.ErrorNotFound(st.Message())
		}
		l.Errorf("delete knowledge collection failed: %v", err)
		return response.ErrorInternalServer(err.Error())
	}

	return nil
}
