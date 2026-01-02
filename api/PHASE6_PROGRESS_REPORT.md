# Phase 6 进度报告：Worker API实现

## 执行时间
**开始时间**: 2025-12-17
**当前状态**: 进行中（60%）

---

## 已完成工作 ✅

### 1. API定义完成

**文件**: `/opt/code/newbee/ops-center/api/desc/worker.api`

**定义的接口** (12个)：
1. ✅ POST `/worker/register` - Worker注册（PSK认证）
2. ✅ POST `/worker/heartbeat` - Worker心跳（PSK认证）
3. ✅ POST `/worker/pick` - Worker选择（JWT认证）
4. ✅ GET `/worker/list` - Worker列表（JWT认证）
5. ✅ GET `/worker/:id` - Worker详情（JWT认证）
6. ✅ POST `/worker/weight` - 更新权重（JWT认证）
7. ✅ POST `/worker/:id/activate` - 手动上线（JWT认证）
8. ✅ POST `/worker/:id/deactivate` - 手动下线（JWT认证）
9. ✅ POST `/worker/delete` - 删除Worker（JWT认证）
10. ✅ GET `/worker/metrics` - 查询指标（JWT认证）
11. ✅ GET `/worker/metrics/stats` - 查询统计（JWT认证）

**定义的类型** (19个)：
- WorkerItem, WorkerDetail
- WorkerRegisterReq, WorkerHeartbeatReq
- WorkerPickReq, WorkerPickResp
- WorkerListReq, WorkerListResp
- UpdateWorkerWeightReq
- WorkerMetricsReq, WorkerMetricsItem, WorkerMetricsResp
- WorkerMetricsStatsReq, WorkerMetricsStatsResp
- BaseResp, BaseIDResp, IDReq, IDsReq

### 2. 代码生成完成

**生成的Handler文件** (11个)：
- `/opt/code/newbee/ops-center/api/internal/handler/worker/*.go`
- 所有handler已生成，包括路由注册

**生成的Logic文件** (11个)：
- `/opt/code/newbee/ops-center/api/internal/logic/worker/*.go`
- 所有logic骨架已生成，待实现业务逻辑

### 3. 数据访问层完成

**文件**: `/opt/code/newbee/ops-center/api/internal/workerclient/worker_client.go`

**实现的方法**：
- `ValidatePSK()` - PSK验证
- `GetWorkerByWorkerID()` - 根据WorkerID查询
- `GetWorkerByID()` - 根据ID查询
- `ListWorkers()` - 列表查询（支持过滤、分页）
- `UpdateWorkerWeight()` - 更新权重
- `ActivateWorker()` - 激活Worker
- `DeactivateWorker()` - 停用Worker
- `DeleteWorkers()` - 批量删除
- `GetWorkerMetrics()` - 查询指标

### 4. 配置更新完成

**文件**: `/opt/code/newbee/ops-center/api/internal/config/config.go`

**新增配置**：
```go
DatabaseConf commoncfg.DatabaseConf // 数据库配置
```

---

## 待完成工作 ❌

### 1. ServiceContext更新

**文件**: `/opt/code/newbee/ops-center/api/internal/svc/service_context.go`

**需要添加**：
```go
type ServiceContext struct {
    // ... existing fields ...

    DB           *ent.Client      // 数据库连接
    WorkerClient *workerclient.WorkerClient  // Worker客户端
}
```

**初始化代码**：
```go
// 初始化数据库连接
db := ent.NewClient(
    ent.Driver(c.DatabaseConf.NewNoCacheDriver()),
    ent.Debug(),
)

// 初始化租户Hook
if err := hooks.QuickSetup(db); err != nil {
    panic("Hook初始化失败: " + err.Error())
}

// 初始化WorkerClient
workerClient := workerclient.NewWorkerClient(db)

svcCtx.DB = db
svcCtx.WorkerClient = workerClient
```

### 2. Logic层实现（11个文件）

#### 2.1 worker_register_logic.go

**需要实现**：
1. PSK验证
2. 检查Worker是否已存在
3. 创建或更新Worker记录
4. 更新内存注册表（如果需要）

