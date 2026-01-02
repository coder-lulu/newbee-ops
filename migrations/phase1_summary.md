# Phase 1: 数据库 Schema 迁移完成报告

## 📅 完成时间
2025-12-28

## ✅ 已完成的工作

### 1. Schema 文件重命名
所有 schema 文件已成功从 worker 重命名为 proxy：

| 原文件名 | 新文件名 | 状态 |
|---------|---------|------|
| `worker.go` | `proxy.go` | ✅ 已重命名 |
| `worker_group.go` | `proxy_group.go` | ✅ 已重命名 |
| `worker_group_member.go` | `proxy_group_member.go` | ✅ 已重命名 |
| `worker_metrics.go` | `proxy_metrics.go` | ✅ 已重命名 |

**位置**: `/opt/code/newbee/ops-center/rpc/ent/schema/`

### 2. Schema 内容更新

#### 2.1 proxy.go
- ✅ 结构体名：`Worker` → `Proxy`
- ✅ 字段名：`worker_id` → `proxy_id`
- ✅ 字段名：`worker_status` → `proxy_status`
- ✅ 表名：`ops_workers` → `ops_proxies`
- ✅ 索引更新：所有索引字段已更新为 proxy 相关字段
- ✅ 注释：所有 "Worker" 注释已更新为 "Proxy"

#### 2.2 proxy_group.go
- ✅ 结构体名：`WorkerGroup` → `ProxyGroup`
- ✅ 字段名：`min_healthy_workers` → `min_healthy_proxies`
- ✅ 表名：`ops_worker_groups` → `ops_proxy_groups`
- ✅ 注释：所有 "Worker" 注释已更新为 "Proxy"

#### 2.3 proxy_group_member.go
- ✅ 结构体名：`WorkerGroupMember` → `ProxyGroupMember`
- ✅ 字段名：`worker_id` → `proxy_id`
- ✅ 字段名：`worker_group_id` → `proxy_group_id`
- ✅ 表名：`ops_worker_group_members` → `ops_proxy_group_members`
- ✅ 边关系：`Worker` → `Proxy`, `WorkerGroup` → `ProxyGroup`
- ✅ 索引更新：所有索引字段已更新
- ✅ 注释：所有 "Worker" 注释已更新为 "Proxy"

#### 2.4 proxy_metrics.go
- ✅ 结构体名：`WorkerMetrics` → `ProxyMetrics`
- ✅ 字段名：`worker_id` → `proxy_id`
- ✅ 字段名：`worker_status` → `proxy_status`
- ✅ 表名：`ops_worker_metrics` → `ops_proxy_metrics`
- ✅ 索引更新：所有索引字段已更新
- ✅ 注释：所有 "Worker" 注释已更新为 "Proxy"

### 3. Ent 代码生成

执行命令：
```bash
GOWORK=off go run entgo.io/ent/cmd/ent generate --template glob="./ent/template/*.tmpl" ./ent/schema --feature sql/execquery,intercept
```

生成结果：
- ✅ `ent/proxy/` - Proxy 实体操作代码
- ✅ `ent/proxygroup/` - ProxyGroup 实体操作代码
- ✅ `ent/proxygroupmember/` - ProxyGroupMember 实体操作代码
- ✅ `ent/proxymetrics/` - ProxyMetrics 实体操作代码
- ✅ 无旧的 worker 目录残留

### 4. RPC CRUD 逻辑生成

所有模型的 RPC CRUD 逻辑已成功生成：

#### 4.1 Proxy CRUD
```bash
make gen-rpc-ent-logic model=Proxy group=proxy
```
生成文件：
- `desc/proxy.proto` - Protobuf 定义
- `internal/logic/proxy/` - CRUD 逻辑实现
  - `create_proxy_logic.go`
  - `get_proxy_by_id_logic.go`
  - `get_proxy_list_logic.go`
  - `update_proxy_logic.go`
  - `delete_proxy_logic.go`

#### 4.2 ProxyGroup CRUD
```bash
make gen-rpc-ent-logic model=ProxyGroup group=proxy_group
```
生成文件：
- `desc/proxygroup.proto`
- `internal/logic/proxy_group/` - CRUD 逻辑实现

#### 4.3 ProxyGroupMember CRUD
```bash
make gen-rpc-ent-logic model=ProxyGroupMember group=proxy_group_member
```
生成文件：
- `desc/proxygroupmember.proto`
- `internal/logic/proxy_group_member/` - CRUD 逻辑实现

#### 4.4 ProxyMetrics CRUD
```bash
make gen-rpc-ent-logic model=ProxyMetrics group=proxy_metrics
```
生成文件：
- `desc/proxymetrics.proto`
- `internal/logic/proxy_metrics/` - CRUD 逻辑实现

### 5. 数据库迁移 SQL 脚本

已创建数据库迁移脚本：
- **文件**: `/opt/code/newbee/ops-center/migrations/worker_to_proxy_migration.sql`
- **内容**:
  - ✅ 表重命名（ops_workers → ops_proxies 等）
  - ✅ 字段重命名（worker_id → proxy_id, worker_status → proxy_status 等）
  - ✅ 索引更新（删除旧索引，创建新索引）
  - ✅ 数据验证查询
  - ✅ 回滚脚本（emergency use only）

