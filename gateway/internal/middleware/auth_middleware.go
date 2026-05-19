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
	// 加入密钥字段
	AccessSecret  string
	AccessExpire  uint64
	RefreshSecret string
	Cache         *redis.Redis
}

func NewAuthMiddleware(c config.Config, Cache *redis.Redis) *AuthMiddleware {
	return &AuthMiddleware{
		AccessSecret:  c.Auth.AccessSecret,
		AccessExpire:  c.Auth.AccessExpire,
		RefreshSecret: c.RefreshSecret,
		Cache:         Cache,
	}
}

// 定义统一的错误响应结构
type ErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 返回 JSON 错误的辅助函数
func jsonError(w http.ResponseWriter, statusCode int, msg string) {
	httpx.WriteJson(w, statusCode, ErrorResponse{
		Code: statusCode,
		Msg:  msg,
	})
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// 1. 提取 token
		AccessToken := utils.GetAccessTokenFromRequest(r)
		RefreshToken, err := utils.GetRefreshTokenFromRequest(r)
		if err != nil {
			logx.Errorf("Get refresh token err: %v", err)
			jsonError(w, http.StatusUnauthorized, "refresh token 获取失败")
			return
		}

		// 2. 判断是不是在黑名单里
		v, err := m.Cache.Get("blacklist:" + RefreshToken)
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

		// 3. 解析 RefreshToken
		Jwt := utils.NewJWT(m.AccessSecret, m.RefreshSecret)
		RefreshClaims, err := Jwt.ParseToken(RefreshToken, m.RefreshSecret)
		if err != nil {
			utils.ClearRefreshToken(w)
			logx.Errorf("Refresh token expired or invalid: %v", err)
			jsonError(w, http.StatusUnauthorized, "refresh token 无效或已过期，请重新登录")
			return
		}

		// 4. 解析和验证 AccessToken
		AccessClaims, err := Jwt.ParseToken(AccessToken, m.AccessSecret)
		if err != nil {
			// Access Token 过期了，自动续期
			if errors.Is(err, utils.TokenExpired) || AccessToken == "" {
				logx.Infof("Access Token 过期，使用 Refresh Token 续期")

				// 刷新 Access Token
				id := RefreshClaims.Id
				NewAccessToken, _ := Jwt.GetAccessToken(id, m.AccessExpire)

				// 写入响应头
				w.Header().Set("Authorization", NewAccessToken)

				// 重新生成 claims 用于后续注入
				AccessClaims, _ = Jwt.ParseToken(NewAccessToken, m.AccessSecret)
			} else {
				// 其他错误（无效、格式错误等）
				logx.Errorf("Access Token 无效: %v", err)
				jsonError(w, http.StatusUnauthorized, "access token 无效")
				return
			}
		}

		// 5. 将用户信息注入上下文
		ctx = context.WithValue(ctx, "X-user-Id", AccessClaims.Id)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
