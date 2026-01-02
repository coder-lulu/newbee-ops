# Ops-Center菜单导入报告

## 执行时间
2025-12-17

## 导入结果
✅ 成功插入 23 条菜单记录

## 菜单结构

### 📁 运维中心 (ID: 300)
顶级目录，path: /ops，icon: ant-design:tool-outlined

#### 📄 会话管理 (ID: 301)
二级菜单，path: session，icon: ant-design:desktop-outlined
- 🔘 创建会话 (ID: 302) - `/ops/session/create`
- 🔘 关闭会话 (ID: 303) - `/ops/session/close`
- 🔘 查看会话详情 (ID: 304) - `/ops/session/:id`
- 🔘 查询会话列表 (ID: 305) - `/ops/session/list`
- 🔘 会话查询 (ID: 306) - `/ops/session/get`

#### 📄 代理管理 (ID: 307)
二级菜单，path: proxy，icon: ant-design:cloud-server-outlined
- 🔘 代理列表 (ID: 308) - `/ops/proxy/list`
- 🔘 代理注册 (ID: 309) - `/ops/proxy/register`
- 🔘 代理心跳 (ID: 310) - `/ops/proxy/heartbeat`
- 🔘 选择代理 (ID: 311) - `/ops/proxy/pick`

#### 📄 访问配置 (ID: 312)
二级菜单，path: profile，icon: ant-design:setting-outlined
- 🔘 查询访问配置 (ID: 313) - `/ops/profile/list`
- 🔘 创建访问配置 (ID: 314) - `/ops/profile/create`
- 🔘 更新访问配置 (ID: 315) - `/ops/profile/update`
- 🔘 删除访问配置 (ID: 316) - `/ops/profile/delete`
- 🔘 查看配置详情 (ID: 317) - `/ops/profile/get`

#### 📄 任务编排 (ID: 318)
二级菜单，path: task，icon: ant-design:schedule-outlined
- 🔘 查询任务列表 (ID: 319) - `/ops/task/list`
- 🔘 创建任务 (ID: 320) - `/ops/task/create`
- 🔘 查看任务状态 (ID: 321) - `/ops/task/status/:taskId`
- 🔘 查看任务结果 (ID: 322) - `/ops/task/result/:taskId`

## 统计信息

| 类型 | 数量 | 说明 |
|------|------|------|
| 顶级目录 (menu_type=0) | 1 | 运维中心 |
| 二级菜单 (menu_type=1) | 4 | 会话管理、代理管理、访问配置、任务编排 |
| 按钮权限 (menu_type=2) | 18 | CRUD操作权限 |
| **总计** | **23** | |

## 菜单ID范围
- 起始ID: 300
- 结束ID: 322
- 预留ID: 323-399 (可用于后续扩展)

## 配置信息
- 服务名称: `ops-center`
- 租户ID: 1 (默认租户)
- parent_id: 100000 (顶级目录的父ID)

## API路径映射

### 会话管理 (Session)
| 按钮权限 | API路径 | 说明 |
|---------|---------|------|
| 创建会话 | POST /ops/session/create | CreateSession |
| 关闭会话 | POST /ops/session/close | CloseSession |
| 查看会话详情 | GET /ops/session/:id | GetSession |
| 查询会话列表 | GET /ops/session/list | ListSession |
| 会话查询 | GET /ops/session/get | GetSessionByQuery |

### 代理管理 (Proxy)
| 按钮权限 | API路径 | 说明 |
|---------|---------|------|
| 代理注册 | POST /ops/proxy/register | RegisterProxy (skipAuth) |
| 代理心跳 | POST /ops/proxy/heartbeat | HeartbeatProxy (skipAuth) |
| 选择代理 | GET /ops/proxy/pick | PickProxy |

### 访问配置 (Profile)
| 按钮权限 | API路径 | 说明 |
|---------|---------|------|
| 查询访问配置 | GET /ops/profile/list | ListProfile |
| 创建访问配置 | POST /ops/profile/create | CreateProfile |
| 更新访问配置 | POST /ops/profile/update | UpdateProfile |
| 删除访问配置 | DELETE /ops/profile/delete | DeleteProfile |
| 查看配置详情 | GET /ops/profile/get | GetProfile |

### 任务编排 (Task)
| 按钮权限 | API路径 | 说明 |
|---------|---------|------|
| 创建任务 | POST /ops/task/create | CreateTask |
| 查看任务状态 | GET /ops/task/status/:taskId | GetTaskStatus |
| 查看任务结果 | GET /ops/task/result/:taskId | GetTaskResult |

## 数据库验证查询

```sql
-- 查询所有ops-center菜单
SELECT id, parent_id, menu_level, menu_type, name, title, path, service_name
FROM sys_menus
WHERE service_name = 'ops-center'
ORDER BY id;

-- 统计菜单数量
SELECT
    menu_type,
    CASE menu_type
        WHEN 0 THEN '目录'
        WHEN 1 THEN '菜单'
        WHEN 2 THEN '按钮'
    END as type_name,
    COUNT(*) as count
FROM sys_menus
WHERE service_name = 'ops-center'
GROUP BY menu_type;
```

## 文件位置
- SQL脚本: `/opt/code/newbee/ops-center/migrations/ops_menus.sql`
- 导入报告: `/opt/code/newbee/ops-center/migrations/ops_menus_import_report.md`

## 下一步操作建议

1. **角色授权**: 需要将这些菜单权限分配给相应的角色
   ```sql
   -- 示例：给管理员角色分配所有ops-center菜单
   INSERT INTO sys_role_menus (role_id, menu_id)
   SELECT 1, id FROM sys_menus WHERE service_name = 'ops-center';
   ```

2. **前端路由配置**: 确保前端项目中存在对应的Vue组件
   - `ops/session/index.vue`
   - `ops/proxy/index.vue`
   - `ops/profile/index.vue`
   - `ops/task/index.vue`

3. **权限验证**: 在ops-center API中启用权限中间件验证

4. **国际化**: 为菜单标题添加多语言支持

## 注意事项

1. ⚠️ 代理注册和心跳接口需要配置为 `skipPaths` (已在配置中设置)
2. ⚠️ 所有菜单权限路径与API定义完全对应
3. ⚠️ 使用了统一的图标库 (ant-design icons)
4. ⚠️ 保持与unified-io菜单结构一致的命名规范
