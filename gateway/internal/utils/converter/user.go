package converter

import (
	core_client "core-rpc/core_client"
	"gateway/internal/types"
)

func ToUserProfile(p *core_client.UserProfile) *types.UserProfile {
	if p == nil {
		return &types.UserProfile{}
	}
	return &types.UserProfile{
		Id:     p.GetId(),
		Name:   p.GetName(),
		Phone:  p.GetPhone(),
		Email:  p.GetEmail(),
		Avatar: p.GetAvatar(),
		Role:   p.GetRole(),
		Sex:    p.GetSex(),
		Age:    p.GetAge(),
	}
}
