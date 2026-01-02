package svc

import (
	"context"
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	commoncasbin "github.com/coder-lulu/newbee-common/v2/casbin"
	commonadapter "github.com/coder-lulu/newbee-common/v2/casbin/adapter"
	"github.com/coder-lulu/newbee-common/v2/i18n"
	"github.com/coder-lulu/newbee-common/v2/middleware/audit"
	"github.com/coder-lulu/newbee-common/v2/middleware/dataperm"
	"github.com/coder-lulu/newbee-common/v2/middleware/integration"
	"github.com/coder-lulu/newbee-common/v2/middleware/keys"
	"github.com/coder-lulu/newbee-common/v2/middleware/logging"
	"github.com/coder-lulu/newbee-common/v2/orm/ent/hooks"
	coreclient "github.com/coder-lulu/newbee-core/rpc/coreclient"
	apicasbin "github.com/coder-lulu/newbee-ops-api/internal/casbin"
	"github.com/coder-lulu/newbee-ops-api/internal/config"
	proxyservice "github.com/coder-lulu/newbee-ops-api/internal/services/proxy"
	rpcclient "github.com/coder-lulu/newbee-ops-api/internal/rpcclient"
	"github.com/coder-lulu/newbee-ops-api/internal/workerclient"
	"github.com/coder-lulu/newbee-ops-rpc/ent"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                 config.Config
	ContextManager         *keys.ContextManager
	IntegrationResult      *integration.Result
	ManagedMiddlewareChain []rest.Middleware
	CoreRpc                coreclient.Core
	Casbin                 *casbin.Enforcer

	// 数据库连接和Worker客户端（Phase 6新增）
	DB           *ent.Client
	WorkerClient *workerclient.WorkerClient

	// Proxy 注册表（内存）
	ProxyRegistry *proxyservice.Registry

	// 临时：AccessProfile 内存存储（后续切换到 RPC/Ent 持久化）
	ProfileStore *AccessProfileStore

	// 会话内存存储（用于管理 Session 生命周期）
	SessionStore *SessionStore

	// 任务内存存储（基础编排占位，后续可切RPC/DB）
	TaskStore *TaskStore

	// 凭证与审计（内存存储，占位，后续切RPC/DB）
	CredentialStore *CredentialStore
	AuditStore      *AuditStore

	// 预置：Ops RPC 客户端（未生成代码前可为nil或stub）
	OpsClient rpcclient.OpsClient
	Trans     *Translator
}

// GetCoreRpcClient implements audit.AuditSvcProvider
func (svc *ServiceContext) GetCoreRpcClient() interface{} { return svc.CoreRpc }

// GetCasbinEnforcer 实现permission.EnforcerProvider接口，为RBAC插件提供Casbin执行器
func (svc *ServiceContext) GetCasbinEnforcer() interface{} { return svc.Casbin }

