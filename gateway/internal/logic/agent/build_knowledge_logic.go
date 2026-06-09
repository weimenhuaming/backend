package agent

import (
	"context"
	"strings"
	"time"

	"gateway/internal/response"
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

func (l *BuildKnowledgeLogic) BuildKnowledge(req *types.BuildKnowledgeReq) (resp *types.BuildKnowledgeData, err error) {
	if code, msg, ok := requireAdmin(l.ctx); !ok {
		return nil, response.NewError(code, msg)
	}

	collection := strings.TrimSpace(req.Collection)
	if collection == "" {
		return nil, response.NewError(400, "collection 名称不能为空")
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
				return nil, response.NewError(409, st.Message())
			case codes.InvalidArgument:
				return nil, response.NewError(400, st.Message())
			}
		}
		l.Errorf("build knowledge failed: %v", err)
		return nil, response.NewError(500, err.Error())
	}

	return &types.BuildKnowledgeData{
		Message:    r.GetMessage(),
		DocCount:   r.GetDocCount(),
		ChunkCount: r.GetChunkCount(),
	}, nil
}
