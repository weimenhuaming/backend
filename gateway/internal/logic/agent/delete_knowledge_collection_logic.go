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

type DeleteKnowledgeCollectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteKnowledgeCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteKnowledgeCollectionLogic {
	return &DeleteKnowledgeCollectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteKnowledgeCollectionLogic) DeleteKnowledgeCollection(req *types.DeleteKnowledgeCollectionReq) (resp *types.DeleteKnowledgeCollectionResp, err error) {
	if code, msg, ok := requireAdmin(l.ctx); !ok {
		return &types.DeleteKnowledgeCollectionResp{Code: code, Msg: msg}, nil
	}

	collection := strings.TrimSpace(req.Collection)
	if collection == "" {
		return &types.DeleteKnowledgeCollectionResp{Code: 400, Msg: "collection 名称不能为空"}, nil
	}

	_, err = l.svcCtx.Agent.DeleteCollection(l.ctx, &agent_client.DeleteCollectionRequest{
		Collection: collection,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return &types.DeleteKnowledgeCollectionResp{Code: 404, Msg: st.Message()}, nil
		}
		l.Errorf("delete knowledge collection failed: %v", err)
		return &types.DeleteKnowledgeCollectionResp{Code: 500, Msg: err.Error()}, nil
	}

	return &types.DeleteKnowledgeCollectionResp{
		Code: 200,
		Msg:  "collection 已删除",
	}, nil
}
