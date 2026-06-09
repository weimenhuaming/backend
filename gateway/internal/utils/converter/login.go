package converter

import (
	"core-rpc/core_client"
	"gateway/internal/types"
)

func ToLoginData(u *core_client.LoginResp, access, refresh string) *types.LoginData {
	return &types.LoginData{
		Id:           u.Id,
		Name:         u.Name,
		Phone:        u.Phone,
		Email:        u.Email,
		Avatar:       u.Avatar,
		Uuid:         u.Uuid,
		Role:         u.Role,
		Sex:          u.Sex,
		Age:          u.Age,
		AccessToken:  access,
		RefreshToken: refresh,
	}
}
