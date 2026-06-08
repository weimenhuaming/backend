package converter

import (
	"core-rpc/internal/model/entity"
	"core-rpc/pb/core"
)

func UserToLoginResp(u *entity.User) *core.LoginResp {
	return &core.LoginResp{
		Id:     u.ID,
		Name:   u.Name,
		Phone:  u.Phone,
		Email:  u.Email,
		Avatar: u.Avatar,
		Role:   u.Role,
		Sex:    u.Sex,
		Age:    u.Age,
	}
}

func UserToProfile(u *entity.User) *core.UserProfile {
	if u == nil {
		return nil
	}
	return &core.UserProfile{
		Id:     u.ID,
		Name:   u.Name,
		Phone:  u.Phone,
		Email:  u.Email,
		Avatar: u.Avatar,
		Role:   u.Role,
		Sex:    u.Sex,
		Age:    u.Age,
	}
}
