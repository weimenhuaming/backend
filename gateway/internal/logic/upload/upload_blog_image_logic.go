package upload

import (
	"context"
	"net/http"

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

func (l *UploadBlogImageLogic) UploadBlogImage(r *http.Request) (resp *types.UploadImageResp, err error) {
	userId, ok := l.ctx.Value("X-user-Id").(uint64)
	if !ok || userId == 0 {
		return &types.UploadImageResp{Code: 401, Msg: "请先登录"}, nil
	}

	urlPath, saveErr := utils.SaveUploadedImage(r, "file", "blog")
	if saveErr != nil {
		return &types.UploadImageResp{Code: 400, Msg: saveErr.Error()}, nil
	}

	return &types.UploadImageResp{
		Code: 200,
		Msg:  "上传成功",
		Data: types.UploadImageData{Url: urlPath},
	}, nil
}
