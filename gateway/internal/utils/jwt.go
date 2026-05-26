package utils

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	TokenExpired     = errors.New("Token 已过期")
	TokenNotValidYet = errors.New("Token 不可用")
	TokenMalformed   = errors.New("Token 格式错误")
	TokenInvalid     = errors.New("Token 无效")
)

type MyClaims struct {
	Id   uint64 `json:"id"`
	Role string `json:"role"`
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
func (j *JWT) GetAccessToken(userId uint64, role string, expireSec uint64) (string, error) {
	claims := MyClaims{
		Id:   userId,
		Role: role,
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
func (j *JWT) GetRefreshToken(userId uint64, role string, expireSec uint64) (string, error) {
	claims := MyClaims{
		Id:   userId,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"TAP"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireSec) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.RefreshTokenSecret)
}

// ParseToken 解析token
func (j *JWT) ParseToken(tokenStr string, secretKey string) (*MyClaims, error) {
	//tokenString := strings.TrimPrefix(tokenStr, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		// 处理具体错误类型
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, TokenMalformed
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, TokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, TokenNotValidYet
		}
		return nil, TokenInvalid
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}

// GetRefreshTokenFromRequest 从请求中获取 refresh_token（标准 net/http）
func GetRefreshTokenFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return "", errors.New("refresh_token not found in cookies")
	}
	return cookie.Value, nil
}

// GetAccessTokenFromRequest 从请求中拿到AccessToken
func GetAccessTokenFromRequest(r *http.Request) string {
	return r.Header.Get("Authorization")
}

// ClearRefreshToken 设置一个立即过期的 Cookie 来清除 refresh_token
func ClearRefreshToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0), // 设置一个已经过去的时间
		MaxAge:   -1,              // 强制删除（有些浏览器只识别 MaxAge）
		SameSite: http.SameSiteLaxMode,
		Secure:   true, // 生产环境建议保留
	})
}
