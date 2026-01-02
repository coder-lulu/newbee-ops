# Phase 4: 完善与优化 - 任务规划

## 📅 规划时间
2025-12-28

## 🎯 Phase 4 目标

在Phase 1-3已经实现核心功能的基础上，Phase 4将完善剩余的API逻辑，并实现高级功能，确保ops-center成为一个完整可用的系统。

## 📋 任务清单

### 🔴 优先级1 - 核心API完善（必须完成）

#### 1.1 会话管理逻辑

| 文件 | 功能 | 复杂度 | 说明 |
|------|------|--------|------|
| `close_session_logic.go` | 关闭会话 | ⭐⭐ | 更新session状态为closed，记录关闭时间 |
| `get_session_logic.go` | 查询单个会话 | ⭐ | 通过ID查询session详情 |
| `get_session_by_query_logic.go` | 条件查询会话 | ⭐⭐ | 支持多条件筛选 |
| `list_session_logic.go` | 会话列表 | ⭐⭐ | 分页查询，支持排序和筛选 |

**实现要点**：
- 使用已有的OpsClient接口（GetSession、GetSessionList、UpdateSession）
- CloseSession需要更新session状态和closed_at时间
- 查询接口支持按状态、协议、ProxyID等筛选

#### 1.2 任务管理逻辑

| 文件 | 功能 | 复杂度 | 说明 |
|------|------|--------|------|
| `get_task_status_logic.go` | 查询任务状态 | ⭐ | 从RPC获取任务状态 |
| `get_task_result_logic.go` | 查询任务结果 | ⭐⭐ | 从RPC获取任务执行结果 |

**实现要点**：
- 使用OpsClient.GetTaskByTaskId接口
- 任务结果可能需要从Proxy的结果存储中获取
- 支持多种任务状态：pending/running/completed/failed

#### 1.3 Proxy接口重定向（兼容性）

| 文件 | 功能 | 复杂度 | 说明 |
|------|------|--------|------|
| `register_proxy_logic.go` | Proxy注册 | ⭐ | 重定向到/proxy/proxy_register |
| `heartbeat_proxy_logic.go` | Proxy心跳 | ⭐ | 重定向到/proxy/proxy_heartbeat |
| `pick_proxy_logic.go` | Proxy选择 | ⭐ | 重定向到/proxy/proxy_pick |

**实现要点**：
- 这些是旧的/ops路径接口，需要重定向到新的/proxy路径
- 保持接口兼容性，避免破坏现有调用

### 🟡 优先级2 - AccessProfile管理（推荐完成）

#### 2.1 AccessProfile CRUD

| 文件 | 功能 | 复杂度 | 说明 |
|------|------|--------|------|
| `create_profile_logic.go` | 创建访问配置 | ⭐⭐ | 创建CI的访问配置文件 |
| `update_profile_logic.go` | 更新访问配置 | ⭐⭐ | 更新访问配置 |
| `get_profile_logic.go` | 查询访问配置 | ⭐ | 查询单个配置 |
| `list_profile_logic.go` | 访问配置列表 | ⭐⭐ | 分页查询配置列表 |
| `delete_profile_logic.go` | 删除访问配置 | ⭐ | 删除配置 |

**AccessProfile数据结构**：
```go
type AccessProfile {
    CiId          string              // CI标识
    Capabilities  []string            // 支持的协议（ssh/rdp/telnet/vnc）
    Ports         map[string]int      // 协议端口映射
    CredentialRef string              // 凭证引用
    PreferProxy   string              // 优先使用的Proxy
    JumpChain     []string            // 跳板机链路
    Tags          map[string]string   // 自定义标签
}
```

**实现要点**：
- AccessProfile存储在数据库中（需要RPC支持）
- 可以在CreateSession时使用AccessProfile获取连接信息
- 支持跳板机链路（JumpChain）配置

### 🟢 优先级3 - 高级功能（可选）

#### 3.1 CMDB集成

**目标**：从CMDB自动获取CI连接信息

**当前问题**：
```go
// create_session_logic.go 和 create_task_logic.go
target := req.Params["target"]      // ❌ 手动传入
username := req.Params["username"]  // ❌ 手动传入
password := req.Params["password"]  // ❌ 手动传入
```

**改进方案**：
```go
// 从CMDB获取CI信息
ciInfo, err := l.svcCtx.CMDBClient.GetCI(req.CiId)
if err != nil {
    return nil, fmt.Errorf("failed to get CI info: %w", err)
}

target := ciInfo.ManagementIP
port := ciInfo.SSHPort
credentials := ciInfo.Credentials
```

