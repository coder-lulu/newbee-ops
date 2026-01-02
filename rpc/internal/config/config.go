package config

import (
	"github.com/coder-lulu/newbee-common/v2/config"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DatabaseConf config.DatabaseConf
	RedisConf    config.RedisConf
	CoreRpc      zrpc.RpcClientConf `json:",optional"` // Core服务RPC配置
}

