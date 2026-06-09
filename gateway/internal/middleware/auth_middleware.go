package middleware

import (
	"context"
	"errors"
	"gateway/internal/config"
	"gateway/internal/utils"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthMiddleware struct {
	AccessSecret  string
	AccessExpire  uint64
	RefreshSecret string
	Cache         *redis.Redis
}

func NewAuthMiddleware(c config.Config, cache *redis.Redis) *AuthMiddleware {
	return &AuthMiddleware{
		AccessSecret:  c.Auth.AccessSecret,
		AccessExpire:  c.Auth.AccessExpire,
		RefreshSecret: c.RefreshSecret,
		Cache:         cache,
	}
}

type ErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func jsonError(w http.ResponseWriter, statusCode int, msg string) {
	httpx.WriteJson(w, statusCode, ErrorResponse{
		Code: statusCode,
		Msg:  msg,
	})
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		accessToken := utils.GetAccessTokenFromRequest(r)
		refreshToken, err := utils.GetRefreshTokenFromRequest(r)
		if err != nil {
			logx.Errorf("Get refresh token err: %v", err)
			jsonError(w, http.StatusUnauthorized, "refresh token 获取失败")
			return
		}

		v, err := m.Cache.Get("blacklist:" + refreshToken)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "服务器内部错误")
			return
		}
		if v == "1" {
			logx.Errorf("RefreshToken 已被拉入黑名单")
			utils.ClearRefreshToken(w)
			jsonError(w, http.StatusUnauthorized, "token 已失效，请重新登录")
			return
		}

		jwt := utils.NewJWT(m.AccessSecret, m.RefreshSecret)
		refreshClaims, err := jwt.ParseToken(refreshToken, m.RefreshSecret)
		if err != nil {
			utils.ClearRefreshToken(w)
			logx.Errorf("Refresh token expired or invalid: %v", err)
			jsonError(w, http.StatusUnauthorized, "refresh token 无效或已过期，请重新登录")
			return
		}

		accessClaims, err := jwt.ParseToken(accessToken, m.AccessSecret)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) || accessToken == "" {
				logx.Infof("Access Token 过期，使用 Refresh Token 续期")

				newAccessToken, _ := jwt.GetAccessToken(refreshClaims.Id, refreshClaims.Role, m.AccessExpire)
				w.Header().Set("Authorization", newAccessToken)
				accessClaims, _ = jwt.ParseToken(newAccessToken, m.AccessSecret)
			} else {
				logx.Errorf("Access Token 无效: %v", err)
				jsonError(w, http.StatusUnauthorized, "access token 无效")
				return
			}
		}

		ctx = context.WithValue(ctx, "X-user-Id", accessClaims.Id)
		ctx = context.WithValue(ctx, "X-user-Role", accessClaims.Role)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
