# Worker 列表加载问题 - 修复完成

## ✅ 已完成的修改

### 1. 前端 API 路径修改

所有 ops-center 相关的 API 调用已添加统一前缀 `/ops-center-api`：

**修改的文件**：
- ✅ `/src/api/ops-center/worker.ts` - Worker 管理 API
- ✅ `/src/api/ops-center/session.ts` - 会话管理 API
- ✅ `/src/api/ops-center/proxy.ts` - 代理管理 API
- ✅ `/src/api/ops-center/profile.ts` - 访问配置 API
- ✅ `/src/api/ops-center/task.ts` - 任务编排 API
- ✅ `/src/api/ops-center/credential.ts` - 凭证管理 API
- ✅ `/src/api/ops-center/audit.ts` - 审计日志 API

**修改示例**：
```typescript
// 修改前
requestClient.get('/worker/list', { params })

// 修改后
requestClient.get('/ops-center-api/worker/list', { params })
```

### 2. Vite 代理配置修改

**文件**: `/opt/code/newbee/ui/apps/web-antd/vite.config.mts`

**添加的配置**：
```typescript
'/ops-center-api': {
  changeOrigin: true,
  rewrite: (path) => path.replace(/^\/ops-center-api/, ''),
  target: 'http://127.0.0.1:9601',  // Ops Center API 服务
  ws: true,
},
```

**工作流程**：
```
前端请求: GET /ops-center-api/worker/list
    ↓
Vite 代理拦截并重写
    ↓
转发到: http://127.0.0.1:9601/worker/list
    ↓
Ops Center API 处理
    ↓
返回 Worker 列表数据
```

## 🧪 测试步骤

### 步骤 1: 重启前端开发服务器

```bash
cd /opt/code/newbee/ui/apps/web-antd

# 如果前端正在运行，先停止 (Ctrl+C)

# 重新启动
pnpm dev
```

**预期输出**：
```
VITE v5.x.x  ready in xxx ms

➜  Local:   http://localhost:3100/
➜  Network: use --host to expose
➜  press h + enter to show help
```

### 步骤 2: 清除浏览器缓存

1. 打开浏览器 (Chrome/Edge)
2. 按 F12 打开开发者工具
3. 右键点击刷新按钮
4. 选择"清空缓存并硬性重新加载"

或者：
- Windows/Linux: `Ctrl + Shift + Delete`
- Mac: `Cmd + Shift + Delete`

### 步骤 3: 登录系统

1. 访问: http://localhost:3100
2. 输入用户名和密码登录
3. 确保登录成功

**验证登录成功**：
- 浏览器控制台 (F12 -> Console)
- 执行: `localStorage.getItem('ACCESS_TOKEN')`
- 应该返回一个 JWT Token 字符串

### 步骤 4: 访问 Worker 管理页面

1. 点击左侧菜单："运维中心"
2. 点击子菜单："Worker管理"
3. 或直接访问: http://localhost:3100/ops-center/workers

### 步骤 5: 验证数据加载

**预期结果**：
- ✅ 页面显示 Worker 列表表格
- ✅ 显示 3 条 Worker 记录：
  - agent-001 (在线)
  - test-worker-001 (离线)
  - test-worker-integration (离线)
- ✅ 显示 CPU、内存使用率等指标
- ✅ 分页器显示"共 3 条"

### 步骤 6: 检查网络请求

**浏览器开发者工具** (F12 -> Network 标签):

1. 刷新页面
2. 查找 `worker/list` 请求
3. 点击该请求查看详情

**验证检查项**：

| 检查项 | 预期值 | 说明 |
|--------|--------|------|
| Request URL | `/ops-center-api/worker/list?page=1&pageSize=10` | 前端请求路径 |
| Status Code | `200 OK` | 请求成功 |
| Request Method | `GET` | GET 请求 |
| Request Headers | 包含 `Authorization: Bearer <token>` | JWT 认证 |
| Response | `{"code":0,"data":[...],"total":3}` | 返回数据结构 |

