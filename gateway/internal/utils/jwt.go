package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func generateAccessToken(secret string, userId int64, role string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"userId": userId,
		"role":   role,
		"iat":    now.Unix(),
		"exp":    now.Add(24 * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(secret))
}

type MyClaims struct {
	Id               uint64 `json:"id"`
	AccessExpireTime uint64 `json:"access_expire_time"`
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

// GetRefreshToken 拿到refreshtoken
func (j *JWT) GetRefreshToken(id uint64, as, rs uint64) (string, error) {
	// 1.先生成accesstoken
	AccessClaims := MyClaims{
		Id:               id,
		AccessExpireTime: rs,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"TAP"}, // 受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(as) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims)
	return token.SignedString(j.AccessTokenSecret)
}

// GetAccessToken 拿到accesstoken
func (j *JWT) GetAccessToken(id uint64, rs uint64) (string, error) {
	RefreshClaims := MyClaims{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"TAP"}, // 受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(rs) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims)
	return token.SignedString(j.RefreshTokenSecret)
}
