package upload

import (
	"context"
	"net/http"

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

func (l *UploadAvatarLogic) UploadAvatar(r *http.Request) (resp *types.UploadImageResp, err error) {
	role, _ := l.ctx.Value("X-user-Role").(string)
	if role != "admin" {
		return &types.UploadImageResp{Code: 403, Msg: "仅管理员可上传头像"}, nil
	}

	urlPath, saveErr := utils.SaveUploadedImage(r, "file", "avatars")
	if saveErr != nil {
		return &types.UploadImageResp{Code: 400, Msg: saveErr.Error()}, nil
	}

	return &types.UploadImageResp{
		Code: 200,
		Msg:  "上传成功",
		Data: types.UploadImageData{Url: urlPath},
	}, nil
}