func NewServiceContext(c config.Config) *ServiceContext {
	// 1. Redis
	rds := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{c.RedisConf.Host},
		Password: c.RedisConf.Pass,
		DB:       c.RedisConf.Db,
	})

	// 预留：如需加载内置多语言资源，请在提供有效 embed.FS 后初始化
	_ = i18n.Success

	// 2. 初始化Core RPC客户端 - 使用SystemContext拦截器支持系统级操作
	coreRpcClient, err := zrpc.NewClient(c.CoreRpc, zrpc.WithUnaryClientInterceptor(hooks.SystemContextClientInterceptor()))
	if err != nil {
		panic("Failed to create Core RPC client: " + err.Error())
	}
	coreRpc := coreclient.NewCore(coreRpcClient)

	// 3. JWT Secret 从中间件配置读取
	jwtSecret := ""
	if c.Middleware.Auth != nil && c.Middleware.Auth.AccessSecret != "" {
		jwtSecret = c.Middleware.Auth.AccessSecret
	}

	// 4. 🔥 初始化Casbin - 使用EntAdapter通过RPC查询规则
	logx.Info("Initializing Casbin with EntAdapter and RPC querier for all tenants")

	// 4.1 创建RPC查询器（仅用于EntAdapter）
	rpcQuerier := apicasbin.NewRpcCasbinRuleQuerier(coreRpc)

	// 4.2 创建EntAdapter
	// 🔥 使用SystemContext绕过租户隔离Hook，加载所有租户的Casbin规则
	systemCtx := hooks.NewSystemContext(context.Background())
	adapter := commonadapter.NewEntAdapter(rpcQuerier, systemCtx)

	// 4.3 创建Casbin模型
	modelText := commoncasbin.GetDefaultRBACWithDomainsModel()
	m, err := model.NewModelFromString(modelText)
	if err != nil {
		logx.Errorf("Failed to create Casbin model: %v", err)
		panic("Casbin model creation failed: " + err.Error())
	}

	// 4.4 创建Enforcer
	cbn, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		logx.Errorf("Failed to create Casbin enforcer: %v", err)
		panic("Casbin enforcer creation failed: " + err.Error())
	}

	// 4.5 加载策略
	err = cbn.LoadPolicy()
	if err != nil {
		logx.Errorf("Failed to load Casbin policy: %v", err)
		panic("Casbin policy loading failed: " + err.Error())
	}

	// 4.6 添加Redis Watcher（用于策略同步）
	w := c.CasbinConf.MustNewOriginalRedisWatcher(c.RedisConf, func(data string) {
		rediswatcher.DefaultUpdateCallback(cbn)(data)
	})
	err = cbn.SetWatcher(w)
	if err != nil {
		logx.Errorf("Failed to set Casbin watcher: %v", err)
		panic("Casbin watcher setup failed: " + err.Error())
	}

	logx.Info("✅ Ops-Center API: Casbin initialized with EntAdapter, connected to sys_casbin_rules table via RPC")

	// 5. 🔥 Phase 6: 初始化数据库连接和WorkerClient
	logx.Info("Initializing database connection and WorkerClient for Worker management")

	// 5.1 创建数据库连接
	db := ent.NewClient(
		ent.Driver(c.DatabaseConf.NewNoCacheDriver()),
		ent.Debug(), // 开启调试模式，打印SQL
	)

	// 5.2 初始化租户Hook（自动租户隔离）
	if err := hooks.QuickSetup(db); err != nil {
		panic("Worker管理 - 租户Hook初始化失败: " + err.Error())
	}

	// 5.3 初始化WorkerClient
	workerClient := workerclient.NewWorkerClient(db)
	logx.Info("✅ Ops-Center API: WorkerClient initialized successfully")

	// 6. 创建服务上下文实例
	svcCtx := &ServiceContext{
		Config:       c,
		CoreRpc:      coreRpc,
		Casbin:       cbn,
		DB:           db,
		WorkerClient: workerClient,
	}

	// 7. 中间件集成 + 审计写入器（可扩展为RPC写入）
	opsAuditWriter := audit.NewBuiltinAuditWriter(svcCtx)

	// 数据权限 CasbinProvider：使用默认提供者，适配 Core RPC
	dpLogger := logging.NewMiddlewareLogger("ops-api-dataperm")
	dpCoreAdapter := apicasbin.NewCoreRPCAdapter(coreRpc)
	var dataPermProvider dataperm.CasbinProvider = dataperm.NewDefaultCasbinProvider(dpCoreAdapter, dpLogger)

	result, err := integration.Setup(&integration.Config{
		Redis:                  rds,
		JWTSecret:              jwtSecret,
		Mode:                   integration.Production,
		ApiResourceProvider:    nil,
		AuditWriter:            opsAuditWriter,
		TenantInfoProvider:     NewRpcTenantInfoProvider(coreRpc),
		RbacProvider:           svcCtx,           // 🔥 提供 Casbin enforcer（通过GetCasbinEnforcer接口）
		DataPermCasbinProvider: dataPermProvider, // 🔥 使用统一数据权限 Casbin 提供者
		Middleware:             &c.Middleware,
	})
	if err != nil {
		panic("Ops Center middleware setup failed: " + err.Error())
	}

	svcCtx.ContextManager = result.ContextManager
	svcCtx.IntegrationResult = result
	svcCtx.ManagedMiddlewareChain = result.Middlewares

	// Proxy 注册表（TTL 2m）
	svcCtx.ProxyRegistry = proxyservice.NewRegistry(2 * 60 * 1e9)
	svcCtx.ProfileStore = NewAccessProfileStore()
	svcCtx.SessionStore = NewSessionStore()
	svcCtx.TaskStore = NewTaskStore()
	svcCtx.CredentialStore = NewCredentialStore()
	svcCtx.AuditStore = NewAuditStore()
	// 初始化 RPC 客户端（若配置了 OpsRpc）
	if c.OpsRpc.Endpoints != nil && len(c.OpsRpc.Endpoints) > 0 {
		// 🔥 添加租户上下文传播拦截器 - 将 tenant_id, user_id, dept_id 注入到 gRPC metadata
		cli := zrpc.MustNewClient(c.OpsRpc,
			zrpc.WithUnaryClientInterceptor(hooks.SystemContextClientInterceptor()),
			zrpc.WithStreamClientInterceptor(hooks.SystemContextStreamClientInterceptor()))
		svcCtx.OpsClient = rpcclient.NewGrpcOpsClient(cli)
		logx.Infow("✅ Ops RPC client initialized successfully (with context propagation interceptor)",
			logx.Field("endpoints", c.OpsRpc.Endpoints))
	} else {
		logx.Info("Ops RPC endpoints not configured, skipping Ops RPC client initialization")
	}
	svcCtx.Trans = &Translator{}
	return svcCtx
}

