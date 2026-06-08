package user

import (
	"context"
	"core-rpc/internal/model/converter"
	"errors"
	"strings"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateUserProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserProfileLogic {
	return &UpdateUserProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserProfileLogic) UpdateUserProfile(in *core.UpdateUserProfileReq) (*core.UpdateUserProfileResp, error) {
	if in.GetUserId() == 0 {
		return nil, errors.New("用户 ID 无效")
	}

	updates := make(map[string]interface{})
	if name := strings.TrimSpace(in.GetName()); name != "" {
		updates["name"] = name
	}
	if phone := strings.TrimSpace(in.GetPhone()); phone != "" {
		updates["phone"] = phone
	}
	if sex := strings.TrimSpace(in.GetSex()); sex != "" {
		switch sex {
		case "男", "女", "未知":
			updates["sex"] = sex
		default:
			return nil, errors.New("性别仅支持：男、女、未知")
		}
	}
	if in.GetAge() > 0 {
		updates["age"] = in.GetAge()
	}
	if avatar := strings.TrimSpace(in.GetAvatar()); avatar != "" {
		updates["avatar"] = avatar
	}

	if len(updates) == 0 {
		return nil, errors.New("没有可更新的字段")
	}

	if err := l.svcCtx.UserRepo.UpdateProfile(in.GetUserId(), updates); err != nil {
		return nil, err
	}

	user, err := l.svcCtx.UserRepo.FindProfileByID(in.GetUserId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &core.UpdateUserProfileResp{
		Profile: converter.UserToProfile(user),
	}, nil
}
