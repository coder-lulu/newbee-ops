# Phase 3: 会话管理与任务下发实现完成报告

## 📅 完成时间
2025-12-28

## ✅ 已完成的工作

### 1. 架构复查与确认

#### 1.1 Proxy架构变更确认

**重要发现**：Proxy已完全移除gRPC通信方式，改为纯HTTP/WebSocket架构。

**通信架构**：
```
┌─────────────┐                    ┌─────────────┐
│             │  HTTP/PSK          │             │
│  Proxy      │─────────────────>  │ Ops-Center  │
│  (newbee-   │  注册/心跳/结果上报  │  (API)      │
│   proxy)    │                    │             │
└─────────────┘                    └─────────────┘
       ▲                                  │
       │                                  │
       │  HTTP调用                         │
       │  任务下发/会话创建                  │
       └──────────────────────────────────┘
```

**文件**: `/opt/code/newbee/newbee-proxy/internal/client/ops_center_client.go`

**Proxy调用Ops-Center的接口**：
- `POST /worker/register` - Proxy注册
- `POST /worker/heartbeat` - Proxy心跳
- `POST /task/result` - 任务结果上报
- 旧版兼容：`/ops/proxy/register`, `/ops/proxy/heartbeat`

#### 1.2 Proxy提供的HTTP接口

**任务执行接口** (`/api/task/*`):
```go
POST /api/task/command      // 执行命令
POST /api/task/script       // 执行脚本
POST /api/task/file         // 文件传输
POST /api/task/cancel       // 取消任务
```

**任务查询接口**:
```go
GET  /api/task/status/:taskId   // 获取任务状态
GET  /api/task/result/:taskId   // 获取任务结果
GET  /api/task/active           // 活跃任务列表
GET  /api/task/stats            // 任务统计
```

**WebSocket会话接口**:
```go
GET  /ws/ssh                    // SSH WebSocket
GET  /ws/telnet                 // Telnet WebSocket
GET  /api/rdp/websocket         // RDP WebSocket (Guacamole)
GET  /api/ssh/websocket         // SSH WebSocket (Guacamole)
GET  /api/vnc/websocket         // VNC WebSocket (Guacamole)
GET  /api/telnet/websocket      // Telnet WebSocket (Guacamole)
```

**数据库操作接口**:
```go
POST /api/db/test               // 测试数据库连接
POST /api/db/connect            // 创建数据库连接
POST /api/db/execute            // 执行SQL
GET  /api/db/databases          // 获取数据库列表
GET  /api/db/tables             // 获取表列表
```

### 2. Proxy HTTP客户端实现

**文件**: `/opt/code/newbee/ops-center/api/internal/client/proxy_http_client.go`

#### 2.1 核心功能

```go
type ProxyHTTPClient struct {
    httpClient *http.Client
    logger     logx.Logger
}
```

**任务执行方法**:
- `ExecuteCommand(ctx, endpoint, req)` - 执行命令
- `ExecuteScript(ctx, endpoint, req)` - 执行脚本
- `ExecuteFileTransfer(ctx, endpoint, req)` - 文件传输
- `CancelTask(ctx, endpoint, taskID)` - 取消任务

**任务查询方法**:
- `GetTaskStatus(ctx, endpoint, taskID)` - 获取任务状态
- `GetTaskResult(ctx, endpoint, taskID)` - 获取任务结果

**WebSocket URL生成方法**:
- `GenerateSSHWebSocketURL(endpoint, target, port, username, password)` - SSH
- `GenerateTelnetWebSocketURL(endpoint, target, port, username, password)` - Telnet
- `GenerateRDPWebSocketURL(endpoint, target, port, username, password)` - RDP
- `GenerateVNCWebSocketURL(endpoint, target, port, password)` - VNC

#### 2.2 请求/响应类型

