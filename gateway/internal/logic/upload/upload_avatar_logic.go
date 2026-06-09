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
	role, _ := l.ctx.Value("X-user-Role").(string)
	if role != "admin" {
		return nil, response.NewError(403, "仅管理员可上传头像")
	}

	urlPath, saveErr := utils.SaveUploadedImage(r, "file", "avatars")
	if saveErr != nil {
		return nil, response.NewError(400, saveErr.Error())
	}

	return &types.UploadImageData{Url: urlPath}, nil
}
