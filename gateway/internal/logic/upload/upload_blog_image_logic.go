package upload

import (
	"context"
	"net/http"

	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils"
	"gateway/internal/utils/vaild"

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
	_, ok, msg := vaild.GetAdminUserID(l.ctx)
	if !ok {
		return nil, response.ErrorAdminAuth(msg)
	}

	urlPath, saveErr := utils.SaveUploadedImage(r, "file", "blog")
	if saveErr != nil {
		return nil, response.ErrorBadRequest(saveErr.Error())
	}

	return &types.UploadImageData{Url: urlPath}, nil
}
