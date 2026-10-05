package main

import (
	"flag"
	"fmt"
	"github.com/coder-lulu/newbee-common/v2/orm/ent/hooks"
	"github.com/coder-lulu/newbee-ops-rpc/internal/config"
	"github.com/coder-lulu/newbee-ops-rpc/internal/server"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/ops.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		newbee_ops_rpc.RegisterOpsServer(grpcServer, server.NewOpsServer(ctx))
	})

	// 添加租户上下文传播拦截器
	s.AddUnaryInterceptors(hooks.ContextPropagationServerInterceptor())
	s.AddStreamInterceptors(hooks.ContextPropagationStreamServerInterceptor())

	defer s.Stop()

	fmt.Printf("Starting ops rpc server at %s...\n", c.ListenOn)
	s.Start()
}
