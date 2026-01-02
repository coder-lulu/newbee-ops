# Worker管理系统 - 架构设计说明

## 📋 问题解答

### 问题1：页面底部显示"权重 (1-1000)"

**原因分析**：
- 这是编辑权重的模态框内容
- 模态框代码本身正确，但可能存在以下原因导致显示异常：
  1. 模态框未正确挂载到body
  2. z-index层级问题
  3. 浏览器缓存问题

**解决方案**：
已优化模态框代码，添加以下属性：
```vue
<a-modal
  v-model:open="weightModalVisible"
  title="编辑Worker权重"
  width="500px"
  :destroy-on-close="true"  <!-- ✅ 关闭时销毁内容 -->
  @ok="handleUpdateWeight"
>
  <a-form-item label="权重值 (范围: 1-1000)">  <!-- ✅ 更清晰的标签 -->
```

**验证方法**：
1. 清除浏览器缓存
2. 刷新页面
3. 点击任意Worker的"权重"按钮
4. 模态框应正常弹出在页面中央

---

### 问题2：Worker管理页面没有"新增"按钮

**结论**: ✅ **这是正确的设计，不需要新增Worker按钮**

**原因说明**：

#### 1️⃣ Worker是自动注册的分布式节点

Worker不是普通的数据记录，而是**分布式系统中的计算节点**。它的生命周期管理遵循以下架构：

```
┌─────────────────────────────────────────────────────┐
│  Worker Agent 程序 (独立进程)                        │
│  ├── 启动时自动注册到 ops-center                     │
│  ├── 定期发送心跳（30秒/次）                        │
│  ├── 上报系统资源指标（CPU/内存/磁盘/网络）          │
│  └── 异常时自动标记为离线                           │
└─────────────────────────────────────────────────────┘
                    ↓ HTTP/PSK认证
┌─────────────────────────────────────────────────────┐
│  OPS-Center API (注册中心)                          │
│  ├── /worker/register   (Worker注册接口)           │
│  ├── /worker/heartbeat  (心跳接口)                  │
│  └── PSK验证 (Pre-Shared Key认证)                  │
└─────────────────────────────────────────────────────┘
                    ↓ 存储
┌─────────────────────────────────────────────────────┐
│  数据库 (sys_workers表)                             │
│  ├── Worker元数据（IP、端口、版本）                 │
│  ├── 资源指标（CPU、内存、活跃会话）                │
│  └── 状态信息（在线、降级、离线）                   │
└─────────────────────────────────────────────────────┘
```

#### 2️⃣ PSK认证机制

Worker注册需要**预共享密钥（PSK）**认证，这是一种自动化的安全机制：

```go
// Worker Agent发送注册请求
POST /worker/register
Headers:
  X-OPS-PSK: dev-psk  // PSK密钥
Body:
{
  "worker_id": "worker-001",
  "name": "Beijing Worker",
  "ip": "192.168.1.100",
  "port": 8890,
  "capabilities": ["ssh", "telnet"],
  "psk": "dev-psk"
}
```

**安全特性**：
- PSK配置在Worker Agent的配置文件中
- 只有持有正确PSK的Agent才能注册
- 前端界面无法直接创建Worker（避免安全风险）

#### 3️⃣ Worker生命周期管理

| 阶段 | 触发方式 | 说明 |
|------|---------|------|
| **创建** | Agent启动时自动注册 | 无需手动创建 |
| **更新** | 心跳时自动更新指标 | 资源使用率实时更新 |
| **降级** | 心跳超时 | 超过阈值自动标记为degraded |
| **离线** | 长期无心跳 | 超过更长阈值标记为offline |
| **删除** | 管理员手动删除 | ✅ 前端可操作 |

#### 4️⃣ 前端管理功能

前端Worker管理页面提供的是**运维管理功能**，而非创建功能：

| 功能 | 说明 | 是否支持 |
|------|------|---------|
| 查看Worker列表 | 显示所有已注册Worker | ✅ 支持 |
| 查看Worker详情 | 查看Worker完整信息 | ✅ 支持 |
| 更新Worker权重 | 调整负载均衡权重 | ✅ 支持 |
| 激活/停用Worker | 控制Worker是否接收任务 | ✅ 支持 |
| 删除Worker | 从系统中移除Worker | ✅ 支持 |
| 查看指标监控 | 实时监控资源使用 | ✅ 支持 |
| 测试Worker选择 | 测试负载均衡算法 | ✅ 支持 |
| **手动创建Worker** | - | ❌ **不支持（设计如此）** |

