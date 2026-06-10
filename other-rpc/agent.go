package main

import (
	"flag"
	"fmt"
	"os"
	"other-rpc/internal/config"
	"other-rpc/internal/server"
	"other-rpc/internal/svc"
	"other-rpc/pb/agent"

	"github.com/joho/godotenv"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/agent.yaml", "the config file")

func main() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		os.Exit(1)
	}

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		agent.RegisterAgentServer(grpcServer, server.NewAgentServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
