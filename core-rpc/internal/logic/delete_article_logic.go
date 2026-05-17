package logic

import (
	"context"
	"errors"
	"strconv"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type DeleteArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleLogic {
	return &DeleteArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteArticleLogic) DeleteArticle(in *core.DeleteArticleReq) (*core.DeleteArticleResp, error) {
	md, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}
	var uidStr string
	if v := md.Get("user-id"); len(v) > 0 {
		uidStr = v[0]
	} else if v := md.Get("user_id"); len(v) > 0 {
		uidStr = v[0]
	} else {
		return nil, errors.New("missing user id in metadata")
	}
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		return nil, err
	}

	a, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	if a.UserId != uid {
		return nil, errors.New("not article owner")
	}

	err = l.svcCtx.ArticleModel.SoftDelete(l.ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &core.DeleteArticleResp{}, nil
}