// Translator 提供生成代码所需的 TransError 接口占位
type Translator struct{}

func (*Translator) TransError(ctx context.Context, err error) error { return err }

// AccessProfile 访问配置信息
type AccessProfile struct {
	CiId          string
	Capabilities  []string
	Ports         map[string]int
	CredentialRef string
	PreferProxy   string
	JumpChain     []string
	Tags          map[string]string
}

// AccessProfileStore 简易内存存储
type AccessProfileStore struct {
	mu   sync.RWMutex
	data map[string]AccessProfile // key=ciId
}

func NewAccessProfileStore() *AccessProfileStore {
	return &AccessProfileStore{data: make(map[string]AccessProfile)}
}

func (s *AccessProfileStore) Upsert(p AccessProfile) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[p.CiId] = p
}

func (s *AccessProfileStore) Get(ciId string) (AccessProfile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[ciId]
	return v, ok
}

func (s *AccessProfileStore) Delete(ciId string) { s.mu.Lock(); delete(s.data, ciId); s.mu.Unlock() }

func (s *AccessProfileStore) List() []AccessProfile {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]AccessProfile, 0, len(s.data))
	for _, v := range s.data {
		out = append(out, v)
	}
	return out
}

// Session 简易会话实体（内存）
type Session struct {
	ID        string
	TenantId  string
	UserId    string
	CiId      string
	Protocol  string
	ProxyId   string
	Endpoint  string
	CreatedAt int64
	ExpiresAt int64
	Status    string // active|closed
	ClosedAt  int64
}

// SessionStore 简易会话存储
type SessionStore struct {
	mu   sync.RWMutex
	data map[string]*Session
}

func NewSessionStore() *SessionStore { return &SessionStore{data: map[string]*Session{}} }

func (s *SessionStore) Start(sess *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[sess.ID] = sess
}

func (s *SessionStore) Get(id string) (*Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[id]
	return v, ok
}

// Close 关闭指定会话；若提供 tenantId 则校验一致性
func (s *SessionStore) Close(id string, tenantId string, ts int64) (bool, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data[id]
	if !ok {
		return false, "not found"
	}
	if tenantId != "" && v.TenantId != "" && tenantId != v.TenantId {
		return false, "tenant mismatch"
	}
	v.Status = "closed"
	v.ClosedAt = ts
	return true, ""
}