**关键逻辑**：
```go
// PSK验证
if !l.svcCtx.WorkerClient.ValidatePSK(req.PSK, l.svcCtx.Config.Ops.Registration.PSK) {
    return &types.BaseResp{Code: 401, Msg: "Invalid PSK"}, nil
}

// 检查是否已存在
existing, err := l.svcCtx.WorkerClient.GetWorkerByWorkerID(l.ctx, req.WorkerID, tenantID)

// 创建或更新
if ent.IsNotFound(err) {
    // 创建新Worker
} else {
    // 更新现有Worker
}
```

#### 2.2 worker_heartbeat_logic.go

**需要实现**：
1. PSK验证
2. 更新Worker心跳时间
3. 更新Worker负载指标

**关键逻辑**：
```go
// 更新心跳和指标
err := l.svcCtx.DB.Worker.Update().
    Where(worker.WorkerIDEQ(req.WorkerID)).
    SetLastHeartbeat(time.Now()).
    SetCPUUsage(req.CPUUsage).
    SetMemoryUsage(req.MemoryUsage).
    // ... other metrics
    Exec(l.ctx)
```

#### 2.3 worker_pick_logic.go

**需要实现**：
1. 获取租户的所有Worker
2. 根据策略选择Worker
3. 返回Worker信息

**重要**：需要实现选择策略逻辑，或者调用RPC服务的WorkerManager.SelectWorker方法

#### 2.4 get_worker_list_logic.go

**需要实现**：
1. 构建查询条件
2. 分页查询
3. 返回列表数据

**已有WorkerClient支持**：`ListWorkers()`方法

#### 2.5 get_worker_by_id_logic.go

**需要实现**：
1. 根据ID查询Worker
2. 返回详细信息

**已有WorkerClient支持**：`GetWorkerByID()`方法

#### 2.6 update_worker_weight_logic.go

**需要实现**：
1. 验证weight范围（1-1000）
2. 更新Worker权重

**已有WorkerClient支持**：`UpdateWorkerWeight()`方法

#### 2.7 activate_worker_logic.go / deactivate_worker_logic.go

**需要实现**：
1. 更新Worker状态

**已有WorkerClient支持**：`ActivateWorker()` / `DeactivateWorker()`方法

#### 2.8 delete_worker_logic.go

**需要实现**：
1. 批量删除Worker

**已有WorkerClient支持**：`DeleteWorkers()`方法

#### 2.9 get_worker_metrics_logic.go

**需要实现**：
1. 查询历史指标
2. 返回指标数据

**已有WorkerClient支持**：`GetWorkerMetrics()`方法

#### 2.10 get_worker_metrics_stats_logic.go

**需要实现**：
1. 查询指标
2. 计算统计信息（平均值、最大值、成功率等）

**逻辑复杂**：需要实现统计计算算法

### 3. PSK认证中间件配置

**文件**: 配置文件（etc/ops-api.yaml）

**需要添加跳过路径**：
```yaml
Middleware:
  auth:
    skipPaths:
      - "/worker/register"
      - "/worker/heartbeat"
  audit:
    skipPaths:
      - "/worker/register"
      - "/worker/heartbeat"
  tenantCheck:
    skipPaths:
      - "/worker/register"
      - "/worker/heartbeat"
```

**注意**：根据统一中间件框架，JWT认证会自动应用，无需在API文件中配置

### 4. PSK验证逻辑

**建议位置**：创建一个自定义中间件 `internal/middleware/psk_auth.go`

**逻辑**：
```go
// PSK认证中间件（仅用于/worker/register和/worker/heartbeat）
func PSKAuthMiddleware(psk string) rest.Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            // 从Header或Body中获取PSK
            requestPSK := r.Header.Get("X-OPS-PSK")
            if requestPSK == "" {
                // 从Body中获取（已在handler中解析）
                requestPSK = ... // 需要从context获取
            }

            if requestPSK != psk {
                httpx.ErrorCtx(r.Context(), w, fmt.Errorf("invalid PSK"))
                return
            }

            next(w, r)
        }
    }
}
```

**或者**：直接在logic层验证PSK（更简单）

---

## 架构说明

### 当前架构选择

