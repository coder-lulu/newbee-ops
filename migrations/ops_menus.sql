-- ================================================
-- Ops-Center服务菜单初始化SQL
-- 生成时间: 2025-12-17
-- 服务名称: ops-center
-- ================================================

-- 说明：
-- 1. 菜单ID从300开始（避免与unified-io的262冲突）
-- 2. parent_id=100000 表示顶级目录
-- 3. menu_type: 0=目录, 1=菜单页面, 2=按钮权限
-- 4. menu_level: 1=顶级目录, 2=二级菜单页面, 3=按钮权限
-- 5. service_name: 'ops-center' 标识所属服务
-- 6. tenant_id=1 为默认租户

-- ================================================
-- 0. 顶级目录：运维中心 (Ops Center)
-- ================================================

-- 0.1 运维中心顶级目录 (ID: 300)
INSERT INTO sys_menus (
    id, created_at, updated_at, sort, tenant_id, menu_level, menu_type,
    path, name, redirect, component, disabled, service_name,
    title, icon, hide_menu, parent_id
) VALUES (
    300, NOW(), NOW(), 5, 1, 1, 0,
    '/ops', '运维中心', '', 'LAYOUT', 0, 'ops-center',
    '运维中心', 'ant-design:tool-outlined', 0, 100000
);

-- ================================================
-- 1. 会话管理 (Session Management)
-- ================================================

-- 1.1 会话管理主菜单 (ID: 301)
INSERT INTO sys_menus (
    id, created_at, updated_at, sort, tenant_id, menu_level, menu_type,
    path, name, redirect, component, disabled, service_name,
    title, icon, hide_menu, parent_id
) VALUES (
    301, NOW(), NOW(), 1, 1, 2, 1,
    'session', '会话管理', '', 'ops/session/index', 0, 'ops-center',
    '会话管理', 'ant-design:desktop-outlined', 0, 300
);

-- 1.2 会话管理权限按钮
INSERT INTO sys_menus (id, created_at, updated_at, sort, tenant_id, menu_level, menu_type, path, name, title, icon, parent_id, service_name) VALUES
(302, NOW(), NOW(), 1, 1, 3, 2, '/ops/session/create', '创建会话', '创建会话', '', 301, 'ops-center'),
(303, NOW(), NOW(), 2, 1, 3, 2, '/ops/session/close', '关闭会话', '关闭会话', '', 301, 'ops-center'),
(304, NOW(), NOW(), 3, 1, 3, 2, '/ops/session/:id', '查看会话详情', '查看详情', '', 301, 'ops-center'),
(305, NOW(), NOW(), 4, 1, 3, 2, '/ops/session/list', '查询会话列表', '查询列表', '', 301, 'ops-center'),
(306, NOW(), NOW(), 5, 1, 3, 2, '/ops/session/get', '会话查询', '会话查询', '', 301, 'ops-center');

-- ================================================
-- 2. 代理管理 (Proxy Management)
-- ================================================

-- 2.1 代理管理主菜单 (ID: 307)
INSERT INTO sys_menus (
    id, created_at, updated_at, sort, tenant_id, menu_level, menu_type,
    path, name, redirect, component, disabled, service_name,
    title, icon, hide_menu, parent_id
) VALUES (
    307, NOW(), NOW(), 2, 1, 2, 1,
    'proxy', '代理管理', '', 'ops/proxy/index', 0, 'ops-center',
    '代理管理', 'ant-design:cloud-server-outlined', 0, 300
);

-- 2.2 代理管理权限按钮
INSERT INTO sys_menus (id, created_at, updated_at, sort, tenant_id, menu_level, menu_type, path, name, title, icon, parent_id, service_name) VALUES
(308, NOW(), NOW(), 1, 1, 3, 2, '/ops/proxy/list', '代理列表', '代理列表', '', 307, 'ops-center'),
(309, NOW(), NOW(), 2, 1, 3, 2, '/ops/proxy/register', '代理注册', '代理注册', '', 307, 'ops-center'),
(310, NOW(), NOW(), 3, 1, 3, 2, '/ops/proxy/heartbeat', '代理心跳', '代理心跳', '', 307, 'ops-center'),
(311, NOW(), NOW(), 4, 1, 3, 2, '/ops/proxy/pick', '选择代理', '选择代理', '', 307, 'ops-center');

-- ================================================
-- 3. 访问配置管理 (Access Profile Management)
-- ================================================

