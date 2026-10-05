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

API 入口启用了环境变量展开；RPC 入口未启用 `conf.UseEnv()`，必须先把配置中的 `${...}` 占位符替换为本地配置值，不能只设置环境变量。真实配置文件保留在本地。

基础服务就绪后，在两个终端依次运行：

```bash
# 在 ops-center/rpc 下
go run . -f etc/ops.yaml
```

```bash
# 在 ops-center/api 下
go run . -f etc/ops.yaml
```

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
