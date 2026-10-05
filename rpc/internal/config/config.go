package config

import (
	"github.com/coder-lulu/newbee-common/v2/config"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DatabaseConf  config.DatabaseConf
	RedisConf     config.RedisConf
	CoreRpc       zrpc.RpcClientConf `json:",optional"` // Core服务RPC配置
	WorkerManager WorkerManagerConf  `json:",optional"`
}

type WorkerManagerConf struct {
	Enabled *bool `json:",optional"`
}

// IsEnabled preserves automatic startup when older configurations omit this section.
func (c WorkerManagerConf) IsEnabled() bool { return c.Enabled == nil || *c.Enabled }
