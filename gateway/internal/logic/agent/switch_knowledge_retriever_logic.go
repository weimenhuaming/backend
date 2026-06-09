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

type SwitchKnowledgeRetrieverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSwitchKnowledgeRetrieverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchKnowledgeRetrieverLogic {
	return &SwitchKnowledgeRetrieverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SwitchKnowledgeRetrieverLogic) SwitchKnowledgeRetriever(req *types.SwitchKnowledgeRetrieverReq) (resp *types.SwitchKnowledgeRetrieverData, err error) {
	if ok := vaild.IsAdmin(l.ctx); !ok {
		return nil, response.ErrorForbidden("非管理员，无权限操作")
	}

	collection := strings.TrimSpace(req.Collection)
	if collection == "" {
		return nil, response.ErrorBadRequest("collection 名称不能为空")
	}

	r, err := l.svcCtx.Agent.SwitchRetriever(l.ctx, &agent_client.SwitchRetrieverRequest{
		Collection: collection,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, response.ErrorNotFound(st.Message())
		}
		l.Errorf("switch knowledge retriever failed: %v", err)
		return nil, response.ErrorInternalServer(err.Error())
	}

	return &types.SwitchKnowledgeRetrieverData{
		Message: r.GetMessage(),
	}, nil
}
