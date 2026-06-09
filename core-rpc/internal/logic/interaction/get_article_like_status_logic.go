package interaction

import (
	"context"
	"errors"

	"core-rpc/internal/model/entity"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleLikeStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleLikeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleLikeStatusLogic {
	return &GetArticleLikeStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetArticleLikeStatusLogic) GetArticleLikeStatus(in *core.GetArticleLikeStatusReq) (*core.GetArticleLikeStatusResp, error) {
	if in.ArticleId == 0 {
		return nil, errors.New("参数无效")
	}

	liked, err := l.svcCtx.InteractionRepo.IsObjectLiked(l.svcCtx.Db, in.UserId, entity.ObjectTypeArticle, in.ArticleId)
	if err != nil {
		return nil, err
	}
	return &core.GetArticleLikeStatusResp{Liked: liked}, nil
}
