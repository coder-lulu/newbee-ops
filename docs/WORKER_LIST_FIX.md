# Worker 列表加载问题修复

## 问题原因

**前端缺少 Vite 代理配置**，导致 Worker API 请求无法正确转发到 ops-center 服务。

### 详细分析

1. **前端 API 调用**: `/worker/list`
2. **后端服务地址**: `http://127.0.0.1:9601` (ops-center)
3. **问题**: Vite 配置中没有 `/worker` 路径的代理

## 修复内容

### 修改文件
`/opt/code/newbee/ui/apps/web-antd/vite.config.mts`

### 添加的代理配置

```typescript
'/worker': {
  changeOrigin: true,
  target: 'http://127.0.0.1:9601',
  ws: true,
},
'/ops/session': {
  changeOrigin: true,
  target: 'http://127.0.0.1:9601',
  ws: true,
},
'/ops/proxy': {
  changeOrigin: true,
  target: 'http://127.0.0.1:9601',
  ws: true,
},
'/ops/profile': {
  changeOrigin: true,
  target: 'http://127.0.0.1:9601',
  ws: true,
},
'/ops/task': {
  changeOrigin: true,
  target: 'http://127.0.0.1:9601',
  ws: true,
},
```

### 工作原理

```
前端请求: GET /worker/list
    ↓
Vite Dev Server (localhost:3100)
    ↓ (代理转发)
Ops Center API (localhost:9601)
    ↓
返回 Worker 列表数据
```

## 应用修复

### 步骤 1: 重启前端开发服务器

```bash
cd /opt/code/newbee/ui/apps/web-antd

# 如果前端正在运行，停止它 (Ctrl+C)
# 然后重新启动
pnpm dev
```

### 步骤 2: 清除浏览器缓存

1. 打开浏览器开发者工具 (F12)
2. 右键点击刷新按钮
3. 选择"清空缓存并硬性重新加载"

### 步骤 3: 验证修复

1. 访问: http://localhost:3100
2. 登录系统
3. 进入"运维中心" -> "Worker管理"
4. 应该能看到 Worker 列表数据

**期望结果**:
- 显示 3 条 Worker 记录
  - agent-001 (online)
  - test-worker-001 (offline)
  - test-worker-integration (offline)

## 验证步骤

### 1. 检查网络请求

浏览器开发者工具 -> Network 标签:

```
请求: /worker/list?page=1&pageSize=10
状态: 200 OK
响应:
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "workerId": "agent-001",
      "name": "agent-001",
      "workerStatus": "online",
      ...
    }
  ],
  "total": 3
}
```

### 2. 检查 Headers

Request Headers 应包含:
```
Authorization: Bearer <token>
Content-Type: application/json
```

### 3. 检查代理日志

终端应显示:
```
[vite] http proxy /worker/list -> http://127.0.0.1:9601/worker/list
```

## 其他相关 API 路径

同时修复了其他 ops-center 相关的 API 路径：

| 前端路径 | 后端服务 | 端口 | 用途 |
|---------|---------|------|------|
| `/worker/*` | ops-center | 9601 | Worker 管理 |
| `/ops/session/*` | ops-center | 9601 | 会话管理 |
| `/ops/proxy/*` | ops-center | 9601 | 代理管理 |
| `/ops/profile/*` | ops-center | 9601 | 访问配置 |
| `/ops/task/*` | ops-center | 9601 | 任务编排 |

## 故障排查

### 问题 1: 仍然返回 404
**原因**: Vite 服务器未重启
**解决**: 停止并重新启动 `pnpm dev`

### 问题 2: 仍然返回 401
**原因**: 用户未登录或 Token 过期
**解决**: 重新登录系统

### 问题 3: 返回空数据
**原因**: 数据库中没有 Worker 记录
**解决**:
```bash
# 检查数据库
mysql -h192.168.26.130 -uroot -p123456 newbee -e \
  "SELECT worker_id, name, worker_status FROM ops_workers;"

# 确保 Worker 进程正在运行
cd /opt/code/newbee/worker
./agent -f etc/agent.yaml
```

### 问题 4: Ops Center 未运行
**症状**: 代理请求超时或连接拒绝
**解决**:
```bash
# 检查 Ops Center 是否运行
ps aux | grep "ops.go" | grep -v grep
netstat -tuln | grep 9601

# 如果未运行，启动它
cd /opt/code/newbee/ops-center/api
go run ops.go -f etc/ops.yaml
```

## 生产环境配置

生产环境通常使用 Nginx 反向代理而不是 Vite Dev Server 代理。

### Nginx 配置示例

```nginx
# ops-center 相关路径
location ~ ^/(worker|ops)/ {
    proxy_pass http://127.0.0.1:9601;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;

    # WebSocket 支持
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
}
```

## 相关文档

- **问题诊断**: `WORKER_LIST_DIAGNOSIS.md`
- **Worker 启动指南**: `/opt/code/newbee/worker/docs/WORKER_STARTUP_GUIDE.md`
- **Ops Center 配置**: `/opt/code/newbee/ops-center/api/etc/ops.yaml`

---

**修复时间**: 2025-12-19
**修复文件**: `vite.config.mts`
**影响范围**: Worker 管理、会话管理、代理管理、访问配置、任务编排等 ops-center 所有功能
**状态**: ✅ 已修复，需要重启前端服务器生效