**API服务直接访问数据库** (类似ProxyRegistry)

**原因**：
1. Worker注册和心跳需要快速响应
2. 避免额外的RPC调用开销
3. 与ProxyRegistry保持一致的架构模式

**优点**：
- 快速实现
- 低延迟

**缺点**：
- 打破了API-RPC分层
- 无法使用RPC服务中的WorkerManager（选择策略等）

### 待优化架构

**长期方案**：API服务通过RPC调用Worker管理功能

**需要做的工作**：
1. 在RPC服务的proto文件中定义Worker管理的RPC方法
2. 生成RPC代码
3. RPC服务的logic层调用WorkerManager
4. API服务的logic层调用OpsRpc客户端

**优点**：
- 清晰的架构分层
- 重用RPC服务中的WorkerManager（选择策略等）
- 更好的可维护性

---

## 关键问题和解决方案

### 问题1：Worker选择策略

**现状**：选择策略在RPC服务的WorkerManager中实现

**当前方案**：
- 简化实现：在API服务的logic层实现基础的选择策略（least_connections）
- 获取所有online状态的Worker
- 按active_sessions升序排序
- 返回第一个Worker

**长期方案**：通过RPC调用WorkerManager.SelectWorker()

### 问题2：PSK认证

**方案1**：自定义中间件
- 优点：统一处理
- 缺点：需要配置路由

**方案2**：Logic层验证
- 优点：简单直接
- 缺点：每个logic都需要验证

**推荐**：方案2（Logic层验证），因为只有2个接口需要PSK认证

### 问题3：内存注册表

**问题**：RPC服务中的WorkerManager有内存注册表，API服务是否需要？

**方案**：
- API服务不需要内存注册表
- 直接从数据库查询（性能足够好）
- Worker选择频率不高（不像session管理）

---

## 测试计划

### 单元测试

1. **WorkerClient测试**
   - PSK验证
   - 数据库CRUD操作
   - 查询过滤和分页

2. **Logic层测试**
   - 注册流程
   - 心跳更新
   - Worker选择

### 集成测试

1. **完整流程测试**
   - Worker注册 → 心跳 → 选择 → 查询
   - Worker权重调整 → 选择结果变化
   - Worker下线 → 选择排除offline Worker

2. **PSK认证测试**
   - 正确的PSK → 成功
   - 错误的PSK → 401错误
   - 无PSK → 401错误

3. **并发测试**
   - 多个Worker同时注册
   - 高频心跳
   - 并发选择

---

## 下一步工作计划

### 立即完成（Phase 6剩余工作）

1. **更新ServiceContext** (30分钟)
   - 添加DB和WorkerClient
   - 更新初始化逻辑

2. **实现关键Logic文件** (2-3小时)
   - worker_register_logic.go ⭐
   - worker_heartbeat_logic.go ⭐
   - worker_pick_logic.go ⭐
   - get_worker_list_logic.go
   - get_worker_by_id_logic.go

3. **实现辅助Logic文件** (1-2小时)
   - update_worker_weight_logic.go
   - activate/deactivate_worker_logic.go
   - delete_worker_logic.go
   - get_worker_metrics_logic.go
   - get_worker_metrics_stats_logic.go

4. **配置更新** (15分钟)
   - 更新etc/ops-api.yaml
   - 添加数据库配置
   - 添加PSK配置

5. **编译测试** (30分钟)
   - 编译API服务
   - 解决编译错误
   - 基础功能测试

### Phase 7: Worker客户端适配

**预计工作量**：1-2天
- 更新Worker客户端配置
- 修改注册/心跳逻辑
- 增加指标采集
- 测试连通性

---

## 当前进度：60% ✅

**完成项**：
- ✅ API定义（100%）
- ✅ 代码生成（100%）
- ✅ 数据访问层（100%）
- ✅ 配置更新（100%）

**进行中**：
- ⏳ ServiceContext更新（0%）
- ⏳ Logic层实现（0%）

**待开始**：
- ❌ PSK认证配置
- ❌ 编译测试
- ❌ 集成测试

---

**报告生成时间**：2025-12-17
**预计完成时间**：需要额外4-6小时开发时间
**状态**：进行中