-- 3.1 访问配置主菜单 (ID: 312)
INSERT INTO sys_menus (
    id, created_at, updated_at, sort, tenant_id, menu_level, menu_type,
    path, name, redirect, component, disabled, service_name,
    title, icon, hide_menu, parent_id
) VALUES (
    312, NOW(), NOW(), 3, 1, 2, 1,
    'profile', '访问配置', '', 'ops/profile/index', 0, 'ops-center',
    '访问配置', 'ant-design:setting-outlined', 0, 300
);

-- 3.2 访问配置权限按钮
INSERT INTO sys_menus (id, created_at, updated_at, sort, tenant_id, menu_level, menu_type, path, name, title, icon, parent_id, service_name) VALUES
(313, NOW(), NOW(), 1, 1, 3, 2, '/ops/profile/list', '查询访问配置', '查询列表', '', 312, 'ops-center'),
(314, NOW(), NOW(), 2, 1, 3, 2, '/ops/profile/create', '创建访问配置', '创建配置', '', 312, 'ops-center'),
(315, NOW(), NOW(), 3, 1, 3, 2, '/ops/profile/update', '更新访问配置', '更新配置', '', 312, 'ops-center'),
(316, NOW(), NOW(), 4, 1, 3, 2, '/ops/profile/delete', '删除访问配置', '删除配置', '', 312, 'ops-center'),
(317, NOW(), NOW(), 5, 1, 3, 2, '/ops/profile/get', '查看配置详情', '查看详情', '', 312, 'ops-center');

-- ================================================
-- 4. 任务编排 (Task Orchestration)
-- ================================================

-- 4.1 任务编排主菜单 (ID: 318)
INSERT INTO sys_menus (
    id, created_at, updated_at, sort, tenant_id, menu_level, menu_type,
    path, name, redirect, component, disabled, service_name,
    title, icon, hide_menu, parent_id
) VALUES (
    318, NOW(), NOW(), 4, 1, 2, 1,
    'task', '任务编排', '', 'ops/task/index', 0, 'ops-center',
    '任务编排', 'ant-design:schedule-outlined', 0, 300
);

-- 4.2 任务编排权限按钮
INSERT INTO sys_menus (id, created_at, updated_at, sort, tenant_id, menu_level, menu_type, path, name, title, icon, parent_id, service_name) VALUES
(319, NOW(), NOW(), 1, 1, 3, 2, '/ops/task/list', '查询任务列表', '查询列表', '', 318, 'ops-center'),
(320, NOW(), NOW(), 2, 1, 3, 2, '/ops/task/create', '创建任务', '创建任务', '', 318, 'ops-center'),
(321, NOW(), NOW(), 3, 1, 3, 2, '/ops/task/status/:taskId', '查看任务状态', '查看状态', '', 318, 'ops-center'),
(322, NOW(), NOW(), 4, 1, 3, 2, '/ops/task/result/:taskId', '查看任务结果', '查看结果', '', 318, 'ops-center');

-- ================================================
-- 菜单汇总统计
-- ================================================
-- 顶级目录: 运维中心 (ID: 300) - menu_level=1, menu_type=0
--
-- 新增菜单页面 (menu_type=1): 4个
--   - 会话管理 (ID: 301)
--   - 代理管理 (ID: 307)
--   - 访问配置 (ID: 312)
--   - 任务编排 (ID: 318)
--
-- 新增权限按钮 (menu_type=2): 18个
--   - 会话管理权限: 5个 (ID: 302-306)
--   - 代理管理权限: 4个 (ID: 308-311)
--   - 访问配置权限: 5个 (ID: 313-317)
--   - 任务编排权限: 4个 (ID: 318-322)
--
-- 总计: 23条菜单记录 (1个顶级目录 + 4个页面 + 18个按钮权限)
-- ================================================

-- ================================================
-- 验证查询
-- ================================================
-- 查看新插入的菜单：
-- SELECT id, parent_id, menu_level, menu_type, name, title, path, service_name
-- FROM sys_menus
-- WHERE service_name = 'ops-center'
-- ORDER BY id;

-- 查看菜单树结构：
-- SELECT
--     p.id as parent_id, p.name as parent_name,
--     c.id as menu_id, c.name as menu_name, c.title, c.menu_type
-- FROM sys_menus p
-- LEFT JOIN sys_menus c ON c.parent_id = p.id
-- WHERE p.service_name = 'ops-center'
-- ORDER BY p.id, c.sort;
