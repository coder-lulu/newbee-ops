# Agent管理功能实现方案

**创建时间**: 2025-12-16
**目标**: 实现ops-center的Agent管理功能，供io服务的发现配置使用
**架构**: ops-center提供Agent管理 → unified-io调用ops-center RPC → 前端调用unified-io API

---

## 🎯 实现目标

### 核心需求
1. ✅ **ops-center RPC服务**已有Agent和AgentGroup的proto定义
2. ❌ **ops-center API服务**需要添加HTTP接口
3. ❌ **unified-io服务**需要配置ops-center RPC客户端
4. ❌ **前端**需要移除mock数据，调用真实API

### 架构流程
```
前端 (AgentSelector.vue)
   ↓ HTTP
unified-io.api (端口9501)
   ↓ gRPC
unified-io.rpc (端口9500)
   ↓ gRPC
ops-center.rpc (端口9600)
   ↓
数据库 (agents表)
```

---

## 📋 任务清单

### Phase 1: ops-center HTTP API实现 (高优先级)

#### 任务1.1: 创建Agent API定义
**文件**: `/opt/code/newbee/ops-center/api/desc/agent.api`

```go
syntax = "api"

info (
	title: "Agent管理API"
	desc: "OPS-Center Agent管理接口"
	version: "v1.0"
)

type (
	// Agent基础信息（简化版，用于列表展示）
	AgentItem {
		Id           uint64  `json:"id"`
		Name         string  `json:"name"`
		AgentId      string  `json:"agent_id"`
		Host         string  `json:"host"`
		Port         int64   `json:"port"`
		AgentStatus  string  `json:"agent_status"`  // online, offline, error
		LastHeartbeat int64  `json:"last_heartbeat"`
		Region       string  `json:"region,optional"`
		Tags         string  `json:"tags,optional"`  // JSON string
		Version      string  `json:"version,optional"`
		CpuUsage     float64 `json:"cpu_usage,optional"`
		MemoryUsage  float64 `json:"memory_usage,optional"`
	}

	// Agent详细信息
	AgentDetail {
		Id                   uint64  `json:"id"`
		CreatedAt            int64   `json:"created_at"`
		UpdatedAt            int64   `json:"updated_at"`
		Status               uint32  `json:"status"`
		Name                 string  `json:"name"`
		AgentId              string  `json:"agent_id"`
		Host                 string  `json:"host"`
		Port                 int64   `json:"port"`
		Description          string  `json:"description,optional"`
		AgentStatus          string  `json:"agent_status"`
		LastHeartbeat        int64   `json:"last_heartbeat"`
		HeartbeatInterval    int64   `json:"heartbeat_interval"`
		LastOnlineAt         int64   `json:"last_online_at,optional"`
		SupportedProviders   string  `json:"supported_providers,optional"`
		Capabilities         string  `json:"capabilities,optional"`
		MaxConcurrentTasks   int64   `json:"max_concurrent_tasks"`
		Tags                 string  `json:"tags,optional"`
		Region               string  `json:"region,optional"`
		LocalIp              string  `json:"local_ip,optional"`
		PublicIp             string  `json:"public_ip,optional"`
		NetworkSegments      string  `json:"network_segments,optional"`
		ActiveSessions       int64   `json:"active_sessions"`
		CpuUsage             float64 `json:"cpu_usage"`
		MemoryUsage          float64 `json:"memory_usage"`
		TotalRequests        int64   `json:"total_requests"`
		SuccessfulRequests   int64   `json:"successful_requests"`
		FailedRequests       int64   `json:"failed_requests"`
		Version              string  `json:"version,optional"`
	}

	// Agent列表请求
	AgentListReq {
		Page       uint64 `form:"page,default=1"`
		PageSize   uint64 `form:"page_size,default=10"`
		Status     uint32 `form:"status,optional"`       // 启用状态: 1=正常 2=禁用
		AgentStatus string `form:"agent_status,optional"` // 运行状态: online, offline
		Region     string `form:"region,optional"`
		Name       string `form:"name,optional"`
		Tags       string `form:"tags,optional"`
	}

	// Agent列表响应
	AgentListResp {
		Total uint64      `json:"total"`
		Data  []AgentItem `json:"data"`
	}

	// 创建/更新Agent请求
	AgentInfo {
		Id                 uint64  `json:"id,optional"`
		Name               string  `json:"name"`
		AgentId            string  `json:"agent_id"`
		Host               string  `json:"host"`
		Port               int64   `json:"port"`
		Description        string  `json:"description,optional"`
		HeartbeatInterval  int64   `json:"heartbeat_interval,default=30"`
		MaxConcurrentTasks int64   `json:"max_concurrent_tasks,default=10"`
		Region             string  `json:"region,optional"`
		Tags               string  `json:"tags,optional"`
	}

	// Agent选择请求（供Discovery使用）
	AgentSelectionReq {
		SelectionMode  string `json:"selection_mode"`   // manual_agent, manual_group, auto_*
		AgentId        string `json:"agent_id,optional"` // manual_agent模式
		AgentGroupId   uint64 `json:"agent_group_id,optional"` // manual_group模式
		Provider       string `json:"provider,optional"`       // 需要的Provider能力
		Region         string `json:"region,optional"`         // 优先区域
		NetworkSegment string `json:"network_segment,optional"` // 需要的网段
	}

	// Agent选择响应
	AgentSelectionResp {
		AgentId string `json:"agent_id"`
		Host    string `json:"host"`
		Port    int64  `json:"port"`
		Msg     string `json:"msg,optional"`
	}

	// 基础响应
	BaseResp {
		Code uint32 `json:"code"`
		Msg  string `json:"msg"`
	}

	// ID响应
	BaseIDResp {
		Id  uint64 `json:"id"`
		Msg string `json:"msg"`
	}

	// ID请求
	IDReq {
		Id uint64 `path:"id"`
	}

	// IDs请求
	IDsReq {
		Ids []uint64 `json:"ids"`
	}
)

@server (
	group: agent
)
service ops-api {
	@doc "获取Agent列表"
	@handler getAgentList
	get /agent/list (AgentListReq) returns (AgentListResp)

	@doc "获取Agent详情"
	@handler getAgentById
	get /agent/:id (IDReq) returns (AgentDetail)

	@doc "创建Agent"
	@handler createAgent
	post /agent/create (AgentInfo) returns (BaseIDResp)

	@doc "更新Agent"
	@handler updateAgent
	post /agent/update (AgentInfo) returns (BaseResp)

	@doc "删除Agent"
	@handler deleteAgent
	post /agent/delete (IDsReq) returns (BaseResp)

	@doc "选择Agent（供Discovery使用）"
	@handler selectAgentForDiscovery
	post /agent/select (AgentSelectionReq) returns (AgentSelectionResp)
}
```

