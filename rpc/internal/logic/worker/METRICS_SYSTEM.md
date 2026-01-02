# Worker监控指标系统说明

## 概述

Worker监控指标系统实现了完整的性能监控、历史数据追踪和统计分析功能。

## 核心功能

### 1. 指标收集（Metrics Collection）

#### 1.1 实时指标收集

**触发方式**：心跳上报时更新
**频率**：Worker发送心跳时（30秒）
**存储位置**：`ops_workers`表 + 内存`WorkerRegistry`

**收集的实时指标**：
- CPU使用率 (`cpu_usage`)
- 内存使用率 (`memory_usage`)
- 磁盘使用率 (`disk_usage`)
- 网络流量 (`network_in`, `network_out`)
- 活跃会话数 (`active_sessions`)
- 累计请求数 (`total_requests`)
- 成功/失败计数 (`success_count`, `failure_count`)

#### 1.2 历史指标快照

**触发方式**：定时任务
**频率**：每分钟一次
**存储位置**：`ops_worker_metrics`表
**保留时间**：7天（自动清理）

**实现函数**：`collectMetricsSnapshot()`

```go
// 每分钟自动执行
func (m *WorkerManager) collectMetricsSnapshot() {
    // 1. 获取所有Worker当前状态
    // 2. 计算增量值（delta）
    // 3. 记录健康检查延迟
    // 4. 保存到ops_worker_metrics表
    // 5. 更新previousMetrics缓存
}
```

### 2. 增量计算（Delta Calculation）

为了准确反映Worker的资源消耗和业务增长，系统会计算相邻两次快照之间的增量值。

#### 2.1 计算逻辑

```go
// 计算增量（处理计数器溢出）
if current >= previous {
    delta = current - previous
} else {
    delta = 0  // 计数器重置或溢出
}
```

#### 2.2 增量字段

| 字段 | 说明 |
|------|------|
| `network_in_delta` | 入站流量增量（字节） |
| `network_out_delta` | 出站流量增量（字节） |
| `request_count_delta` | 请求数增量 |
| `success_count_delta` | 成功请求增量 |
| `failure_count_delta` | 失败请求增量 |

### 3. 延迟追踪（Latency Tracking）

#### 3.1 健康检查延迟

**数据来源**：健康检查HTTP请求的响应时间
**更新频率**：每次健康检查（30秒）
**存储方式**：内存map `healthCheckLatency`

**实现逻辑**：
```go
func (m *WorkerManager) performHealthChecks() {
    stats, results := m.healthChecker.BatchDetailedCheck(...)

    for _, result := range results {
        if result.Success {
            // 记录成功的健康检查延迟
            m.healthCheckLatency[result.WorkerID] = float64(result.LatencyMs)
        }
    }
}
```

#### 3.2 平均延迟计算

**记录到指标表**：
```go
avgLatency, _ := m.healthCheckLatency[w.WorkerID]
metric.SetAvgLatencyMs(avgLatency)
```

### 4. 指标查询API

#### 4.1 查询历史指标

**方法**：`GetWorkerMetrics()`

**参数**：
```go
type WorkerMetricsQueryOptions struct {
    WorkerID  string     // 必须：Worker ID
    TenantID  uint64     // 必须：租户ID
    StartTime *time.Time // 可选：开始时间
    EndTime   *time.Time // 可选：结束时间
    Limit     int        // 可选：限制返回数量（默认100，最大1000）
}
```

**返回**：按时间降序排列的指标快照列表

**示例**：
```go
startTime := time.Now().Add(-24 * time.Hour)
endTime := time.Now()

metrics, err := workerManager.GetWorkerMetrics(ctx, &worker.WorkerMetricsQueryOptions{
    WorkerID:  "worker-001",
    TenantID:  1,
    StartTime: &startTime,
    EndTime:   &endTime,
    Limit:     200,
})

// metrics: 最近24小时的200条指标记录
```

#### 4.2 查询指标统计

**方法**：`GetWorkerMetricsStats()`