**实现步骤**：
1. 创建CMDBClient接口
2. 实现GetCI方法（HTTP调用CMDB服务）
3. 在ServiceContext中注入CMDBClient
4. 修改CreateSession和CreateTask使用CMDB数据

#### 3.2 任务结果自动收集

**目标**：定期从Proxy拉取任务执行结果

**当前问题**：
- 任务下发后没有收集结果
- Proxy执行完成后，ops-center不知道

**改进方案1 - 定时轮询**：
```go
// 启动结果收集器
go l.svcCtx.TaskResultCollector.Start()

// 定时轮询逻辑
func (c *TaskResultCollector) Start() {
    ticker := time.NewTicker(10 * time.Second)
    for range ticker.C {
        tasks := c.getPendingTasks()
        for _, task := range tasks {
            result, err := c.proxyClient.GetTaskResult(
                ctx, task.ProxyEndpoint, task.ProxyTaskID)
            if err == nil && result.HasResult {
                c.saveTaskResult(task.TaskID, result)
            }
        }
    }
}
```

**改进方案2 - Proxy主动上报**（推荐）：
- Proxy已经实现了结果上报接口：`POST /task/result`
- 在ReportTaskResultLogic中处理上报的结果
- 更新任务状态到数据库

**实现步骤**：
1. 完善ReportTaskResultLogic逻辑
2. 更新Task状态到数据库
3. 可选：发送WebSocket通知给前端

#### 3.3 会话超时清理

**目标**：定期清理过期会话

**实现方案**：
```go
// 启动会话清理器
go l.svcCtx.SessionCleaner.Start()

func (c *SessionCleaner) Start() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        now := time.Now().Unix()

        // 查询过期会话
        sessions, err := c.opsClient.GetSessionList(ctx, &ops.SessionListReq{
            Status: pointy("active"),
        })

        for _, session := range sessions.Data {
            if session.ExpiresAt < now {
                // 关闭过期会话
                c.opsClient.UpdateSession(ctx, &ops.SessionInfo{
                    Id:       session.Id,
                    Status:   pointy("expired"),
                    ClosedAt: pointy(now),
                })

                c.logger.Infow("Session expired",
                    logx.Field("session_id", session.SessionId),
                    logx.Field("expires_at", session.ExpiresAt))
            }
        }
    }
}
```

**实现步骤**：
1. 创建SessionCleaner组件
2. 在ServiceContext中启动
3. 配置清理间隔和过期阈值

#### 3.4 错误重试机制

**目标**：任务下发失败时自动重试

**实现方案**：
```go
// 任务下发失败时
if err != nil {
    // 记录失败信息
    failureInfo := &TaskFailure{
        TaskID:     taskID,
        CiID:       ciId,
        Error:      err.Error(),
        RetryCount: 0,
        NextRetry:  time.Now().Add(1 * time.Minute),
    }

    // 保存到重试队列
    l.svcCtx.RetryQueue.Add(failureInfo)
    continue
}

// 重试处理器
func (r *RetryHandler) Start() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        tasks := r.getRetryableTasks()
        for _, task := range tasks {
            if task.RetryCount >= 3 {
                r.markAsFailed(task)
                continue
            }

            // 重试任务下发
            err := r.retryTask(task)
            if err == nil {
                r.removeFromQueue(task)
            } else {
                task.RetryCount++
                task.NextRetry = time.Now().Add(
                    time.Duration(task.RetryCount) * time.Minute)
                r.updateRetryInfo(task)
            }
        }
    }
}
```

**实现步骤**：
1. 创建RetryQueue和RetryHandler
2. 在任务下发失败时记录到队列
3. 定时重试，使用指数退避算法
4. 最大重试3次

#### 3.5 监控与告警

**目标**：监控系统运行状态

**监控指标**：
- Proxy健康状态（在线/离线/降级）
- 活跃会话数量
- 任务执行统计（成功/失败率）
- API响应时间
- 资源使用情况

**实现方案**：
```go
// 集成Prometheus指标
import "github.com/prometheus/client_golang/prometheus"

var (
    activeSessionsGauge = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "ops_active_sessions",
            Help: "Number of active sessions",
        })

    taskSuccessCounter = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "ops_task_success_total",
            Help: "Total number of successful tasks",
        })
)

// 在API逻辑中更新指标
func (l *CreateSessionLogic) CreateSession(...) {
    // ...
    activeSessionsGauge.Inc()
    // ...
}
```

