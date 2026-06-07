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

type DeleteCollectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCollectionLogic {
	return &DeleteCollectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCollectionLogic) DeleteCollection(in *agent.DeleteCollectionRequest) (*agent.DeleteCollectionResponse, error) {
	name := in.GetCollection()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "collection 名称不能为空")
	}

	if err := l.svcCtx.Chroma.DeleteCollection(l.ctx, name); err != nil {
		if errors.Is(err, vector.ErrCollectionNotFound) {
			return nil, status.Errorf(codes.NotFound, "collection %q 不存在", name)
		}
		l.Errorf("删除 collection 失败: %v", err)
		return nil, err
	}

	msg := fmt.Sprintf("collection %q 已删除", name)
	l.Infof(msg)
	return &agent.DeleteCollectionResponse{Message: msg}, nil
}