#### 任务1.2: 创建AgentGroup API定义
**文件**: `/opt/code/newbee/ops-center/api/desc/agentgroup.api`

```go
syntax = "api"

info (
	title: "Agent分组管理API"
	desc: "OPS-Center Agent分组管理接口"
	version: "v1.0"
)

type (
	// AgentGroup基础信息
	AgentGroupItem {
		Id                  uint64 `json:"id"`
		Name                string `json:"name"`
		Description         string `json:"description,optional"`
		SelectionStrategy   string `json:"selection_strategy"` // round_robin, least_load, random
		HealthCheckInterval int64  `json:"health_check_interval"`
		AutoFailover        bool   `json:"auto_failover"`
		MemberCount         int    `json:"member_count,optional"` // 成员数量
	}

	// AgentGroup详细信息
	AgentGroupDetail {
		Id                  uint64 `json:"id"`
		CreatedAt           int64  `json:"created_at"`
		UpdatedAt           int64  `json:"updated_at"`
		Status              uint32 `json:"status"`
		Name                string `json:"name"`
		Description         string `json:"description,optional"`
		SelectionStrategy   string `json:"selection_strategy"`
		HealthCheckInterval int64  `json:"health_check_interval"`
		AutoFailover        bool   `json:"auto_failover"`
		MaxRetryCount       int64  `json:"max_retry_count"`
	}

	// AgentGroup列表请求
	AgentGroupListReq {
		Page     uint64 `form:"page,default=1"`
		PageSize uint64 `form:"page_size,default=10"`
		Status   uint32 `form:"status,optional"`
		Name     string `form:"name,optional"`
	}

	// AgentGroup列表响应
	AgentGroupListResp {
		Total uint64           `json:"total"`
		Data  []AgentGroupItem `json:"data"`
	}

	// 创建/更新AgentGroup请求
	AgentGroupInfo {
		Id                  uint64 `json:"id,optional"`
		Name                string `json:"name"`
		Description         string `json:"description,optional"`
		SelectionStrategy   string `json:"selection_strategy,default=round_robin"`
		HealthCheckInterval int64  `json:"health_check_interval,default=60"`
		AutoFailover        bool   `json:"auto_failover,default=true"`
		MaxRetryCount       int64  `json:"max_retry_count,default=3"`
	}
)

@server (
	group: agentgroup
)
service ops-api {
	@doc "获取Agent分组列表"
	@handler getAgentGroupList
	get /agentgroup/list (AgentGroupListReq) returns (AgentGroupListResp)

	@doc "获取Agent分组详情"
	@handler getAgentGroupById
	get /agentgroup/:id (IDReq) returns (AgentGroupDetail)

	@doc "创建Agent分组"
	@handler createAgentGroup
	post /agentgroup/create (AgentGroupInfo) returns (BaseIDResp)

	@doc "更新Agent分组"
	@handler updateAgentGroup
	post /agentgroup/update (AgentGroupInfo) returns (BaseResp)

	@doc "删除Agent分组"
	@handler deleteAgentGroup
	post /agentgroup/delete (IDsReq) returns (BaseResp)
}
```

