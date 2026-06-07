package logic

import (
	"context"
	"errors"
	"fmt"

	"other-rpc/internal/agent/vector"
	"other-rpc/internal/svc"
	"other-rpc/pb/agent"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// Build 从知识库目录构建向量索引并持久化（名称重复时返回已存在）。
func (l *BuildLogic) Build(in *agent.BuildRequest) (*agent.BuildResponse, error) {
	name := in.GetCollection()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "collection 名称不能为空")
	}

	_, docCount, chunkCount, err := l.svcCtx.Chroma.Build(l.ctx, name, l.svcCtx.Embedder)
	if err != nil {
		if errors.Is(err, vector.ErrCollectionExists) {
			return nil, status.Errorf(codes.AlreadyExists, "collection %q 已存在", name)
		}
		l.Errorf("构建向量索引失败: %v", err)
		return nil, err
	}

	msg := fmt.Sprintf("知识库 %q 已构建，文档数: %d, 切片数: %d", name, docCount, chunkCount)
	l.Infof(msg)
	return &agent.BuildResponse{
		Message:    msg,
		DocCount:   int32(docCount),
		ChunkCount: int32(chunkCount),
	}, nil
}
