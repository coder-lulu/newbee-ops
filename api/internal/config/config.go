package config

import (
	commoncfg "github.com/coder-lulu/newbee-common/v2/config"
	"github.com/coder-lulu/newbee-common/v2/i18n"
	"github.com/coder-lulu/newbee-common/v2/middleware/framework"
	"github.com/coder-lulu/newbee-common/v2/plugins/casbin"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Middleware   framework.UnifiedConfig `json:",optional"`
	RedisConf    commoncfg.RedisConf
	DatabaseConf commoncfg.DatabaseConf // 添加数据库配置
	CoreRpc      zrpc.RpcClientConf
	OpsRpc       zrpc.RpcClientConf
	I18nConf     i18n.Conf
	CROSConf     commoncfg.CROSConf
	CasbinConf   casbin.CasbinConf

	Ops   OpsConf   `json:",optional"`
	Agent AgentConf `json:",optional"` // nb-agent 对接配置
}

type OpsConf struct {
	DefaultProxyEndpoints []string         `json:",optional"`              // e.g., ["http://127.0.0.1:8889"]
	SessionTokenTTL       int64            `json:",optional,default=1800"` // seconds
	Registration          RegistrationConf `json:",optional"`
}

type RegistrationConf struct {
	TenantID uint64 `json:",optional"` // Explicit tenant authorized by this PSK; zero disables registration.
	PSK      string `json:",optional"` // Pre-shared key for proxy register/heartbeat
}

// AgentConf nb-agent 接口调用配置
type AgentConf struct {
	Endpoints []string `json:",optional"` // e.g., ["http://127.0.0.1:8090"]
	TimeoutMS int      `json:",optional"` // 默认 5000
	Token     string   `json:",optional"` // 可选：鉴权 Token
}