// ListBy 过滤并分页返回会话列表
func (s *SessionStore) ListBy(tenantId, status, ciId string, page, size int) ([]*Session, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	all := make([]*Session, 0, len(s.data))
	for _, v := range s.data {
		if tenantId != "" && v.TenantId != tenantId {
			continue
		}
		if status != "" && v.Status != status {
			continue
		}
		if ciId != "" && v.CiId != ciId {
			continue
		}
		all = append(all, v)
	}
	total := len(all)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	start := (page - 1) * size
	if start >= total {
		return []*Session{}, total
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total
}

// Task 简易任务实体（内存）
type Task struct {
	ID        string
	TenantId  string
	CreatorId string
	CiIds     []string
	Executor  string // agent|ssh|telnet (预留)
	Command   string
	Timeout   int
	Status    string // pending|running|success|failed
	Result    string // 简化：stdout/错误消息
	CreatedAt int64
	UpdatedAt int64
}

// TaskStore 简易任务存储
type TaskStore struct {
	mu   sync.RWMutex
	data map[string]*Task
}

func NewTaskStore() *TaskStore { return &TaskStore{data: map[string]*Task{}} }

func (s *TaskStore) Create(t *Task) { s.mu.Lock(); defer s.mu.Unlock(); s.data[t.ID] = t }
func (s *TaskStore) Get(id string) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[id]
	return v, ok
}
func (s *TaskStore) Update(id string, fn func(*Task)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data[id]
	if !ok {
		return false
	}
	fn(v)
	v.UpdatedAt = time.Now().Unix()
	return true
}

// Count returns number of tasks in store (for tests/metrics)
func (s *TaskStore) Count() int { s.mu.RLock(); defer s.mu.RUnlock(); return len(s.data) }

// First returns an arbitrary task (for tests)
func (s *TaskStore) First() (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.data {
		return v, true
	}
	return nil, false
}

// CredentialRef 简易实体
type CredentialRef struct {
	ID        string
	Provider  string
	Ref       string
	Scope     string
	CreatedBy string
	Tags      map[string]string
	CreatedAt int64
	UpdatedAt int64
}

type CredentialStore struct {
	mu   sync.RWMutex
	data map[string]*CredentialRef
}

func NewCredentialStore() *CredentialStore {
	return &CredentialStore{data: map[string]*CredentialRef{}}
}
func (s *CredentialStore) Upsert(c *CredentialRef) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c.CreatedAt == 0 {
		c.CreatedAt = time.Now().Unix()
	}
	c.UpdatedAt = time.Now().Unix()
	s.data[c.ID] = c
}
func (s *CredentialStore) Get(id string) (*CredentialRef, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[id]
	return v, ok
}
func (s *CredentialStore) Delete(id string) { s.mu.Lock(); delete(s.data, id); s.mu.Unlock() }
func (s *CredentialStore) List() []*CredentialRef {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*CredentialRef, 0, len(s.data))
	for _, v := range s.data {
		out = append(out, v)
	}
	return out
}

// AuditEvent 简易实体
type AuditEvent struct {
	EventID string
	Type    string
	Subject struct {
		User  string
		CI    string
		Proxy string
	}
	OccurAt int64
	Meta    map[string]string
}

type AuditStore struct {
	mu   sync.RWMutex
	data map[string]*AuditEvent
}

func NewAuditStore() *AuditStore           { return &AuditStore{data: map[string]*AuditEvent{}} }
func (s *AuditStore) Create(e *AuditEvent) { s.mu.Lock(); defer s.mu.Unlock(); s.data[e.EventID] = e }
func (s *AuditStore) Get(id string) (*AuditEvent, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[id]
	return v, ok
}
func (s *AuditStore) List() []*AuditEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*AuditEvent, 0, len(s.data))
	for _, v := range s.data {
		out = append(out, v)
	}
	return out
}