#### 任务1.3: 生成API代码
```bash
cd /opt/code/newbee/ops-center/api

# 生成Agent API
goctl api go -api desc/agent.api -dir . --style=go_zero

# 生成AgentGroup API
goctl api go -api desc/agentgroup.api -dir . --style=go_zero

# 编译验证
go build -v .
```

---

### Phase 2: unified-io集成ops-center RPC (高优先级)

#### 任务2.1: 添加OpsRpc配置
**文件**: `/opt/code/newbee/unified-io/rpc/etc/io.yaml`

```yaml
# 在CoreRpc配置后添加
# Ops-Center服务RPC连接配置
OpsRpc:
  Endpoints:
    - 127.0.0.1:9600
  Timeout: 30000  # 30秒超时，单位：毫秒
```

#### 任务2.2: 修改Config结构体
**文件**: `/opt/code/newbee/unified-io/rpc/internal/config/config.go`

```go
type Config struct {
	zrpc.RpcServerConf
	DatabaseConf    database.Config
	RedisConf       redis.Config
	CasbinConf      casbin.CasbinConf
	CoreRpc         CoreRpcConf   // 已存在
	OpsRpc          OpsRpcConf    // 新增
	TaskWorker      TaskWorkerConf
	InputAdapterConf input_adapter.Config
	Middleware      middleware.MiddlewareConf
}

// OpsRpc客户端配置（新增）
type OpsRpcConf struct {
	Endpoints []string
	Timeout   int64 `json:",optional"`
}
```

#### 任务2.3: 修改ServiceContext
**文件**: `/opt/code/newbee/unified-io/rpc/internal/svc/service_context.go`

