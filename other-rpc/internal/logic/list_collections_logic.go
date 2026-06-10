package logic

import (
	"context"

	"other-rpc/internal/svc"
	"other-rpc/pb/agent"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCollectionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCollectionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCollectionsLogic {
	return &ListCollectionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCollectionsLogic) ListCollections(in *agent.ListCollectionsRequest) (*agent.ListCollectionsResponse, error) {
	collections, err := l.svcCtx.Chroma.ListCollections(l.ctx)
	if err != nil {
		l.Errorf("获取 collection 列表失败: %v", err)
		return nil, err
	}

	out := make([]*agent.CollectionInfo, 0, len(collections))
	for _, item := range collections {
		out = append(out, &agent.CollectionInfo{
			Name:       item.Name,
			DocCount:   int32(item.DocCount),
			ChunkCount: int32(item.ChunkCount),
			Count:      int32(item.Count),
		})
	}

	return &agent.ListCollectionsResponse{Collections: out}, nil
}
