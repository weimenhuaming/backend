package logic

import (
	"context"
	"fmt"

	"other-rpc/internal/svc"
	"other-rpc/pb/agent"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SwitchRetrieverLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSwitchRetrieverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchRetrieverLogic {
	return &SwitchRetrieverLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SwitchRetrieverLogic) SwitchRetriever(in *agent.SwitchRetrieverRequest) (*agent.SwitchRetrieverResponse, error) {
	name := in.GetCollection()

	retriever, err := l.svcCtx.Chroma.OpenRetriever(l.ctx, name, l.svcCtx.Embedder)
	if err != nil {
		l.Errorf("切换检索器失败: %v", err)
		return nil, status.Errorf(codes.NotFound, "%v", err)
	}

	l.svcCtx.Agent.SetRetriever(retriever)

	msg := fmt.Sprintf("已切换检索器到 collection %q", name)
	l.Infof(msg)
	return &agent.SwitchRetrieverResponse{Message: msg}, nil
}