```go
import (
	"github.com/coder-lulu/newbee-ops-center/rpc/opsclient"  // 新增
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config           config.Config
	DB               *ent.Client
	Redis            redis.UniversalClient
	CoreRpc          core.Core
	OpsRpc           ops.Ops  // 新增
	// ... 其他字段
}

func NewServiceContext(c config.Config) *ServiceContext {
	// ... 数据库、Redis初始化

	// 初始化Core RPC客户端
	coreRpcConn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: c.CoreRpc.Endpoints,
		Timeout:   c.CoreRpc.Timeout,
	})

	// 初始化Ops RPC客户端（新增）
	opsRpcConn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: c.OpsRpc.Endpoints,
		Timeout:   c.OpsRpc.Timeout,
	})

	return &ServiceContext{
		Config:  c,
		DB:      db,
		Redis:   rds,
		CoreRpc: core.NewCore(coreRpcConn),
		OpsRpc:  ops.NewOps(opsRpcConn),  // 新增
	}
}
```

#### 任务2.4: 在Discovery相关Logic中使用OpsRpc
**示例**: 修改发现池创建逻辑，验证Agent存在性

**文件**: `/opt/code/newbee/unified-io/rpc/internal/logic/discovery_pool/create_discovery_pool_logic.go`

```go
func (l *CreateDiscoveryPoolLogic) CreateDiscoveryPool(in *io.DiscoveryPoolInfo) (*io.BaseIDResp, error) {
	// 验证Agent存在且在线
	if in.AgentId != nil && *in.AgentId != "" {
		agentResp, err := l.svcCtx.OpsRpc.GetAgentList(l.ctx, &ops.AgentListReq{
			Page:     1,
			PageSize: 1,
			AgentId:  in.AgentId,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to verify agent: %w", err)
		}
		if agentResp.Total == 0 {
			return nil, fmt.Errorf("agent not found: %s", *in.AgentId)
		}
		agent := agentResp.Data[0]
		if agent.AgentStatus != nil && *agent.AgentStatus != "online" {
			return nil, fmt.Errorf("agent is not online: %s (status: %s)",
				*in.AgentId, *agent.AgentStatus)
		}
	}

	// ... 继续创建发现池逻辑
}
```

---

### Phase 3: 前端集成 (中优先级)

#### 任务3.1: 创建前端Agent API
**文件**: `/opt/code/newbee/ui/apps/web-antd/src/api/ops/agent.ts`

```typescript
import { requestClient } from '#/api/request';

export namespace AgentAPI {
  // Agent基础信息
  export interface AgentItem {
    id: number;
    name: string;
    agent_id: string;
    host: string;
    port: number;
    agent_status: 'online' | 'offline' | 'error';
    last_heartbeat: number;
    region?: string;
    tags?: string; // JSON string
    version?: string;
    cpu_usage?: number;
    memory_usage?: number;
  }

  // Agent列表请求
  export interface AgentListReq {
    page?: number;
    page_size?: number;
    status?: number;
    agent_status?: string;
    region?: string;
    name?: string;
    tags?: string;
  }

  // Agent列表响应
  export interface AgentListResp {
    total: number;
    data: AgentItem[];
  }

  // Agent选择请求
  export interface AgentSelectionReq {
    selection_mode: string;
    agent_id?: string;
    agent_group_id?: number;
    provider?: string;
    region?: string;
    network_segment?: string;
  }

  // Agent选择响应
  export interface AgentSelectionResp {
    agent_id: string;
    host: string;
    port: number;
    msg?: string;
  }

  /**
   * 获取Agent列表
   */
  export function listAgents(params: AgentListReq = {}) {
    return requestClient.get<AgentListResp>('/ops-api/agent/list', { params });
  }

  /**
   * 选择Agent（供Discovery使用）
   */
  export function selectAgent(data: AgentSelectionReq) {
    return requestClient.post<AgentSelectionResp>('/ops-api/agent/select', data);
  }
}
```

#### 任务3.2: 修改AgentSelector组件
**文件**: `/opt/code/newbee/ui/apps/web-antd/src/views/cmdb/ci_types/components/discovery/AgentSelector.vue`

