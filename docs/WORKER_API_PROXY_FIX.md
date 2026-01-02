# Worker API 代理配置修复方案

## 问题分析

### 当前问题
Worker 管理页面无法加载数据，前端调用 `/worker/list` 但 Vite 配置中缺少对应的代理。

### 现有配置模式分析

查看 `vite.config.mts` 中的现有配置：

```typescript
'/io-api': {
  changeOrigin: true,
  rewrite: (path) => path.replace(/^\/io-api/, ''),
  target: 'http://127.0.0.1:9501',  // IO 服务
  ws: true,
},
'/cmdb-api': {
  changeOrigin: true,
  rewrite: (path) => path.replace(/^\/cmdb-api/, ''),
  target: 'http://127.0.0.1:9201',  // CMDB 服务
  ws: true,
},
'/sys-api': {
  changeOrigin: true,
  rewrite: (path) => path.replace(/^\/sys-api/, ''),
  target: 'http://127.0.0.1:9101',  // Core API 服务
  ws: true,
},
```

**配置模式**：
1. 前端API路径：`/{service}-api/{endpoint}`
2. Vite 代理：移除 `/{service}-api` 前缀
3. 后端接收：`/{endpoint}`

### 问题根源

1. **前端 API 路径不匹配**：
   - Worker API 使用：`/worker/list`
   - Session API 使用：`/ops/session/xxx`
   - 但配置中没有对应的代理

2. **ops-api 配置错误**：
   ```typescript
   '/ops-api': {
     target: 'http://127.0.0.1:9410',  // ❌ 此端口无服务运行！
   }
   ```
   - 9410 端口没有服务
   - ops-center 实际运行在 9600 (RPC) 和 9601 (API)

## 解决方案

### 方案 1：修改前端 API 调用路径（推荐）✅

修改前端代码，使用统一的 API 前缀模式。

#### 步骤 1: 修改 Worker API 路径

**文件**: `/opt/code/newbee/ui/apps/web-antd/src/api/ops-center/worker.ts`

```typescript
// 修改前
export async function listWorkers(params: WorkerListParams) {
  const ret = await requestClient.get<any>('/worker/list', { params });
  // ...
}

// 修改后
export async function listWorkers(params: WorkerListParams) {
  const ret = await requestClient.get<any>('/ops-center-api/worker/list', { params });
  // ...
}
```

#### 步骤 2: 修改所有 ops-center API 路径

**需要修改的文件**：
- `worker.ts` - 所有 Worker API
- `session.ts` - 所有 Session API
- `proxy.ts` - 所有 Proxy API
- `profile.ts` - 所有 Profile API
- `task.ts` - 所有 Task API
- `credential.ts` - 所有 Credential API

**统一替换规则**：
- `/worker/` → `/ops-center-api/worker/`
- `/ops/session/` → `/ops-center-api/ops/session/`
- `/ops/proxy/` → `/ops-center-api/ops/proxy/`
- `/ops/profile/` → `/ops-center-api/ops/profile/`
- `/ops/task/` → `/ops-center-api/ops/task/`

#### 步骤 3: 修改 Vite 代理配置

**文件**: `/opt/code/newbee/ui/apps/web-antd/vite.config.mts`

```typescript
'/ops-center-api': {
  changeOrigin: true,
  rewrite: (path) => path.replace(/^\/ops-center-api/, ''),
  target: 'http://127.0.0.1:9601',  // Ops Center API
  ws: true,
},
```

**替换原有的错误配置**：
```typescript
// 删除或修正这个配置
'/ops-api': {
  changeOrigin: true,
  rewrite: (path) => path.replace(/^\/ops-api/, ''),
  target: 'http://127.0.0.1:9410',  // ❌ 删除此配置
  ws: true,
},
```

### 方案 2：修改 Vite 代理配置（不推荐）⚠️

保持前端 API 路径不变，只修改 Vite 配置。

**问题**：需要为每个路径单独配置代理，不符合现有的统一模式。

