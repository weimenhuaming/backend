package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type MyClaims struct {
	Id uint64 `json:"id"`
	jwt.RegisteredClaims
}

type JWT struct {
	AccessTokenSecret  []byte
	RefreshTokenSecret []byte
}

func NewJWT(ats, rts string) *JWT {
	return &JWT{
		AccessTokenSecret:  []byte(ats),
		RefreshTokenSecret: []byte(rts),
	}
}

// GetAccessToken 用 AccessTokenSecret 签发 access token，expireSec 单位为秒
func (j *JWT) GetAccessToken(id uint64, expireSec uint64) (string, error) {
	claims := MyClaims{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"TAP"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireSec) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.AccessTokenSecret)
}

// GetRefreshToken 用 RefreshTokenSecret 签发 refresh token，expireSec 单位为秒
func (j *JWT) GetRefreshToken(id uint64, expireSec uint64) (string, error) {
	claims := MyClaims{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"TAP"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireSec) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.RefreshTokenSecret)
}
