# Worker健康检查机制说明

## 概述

Worker健康检查机制通过主动HTTP探测和被动TTL监控双重机制，实现自动故障检测和状态转换。

## 核心组件

### 1. HealthChecker (health_checker.go)

**职责**: 执行HTTP健康检查

**主要方法**:
- `CheckWorker()` - 单个Worker健康检查
- `CheckWorkers()` - 批量并发检查（最多20个并发）
- `BatchDetailedCheck()` - 批量详细检查（返回统计信息）
- `CheckWorkerWithRetry()` - 带重试的检查

**检查逻辑**:
```
1. 发送HTTP GET请求到 worker.HealthCheckURL
2. 超时时间: 5秒（可配置）
3. 成功条件: HTTP状态码 200-299
4. 记录响应延迟
5. 返回成功/失败结果
```

### 2. WorkerRegistry (worker_registry.go)

**职责**: 管理Worker状态转换

**UpdateHealthCheck() 方法**:
```go
// 健康检查成功
if success {
    HealthCheckFailures = 0
    Status = "online"  // 恢复为在线
}

// 健康检查失败
if !success {
    HealthCheckFailures++

    // 状态转换规则
    if HealthCheckFailures >= 5 {
        Status = "offline"   // 连续5次失败 → 离线
    } else if HealthCheckFailures >= 3 {
        Status = "degraded"  // 连续3次失败 → 降级
    }
}
```

### 3. WorkerManager (worker_manager.go)

**职责**: 协调健康检查流程

**healthCheckLoop() 方法**:
- 每30秒执行一次健康检查循环
- 调用 `performHealthChecks()`

**performHealthChecks() 方法**:
1. 获取所有Worker列表
2. 使用 `BatchDetailedCheck()` 并发检查
3. 更新每个Worker的健康状态
4. 记录状态变化和统计信息

## 状态转换流程

```
┌─────────┐
│  online │ ◄──────────────────────┐
└────┬────┘                        │
     │                              │
     │ 连续3次健康检查失败           │ 健康检查成功
     │                              │
     ▼                              │
┌──────────┐                       │
│ degraded │                       │
└────┬─────┘                       │
     │                              │
     │ 连续5次健康检查失败           │
     │                              │
     ▼                              │
┌─────────┐                        │
│ offline │ ───────────────────────┘
└─────────┘
```

## 双重监控机制

### 主动健康探测 (Active Health Check)

**触发条件**:
- Worker配置了 `health_check_url`
- 每30秒自动执行

**检查方式**:
- HTTP GET请求
- 5秒超时
- 并发检查（信号量限制20个并发）

**适用场景**:
- Worker提供HTTP健康端点
- 需要快速检测故障

### 被动TTL监控 (Passive TTL Monitoring)

**触发条件**:
- Worker超过3分钟未发送心跳

**检查方式**:
- 每30秒扫描一次所有Worker
- 检查 `last_heartbeat` 时间
- 超过TTL自动标记为offline

**适用场景**:
- Worker网络断开
- Worker进程崩溃
- 心跳机制失效

## 自动故障转移

**选择策略自动排除offline Worker**:
```go
func filterByStatus(workers []*WorkerInfo) []*WorkerInfo {
    result := make([]*WorkerInfo, 0)
    for _, w := range workers {
        if w.Status != "offline" {
            result = append(result, w)
        }
    }
    return result
}
```

**加权策略自动降低degraded Worker权重**:
```go
switch w.Status {
case "online":
    // 保持原权重
case "degraded":
    weight = weight / 2  // 降级状态权重减半
case "offline":
    weight = 0           // 离线状态权重为0
}
```

## 配置参数

```go
type WorkerManagerConfig struct {
    TTL                 time.Duration  // 3分钟（心跳TTL）
    HealthCheckInterval time.Duration  // 30秒（检查间隔）
    HealthCheckTimeout  time.Duration  // 5秒（请求超时）
}
```

## 健康检查统计

每次健康检查完成后记录:
- **total_checked**: 检查的Worker总数
- **success**: 成功次数
- **failure**: 失败次数
- **avg_latency_ms**: 平均响应延迟
- **status_changed**: 状态变化次数
- **status_changes**: 状态变化详情 (如: online→degraded: 2次)

## 示例日志输出

```
[INFO] Starting health check cycle worker_count=5
[INFO] Health check passed worker_id=worker-001 status_code=200 latency_ms=23
[ERROR] Health check failed worker_id=worker-002 status_code=503 latency_ms=156
[INFO] Worker status changed worker_id=worker-002 previous_status=online new_status=degraded failure_count=3
[INFO] Health check cycle completed total_checked=5 success=4 failure=1 avg_latency_ms=45 status_changed=1
[INFO] Status change summary change=online->degraded count=1
```

## 性能特点

✅ **并发执行**: 最多20个并发HTTP请求
✅ **超时控制**: 5秒超时，避免阻塞
✅ **连接复用**: HTTP连接池，减少开销
✅ **无阻塞**: 健康检查不阻塞Worker选择
✅ **自动恢复**: 健康检查成功自动恢复online状态

## Worker客户端要求

Worker需要提供健康检查端点:

```go
// 示例: Worker HTTP健康端点
// GET /health
// 响应: HTTP 200 OK
{
    "status": "healthy",
    "timestamp": "2025-12-17T10:30:00Z"
}
```

**配置示例**:
```yaml
Worker:
  HealthCheckURL: "http://127.0.0.1:8889/health"
```

## 最佳实践

1. **设置合理的健康检查端点** - 轻量级检查，避免耗时操作
2. **配置适当的超时** - 5秒超时适合大多数场景
3. **监控状态变化** - 关注频繁的状态切换
4. **结合TTL监控** - 双重保障，提高可靠性
5. **定期审查日志** - 识别不稳定的Worker
