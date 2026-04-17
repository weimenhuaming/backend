package logic

import (
	"context"

	"core-rpc/internal/svc"
	"core-rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type TestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestLogic {
	return &TestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 测试接口
func (l *TestLogic) Test(in *pb.TestRequest) (*pb.TestResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.TestResponse{}, nil
}
