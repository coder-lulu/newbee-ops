# Agent管理功能实施快速清单

**创建时间**: 2025-12-16
**预计耗时**: 3.5小时
**依赖**: ops-center RPC已有完整proto定义

---

## 📝 需要创建的文件 (5个)

### 1. ops-center API定义文件

#### `/opt/code/newbee/ops-center/api/desc/agent.api`
- Agent列表、详情、创建、更新、删除接口
- Agent选择接口（供Discovery使用）
- **参考**: 完整代码见 `AGENT_MANAGEMENT_IMPLEMENTATION_PLAN.md` 任务1.1

#### `/opt/code/newbee/ops-center/api/desc/agentgroup.api`
- AgentGroup列表、详情、创建、更新、删除接口
- **参考**: 完整代码见 `AGENT_MANAGEMENT_IMPLEMENTATION_PLAN.md` 任务1.2

### 2. 前端API文件

#### `/opt/code/newbee/ui/apps/web-antd/src/api/ops/agent.ts`
- AgentAPI命名空间
- listAgents() - 获取Agent列表
- selectAgent() - 选择Agent
- **参考**: 完整代码见 `AGENT_MANAGEMENT_IMPLEMENTATION_PLAN.md` 任务3.1

### 3. 创建ops API目录
```bash
mkdir -p /opt/code/newbee/ui/apps/web-antd/src/api/ops
```

---

## 🔧 需要修改的文件 (5个)

### 1. unified-io RPC配置

#### `/opt/code/newbee/unified-io/rpc/etc/io.yaml`
**添加内容**（在CoreRpc后面）：
```yaml
# Ops-Center服务RPC连接配置
OpsRpc:
  Endpoints:
    - 127.0.0.1:9600
  Timeout: 30000  # 30秒超时，单位：毫秒
```

### 2. unified-io Config结构体

#### `/opt/code/newbee/unified-io/rpc/internal/config/config.go`
**添加内容**：
```go
// 在Config结构体中添加
type Config struct {
	// ... 已有字段
	OpsRpc          OpsRpcConf    // 新增
}

// 在文件末尾添加
type OpsRpcConf struct {
	Endpoints []string
	Timeout   int64 `json:",optional"`
}
```

### 3. unified-io ServiceContext

#### `/opt/code/newbee/unified-io/rpc/internal/svc/service_context.go`
**修改步骤**：
1. 添加import: `"github.com/coder-lulu/newbee-ops-center/rpc/opsclient"`
2. 在ServiceContext结构体添加: `OpsRpc ops.Ops`
3. 在NewServiceContext函数中初始化OpsRpc客户端

**参考**: 完整代码见 `AGENT_MANAGEMENT_IMPLEMENTATION_PLAN.md` 任务2.3

### 4. discovery_pool创建逻辑

#### `/opt/code/newbee/unified-io/rpc/internal/logic/discovery_pool/create_discovery_pool_logic.go`
**添加内容**：在创建发现池前验证Agent存在且在线
- 调用 `l.svcCtx.OpsRpc.GetAgentList()` 验证Agent
- 检查Agent状态为online

**参考**: 完整代码见 `AGENT_MANAGEMENT_IMPLEMENTATION_PLAN.md` 任务2.4

### 5. 前端AgentSelector组件

#### `/opt/code/newbee/ui/apps/web-antd/src/views/cmdb/ci_types/components/discovery/AgentSelector.vue`
**修改步骤**：
1. 导入 `AgentAPI` from `#/api/ops/agent`
2. 修改agents类型为 `AgentAPI.AgentItem[]`
3. 移除loadAgents()中的mock数据
4. 替换为真实API调用: `AgentAPI.listAgents()`
5. 更新模板中的字段名（agent.id → agent.agent_id）

**参考**: 完整代码见 `AGENT_MANAGEMENT_IMPLEMENTATION_PLAN.md` 任务3.2

---

## 🚀 执行命令清单

### Step 1: 创建ops-center API文件
```bash
# 创建agent.api（复制完整内容）
vim /opt/code/newbee/ops-center/api/desc/agent.api

# 创建agentgroup.api（复制完整内容）
vim /opt/code/newbee/ops-center/api/desc/agentgroup.api
```

### Step 2: 生成ops-center API代码
```bash
cd /opt/code/newbee/ops-center/api

# 生成Agent API代码
goctl api go -api desc/agent.api -dir . --style=go_zero

# 生成AgentGroup API代码
goctl api go -api desc/agentgroup.api -dir . --style=go_zero

# 验证编译
go build -v .
```

### Step 3: 修改unified-io配置和代码
```bash
# 1. 修改配置文件
vim /opt/code/newbee/unified-io/rpc/etc/io.yaml

# 2. 修改Config结构体
vim /opt/code/newbee/unified-io/rpc/internal/config/config.go

# 3. 修改ServiceContext
vim /opt/code/newbee/unified-io/rpc/internal/svc/service_context.go

# 4. 修改discovery_pool逻辑
vim /opt/code/newbee/unified-io/rpc/internal/logic/discovery_pool/create_discovery_pool_logic.go

# 5. 验证编译
cd /opt/code/newbee/unified-io/rpc
go build -v .
```