**Response 示例**：
```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "workerId": "agent-001",
      "name": "agent-001",
      "workerStatus": "online",
      "cpuUsage": 3.44,
      "memoryUsage": 65.31,
      "activeSessions": 0,
      "lastHeartbeat": 1734567896
    }
  ],
  "total": 3
}
```

### 步骤 7: 检查代理日志

**终端输出** (运行 `pnpm dev` 的终端):

应该能看到代理日志：
```
[vite] http proxy /ops-center-api/worker/list -> http://127.0.0.1:9601/worker/list
```

如果没有看到，尝试在浏览器中访问其他 Worker 功能触发更多请求。

## 🔍 故障排查

### 问题 1: 前端启动失败

**症状**：
```
Error: Cannot find module '@vben/...'
```

**解决**：
```bash
# 重新安装依赖
pnpm install

# 清除缓存
rm -rf node_modules/.vite
```

### 问题 2: 仍然返回 404

**症状**: Network 标签显示 404 Not Found

**检查**：
1. 确认前端已重启
2. 检查 Vite 配置是否正确保存
3. 清除浏览器缓存

**验证 Vite 配置**：
```bash
grep -A5 "ops-center-api" /opt/code/newbee/ui/apps/web-antd/vite.config.mts
```

应该显示：
```typescript
'/ops-center-api': {
  changeOrigin: true,
  rewrite: (path) => path.replace(/^\/ops-center-api/, ''),
  target: 'http://127.0.0.1:9601',
  ws: true,
},
```

### 问题 3: 返回 401 未授权

**症状**:
```json
{"code": 40001, "message": "认证Token缺失"}
```

**原因**: 未登录或 Token 过期

**解决**：
1. 重新登录系统
2. 确保登录成功后再访问 Worker 页面

**验证 Token**：
```javascript
// 浏览器控制台
localStorage.getItem('ACCESS_TOKEN')
```

### 问题 4: 返回空数据

**症状**:
```json
{"code": 0, "data": [], "total": 0}
```

**原因**: 数据库中没有 Worker 记录

**检查数据库**：
```bash
mysql -h192.168.26.130 -uroot -p123456 newbee -e \
  "SELECT worker_id, name, worker_status FROM ops_workers;"
```

**启动 Worker**：
```bash
cd /opt/code/newbee/worker
./agent -f etc/agent.yaml
```

### 问题 5: Ops Center 未运行

**症状**:
```
[vite] http proxy error: connect ECONNREFUSED 127.0.0.1:9601
```

**检查**：
```bash
# 检查端口
netstat -tuln | grep 9601

# 检查进程
ps aux | grep "ops.go" | grep -v grep
```

**启动 Ops Center**：
```bash
cd /opt/code/newbee/ops-center/api
go run ops.go -f etc/ops.yaml
```

### 问题 6: TypeScript 类型错误

**症状**: 前端编译时报类型错误

**检查**：
```bash
cd /opt/code/newbee/ui/apps/web-antd
npm run typecheck
```

**常见错误**：通常是其他模块的已有错误，不影响 Worker 功能。

## 📊 功能验证清单

测试所有 Worker 管理功能：

- [ ] Worker 列表加载
- [ ] 分页功能
- [ ] 筛选功能（按状态、区域等）
- [ ] Worker 详情页面
- [ ] Worker 指标监控页面
- [ ] Worker 选择测试页面
- [ ] 更新 Worker 权重
- [ ] 激活/停用 Worker
- [ ] 删除 Worker

### 快速验证脚本