**命令执行请求**:
```go
type CommandExecuteRequest struct {
    Target       string            `json:"target"`
    Port         int32             `json:"port"`
    Protocol     string            `json:"protocol"`      // ssh/telnet
    Username     string            `json:"username"`
    Password     string            `json:"password"`
    PrivateKey   string            `json:"private_key"`
    Command      string            `json:"command"`
    WorkingDir   string            `json:"working_dir"`
    Timeout      int32             `json:"timeout"`
    UseSudo      bool              `json:"use_sudo"`
    SudoPassword string            `json:"sudo_password"`
    Environment  map[string]string `json:"environment"`
}
```

**任务响应**:
```go
type TaskResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    TaskID  string `json:"task_id"`
    Error   string `json:"error,omitempty"`
}
```

**任务结果响应**:
```go
type TaskResultResponse struct {
    TaskID         string                 `json:"task_id"`
    TaskType       string                 `json:"task_type"`
    Status         string                 `json:"status"`
    StartTime      string                 `json:"start_time"`
    UpdateTime     string                 `json:"update_time"`
    Duration       int64                  `json:"duration"`
    ResultStatus   string                 `json:"result_status"`
    ErrorMessage   string                 `json:"error_message"`
    HasResult      bool                   `json:"has_result"`
    DetailedResult map[string]interface{} `json:"detailed_result"`
    ParsedResult   map[string]interface{} `json:"parsed_result"`
}
```

### 3. 会话创建逻辑实现

**文件**: `/opt/code/newbee/ops-center/api/internal/logic/ops/create_session_logic.go`

#### 3.1 实现流程

```
1. 验证协议 (ssh|telnet|rdp|vnc)
   ↓
2. 选择合适的Proxy (使用ProxyPick)
   ↓
3. 获取连接凭证 (从params参数)
   ↓
4. 生成WebSocket URL
   ↓
5. 创建Session记录到数据库
   ↓
6. 返回SessionID、WebSocket URL、Token
```

#### 3.2 核心代码

```go
func (l *CreateSessionLogic) CreateSession(req *types.CreateSessionReq)
    (resp *types.CreateSessionResp, err error) {

    // 1. 验证协议
    if !isValidProtocol(req.Protocol) {
        return nil, fmt.Errorf("invalid protocol: %s", req.Protocol)
    }

    // 2. 选择Proxy
    pickReq := &types.ProxyPickReq{
        Strategy:             "least_connections",
        RequiredCapabilities: []string{req.Protocol},
    }
    pickLogic := NewProxyPickLogic(l.ctx, l.svcCtx)
    proxyPickResp, err := pickLogic.ProxyPick(pickReq)

    // 3. 获取连接信息
    target := req.Params["target"]
    username := req.Params["username"]
    password := req.Params["password"]
    port := req.Options.Port
    if port == 0 {
        port = getDefaultPort(req.Protocol)
    }

    // 4. 生成WebSocket URL
    proxyClient := client.NewProxyHTTPClient()
    proxyHTTPEndpoint := proxyPickResp.Endpoints["http"]

    switch req.Protocol {
    case "ssh":
        wsURL = proxyClient.GenerateSSHWebSocketURL(
            proxyHTTPEndpoint, target, port, username, password)
    case "rdp":
        wsURL = proxyClient.GenerateRDPWebSocketURL(
            proxyHTTPEndpoint, target, port, username, password)
    // ... 其他协议
    }

    // 5. 创建Session记录
    sessionInfo := &ops.SessionInfo{
        CiId:      &req.CiId,
        Protocol:  &req.Protocol,
        ProxyId:   &proxyPickResp.ProxyID,
        Endpoint:  &proxyHTTPEndpoint,
        CreatedAt: pointy(now),
        ExpiresAt: pointy(expiresAt),
        Status:    pointy("active"),
    }
    createResp, err := l.svcCtx.OpsClient.CreateSession(l.ctx, sessionInfo)

    // 6. 返回响应
    return &types.CreateSessionResp{
        SessionId: sessionID,
        ProxyId:   proxyPickResp.ProxyID,
        WsUrl:     wsURL,
        Token:     token,
        ExpiresAt: expiresAt,
        Handshake: handshake,
    }, nil
}
```

