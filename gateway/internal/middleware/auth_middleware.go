package middleware

import (
	"context"
	"errors"
	"gateway/internal/config"
	"gateway/internal/utils"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthMiddleware struct {
	// 加入密钥字段
	AccessSecret  string
	AccessExpire  uint64
	RefreshSecret string
}

func NewAuthMiddleware(c config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		AccessSecret:  c.Auth.AccessSecret,
		AccessExpire:  c.Auth.AccessExpire,
		RefreshSecret: c.RefreshSecret,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// 1.提取token
		AccessToken := utils.GetAccessTokenFromRequest(r)
		RefreshToken, err := utils.GetRefreshTokenFromRequest(r)
		if err != nil {
			logx.Errorf("Get refresh token err: %v", err)
			return
		}

		// 2.判断是不是在黑名单里

		// 3.解析token
		Jwt := utils.NewJWT(m.AccessSecret, m.RefreshSecret)
		RefreshClaims, err := Jwt.ParseToken(RefreshToken, m.RefreshSecret)
		if err != nil {
			// 判断这个RefreshToken是否有效
			utils.ClearRefreshToken(w)
			logx.Errorf("Refresh token expired or invalid")
			return
		}

		// 4. 解析和验证 AccessToken
		// 判断AccessClaims有没有问题,如果有，就需要更具Refresh刷新
		AccessClaims, err := Jwt.ParseToken(AccessToken, m.AccessSecret)
		if err != nil {
			// 1.过期了
			if errors.Is(err, utils.TokenExpired) || AccessToken == "" {
				logx.Errorf("JWT 过期: %v", err)
				// 刷新令牌
				id := RefreshClaims.Id
				// 这里可以根据需要调整过期时间
				NewAccessToken, _ := Jwt.GetAccessToken(id, m.AccessExpire)

				// 写入响应头（自定义 Header）
				w.Header().Set("Authorization", NewAccessToken)
			} else {
				// 这里就是说，无效或者其它
				httpx.Error(w, err)
				return
			}
		}
		// 5.继续处理请求,将用户信息注入上下文
		ctx = context.WithValue(ctx, "X-user-Id", AccessClaims.Id)

		r = r.WithContext(ctx)
		next(w, r)
	}
}
