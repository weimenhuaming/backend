package main

import (
	"flag"
	"fmt"
	"gateway/internal/middleware"
	"net/http"
	"path/filepath"

	"gateway/internal/config"
	"gateway/internal/handler"
	"gateway/internal/svc"
	"gateway/internal/utils"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	staticDir := resolveStaticDir(*configFile)
	utils.InitStaticRoot(staticDir)

	server := rest.MustNewServer(
		c.RestConf,
		rest.WithFileServer("/static", http.Dir(staticDir)),
	)
	// 注册中间件（在 RegisterHandlers 之前）
	server.Use(middleware.CorsMiddleware())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d..., static dir: %s\n", c.Host, c.Port, staticDir)
	server.Start()
}

func resolveStaticDir(configFile string) string {
	absConfig, err := filepath.Abs(configFile)
	if err != nil {
		return "static"
	}
	return filepath.Clean(filepath.Join(filepath.Dir(absConfig), "..", "static"))
}