#### 3.3 辅助函数

```go
// 协议验证
func isValidProtocol(protocol string) bool {
    validProtocols := []string{"ssh", "telnet", "rdp", "vnc"}
    for _, p := range validProtocols {
        if p == protocol {
            return true
        }
    }
    return false
}

// 默认端口
func getDefaultPort(protocol string) int {
    switch protocol {
    case "ssh":    return 22
    case "telnet": return 23
    case "rdp":    return 3389
    case "vnc":    return 5900
    default:       return 22
    }
}
```

### 4. 任务下发逻辑实现

**文件**: `/opt/code/newbee/ops-center/api/internal/logic/ops/create_task_logic.go`

#### 4.1 实现流程

```
1. 验证任务参数 (CiIds, Command, Executor)
   ↓
2. 生成任务ID
   ↓
3. 创建Task记录到数据库
   ↓
4. 为每个CI下发任务:
   4.1 选择Proxy
   4.2 构建HTTP endpoint
   4.3 根据执行器类型调用Proxy接口
   ↓
5. 统计执行结果
   ↓
6. 返回任务ID和状态
```

#### 4.2 核心代码

```go
func (l *CreateTaskLogic) CreateTask(req *types.CreateTaskReq)
    (resp *types.CreateTaskResp, err error) {

    // 1. 验证参数
    if len(req.CiIds) == 0 {
        return nil, fmt.Errorf("ci_ids is required")
    }
    if req.Command.Content == "" {
        return nil, fmt.Errorf("command content is required")
    }

    // 2. 生成任务ID
    taskID := fmt.Sprintf("task_%d", time.Now().Unix())

    // 3. 创建Task记录
    taskInfo := &ops.TaskInfo{
        TaskId:    &taskID,
        Name:      &req.Name,
        Status:    pointy("pending"),
        CreatedAt: pointy(now),
    }
    createResp, err := l.svcCtx.OpsClient.CreateTask(l.ctx, taskInfo)

    // 4. 为每个CI下发任务
    proxyClient := client.NewProxyHTTPClient()
    for _, ciId := range req.CiIds {
        // 4.1 选择Proxy
        pickReq := &types.ProxyPickReq{
            Strategy:             "least_connections",
            RequiredCapabilities: []string{req.Executor},
        }
        proxyPickResp, err := pickLogic.ProxyPick(pickReq)

        // 4.2 构建endpoint
        proxyHTTPEndpoint := proxyPickResp.Endpoints["http"]

        // 4.3 根据执行器类型下发
        if req.Executor == "ssh" || req.Executor == "telnet" {
            cmdReq := &client.CommandExecuteRequest{
                Target:      target,
                Port:        22,
                Protocol:    req.Executor,
                Username:    username,
                Password:    password,
                Command:     req.Command.Content,
                WorkingDir:  req.Command.WorkDir,
                Timeout:     int32(req.Command.Timeout),
                Environment: req.Command.Env,
            }
            taskResp, err = proxyClient.ExecuteCommand(
                l.ctx, proxyHTTPEndpoint, cmdReq)
        } else if req.Executor == "agent" {
            scriptReq := &client.ScriptExecuteRequest{
                Target:        target,
                ScriptContent: req.Command.Content,
                ScriptType:    "bash",
                // ...
            }
            taskResp, err = proxyClient.ExecuteScript(
                l.ctx, proxyHTTPEndpoint, scriptReq)
        }
    }

    // 5. 统计结果
    status := "running"
    if successCount == 0 {
        status = "failed"
    } else if len(failedCIs) > 0 {
        status = "partial"
    }

    return &types.CreateTaskResp{
        TaskId: taskID,
        Status: status,
    }, nil
}
```

#### 4.3 支持的执行器类型

