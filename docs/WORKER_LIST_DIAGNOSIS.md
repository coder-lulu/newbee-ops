# Worker 列表加载问题诊断

## 问题描述
前端 Worker 管理页面加载 Worker 列表数据为空。

## 诊断步骤

### 1. 数据库检查 ✅
```bash
mysql -h192.168.26.130 -uroot -p123456 newbee -e \
  "SELECT id, worker_id, name, worker_status, tenant_id FROM ops_workers;"
```

**结果**：数据库中有 3 条 Worker 记录，tenant_id 都是 1
- agent-001 (online)
- test-worker-001 (offline)
- test-worker-integration (offline)

### 2. API 路由检查 ✅
- **前端调用**: `GET /worker/list`
- **后端定义**: `/opt/code/newbee/ops-center/api/desc/worker.api` ✅
- **Logic 实现**: `get_worker_list_logic.go` ✅

### 3. 认证要求检查 ⚠️

API 需要 JWT 认证：
```bash
curl 'http://127.0.0.1:9601/worker/list?page=1&pageSize=10'
# 返回: {"code":40001,"message":"认证Token缺失"}
```

## 可能的问题

### 问题 1: 前端未传递 JWT Token
**症状**: API 返回 40001 "认证Token缺失"

**检查方法**（浏览器开发者工具）:
1. 打开浏览器控制台 (F12)
2. 切换到 Network 标签
3. 刷新 Worker 列表页面
4. 查找 `worker/list` 请求
5. 检查 Request Headers 中是否有 `Authorization: Bearer <token>`

**原因**:
- 用户未登录或 Token 已过期
- Request interceptor 未正确添加 Token

### 问题 2: API baseURL 配置错误
**症状**: 请求发送到了错误的地址

**检查方法**:
1. 浏览器控制台 -> Network 标签
2. 查看 `worker/list` 请求的完整 URL
3. 应该是: `http://localhost:9100/worker/list` (通过 Core API 网关)
4. 或者: `http://localhost:9601/worker/list` (直接访问 Ops Center)

**Ops Center API 地址**:
- 直接访问: `http://127.0.0.1:9601`
- 通过网关: `http://127.0.0.1:9100` (Core API 转发到 Ops Center)

### 问题 3: CORS 跨域问题
**症状**: 浏览器控制台显示 CORS 错误

**检查方法**:
1. 浏览器控制台查看是否有 CORS 错误
2. 错误信息: "Access-Control-Allow-Origin"

**解决方案**:
检查 Ops Center 配置 `etc/ops.yaml`:
```yaml
CROSConf:
  Address: '*'  # 允许所有来源（开发环境）
```

## 解决方案

### 方案 1: 确认用户已登录

**步骤**:
1. 访问前端: http://localhost:3100
2. 使用正确的账号密码登录
3. 登录成功后，JWT Token 会自动保存到 localStorage
4. 刷新 Worker 管理页面

**验证**:
浏览器控制台执行:
```javascript
localStorage.getItem('ACCESS_TOKEN')
// 应该返回一个长字符串
```

### 方案 2: 检查 API 代理配置

**检查前端配置**:
```bash
# 查找 Vite 配置
cat /opt/code/newbee/ui/apps/web-antd/vite.config.mts | grep proxy -A 20
```

前端可能通过代理访问后端，需要确认代理配置正确。

### 方案 3: 检查中间件配置

**检查 Ops Center 中间件**:
```bash
grep -A 5 "worker/list" /opt/code/newbee/ops-center/api/etc/ops.yaml
```

确认 `/worker/list` 不在 `skipPaths` 中（需要认证）。

### 方案 4: 使用 SystemContext 查询（临时方案）

如果 Worker 是系统级资源，可能需要使用 SystemContext 绕过租户隔离。

**修改 `get_worker_list_logic.go`**:
```go
// 在查询前使用 SystemContext
import "github.com/coder-lulu/newbee-common/orm/ent/hooks"

func (l *GetWorkerListLogic) GetWorkerList(req *types.WorkerListReq) (resp *types.WorkerListResp, err error) {
    // 使用 SystemContext 绕过租户隔离
    ctx := hooks.NewSystemContext(l.ctx)

    // 查询所有 Worker（不过滤 tenant_id）
    workers, total, err := l.svcCtx.WorkerClient.ListWorkersWithoutTenant(
        ctx,
        req.Page,
        req.PageSize,
        filters,
    )
    // ...
}
```

## 调试命令

### 1. 查看 Ops Center 日志
```bash
# 如果使用 go run 启动
# 查看终端输出

# 如果使用 nohup 后台运行
tail -f ops-center.log
```

### 2. 测试 API（需要 Token）
```bash
# 从浏览器 localStorage 获取 token
TOKEN="<your_jwt_token_here>"

# 测试 API
curl -H "Authorization: Bearer $TOKEN" \
  'http://127.0.0.1:9601/worker/list?page=1&pageSize=10' | jq
```

### 3. 检查数据库连接
```bash
# 测试数据库
mysql -h192.168.26.130 -uroot -p123456 newbee -e "SELECT 1;"
```

### 4. 检查 Ops Center 进程
```bash
ps aux | grep "ops.go\|ops-center" | grep -v grep
netstat -tuln | grep 9601
```

## 推荐的排查顺序

1. **确认用户已登录** ⭐ (最可能)
   - 检查浏览器 localStorage 是否有 ACCESS_TOKEN
   - 检查 Network 标签中请求是否带 Authorization header

2. **检查 API 请求地址**
   - Network 标签查看完整 URL
   - 确认请求到了正确的端口

3. **查看浏览器控制台错误**
   - Console 标签查看 JavaScript 错误
   - Network 标签查看 HTTP 状态码

4. **检查后端日志**
   - 查看 Ops Center 日志输出
   - 确认请求是否到达后端

5. **数据库验证**
   - 确认数据库中有 Worker 记录
   - 确认 tenant_id 匹配

## 临时解决方案

如果急需查看数据，可以临时修改 API 跳过认证：

**修改 `etc/ops.yaml`**:
```yaml
Middleware:
  auth:
    skipPaths:
      - "/worker/list"  # 临时添加，仅用于调试
```

**⚠️ 注意**: 这是临时方案，调试完成后请移除！

## 需要提供的信息

如果以上方法都无法解决，请提供：

1. 浏览器控制台截图（Console + Network 标签）
2. Ops Center 日志输出
3. 前端启动日志
4. 用户登录状态确认

---

**更新时间**: 2025-12-19
**诊断工具**: `/opt/code/newbee/ops-center/docs/WORKER_LIST_DIAGNOSIS.md`