**返回统计信息**：
```go
type WorkerMetricsStats struct {
    WorkerID        string    // Worker ID
    StartTime       time.Time // 统计开始时间
    EndTime         time.Time // 统计结束时间
    AvgCPU          float64   // 平均CPU使用率
    MaxCPU          float64   // 最大CPU使用率
    AvgMemory       float64   // 平均内存使用率
    MaxMemory       float64   // 最大内存使用率
    AvgSessions     float64   // 平均会话数
    MaxSessions     int       // 最大会话数
    TotalRequests   uint64    // 总请求数（增量累加）
    TotalSuccess    uint64    // 总成功数
    TotalFailure    uint64    // 总失败数
    SuccessRate     float64   // 成功率（%）
    AvgLatency      float64   // 平均延迟(ms)
    DataPoints      int       // 数据点数量
}
```

**示例**：
```go
stats, err := workerManager.GetWorkerMetricsStats(ctx, &worker.WorkerMetricsQueryOptions{
    WorkerID:  "worker-001",
    TenantID:  1,
    StartTime: &startTime,
    EndTime:   &endTime,
})

// stats.AvgCPU: 15.3%
// stats.SuccessRate: 99.8%
// stats.AvgLatency: 23.5ms
```

#### 4.3 查询所有Worker汇总

**方法**：`GetAllWorkersMetricsSummary()`

**返回**：当前所有Worker的实时指标汇总

**示例输出**：
```json
{
    "total_workers": 10,
    "online_count": 8,
    "degraded_count": 1,
    "offline_count": 1,
    "total_sessions": 156,
    "total_requests": 1523456,
    "total_success": 1519234,
    "total_failure": 4222,
    "success_rate": 99.72,
    "avg_cpu_usage": 23.5,
    "avg_memory_usage": 45.8
}
```

### 5. 自动清理（Auto Cleanup）

**触发方式**：后台定时任务
**频率**：每小时一次
**保留策略**：删除7天前的指标数据

**实现函数**：`cleanOldMetrics()`

```go
func (m *WorkerManager) cleanOldMetrics() {
    threshold := time.Now().Add(-m.config.MetricsRetention) // 7天前

    deleted, err := m.db.WorkerMetrics.Delete().
        Where(workermetrics.TimestampLT(threshold)).
        Exec(systemCtx)

    // 日志记录删除数量
}
```

**清理日志示例**：
```
[INFO] Old metrics cleaned deleted_count=10080 threshold=2025-12-10T02:34:18Z
```

## 数据结构设计

### 上一次指标快照（WorkerMetricsSnapshot）

用于计算增量值的内存缓存：

```go
type WorkerMetricsSnapshot struct {
    Timestamp     time.Time
    NetworkIn     uint64
    NetworkOut    uint64
    TotalRequests uint64
    SuccessCount  uint64
    FailureCount  uint64
}
```

**存储位置**：`WorkerManager.previousMetrics` (map[string]*WorkerMetricsSnapshot)

### 健康检查延迟缓存

```go
// WorkerManager
healthCheckLatency map[string]float64 // workerID -> 最近健康检查延迟(ms)
```

**更新时机**：每次健康检查成功后

## 性能特点

✅ **增量计算** - 准确反映资源消耗变化
✅ **时序存储** - 支持历史数据趋势分析
✅ **自动清理** - 防止数据无限增长
✅ **并发安全** - 使用读写锁保护共享数据
✅ **低开销** - 每分钟一次快照，不影响主流程
✅ **灵活查询** - 支持时间范围、数量限制

## 监控指标类型

### 资源指标（Resource Metrics）
- CPU使用率
- 内存使用率
- 磁盘使用率
- 网络流量

### 业务指标（Business Metrics）
- 活跃会话数
- 请求总数
- 成功/失败次数
- 成功率

### 性能指标（Performance Metrics）
- 健康检查延迟
- 平均响应时间

### 状态指标（Status Metrics）
- Worker状态快照
- 状态变化次数

## 使用场景

