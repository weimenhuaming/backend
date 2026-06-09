package login

import (
	"context"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/mojocn/base64Captcha"
	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaLogic {
	return &CaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaLogic) Captcha() (resp *types.CaptchaData, err error) {
	// 高度、宽度、位数、干扰、字体大小
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)

	id, b64s, _, err := captcha.Generate()
	if err != nil {
		return nil, err
	}

	return &types.CaptchaData{
		CaptchaId: id,
		PicBase64: b64s,
	}, nil
}
