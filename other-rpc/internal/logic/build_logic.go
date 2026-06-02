package logic

import (
	"context"
	"fmt"

	"other-rpc/internal/svc"
	"other-rpc/pb/agent"

	"github.com/zeromicro/go-zero/core/logx"
)

type BuildLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBuildLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BuildLogic {
	return &BuildLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Build 从知识库目录构建向量索引并持久化（启动前或更新文档后调用）。
func (l *BuildLogic) Build(in *agent.BuildRequest) (*agent.BuildResponse, error) {
	_ = in

	docCount, chunkCount, err := l.svcCtx.Agent.Build(l.ctx)
	if err != nil {
		l.Errorf("构建向量索引失败: %v", err)
		return nil, err
	}

	msg := fmt.Sprintf("向量索引已构建并保存，文档数: %d, 切片数: %d", docCount, chunkCount)
	l.Infof(msg)
	return &agent.BuildResponse{
		Message:    msg,
		DocCount:   int32(docCount),
		ChunkCount: int32(chunkCount),
	}, nil
}
