package agent

import (
	"context"

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

func (l *ListKnowledgeCollectionsLogic) ListKnowledgeCollections(req *types.ListKnowledgeCollectionsReq) (resp *types.ListKnowledgeCollectionsResp, err error) {
	_ = req

	if code, msg, ok := requireAdmin(l.ctx); !ok {
		return &types.ListKnowledgeCollectionsResp{Code: code, Msg: msg}, nil
	}

	r, err := l.svcCtx.Agent.ListCollections(l.ctx, &agent_client.ListCollectionsRequest{})
	if err != nil {
		l.Errorf("list knowledge collections failed: %v", err)
		return &types.ListKnowledgeCollectionsResp{Code: 500, Msg: err.Error()}, nil
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

	return &types.ListKnowledgeCollectionsResp{
		Code: 200,
		Msg:  "ok",
		Data: types.ListKnowledgeCollectionsData{
			Collections: collections,
		},
	}, nil
}
