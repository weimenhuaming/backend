package article

import (
	"context"
	"core-rpc/internal/model/converter"
	"errors"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetArticleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleDetailLogic {
	return &GetArticleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetArticleDetailLogic) GetArticleDetail(in *core.GetArticleDetailReq) (*core.GetArticleDetailResp, error) {
	a, err := l.svcCtx.ArtRepo.FindByID(in.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	// fetch author info
	authorName, authorAvatar := "", ""
	if u, err := l.svcCtx.UserRepo.FindByID(a.UserID); err == nil {
		authorName = u.Name
		authorAvatar = u.Avatar
	}

	return &core.GetArticleDetailResp{
		Article: converter.ArticleToProto(a, authorName, authorAvatar),
	}, nil
}