```vue
<script lang="ts" setup>
import type { DiscoveryMethod } from './DiscoveryMethodSelector.vue';
import { AgentAPI } from '#/api/ops/agent';  // 新增导入
import { onMounted, ref } from 'vue';
import { CheckCircleOutlined, DesktopOutlined, InfoCircleOutlined } from '@ant-design/icons-vue';
import { Alert, Card, Col, Empty, message, Radio, Row, Spin, Tag } from 'ant-design-vue';

interface Props {
  selectedAgentId: string | null;
  method: DiscoveryMethod | null;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  agentSelected: [agentId: string, agentInfo: any];
}>();

// 代理列表
const loading = ref(false);
const agents = ref<AgentAPI.AgentItem[]>([]);  // 修改类型
const selectedAgent = ref<string | null>(props.selectedAgentId);

// 加载代理列表（移除Mock数据）
const loadAgents = async () => {
  try {
    loading.value = true;

    // ✅ 调用真实API
    const response = await AgentAPI.listAgents({
      agent_status: 'online',
      page: 1,
      page_size: 100
    });

    agents.value = response.data;

    if (agents.value.length === 0) {
      message.warning('暂无在线Agent，请先启动Agent服务');
    }
  } catch (error) {
    console.error('加载代理失败:', error);
    message.error('加载代理列表失败');
    agents.value = [];
  } finally {
    loading.value = false;
  }
};

// 选择代理
const handleSelectAgent = (agentId: string) => {
  selectedAgent.value = agentId;
  const agentInfo = agents.value.find((a) => a.agent_id === agentId);
  emit('agentSelected', agentId, agentInfo);
};

onMounted(() => {
  loadAgents();
});
</script>

<template>
  <div class="agent-selector">
    <div class="selector-header">
      <h3>选择执行代理</h3>
      <p>选择一个OPS-Center代理来执行发现任务</p>
    </div>

    <Alert
      message="代理说明"
      description="代理需要能够访问目标系统。请选择网络连通性好、负载较低的代理。"
      type="info"
      show-icon
      closable
      style="margin-bottom: 24px"
    />

    <Spin :spinning="loading">
      <div v-if="agents.length === 0 && !loading" class="empty-state">
        <Empty description="暂无可用的在线Agent">
          <template #image>
            <DesktopOutlined style="font-size: 64px; color: #d9d9d9" />
          </template>
        </Empty>
      </div>

      <div v-else class="agents-list">
        <Row :gutter="[16, 16]">
          <Col v-for="agent in agents" :key="agent.agent_id" :xs="24" :sm="12" :md="8">
            <Card
              :class="['agent-card', { 'agent-card-selected': selectedAgent === agent.agent_id }]"
              hoverable
              @click="handleSelectAgent(agent.agent_id)"
            >
              <div class="agent-header">
                <div class="agent-name">
                  <DesktopOutlined />
                  <span>{{ agent.name }}</span>
                </div>
                <Tag v-if="agent.agent_status === 'online'" color="success">
                  <CheckCircleOutlined />
                  在线
                </Tag>
                <Tag v-else color="error">离线</Tag>
              </div>

              <div class="agent-info">
                <div class="info-row">
                  <span class="label">地址:</span>
                  <span class="value">{{ agent.host }}:{{ agent.port }}</span>
                </div>
                <div v-if="agent.region" class="info-row">
                  <span class="label">区域:</span>
                  <span class="value">{{ agent.region }}</span>
                </div>
                <div v-if="agent.version" class="info-row">
                  <span class="label">版本:</span>
                  <span class="value">{{ agent.version }}</span>
                </div>
                <div v-if="agent.cpu_usage !== undefined" class="info-row">
                  <span class="label">CPU:</span>
                  <span class="value">{{ agent.cpu_usage.toFixed(1) }}%</span>
                </div>
                <div v-if="agent.memory_usage !== undefined" class="info-row">
                  <span class="label">内存:</span>
                  <span class="value">{{ agent.memory_usage.toFixed(1) }}%</span>
                </div>
              </div>

              <div v-if="agent.tags" class="agent-tags">
                <Tag
                  v-for="tag in JSON.parse(agent.tags || '[]')"
                  :key="tag"
                  size="small"
                >
                  {{ tag }}
                </Tag>
              </div>

              <Radio
                :checked="selectedAgent === agent.agent_id"
                class="agent-radio"
              >
                选择此代理
              </Radio>
            </Card>
          </Col>
        </Row>
      </div>
    </Spin>
  </div>
</template>

<style scoped lang="less">
.agent-selector {
  .selector-header {
    margin-bottom: 24px;

    h3 {
      font-size: 18px;
      font-weight: 600;
      margin-bottom: 8px;
    }

    p {
      color: #666;
      font-size: 14px;
    }
  }

  .empty-state {
    padding: 48px 0;
    text-align: center;
  }

  .agents-list {
    .agent-card {
      transition: all 0.3s;
      border: 2px solid transparent;

      &:hover {
        border-color: #1890ff;
      }

      &-selected {
        border-color: #1890ff;
        background-color: #e6f7ff;
      }

      .agent-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;

        .agent-name {
          display: flex;
          align-items: center;
          gap: 8px;
          font-weight: 600;
          font-size: 16px;
        }
      }

      .agent-info {
        margin-bottom: 12px;

        .info-row {
          display: flex;
          justify-content: space-between;
          padding: 4px 0;
          font-size: 14px;

          .label {
            color: #666;
          }

          .value {
            font-weight: 500;
          }
        }
      }

      .agent-tags {
        margin-bottom: 12px;
        display: flex;
        gap: 4px;
        flex-wrap: wrap;
      }

      .agent-radio {
        width: 100%;
        display: block;
      }
    }
  }
}
</style>
```

