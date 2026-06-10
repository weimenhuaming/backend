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

type UploadAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAvatarLogic {
	return &UploadAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadAvatarLogic) UploadAvatar(r *http.Request) (resp *types.UploadImageData, err error) {
	_, ok := vaild.GetUserID(l.ctx)
	if !ok {
		return nil, response.ErrorUnauthorized("请先登录")
	}

	if !vaild.IsAdmin(l.ctx) {
		return nil, response.ErrorForbidden("仅管理员可上传头像")
	}

	urlPath, saveErr := utils.SaveUploadedImage(r, "file", "avatars")
	if saveErr != nil {
		return nil, response.ErrorBadRequest(saveErr.Error())
	}

	return &types.UploadImageData{Url: urlPath}, nil
}