### Step 4: 创建前端API和修改组件
```bash
# 1. 创建ops API目录
mkdir -p /opt/code/newbee/ui/apps/web-antd/src/api/ops

# 2. 创建agent.ts
vim /opt/code/newbee/ui/apps/web-antd/src/api/ops/agent.ts

# 3. 修改AgentSelector.vue
vim /opt/code/newbee/ui/apps/web-antd/src/views/cmdb/ci_types/components/discovery/AgentSelector.vue
```

### Step 5: 插入测试数据
```sql
-- 连接到MySQL
mysql -h 192.168.26.130 -P 3306 -uroot -p123456 -D newbee

-- 插入测试Agent
INSERT INTO agents (
  name, agent_id, host, port, agent_status,
  last_heartbeat, heartbeat_interval, region,
  version, status, tenant_id, created_at, updated_at
) VALUES (
  'Test-Agent-Beijing-01',
  'agent-test-001',
  '127.0.0.1',
  8080,
  'online',
  UNIX_TIMESTAMP(),
  30,
  'beijing',
  '1.0.0',
  1,
  1,
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP()
);

-- 验证数据
SELECT * FROM agents WHERE tenant_id = 1;
```

---

## ✅ 测试验证清单

### 测试1: ops-center API
```bash
# 启动ops-center服务
cd /opt/code/newbee/ops-center/rpc && go run ops.go &
cd /opt/code/newbee/ops-center/api && go run ops.go &

# 测试Agent列表接口
curl -X GET "http://localhost:9410/agent/list?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 预期：返回JSON，包含刚插入的测试Agent
```

### 测试2: unified-io调用ops-center
```bash
# 启动unified-io服务
cd /opt/code/newbee/unified-io/rpc && go run io.go

# 检查日志，应该看到：
# "Initializing OpsRpc client..."
# "OpsRpc client initialized successfully"
```

### 测试3: 前端完整流程
```bash
# 1. 启动前端
cd /opt/code/newbee/ui/apps/web-antd
pnpm dev

# 2. 浏览器访问
http://localhost:5173

# 3. 操作步骤
# - 登录系统
# - 进入CMDB → CI类型管理
# - 选择任意CI类型 → 新建发现配置
# - 在"选择执行代理"步骤，应该能看到真实的Agent列表
# - 选择Agent后继续，验证能否成功创建发现池
```

---

## 🎯 完成标准

### ✅ Phase 1: ops-center API (必须)
- [ ] agent.api文件创建完成
- [ ] agentgroup.api文件创建完成
- [ ] API代码生成成功（无错误）
- [ ] ops-center.api编译通过
- [ ] curl测试 `/agent/list` 返回200

### ✅ Phase 2: unified-io集成 (必须)
- [ ] io.yaml添加OpsRpc配置
- [ ] config.go添加OpsRpcConf结构体
- [ ] service_context.go初始化OpsRpc客户端
- [ ] discovery_pool逻辑添加Agent验证
- [ ] unified-io.rpc编译通过

### ✅ Phase 3: 前端集成 (必须)
- [ ] agent.ts API文件创建完成
- [ ] AgentSelector.vue移除mock数据
- [ ] 前端能成功调用 `/ops-api/agent/list`
- [ ] 浏览器Network面板看到真实API请求
- [ ] 发现配置页面显示真实Agent列表

### ✅ Phase 4: 数据验证 (可选)
- [ ] agents表存在且结构正确
- [ ] 至少插入1条测试Agent数据
- [ ] Agent状态为online
- [ ] last_heartbeat时间正确

---

## 📦 文件依赖关系图

```
ops-center/api/desc/agent.api
        ↓ (goctl生成)
ops-center/api/internal/handler/agent/*.go
        ↓ (调用)
ops-center/rpc (agent.proto已存在)
        ↑ (gRPC调用)
unified-io/rpc/internal/svc/service_context.go (OpsRpc客户端)
        ↑ (调用)
unified-io/rpc/internal/logic/discovery_pool/*.go
        ↑ (gRPC调用)
unified-io/api/internal/handler/*.go
        ↑ (HTTP调用)
前端 src/api/ops/agent.ts
        ↑ (导入)
前端 src/views/.../AgentSelector.vue
```

---

## 🔍 常见问题

### Q1: 生成API代码报错
**A**: 检查goctl版本，确保>=1.6.0，并且当前目录在ops-center/api

### Q2: unified-io启动报错"unknown field OpsRpc"
**A**: 确认config.go中已添加OpsRpcConf结构体定义

### Q3: 前端调用Agent API返回404
**A**: 检查vite.config.mts中ops-api代理配置，确认端口为9410

### Q4: Agent列表为空
**A**: 检查agents表中是否有数据，且agent_status='online'

### Q5: unified-io调用ops-center超时
**A**: 确认ops-center.rpc服务已启动（端口9600）

---

## 📚 参考文档

- **完整实施方案**: `AGENT_MANAGEMENT_IMPLEMENTATION_PLAN.md`
- **Agent Proto定义**: `/opt/code/newbee/ops-center/rpc/desc/agent.proto`
- **Discovery配置流程**: `/opt/code/newbee/unified-io/migrations/discovery_configuration_service_dependencies.md`

---

**快速开始**: 按顺序执行"执行命令清单"，每步完成后验证无误再继续下一步。

**预计完成时间**: 3.5小时
**文档创建时间**: 2025-12-16