#### 任务3.3: 修改vite代理配置（确认端口）
**文件**: `/opt/code/newbee/ui/apps/web-antd/vite.config.mts`

确认ops-api代理配置正确：
```typescript
'/ops-api': {
  target: 'http://127.0.0.1:9410',  // ⚠️ 确认ops.yaml中实际端口
  changeOrigin: true,
  rewrite: (path) => path.replace(/^\/ops-api/, ''),
  ws: true,
},
```

---

### Phase 4: 数据库Schema验证 (低优先级)

#### 任务4.1: 检查agents表结构
```sql
-- 验证agents表是否存在且字段完整
SHOW CREATE TABLE agents;

-- 如果不存在，需要创建（通常通过ent migration）
```

#### 任务4.2: 初始化测试数据
```sql
-- 插入测试Agent数据
INSERT INTO agents (
  name, agent_id, host, port, agent_status,
  last_heartbeat, heartbeat_interval, region,
  version, status, tenant_id
) VALUES (
  'Test-Agent-01',
  'agent-test-001',
  '127.0.0.1',
  8080,
  'online',
  UNIX_TIMESTAMP(),
  30,
  'local',
  '1.0.0',
  1,
  1
);
```

---

## 🔧 实施步骤顺序

### 推荐实施顺序
```
1️⃣ Phase 1.1 - 创建agent.api (30分钟)
   ↓
2️⃣ Phase 1.2 - 创建agentgroup.api (20分钟)
   ↓
3️⃣ Phase 1.3 - 生成API代码并测试 (20分钟)
   ↓
4️⃣ Phase 2.1-2.3 - unified-io配置OpsRpc (30分钟)
   ↓
5️⃣ Phase 2.4 - 修改discovery_pool logic (20分钟)
   ↓
6️⃣ Phase 3.1 - 创建前端Agent API (15分钟)
   ↓
7️⃣ Phase 3.2 - 修改AgentSelector组件 (30分钟)
   ↓
8️⃣ Phase 4 - 数据库验证和测试数据 (20分钟)
   ↓
9️⃣ 集成测试 (30分钟)
```