| 执行器类型 | 说明 | Proxy接口 |
|----------|------|----------|
| `ssh` | SSH命令执行 | `/api/task/command` |
| `telnet` | Telnet命令执行 | `/api/task/command` |
| `agent` | 脚本执行（通过SSH上传脚本） | `/api/task/script` |

## 📊 代码统计

### 新增文件

| 文件 | 行数 | 说明 |
|------|------|------|
| `proxy_http_client.go` | 308 | Proxy HTTP客户端 |
| `create_session_logic.go` | 197 | 会话创建逻辑 |
| `create_task_logic.go` | 192 | 任务下发逻辑 |
| **总计** | **697** | **3个核心文件** |

### 功能实现

| 功能模块 | 方法数 | 说明 |
|---------|-------|------|
| HTTP客户端 | 10 | 任务执行、查询、WebSocket URL生成 |
| 会话创建 | 3 | 主逻辑 + 2个辅助函数 |
| 任务下发 | 1 | 主逻辑（含多CI批量下发） |
| **总计** | **14** | **完整的HTTP集成** |

## 🔍 技术亮点

### 1. 纯HTTP架构

- ✅ **完全移除gRPC依赖** - Proxy与Ops-Center之间仅通过HTTP通信
- ✅ **WebSocket支持** - 实时会话通过WebSocket实现
- ✅ **RESTful设计** - 清晰的资源路径和HTTP动词

### 2. 智能Proxy选择

- ✅ **复用Phase 2的ProxyPick** - 会话和任务都使用同一套负载均衡
- ✅ **能力匹配** - 根据协议/执行器类型筛选Proxy
- ✅ **最少连接数策略** - 默认选择负载最轻的Proxy

### 3. 多协议支持

- ✅ **4种远程协议** - SSH、Telnet、RDP、VNC
- ✅ **2种WebSocket模式** - 基础WebSocket + Guacamole协议
- ✅ **灵活的端口配置** - 支持自定义端口和默认端口

### 4. 批量任务下发

- ✅ **多CI支持** - 一次请求可下发到多个CI
- ✅ **失败容错** - 部分CI失败不影响其他CI
- ✅ **状态追踪** - 区分running/failed/partial状态

### 5. WebSocket URL生成

- ✅ **协议自动转换** - HTTP→WS、HTTPS→WSS
- ✅ **参数编码** - 自动处理URL参数
- ✅ **多endpoint支持** - 优先使用Proxy注册的HTTP endpoint

## 🔗 API使用示例

### 会话创建

**请求示例**:
```bash
curl -X POST http://localhost:9601/ops/session/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "ci_id": "server-001",
    "protocol": "ssh",
    "params": {
      "target": "192.168.1.100",
      "username": "root",
      "password": "password123"
    },
    "options": {
      "port": 22
    }
  }'
```

**响应示例**:
```json
{
  "session_id": "session_1735363200",
  "proxy_id": "proxy-001",
  "ws_url": "ws://192.168.1.50:8080/ws/ssh?target=192.168.1.100&port=22&username=root&password=password123",
  "token": "token_12345_1735363200",
  "expires_at": 1735370400,
  "handshake": {
    "protocol": "ssh",
    "proxy_id": "proxy-001",
    "target": "192.168.1.100"
  }
}
```

### 任务创建

**请求示例**:
```bash
curl -X POST http://localhost:9601/ops/task/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "批量检查磁盘空间",
    "ci_ids": ["server-001", "server-002", "server-003"],
    "executor": "ssh",
    "command": {
      "content": "df -h",
      "timeout": 60
    }
  }'
```

**响应示例**:
```json
{
  "task_id": "task_1735363200",
  "status": "running"
}
```

### 任务状态查询

**请求示例**:
```bash
curl -X GET http://localhost:9601/ops/task/status/task_1735363200 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**响应示例**:
```json
{
  "task_id": "task_1735363200",
  "status": "running"
}
```

## 🎯 后续优化建议

### 1. CMDB集成

**当前状态**: 连接信息从请求参数获取
**优化方向**: 从CMDB自动获取CI的连接信息

```go
// 当前代码 (TODO)
target := req.Params["target"]
username := req.Params["username"]
password := req.Params["password"]