```typescript
// 不推荐 - 过于分散
'/worker': {
  target: 'http://127.0.0.1:9601',
},
'/ops/session': {
  target: 'http://127.0.0.1:9601',
},
// ... 需要配置很多路径
```

## 修复脚本

### 自动修改 API 路径的脚本

```bash
#!/bin/bash
# fix-ops-center-api-paths.sh

API_DIR="/opt/code/newbee/ui/apps/web-antd/src/api/ops-center"

echo "修复 ops-center API 路径..."

# Worker API
sed -i "s|'/worker/|'/ops-center-api/worker/|g" "$API_DIR/worker.ts"
sed -i 's|`/worker/|`/ops-center-api/worker/|g' "$API_DIR/worker.ts"

# Session API
sed -i "s|'/ops/session/|'/ops-center-api/ops/session/|g" "$API_DIR/session.ts"

# Proxy API (如果存在)
if [ -f "$API_DIR/proxy.ts" ]; then
  sed -i "s|'/ops/proxy/|'/ops-center-api/ops/proxy/|g" "$API_DIR/proxy.ts"
fi

# Profile API
sed -i "s|'/ops/profile/|'/ops-center-api/ops/profile/|g" "$API_DIR/profile.ts"

# Task API
sed -i "s|'/ops/task/|'/ops-center-api/ops/task/|g" "$API_DIR/task.ts"

# Credential API
sed -i "s|'/ops/credential/|'/ops-center-api/ops/credential/|g" "$API_DIR/credential.ts"

echo "✅ API 路径修复完成"
echo "请检查修改后的文件，然后重启前端服务"
```

### 使用方法

```bash
chmod +x fix-ops-center-api-paths.sh
./fix-ops-center-api-paths.sh
```

## 验证步骤

### 1. 修改后重启前端

```bash
cd /opt/code/newbee/ui/apps/web-antd

# 停止前端 (Ctrl+C)
# 重新启动
pnpm dev
```

### 2. 检查代理日志

终端应显示：
```
[vite] http proxy /ops-center-api/worker/list -> http://127.0.0.1:9601/worker/list
```

### 3. 浏览器验证

1. 打开开发者工具 (F12)
2. Network 标签
3. 刷新 Worker 管理页面
4. 检查请求：
   - URL: `/ops-center-api/worker/list?page=1&pageSize=10`
   - Status: 200 OK
   - Response: 包含 Worker 数据

## 端口映射总结

| 服务 | RPC 端口 | API 端口 | Vite 代理前缀 | 配置状态 |
|------|---------|----------|--------------|---------|
| Core | 9100 | 9101 | `/sys-api` | ✅ 正确 |
| CMDB | 9200 | 9201 | `/cmdb-api` | ✅ 正确 |
| IO | 9500 | 9501 | `/io-api` | ✅ 正确 |
| FMS | - | 9102 | `/fms-api` | ✅ 正确 |
| MMS | - | 9104 | `/mms-api` | ✅ 正确 |
| IPAM | - | 9302 | `/ipam-api` | ✅ 正确 |
| Ops Center | 9600 | 9601 | `/ops-center-api` | ❌ 需要添加 |

## 常见问题

### Q1: 为什么不直接修改 Vite 配置？
**A**: 为了保持配置的一致性和可维护性，所有服务都应该使用 `/{service}-api/` 前缀模式。

### Q2: 9410 端口的服务是什么？
**A**: 9410 端口目前没有服务，可能是配置错误或历史遗留。应该删除或修正。

### Q3: 修改后是否影响其他功能？
**A**: 只要统一使用新的路径前缀，不会影响其他功能。确保前后端路径一致即可。

### Q4: 生产环境如何配置？
**A**: 生产环境通常使用 Nginx，需要配置相应的 location 规则：

```nginx
location /ops-center-api/ {
    rewrite ^/ops-center-api/(.*) /$1 break;
    proxy_pass http://127.0.0.1:9601;
}
```

---

**文档版本**: v1.0
**更新时间**: 2025-12-19
**状态**: 待执行方案 1