### 1. 性能监控
```go
// 查询最近1小时的指标
metrics, _ := workerManager.GetWorkerMetrics(ctx, &worker.WorkerMetricsQueryOptions{
    WorkerID:  "worker-001",
    TenantID:  1,
    StartTime: &oneHourAgo,
    Limit:     60,
})

// 绘制CPU/内存使用率曲线
```

### 2. 容量规划
```go
// 查询最近7天的统计信息
stats, _ := workerManager.GetWorkerMetricsStats(ctx, &worker.WorkerMetricsQueryOptions{
    WorkerID:  "worker-001",
    TenantID:  1,
    StartTime: &sevenDaysAgo,
})

// 分析：
// - 平均会话数：评估负载
// - 峰值会话数：评估容量上限
// - 成功率：评估稳定性
```

### 3. 故障排查
```go
// 查询故障时间段的指标
metrics, _ := workerManager.GetWorkerMetrics(ctx, &worker.WorkerMetricsQueryOptions{
    WorkerID:  "worker-001",
    TenantID:  1,
    StartTime: &incidentStart,
    EndTime:   &incidentEnd,
})

// 分析资源异常、请求失败等
```

### 4. 实时监控大屏
```go
// 获取所有Worker实时汇总
summary := workerManager.GetAllWorkersMetricsSummary(tenantID)

// 显示：
// - 在线Worker数量
// - 总会话数
// - 实时成功率
// - 平均资源使用率
```

## 数据库表结构

### ops_worker_metrics

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint64 | 主键 |
| worker_id | string | Worker ID |
| timestamp | time | 采集时间 |
| cpu_usage | float64 | CPU使用率 |
| memory_usage | float64 | 内存使用率 |
| disk_usage | float64 | 磁盘使用率 |
| network_in_delta | uint64 | 入站流量增量 |
| network_out_delta | uint64 | 出站流量增量 |
| active_sessions | int | 活跃会话数 |
| request_count_delta | uint64 | 请求增量 |
| success_count_delta | uint64 | 成功增量 |
| failure_count_delta | uint64 | 失败增量 |
| avg_latency_ms | float64 | 平均延迟 |
| worker_status | string | Worker状态快照 |
| tenant_id | uint64 | 租户ID |

**索引**：
- `idx_worker_time` - (worker_id, timestamp) - 时序查询
- `idx_tenant_time` - (tenant_id, timestamp) - 租户查询
- `idx_timestamp` - (timestamp) - TTL清理

## 配置参数

```go
type WorkerManagerConfig struct {
    MetricsRetention time.Duration // 指标保留时间（默认7天）
}
```

## 最佳实践

1. **合理设置查询限制** - 避免一次查询过多数据（建议≤1000条）
2. **使用时间范围过滤** - 减少数据库扫描范围
3. **定期监控清理日志** - 确认自动清理正常工作
4. **结合实时指标和历史指标** - 实时用Registry，历史用Metrics表
5. **监控指标表增长** - 确保不会无限增长

## 故障排查

### 问题1：指标数据缺失

**原因**：
- Worker未正常发送心跳
- 指标收集任务未启动
- 数据库写入失败

**检查**：
```bash
# 查看日志
grep "Metrics snapshot collected" ops-center.log

# 检查指标表
SELECT COUNT(*) FROM ops_worker_metrics WHERE timestamp > NOW() - INTERVAL 1 HOUR;
```

### 问题2：增量值为0

**原因**：
- Worker计数器未更新
- 首次快照（没有previous值）
- 计数器溢出重置

**正常**：首次快照的delta值为0是正常的

### 问题3：指标表快速增长

**原因**：
- Worker数量过多
- 自动清理失败

**检查**：
```bash
# 查看清理日志
grep "Old metrics cleaned" ops-center.log

# 检查指标数量
SELECT DATE(timestamp), COUNT(*) FROM ops_worker_metrics GROUP BY DATE(timestamp);
```

---

**实施日期**：2025-12-17
**Phase 5状态**：✅ 完成
**测试覆盖率**：待后续集成测试