```bash
#!/bin/bash
# verify-worker-ui.sh

echo "=== Worker UI 功能验证 ==="
echo ""

# 1. 检查前端进程
echo "1. 检查前端服务:"
if ps aux | grep -E "vite|pnpm dev" | grep web-antd | grep -v grep > /dev/null; then
    echo "✅ 前端服务运行中"
else
    echo "❌ 前端服务未运行"
    echo "   启动: cd /opt/code/newbee/ui/apps/web-antd && pnpm dev"
fi
echo ""

# 2. 检查 Ops Center
echo "2. 检查 Ops Center:"
if netstat -tuln 2>/dev/null | grep ":9601 " > /dev/null; then
    echo "✅ Ops Center 运行中 (9601)"
else
    echo "❌ Ops Center 未运行"
    echo "   启动: cd /opt/code/newbee/ops-center/api && go run ops.go -f etc/ops.yaml"
fi
echo ""

# 3. 检查 Worker 数据
echo "3. 检查 Worker 数据:"
WORKER_COUNT=$(mysql -h192.168.26.130 -uroot -p123456 newbee -se "SELECT COUNT(*) FROM ops_workers;" 2>/dev/null)
if [ "$WORKER_COUNT" -gt 0 ]; then
    echo "✅ 数据库中有 $WORKER_COUNT 条 Worker 记录"
else
    echo "❌ 数据库中没有 Worker 记录"
    echo "   启动 Worker: cd /opt/code/newbee/worker && ./agent -f etc/agent.yaml"
fi
echo ""

# 4. 检查 Vite 配置
echo "4. 检查 Vite 配置:"
if grep -q "ops-center-api" /opt/code/newbee/ui/apps/web-antd/vite.config.mts; then
    echo "✅ Vite 代理配置正确"
else
    echo "❌ Vite 代理配置缺失"
fi
echo ""

echo "=== 验证完成 ==="
echo ""
echo "如果所有检查都通过，访问: http://localhost:3100/ops-center/workers"
```

## 📝 修改总结

### 文件修改统计

| 文件 | 修改类型 | 修改行数 |
|------|---------|---------|
| `vite.config.mts` | 新增代理配置 | +6 行 |
| `worker.ts` | 路径前缀修改 | ~6 行 |
| `session.ts` | 路径前缀修改 | ~5 行 |
| `proxy.ts` | 路径前缀修改 | ~3 行 |
| `profile.ts` | 路径前缀修改 | ~4 行 |
| `task.ts` | 路径前缀修改 | ~3 行 |
| `credential.ts` | 路径前缀修改 | ~2 行 |
| `audit.ts` | 路径前缀修改 | ~2 行 |
| **总计** | | **~31 行** |

### 配置模式对比

**修改前** ❌:
```
前端: /worker/list
     ↓ (无代理)
     ❌ 404 Not Found
```

**修改后** ✅:
```
前端: /ops-center-api/worker/list
     ↓ (Vite 代理)
后端: http://127.0.0.1:9601/worker/list
     ↓
     ✅ 返回数据
```

### 符合现有模式

所有服务现在使用统一的代理模式：

| 服务 | 前端前缀 | 后端端口 | 状态 |
|------|---------|---------|------|
| Core | `/sys-api` | 9101 | ✅ |
| CMDB | `/cmdb-api` | 9201 | ✅ |
| IO | `/io-api` | 9501 | ✅ |
| FMS | `/fms-api` | 9102 | ✅ |
| MMS | `/mms-api` | 9104 | ✅ |
| IPAM | `/ipam-api` | 9302 | ✅ |
| **Ops Center** | **`/ops-center-api`** | **9601** | **✅ 已修复** |

## 🎉 预期成果

修复完成后，你应该能够：

1. ✅ 正常加载 Worker 列表页面
2. ✅ 看到 3 条 Worker 记录
3. ✅ 查看 Worker 详情
4. ✅ 监控 Worker 指标
5. ✅ 测试 Worker 选择算法
6. ✅ 管理 Worker 权重
7. ✅ 激活/停用/删除 Worker

---

**修复完成时间**: 2025-12-19
**修复方式**: 统一 API 前缀模式
**影响范围**: 所有 Ops Center 功能
**向后兼容**: 是（符合现有架构模式）
**状态**: ✅ 已完成，等待测试验证