**预计总时长**: 3.5小时

---

## 🧪 测试验证

### 测试1: ops-center API验证
```bash
# 测试获取Agent列表
curl -X GET "http://localhost:9410/agent/list?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 预期响应
{
  "total": 1,
  "data": [
    {
      "id": 1,
      "name": "Test-Agent-01",
      "agent_id": "agent-test-001",
      "host": "127.0.0.1",
      "port": 8080,
      "agent_status": "online",
      "last_heartbeat": 1734336000,
      "region": "local",
      "version": "1.0.0"
    }
  ]
}
```

### 测试2: unified-io调用ops-center RPC
```bash
# 在unified-io.rpc中添加临时测试代码
func TestOpsRpcConnection(ctx context.Context, svcCtx *svc.ServiceContext) {
	resp, err := svcCtx.OpsRpc.GetAgentList(ctx, &ops.AgentListReq{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Printf("❌ OpsRpc调用失败: %v", err)
		return
	}
	log.Printf("✅ OpsRpc调用成功: total=%d", resp.Total)
}
```

### 测试3: 前端完整流程
1. 启动所有服务（ops-center.rpc, ops-center.api, unified-io.rpc, unified-io.api, web-antd）
2. 登录系统
3. 进入CMDB → CI类型管理
4. 选择CI类型 → 新建发现配置
5. 在"选择执行代理"步骤，验证能否看到真实Agent列表
6. 选择Agent后继续流程，验证能否成功创建发现池

---

## ⚠️ 注意事项

### 关键依赖关系
1. **ops-center.rpc 必须先启动**，否则unified-io无法调用
2. **Agent数据必须存在**，否则前端显示空列表
3. **JWT认证必须正确配置**，否则API调用401

### 端口冲突检查
- ops-center.api: 确认是9410还是9402（参考vite.config.mts）
- 建议统一使用9410（ops.yaml中的配置）

### 错误处理
在所有RPC调用处添加超时和重试逻辑：
```go
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

resp, err := l.svcCtx.OpsRpc.GetAgentList(ctx, req)
if err != nil {
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, errors.New("ops-center服务超时")
	}
	return nil, fmt.Errorf("调用ops-center失败: %w", err)
}
```

---

## 📚 相关文档

1. **Agent Proto定义**: `/opt/code/newbee/ops-center/rpc/desc/agent.proto`
2. **AgentGroup Proto定义**: `/opt/code/newbee/ops-center/rpc/desc/agentgroup.proto`
3. **Discovery配置流程**: `/opt/code/newbee/unified-io/migrations/discovery_configuration_service_dependencies.md`
4. **前端组件**: `/opt/code/newbee/ui/apps/web-antd/src/views/cmdb/ci_types/components/discovery/AgentSelector.vue`

---

## 🎯 完成标准

### Phase 1完成标准
- [ ] agent.api和agentgroup.api文件创建完成
- [ ] API代码生成无错误
- [ ] ops-center.api服务编译通过
- [ ] 使用curl能访问 `/agent/list` 接口

### Phase 2完成标准
- [ ] unified-io配置文件添加OpsRpc配置
- [ ] ServiceContext成功初始化OpsRpc客户端
- [ ] discovery_pool创建时能调用OpsRpc验证Agent
- [ ] unified-io.rpc服务启动无报错

### Phase 3完成标准
- [ ] 前端agent.ts API文件创建完成
- [ ] AgentSelector.vue移除所有mock数据
- [ ] 前端能成功调用 `/ops-api/agent/list`
- [ ] 发现配置页面能显示真实Agent列表

### Phase 4完成标准
- [ ] agents表结构验证完成
- [ ] 至少有1条测试Agent数据
- [ ] Agent状态为online
- [ ] 心跳时间正常更新

---

**文档创建时间**: 2025-12-16
**预计完成时间**: 3.5小时
**优先级**: 高
**维护人**: NewBee Team
