//	ops-center
//
//	Description: Ops Center service
//
// swagger:meta
package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/coder-lulu/newbee-common/v2/middleware/integration"
	"github.com/coder-lulu/newbee-ops-api/internal/config"
	"github.com/coder-lulu/newbee-ops-api/internal/handler"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/ops.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	server := rest.MustNewServer(c.RestConf, rest.WithCors(c.CROSConf.Address))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	defer func() {
		if ctx.IntegrationResult != nil && ctx.IntegrationResult.Manager != nil {
			if err := ctx.IntegrationResult.Manager.Shutdown(); err != nil {
				logx.Errorf("failed to shutdown middleware manager: %v", err)
			}
		}
	}()

	// 统一中间件链
	integration.ApplyToServer(server, ctx.IntegrationResult)

	handler.RegisterHandlers(server, ctx)

	// metrics endpoint
	server.AddRoute(rest.Route{Method: http.MethodGet, Path: "/metrics", Handler: func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	}})

	fmt.Printf("Starting Ops Center at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// registerTaskHandlers is a thin wrapper to avoid import cycles in main.
// 保留空壳（兼容历史调用形式）
func registerTaskHandlers(server *rest.Server, ctx *svc.ServiceContext) {}