**告警规则**：
- Proxy离线超过5分钟
- 任务失败率超过20%
- 会话数量超过阈值

## 📊 工作量估算

### 优先级1 - 核心API完善

| 模块 | 文件数 | 预计工时 | 说明 |
|------|--------|---------|------|
| 会话管理 | 4 | 4小时 | 简单的CRUD逻辑 |
| 任务管理 | 2 | 2小时 | 查询逻辑 |
| Proxy重定向 | 3 | 1小时 | 简单重定向 |
| **小计** | **9** | **7小时** | |

### 优先级2 - AccessProfile管理

| 模块 | 文件数 | 预计工时 | 说明 |
|------|--------|---------|------|
| AccessProfile CRUD | 5 | 5小时 | 需要RPC支持 |

### 优先级3 - 高级功能

| 模块 | 预计工时 | 说明 |
|------|---------|------|
| CMDB集成 | 3小时 | 需要CMDB服务配合 |
| 任务结果收集 | 2小时 | 完善ReportTaskResult |
| 会话超时清理 | 2小时 | 后台定时任务 |
| 错误重试 | 3小时 | 重试队列和处理器 |
| 监控告警 | 2小时 | Prometheus集成 |
| **小计** | **12小时** | |

### 总计

| 优先级 | 工时 | 说明 |
|--------|------|------|
| P1 | 7小时 | 必须完成 |
| P2 | 5小时 | 推荐完成 |
| P3 | 12小时 | 可选 |
| **总计** | **24小时** | **约3个工作日** |

## 🎯 分阶段实施建议

### 第一阶段（1天）- 完成P1

**目标**：实现所有核心API逻辑，确保系统可用

**任务**：
1. ✅ 实现会话管理逻辑（4个文件）
2. ✅ 实现任务查询逻辑（2个文件）
3. ✅ 实现Proxy接口重定向（3个文件）
4. ✅ 端到端测试

**交付物**：
- 完整的会话管理API
- 完整的任务管理API
- 所有API逻辑无TODO

### 第二阶段（0.5天）- 完成P2

**目标**：实现AccessProfile管理

**任务**：
1. ✅ 实现AccessProfile CRUD（5个文件）
2. ✅ 集成到CreateSession逻辑
3. ✅ 测试验证

**交付物**：
- AccessProfile管理功能
- 支持从AccessProfile获取连接信息

### 第三阶段（1.5天）- 完成P3

**目标**：实现高级功能

**任务**：
1. ✅ CMDB集成
2. ✅ 任务结果自动收集
3. ✅ 会话超时清理
4. ⚠️ 错误重试（可选）
5. ⚠️ 监控告警（可选）

**交付物**：
- 完整的自动化运维系统
- 监控面板

## 📝 验收标准

### P1验收标准

- [ ] 所有会话管理接口可正常调用
- [ ] 所有任务管理接口可正常调用
- [ ] CloseSession正确更新会话状态
- [ ] GetTaskResult返回完整的任务结果
- [ ] 无TODO代码

### P2验收标准

- [ ] AccessProfile CRUD接口可正常调用
- [ ] CreateSession支持从AccessProfile获取配置
- [ ] 支持跳板机链路配置

### P3验收标准

- [ ] CMDB集成成功，自动获取CI信息
- [ ] 任务结果自动收集，数据库有记录
- [ ] 会话自动过期并清理
- [ ] Prometheus指标正常上报

## 🔗 依赖关系

```
Phase 4
  ├── P1: 核心API完善（无依赖）
  │   ├── 会话管理逻辑
  │   ├── 任务管理逻辑
  │   └── Proxy接口重定向
  │
  ├── P2: AccessProfile管理（依赖P1）
  │   └── 需要RPC支持
  │
  └── P3: 高级功能（依赖P1+P2）
      ├── CMDB集成（需要CMDB服务）
      ├── 任务结果收集（依赖ReportTaskResult）
      ├── 会话超时清理（依赖会话管理）
      └── 错误重试（依赖任务管理）
```

## 📚 相关文档

- **Phase 1**: `/opt/code/newbee/ops-center/migrations/phase1_summary.md`
- **Phase 2**: `/opt/code/newbee/ops-center/migrations/phase2_summary.md`
- **Phase 3**: `/opt/code/newbee/ops-center/migrations/phase3_summary.md`
- **API定义**: `/opt/code/newbee/ops-center/api/desc/ops.api`

---

**规划时间**: 2025-12-28
**预计完成**: P1完成后系统即可基本使用，P2+P3为增强功能