// 优化后
ciInfo, err := l.svcCtx.CMDBClient.GetCI(req.CiId)
target := ciInfo.IPAddress
credentials := ciInfo.Credentials
```

### 2. AccessProfile集成

**当前状态**: 未使用AccessProfile
**优化方向**: 支持从AccessProfile获取访问策略

```go
// 获取AccessProfile
profile, err := l.svcCtx.OpsClient.GetAccessProfile(l.ctx, req.CiId)

// 使用AccessProfile中的配置
credentials := profile.CredentialRef
preferProxy := profile.PreferProxy
jumpChain := profile.JumpChain
```

### 3. Token鉴权

**当前状态**: 简单的字符串拼接
**优化方向**: 使用JWT生成临时token

```go
// 生成JWT token
token, err := jwt.Generate(&jwt.Claims{
    SessionID: sessionID,
    UserID:    userID,
    ExpiresAt: expiresAt,
})
```

### 4. 会话超时管理

**当前状态**: 固定2小时过期
**优化方向**:
- 可配置的过期时间
- 定时清理过期会话
- 会话活跃度检测

### 5. 任务结果收集

**当前状态**: 仅下发任务，未实现结果收集
**优化方向**:
- 定时轮询Proxy获取任务结果
- WebSocket实时推送结果
- 结果写入数据库

### 6. 错误重试机制

**当前状态**: 失败直接记录
**优化方向**:
- 自动重试失败的任务
- 指数退避算法
- 最大重试次数限制

## 📝 配置要求

### Ops-Center配置

```yaml
# etc/ops.yaml
Ops:
  Registration:
    PSK: "dev-psk"

  ProxySelection:
    DefaultStrategy: "least_connections"
    HealthCheckInterval: 60
    OfflineThreshold: 180

  Session:
    DefaultTimeout: 7200  # 2小时
    MaxConcurrent: 1000

  Task:
    DefaultTimeout: 300   # 5分钟
    MaxConcurrent: 100
```

## ✅ 验收标准

### Phase 3完成标准

- [x] Proxy HTTP客户端已实现
- [x] 会话创建逻辑已实现（含Proxy选择）
- [x] 任务下发逻辑已实现（支持多CI）
- [x] WebSocket URL生成已实现（4种协议）
- [x] 所有代码已添加详细注释
- [x] 代码符合Go最佳实践
- [ ] 端到端测试通过（待测试）
- [ ] CMDB集成（待实现）
- [ ] 任务结果收集（待实现）

## 📚 参考文档

- **Phase 1 总结**: `/opt/code/newbee/ops-center/migrations/phase1_summary.md`
- **Phase 2 总结**: `/opt/code/newbee/ops-center/migrations/phase2_summary.md`
- **Proxy集成说明**: `/opt/code/newbee/newbee-proxy/README_OPS_INTEGRATION.md`
- **API定义**: `/opt/code/newbee/ops-center/api/desc/ops.api`

## 🔄 架构演进

### Phase 1: 数据库Schema
- ✅ Worker → Proxy重命名
- ✅ 数据库迁移SQL

### Phase 2: API Logic
- ✅ Proxy注册/心跳逻辑
- ✅ Proxy选择逻辑（5种策略）

### Phase 3: 会话与任务（当前）
- ✅ HTTP客户端工具
- ✅ 会话创建（WebSocket URL生成）
- ✅ 任务下发（多CI批量）

### Phase 4: 完善与优化（建议）
1. CMDB集成
2. AccessProfile集成
3. 任务结果收集
4. 会话超时管理
5. 错误重试机制
6. 性能监控和告警

---

**Phase 3 完成时间**: 2025-12-28
**状态**: ✅ 核心功能全部完成
**准备进入**: Phase 4 - 完善与优化 或 直接测试验证