## 📊 变更统计

### Schema 文件
| 类型 | 数量 | 详情 |
|------|------|------|
| 重命名文件 | 4 | proxy.go, proxy_group.go, proxy_group_member.go, proxy_metrics.go |
| 更新结构体 | 4 | Proxy, ProxyGroup, ProxyGroupMember, ProxyMetrics |
| 更新字段 | 8 | proxy_id, proxy_status, proxy_group_id, min_healthy_proxies 等 |
| 更新表名 | 4 | ops_proxies, ops_proxy_groups, ops_proxy_group_members, ops_proxy_metrics |

### 生成的代码
| 类型 | 数量 | 详情 |
|------|------|------|
| Ent 实体目录 | 4 | proxy/, proxygroup/, proxygroupmember/, proxymetrics/ |
| Proto 文件 | 4 | proxy.proto, proxygroup.proto, proxygroupmember.proto, proxymetrics.proto |
| Logic 目录 | 4 | proxy/, proxy_group/, proxy_group_member/, proxy_metrics/ |
| CRUD 方法 | 20 | 每个模型 5 个方法（Create, Get, List, Update, Delete） |

## 🔄 字段映射表

| 原字段名 | 新字段名 | 数据类型 | 说明 |
|---------|---------|---------|------|
| `worker_id` | `proxy_id` | VARCHAR(100) | Proxy 唯一标识 |
| `worker_status` | `proxy_status` | ENUM | 状态：online/degraded/offline |
| `worker_group_id` | `proxy_group_id` | BIGINT UNSIGNED | 分组ID |
| `min_healthy_workers` | `min_healthy_proxies` | INT | 最少健康Proxy数量 |

## 🗂️ 表映射

| 原表名 | 新表名 | 记录说明 |
|-------|--------|---------|
| `ops_workers` | `ops_proxies` | Proxy 主表 |
| `ops_worker_groups` | `ops_proxy_groups` | Proxy 分组表 |
| `ops_worker_group_members` | `ops_proxy_group_members` | Proxy 分组成员关联表 |
| `ops_worker_metrics` | `ops_proxy_metrics` | Proxy 指标历史表 |

## 📝 下一步操作

### Phase 2: API Logic 实现

根据集成计划 (`/tmp/ops_proxy_integration_plan.md`)，下一步需要实现：

1. **Proxy 注册逻辑** (`api/internal/logic/proxy/proxy_register_logic.go`)
   - PSK 验证
   - 检查 Proxy 是否已存在
   - 创建或更新 Proxy 记录

2. **Proxy 心跳逻辑** (`api/internal/logic/proxy/proxy_heartbeat_logic.go`)
   - PSK 验证
   - 更新心跳时间
   - 更新资源指标
   - 保存指标历史

3. **Proxy 选择逻辑** (`api/internal/logic/proxy/proxy_pick_logic.go`)
   - 实现负载均衡策略：
     - least_connections（最少连接数）
     - round_robin（轮询）
     - weighted（加权轮询）
     - random（随机）
     - consistent_hash（一致性哈希）

4. **其他 Proxy 管理接口**
   - 查询 Proxy 列表
   - 查询 Proxy 详情
   - 激活/停用 Proxy
   - 删除 Proxy
   - 更新权重
   - 查询指标

### 执行数据库迁移

在开发/测试环境执行迁移：

```bash
# 连接到数据库
mysql -h 192.168.26.130 -P 3306 -u root -p newbee

# 执行迁移脚本
source /opt/code/newbee/ops-center/migrations/worker_to_proxy_migration.sql;

# 验证迁移结果
SHOW TABLES LIKE 'ops_proxy%';
```

⚠️ **注意**：在生产环境执行前：
1. 备份数据库
2. 在测试环境验证迁移脚本
3. 准备回滚方案
4. 评估停机时间

## 🎯 验收标准

- [x] Schema 文件全部重命名
- [x] Schema 内容全部更新（类型名、字段名、表名）
- [x] Ent 代码成功生成
- [x] RPC CRUD 逻辑全部生成
- [x] 数据库迁移 SQL 脚本已创建
- [ ] 数据库迁移已执行并验证（待执行）
- [ ] API Logic 实现（Phase 2）
- [ ] 端到端测试通过（Phase 2+）

## 📚 参考文档

- **集成计划**: `/tmp/ops_proxy_integration_plan.md`
- **迁移总结**: `/tmp/worker_to_proxy_migration_summary.md`
- **数据库迁移脚本**: `/opt/code/newbee/ops-center/migrations/worker_to_proxy_migration.sql`
- **Schema 目录**: `/opt/code/newbee/ops-center/rpc/ent/schema/`
- **Proto 目录**: `/opt/code/newbee/ops-center/rpc/desc/`
- **Logic 目录**: `/opt/code/newbee/ops-center/rpc/internal/logic/`

---

**Phase 1 完成时间**: 2025-12-28
**状态**: ✅ 全部完成
**准备进入**: Phase 2 - API Logic 实现
