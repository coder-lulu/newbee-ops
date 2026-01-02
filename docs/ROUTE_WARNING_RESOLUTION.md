# 前端路由警告说明

## 1. 已解决的警告（Ops 服务）

### 问题
路由生成脚本报告以下组件未找到：
- ❌ `/views/ops/session/index.vue`
- ❌ `/views/ops/proxy/index.vue`
- ❌ `/views/ops/profile/index.vue`
- ❌ `/views/ops/task/index.vue`

### 原因
数据库中的菜单配置路径使用旧的 `ops/` 前缀，但实际文件使用 `ops-center/` 前缀。

### 解决方案
已执行 SQL 更新脚本 `/opt/code/newbee/ops-center/migrations/fix_ops_menu_components.sql`

**更新内容**：
| 菜单ID | 旧路径 | 新路径 | 文件状态 |
|--------|--------|--------|----------|
| 301 | `ops/session/index` | `ops-center/sessions/index` | ✅ 已存在 |
| 307 | `ops/proxy/index` | `ops-center/proxy/index` | ✅ 已存在 |
| 312 | `ops/profile/index` | `ops-center/profiles/index` | ✅ 已存在 |
| 318 | `ops/task/index` | `ops-center/tasks/index` | ✅ 已存在 |

---

## 2. 待实现的功能（IO 服务）

### 问题
以下 unified-io 服务的前端页面尚未实现：

| 菜单ID | 组件路径 | 菜单名称 | 状态 |
|--------|---------|---------|------|
| 221 | `io/discovery-template/index` | 发现模板 | 📋 待实现 |
| 227 | `io/discovery-provider/index` | Provider管理 | 📋 待实现 |
| 231 | `io/data-target/index` | 数据目标 | 📋 待实现 |

### 说明
这些是 **unified-io** 服务的功能菜单，对应的后端 API 已实现，但前端页面尚未开发。

### 建议实现顺序

#### Phase 1: 数据目标管理
```
文件路径: /opt/code/newbee/ui/apps/web-antd/src/views/io/data-target/index.vue
功能: 管理数据目标（数据库、API端点等）
参考: /opt/code/newbee/unified-io/rpc/internal/logic/data_target/
```

#### Phase 2: 发现模板管理
```
文件路径: /opt/code/newbee/ui/apps/web-antd/src/views/io/discovery-template/index.vue
功能: 管理CI发现模板
参考: /opt/code/newbee/unified-io/rpc/internal/logic/discovery_template/
```

#### Phase 3: Provider管理
```
文件路径: /opt/code/newbee/ui/apps/web-antd/src/views/io/discovery-provider/index.vue
功能: 管理发现提供者（云厂商、网络设备等）
参考: /opt/code/newbee/unified-io/rpc/internal/logic/discovery_provider/
```

---

## 3. 验证步骤

### 验证 Ops 菜单修复
1. 清除浏览器缓存
2. 重新登录系统
3. 访问"运维中心"菜单
4. 确认以下页面正常加载：
   - 会话管理 (`/ops-center/sessions`)
   - 代理管理 (`/ops-center/proxy`)
   - 访问配置 (`/ops-center/profiles`)
   - 任务编排 (`/ops-center/tasks`)

### 临时屏蔽 IO 警告
如果暂时不实现 IO 服务页面，可以在数据库中将这些菜单设置为隐藏：

```sql
-- 临时隐藏未实现的菜单
UPDATE sys_menus
SET hide_menu = 1,
    updated_at = NOW()
WHERE id IN (221, 227, 231);
```

---

**文档更新时间**: 2025-12-17
**相关脚本**: `/opt/code/newbee/ops-center/migrations/fix_ops_menu_components.sql`
