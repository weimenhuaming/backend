package upload

import (
	"context"
	"net/http"

	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadBlogImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadBlogImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadBlogImageLogic {
	return &UploadBlogImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadBlogImageLogic) UploadBlogImage(r *http.Request) (resp *types.UploadImageData, err error) {
	userId, ok := l.ctx.Value("X-user-Id").(uint64)
	if !ok || userId == 0 {
		return nil, response.NewError(401, "请先登录")
	}

	urlPath, saveErr := utils.SaveUploadedImage(r, "file", "blog")
	if saveErr != nil {
		return nil, response.NewError(400, saveErr.Error())
	}

	return &types.UploadImageData{Url: urlPath}, nil
}
