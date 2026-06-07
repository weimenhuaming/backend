package agent

import (
	"context"
	"strings"
	"time"

	"gateway/internal/svc"
	"gateway/internal/types"
	agent_client "other-rpc/agent_client"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BuildKnowledgeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBuildKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BuildKnowledgeLogic {
	return &BuildKnowledgeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BuildKnowledgeLogic) BuildKnowledge(req *types.BuildKnowledgeReq) (resp *types.BuildKnowledgeResp, err error) {
	if code, msg, ok := requireAdmin(l.ctx); !ok {
		return &types.BuildKnowledgeResp{Code: code, Msg: msg}, nil
	}

	collection := strings.TrimSpace(req.Collection)
	if collection == "" {
		return &types.BuildKnowledgeResp{Code: 400, Msg: "collection 名称不能为空"}, nil
	}

	buildCtx, cancel := context.WithTimeout(context.WithoutCancel(l.ctx), 2*time.Minute)
	defer cancel()

	r, err := l.svcCtx.Agent.Build(buildCtx, &agent_client.BuildRequest{
		Collection: collection,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				return &types.BuildKnowledgeResp{Code: 409, Msg: st.Message()}, nil
			case codes.InvalidArgument:
				return &types.BuildKnowledgeResp{Code: 400, Msg: st.Message()}, nil
			}
		}
		l.Errorf("build knowledge failed: %v", err)
		return &types.BuildKnowledgeResp{Code: 500, Msg: err.Error()}, nil
	}

	return &types.BuildKnowledgeResp{
		Code: 200,
		Msg:  "ok",
		Data: types.BuildKnowledgeData{
			Message:    r.GetMessage(),
			DocCount:   r.GetDocCount(),
			ChunkCount: r.GetChunkCount(),
		},
	}, nil
}
