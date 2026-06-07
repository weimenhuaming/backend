package agent

import (
	"context"
	"strings"

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

func (l *SwitchKnowledgeRetrieverLogic) SwitchKnowledgeRetriever(req *types.SwitchKnowledgeRetrieverReq) (resp *types.SwitchKnowledgeRetrieverResp, err error) {
	if code, msg, ok := requireAdmin(l.ctx); !ok {
		return &types.SwitchKnowledgeRetrieverResp{Code: code, Msg: msg}, nil
	}

	collection := strings.TrimSpace(req.Collection)
	if collection == "" {
		return &types.SwitchKnowledgeRetrieverResp{Code: 400, Msg: "collection 名称不能为空"}, nil
	}

	r, err := l.svcCtx.Agent.SwitchRetriever(l.ctx, &agent_client.SwitchRetrieverRequest{
		Collection: collection,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return &types.SwitchKnowledgeRetrieverResp{Code: 404, Msg: st.Message()}, nil
		}
		l.Errorf("switch knowledge retriever failed: %v", err)
		return &types.SwitchKnowledgeRetrieverResp{Code: 500, Msg: err.Error()}, nil
	}

	return &types.SwitchKnowledgeRetrieverResp{
		Code: 200,
		Msg:  "ok",
		Data: types.SwitchKnowledgeRetrieverData{
			Message: r.GetMessage(),
		},
	}, nil
}
