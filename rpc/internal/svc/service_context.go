package svc

import (
	"github.com/coder-lulu/newbee-core/rpc/coreclient"
	"github.com/coder-lulu/newbee-ops-rpc/ent"
	_ "github.com/coder-lulu/newbee-ops-rpc/ent/runtime"
	"github.com/coder-lulu/newbee-ops-rpc/internal/config"
	// TODO: "github.com/coder-lulu/newbee-ops-rpc/internal/dispatcher" - 需要更新为Proxy架构
	"github.com/coder-lulu/newbee-ops-rpc/internal/logic/worker"
	// TODO: "github.com/coder-lulu/newbee-ops-rpc/internal/services/executor" - 需要更新为不依赖dispatcher
	"github.com/coder-lulu/newbee-ops-rpc/internal/services/template"

	"github.com/coder-lulu/newbee-common/v2/orm/ent/hooks"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	DB             *ent.Client
	Redis          redis.UniversalClient
	CoreRpc        coreclient.Core       // Core服务RPC客户端
	WorkerManager  *worker.WorkerManager // Worker管理器
	// TODO: TaskDispatcher *dispatcher.TaskDispatcher // 任务分配器 (Phase 2) - 需要更新为Proxy架构

	// Phase 2 - 模板引擎和参数验证
	TemplateEngine     template.TemplateEngine     // 模板引擎
	ParameterValidator template.ParameterValidator // 参数验证器

	// TODO: Phase 3 - 脚本执行引擎 - 需要更新为不依赖TaskDispatcher
	// ScriptExecutor *executor.ScriptExecutor // 脚本执行器
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := ent.NewClient(
		ent.Log(logx.Info), // logger
		ent.Driver(c.DatabaseConf.NewNoCacheDriver()),
		ent.Debug(), // debug mode
	)

	// 🎯 使用统一Hook系统 - 一键设置租户和部门Hook
	// 配置Ops服务的租户过滤规则 - 添加Ops特有的系统表（如果有）
	// 注：根据当前schema，所有表都应该是租户级别的，暂无系统级表需要排除

	// 一键设置：初始化配置 + 注册所有hooks (租户Hook + 部门Hook)
	if err := hooks.QuickSetup(db); err != nil {
		logx.Errorw("Failed to setup unified hooks", logx.Field("error", err.Error()))
		panic("统一Hook初始化失败: " + err.Error())
	}
	logx.Infow("✅ Ops-Center service: Unified hooks initialized successfully")

	// RPC服务层不注册数据权限拦截器
	// 数据权限控制应该在API层通过中间件处理，RPC层作为数据访问层不承担权限职责
	// 这样可以保持清晰的层次分离，避免跨服务的上下文传递问题

	// 初始化Redis连接
	rds := c.RedisConf.MustNewUniversalRedis()

	// 初始化Core RPC客户端 - 参考CMDB服务的实现模式
	var coreRpc coreclient.Core
	if c.CoreRpc.Endpoints != nil && len(c.CoreRpc.Endpoints) > 0 {
		// 创建RPC客户端，使用SystemContext拦截器支持系统级操作
		rpcClient, err := zrpc.NewClient(c.CoreRpc, zrpc.WithUnaryClientInterceptor(hooks.SystemContextClientInterceptor()))
		if err != nil {
			logx.Errorf("Failed to create Core RPC client: %v", err)
		} else {
			coreRpc = coreclient.NewCore(rpcClient)
			logx.Infow("✅ Core RPC client initialized successfully",
				logx.Field("endpoints", c.CoreRpc.Endpoints))
		}
	} else {
		logx.Info("Core RPC endpoints not configured, skipping Core RPC client initialization")
	}

	// 初始化Worker管理器
	workerManager := worker.NewWorkerManager(db, worker.DefaultWorkerManagerConfig())
	if err := workerManager.Start(); err != nil {
		logx.Errorf("Failed to start WorkerManager: %v", err)
		panic("WorkerManager启动失败: " + err.Error())
	}
	logx.Info("✅ WorkerManager started successfully")

	// TODO: 初始化任务分配器 (Phase 2) - 需要更新为Proxy架构
	// taskDispatcher := dispatcher.NewTaskDispatcher()
	// logx.Info("✅ TaskDispatcher initialized successfully")

	// 初始化模板引擎 (Phase 2)
	templateEngine := template.NewTemplateEngine()
	logx.Info("✅ TemplateEngine initialized successfully")

	// 初始化参数验证器 (Phase 2)
	parameterValidator := template.NewParameterValidator()
	logx.Info("✅ ParameterValidator initialized successfully")

	// TODO: 初始化脚本执行器 (Phase 3) - 需要更新为不依赖TaskDispatcher
	// scriptExecutor := executor.NewScriptExecutor(db, templateEngine, taskDispatcher, workerManager)
	// logx.Info("✅ ScriptExecutor initialized successfully")

	return &ServiceContext{
		Config:        c,
		DB:            db,
		Redis:         rds,
		CoreRpc:       coreRpc,
		WorkerManager: workerManager,
		// TODO: TaskDispatcher: taskDispatcher,

		// Phase 2
		TemplateEngine:     templateEngine,
		ParameterValidator: parameterValidator,

		// TODO: Phase 3
		// ScriptExecutor: scriptExecutor,
	}
}
