# Ops Center RPC（计划）

本模块用于落地 AccessProfile 等持久化实体，遵循 CLAUDE.md 规范：

1) 修改 ent Schema：`ent/schema/accessprofile.go`
2) 生成 ent 代码：
   go run entgo.io/ent/cmd/ent generate ./ent/schema --feature sql/execquery,intercept,sql/modifier
3) 生成 RPC 代码：
   make gen-rpc

实体要求：
- AccessProfile:
  - ci_id（唯一）
  - capabilities[]、ports(json)、credential_ref、prefer_proxy、jump_chain[]、tags(json)
  - mixins: IDMixin、StatusMixin（可选）、TenantMixin
- 启动时注册 TenantHook/DataPerm 拦截器；API 层挂 Authority/TenantCheck/DataPerm 中间件。

当前阶段：为 API 验证闭环提供内存版画像；RPC 将在下一阶段替换内存存储。

