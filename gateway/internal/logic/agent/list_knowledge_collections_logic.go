package agent

import (
	"context"
	"gateway/internal/utils/vaild"

	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	agent_client "other-rpc/agent_client"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListKnowledgeCollectionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListKnowledgeCollectionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListKnowledgeCollectionsLogic {
	return &ListKnowledgeCollectionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListKnowledgeCollectionsLogic) ListKnowledgeCollections(req *types.ListKnowledgeCollectionsReq) (resp *types.ListKnowledgeCollectionsData, err error) {
	_ = req

	if _, ok, msg := vaild.GetAdminUserID(l.ctx); !ok {
		return nil, response.ErrorAdminAuth(msg)
	}

	r, err := l.svcCtx.Agent.ListCollections(l.ctx, &agent_client.ListCollectionsRequest{})
	if err != nil {
		l.Errorf("list knowledge collections failed: %v", err)
		return nil, response.ErrorInternalServer(err.Error())
	}

	collections := make([]types.KnowledgeCollectionInfo, 0, len(r.GetCollections()))
	for _, item := range r.GetCollections() {
		collections = append(collections, types.KnowledgeCollectionInfo{
			Name:       item.GetName(),
			DocCount:   item.GetDocCount(),
			ChunkCount: item.GetChunkCount(),
			Count:      item.GetCount(),
		})
	}

	return &types.ListKnowledgeCollectionsData{
		Collections: collections,
	}, nil
}
