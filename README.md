# 新蜂资产管理平台 — 运维中心

负责平台 Worker/Agent 管理、接入画像和任务调度，与 Proxy、主机 Agent 及核心权限服务协作。单个仓库内包含 `api` 和 `rpc` 两个独立 Go 模块。

仓库：[coder-lulu/newbee-ops](https://github.com/coder-lulu/newbee-ops) · [平台工作区](https://github.com/coder-lulu/newbee)

## 获取代码

推荐通过完整工作区开发，保留兄弟模块目录及本地 `replace` 依赖。以下命令使用 Bash；Go 工作区要求 Go 1.25.1 或更高版本。

```bash
git clone --recurse-submodules https://github.com/coder-lulu/newbee.git
cd newbee/ops-center
```

已有工作区执行 `git submodule update --init --recursive`。单独克隆模块时，需要自行补齐 `go.mod` 中的本地依赖路径。

## 目录导航

| 目录 | 用途 |
| --- | --- |
| `api/` | HTTP 接口、Worker 客户端与认证接入 |
| `rpc/` | RPC 服务、任务分发与持久化逻辑 |
| `migrations/` | 数据库迁移文件 |
| `docs/` | Worker 管理与 Agent 接入文档 |

## 配置与本地运行

首次分别复制 `api/etc/ops.yaml.example` 和 `rpc/etc/ops.yaml.example` 到同目录的 `ops.yaml`，填写数据库、Redis、Core RPC 及 Worker 接入参数。示例中的地址和端口需要按实际环境调整；示例 API 端口为 `9601`，RPC 为 `9600`。

API 和 RPC 入口均启用了 `conf.UseEnv()` 环境变量展开，可为配置中的 `${...}` 占位符设置对应环境变量。真实配置文件保留在本地。

基础服务就绪后，在两个终端依次运行：

```bash
# 在 ops-center/rpc 下
go run . -f etc/ops.yaml
```

```bash
# 在 ops-center/api 下
go run . -f etc/ops.yaml
```

## 首次数据库初始化

先完成 Core 的数据库初始化，并配置 RPC 的 `CoreRpc.Endpoints`。首次空库启动必须设置 `WorkerManager.Enabled: false`（示例配置已关闭），使 RPC 启动时暂不读取尚未创建的 Worker 表。Ops RPC 启动后，在平台工作区根目录运行 `bash init-databases.sh -s ops-center`；也可在本目录使用已安装的 grpcurl：

```bash
grpcurl -plaintext -import-path rpc/desc -proto ops.proto \
  -d '{}' 127.0.0.1:9600 ops.Ops/initDatabase
```

此操作按当前 Ent 模型创建 Ops 表，不删除现有列、索引或其他模块的表；重复调用保留已有业务记录。随后通过 Core RPC 为默认租户登记 Ops 菜单和 API 目录，合并默认 `superadmin` 的菜单授权，保留其原有菜单。菜单使用数据库分配的 ID，不依赖历史 SQL 中的固定编号；普通角色仍需管理员单独授权。Core 不可用、未初始化或目录登记失败会返回错误，修复后可重试。

初始化成功后设置 `WorkerManager.Enabled: true` 并重启 Ops RPC，恢复 Worker 注册表、心跳及指标后台处理，然后启动 API。关闭期间保留管理器对象，但不运行数据恢复或后台任务；省略此配置时沿用自动启用行为，真实数据库连接或查询错误仍会导致启动失败。

初始化无需导入本机数据库备份，不会创建业务任务、凭证或代理。RPC 端口应仅用于部署内网；请在启动 API 前完成初始化，以便 API 加载完整的菜单和权限配置。`migrations/` 中的历史菜单 SQL 不作为首次部署的必需步骤。

## 构建与验证

本仓库顶层没有 `go.mod`；分别进入 `rpc/`、`api/` 执行：

```bash
go build .
go test ./...
go vet ./...
```

数据库迁移应按目标环境审核后执行；这里不将历史设计文档中的完成状态视为当前部署验证结果。

## 文档

- [Worker 管理常见问题](docs/WORKER_MANAGEMENT_FAQ.md)
- [Agent 接入计划](docs/AGENT_MANAGEMENT_IMPLEMENTATION_PLAN.md)
- [Worker 项目说明](WORKER_PROJECT_SUMMARY.md)

## 许可证与来源

本仓库采用 [Apache-2.0](LICENSE)。沿用仓库现有上游许可和文件版权声明。第三方依赖遵循各自许可证，保留原有版权与许可声明。


### Proxy 机器注册的租户绑定

Ops API 的 `Ops.Registration.PSK` 绑定单个租户，必须同时设置 `Ops.Registration.TenantID` 为已存在且授权该 Proxy 的租户 ID。示例值 `0` 表示未配置，注册和心跳会返回业务码 `503`；空 PSK 或不匹配的 PSK 返回 `401`。客户端提交的租户信息不能覆盖此绑定。注册、查询与心跳均携带该租户上下文调用 RPC，继续执行 Ent 租户隔离。不同租户应使用独立配置的注册入口和 PSK。
