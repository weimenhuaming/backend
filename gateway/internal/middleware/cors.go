package middleware

import (
	"net/http"
)

// CorsMiddleware 创建一个CORS中间件，处理跨域请求
// 允许的前端源为 http://localhost:3000
func CorsMiddleware() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 获取请求来源
			origin := r.Header.Get("Origin")

			// 只允许指定的前端源跨域请求
			if origin == "http://localhost:3000" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			// 告诉浏览器Vary: Origin 用于缓存优化，告诉缓存服务器根据Origin不同而缓存不同版本
			w.Header().Set("Vary", "Origin")

			// 允许携带Cookie等凭证信息
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// 允许前端发送的请求头
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

			// 允许的HTTP方法
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			// 允许前端读取的响应头，特别是Authorization（用于传递token）
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")

			// 处理预检请求（OPTIONS方法）
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// 继续处理下一个中间件或处理器
			next(w, r)
		}
	}
}