---

## 🎯 如何部署Worker？

### 方案1：使用Worker Agent程序

**步骤**：
```bash
# 1. 配置Worker Agent
cd /opt/code/newbee/worker
vim etc/agent.yaml

# 配置示例：
OpsCenter:
  Endpoints:
    - http://ops-center-api:9601
  PSK: "your-secure-psk"  # 与ops-center配置一致
Worker:
  ID: "worker-beijing-01"
  Name: "北京机房Worker"
  IP: "192.168.1.100"
  Port: 8890

# 2. 启动Worker Agent
./worker -f etc/agent.yaml

# 3. 验证注册成功
# Worker会自动注册到ops-center，可在前端页面查看
```

### 方案2：使用Docker部署

```bash
docker run -d \
  --name worker-agent \
  -e OPS_CENTER_ENDPOINT=http://ops-center:9601 \
  -e OPS_CENTER_PSK=your-secure-psk \
  -e WORKER_ID=worker-docker-01 \
  newbee/worker-agent:latest
```

### 方案3：Kubernetes部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: worker-agent
spec:
  replicas: 3  # 自动创建3个Worker
  template:
    spec:
      containers:
      - name: worker
        image: newbee/worker-agent:latest
        env:
        - name: OPS_CENTER_ENDPOINT
          value: "http://ops-center-api:9601"
        - name: OPS_CENTER_PSK
          valueFrom:
            secretKeyRef:
              name: ops-center-secret
              key: psk
```

---

## 📊 前端页面优化

### 优化1：添加说明提示

在Worker列表页面顶部添加了Alert提示：

```vue
<a-alert
  message="Worker自动注册"
  description="Worker由独立的Agent程序通过PSK认证自动注册到系统，无需手动创建。您可以通过此页面管理已注册的Worker节点。"
  type="info"
  show-icon
  closable
  style="margin-bottom: 16px"
/>
```

**效果**：
- ✅ 用户一眼就能看到Worker是自动注册的
- ✅ 避免用户寻找"新增"按钮
- ✅ 可关闭，不会干扰正常使用

### 优化2：模态框改进

```vue
<a-modal
  v-model:open="weightModalVisible"
  title="编辑Worker权重"
  width="500px"
  :destroy-on-close="true"  <!-- 关闭时销毁，避免显示异常 -->
  @ok="handleUpdateWeight"
>
```

---

## 🔍 验证Worker注册

### 1. 检查Worker是否注册成功

**前端验证**：
- 访问 `/ops/workers`
- 查看Worker列表是否有数据
- 检查Worker状态（应为"在线"）

**后端验证**：
```sql
-- 查询所有已注册Worker
SELECT
    worker_id,
    name,
    ip,
    port,
    worker_status,
    last_heartbeat,
    cpu_usage,
    memory_usage
FROM sys_workers
WHERE status = 1  -- 启用状态
ORDER BY last_heartbeat DESC;
```

### 2. 检查Worker心跳

**日志验证**：
```bash
# Worker Agent日志
tail -f /var/log/worker/agent.log | grep "heartbeat"

# OPS-Center API日志
tail -f /var/log/ops-center/api.log | grep "WorkerHeartbeat"
```

**数据库验证**：
```sql
-- 检查最近心跳时间
SELECT
    worker_id,
    name,
    TIMESTAMPDIFF(SECOND, FROM_UNIXTIME(last_heartbeat), NOW()) as seconds_ago
FROM sys_workers
WHERE status = 1
ORDER BY last_heartbeat DESC;
```

---

## 🎯 总结

### 问题1解决方案
- ✅ 模态框代码已优化（添加`destroy-on-close`和`width`）
- ✅ 标签文本更清晰："权重值 (范围: 1-1000)"
- ✅ 添加placeholder提示

### 问题2设计解释
- ✅ Worker是分布式系统的计算节点，由Agent自动注册
- ✅ 采用PSK认证机制，保证安全性
- ✅ 前端提供运维管理功能，不提供手动创建
- ✅ 添加了Alert说明，避免用户困惑

### 架构优势
1. **自动化**: Worker自动注册和心跳，无需人工干预
2. **安全性**: PSK认证，防止未授权Worker注册
3. **可扩展**: 支持动态扩容，启动新Agent即可
4. **易部署**: 支持容器化和Kubernetes部署
5. **高可用**: 心跳机制自动检测Worker状态

---

**文档更新时间**: 2025-12-17
**文档版本**: v1.0.0
