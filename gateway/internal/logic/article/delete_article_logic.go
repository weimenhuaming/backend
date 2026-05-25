package article

import (
	"context"
	"core-rpc/core_client"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleLogic {
	return &DeleteArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteArticleLogic) DeleteArticle(req *types.DeleteArticleReq) (resp *types.DeleteArticleResp, err error) {
	// extract user id from context (set by auth middleware)
	uidVal := l.ctx.Value("X-user-Id")
	if uidVal == nil {
		return &types.DeleteArticleResp{
			Code: 401,
			Msg:  "missing user id",
		}, nil
	}

	uid, ok := uidVal.(uint64)
	if !ok {
		// try if it was stored as int64 or int
		switch v := uidVal.(type) {
		case int:
			uid = uint64(v)
			ok = true
		case int64:
			uid = uint64(v)
			ok = true
		case float64:
			uid = uint64(v)
			ok = true
		default:
			ok = false
		}
	}
	if !ok {
		return &types.DeleteArticleResp{
			Code: 401,
			Msg:  "invalid user id",
		}, nil
	}

	// call core rpc and pass user id explicitly
	_, err = l.svcCtx.Core.DeleteArticle(l.ctx, &core_client.DeleteArticleReq{Id: req.Id, UserId: uid})
	if err != nil {
		return &types.DeleteArticleResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &types.DeleteArticleResp{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}
