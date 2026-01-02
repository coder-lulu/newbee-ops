-- ================================================
-- OPS-CENTER Proxy管理菜单初始化SQL
-- 生成时间: 2025-12-17
-- 服务名称: ops-center
-- ================================================

-- 说明：
-- 1. 菜单ID从323开始（当前最大ID是322）
-- 2. parent_id=300 表示挂载在"运维中心"模块下
-- 3. menu_type: 0=目录, 1=菜单页面, 2=按钮权限
-- 4. service_name: 'ops-center' 标识所属服务
-- 5. tenant_id=1 为默认租户
-- 6. 先删除已存在的记录，避免重复插入

-- ================================================
-- 0. 清理已存在的记录（如果需要重新插入）
-- ================================================
-- DELETE FROM sys_menus WHERE id BETWEEN 323 AND 337;
-- DELETE FROM role_menus WHERE menu_id BETWEEN 323 AND 337;

-- ================================================
-- 1. Dashboard概览 (ID: 323)
-- ================================================

INSERT INTO sys_menus (
    id, created_at, updated_at, sort, tenant_id, menu_level, menu_type,
    path, name, redirect, component, disabled, service_name,
    title, icon, hide_menu, parent_id
) VALUES (
    323, NOW(), NOW(), 5, 1, 2, 1,
    'dashboard', 'OPS Dashboard概览', '', 'ops-center/dashboard/index', 0, 'ops-center',
    'Dashboard概览', 'ant-design:dashboard-outlined', 0, 300
)
ON DUPLICATE KEY UPDATE
    updated_at = NOW(),
    title = 'Dashboard概览',
    component = 'ops-center/dashboard/index';

-- 1.1 Dashboard概览权限按钮
INSERT INTO sys_menus (id, created_at, updated_at, sort, tenant_id, menu_level, menu_type, path, name, title, icon, parent_id, service_name) VALUES
(324, NOW(), NOW(), 1, 1, 3, 2, '/dashboard/view', 'OPS 查看Dashboard', '查看Dashboard', '', 323, 'ops-center')
ON DUPLICATE KEY UPDATE
    updated_at = NOW(),
    title = '查看Dashboard';

-- ================================================
-- 2. Proxy管理 (ID: 325)
-- ================================================

INSERT INTO sys_menus (
    id, created_at, updated_at, sort, tenant_id, menu_level, menu_type,
    path, name, redirect, component, disabled, service_name,
    title, icon, hide_menu, parent_id
) VALUES (
    325, NOW(), NOW(), 6, 1, 2, 1,
    'proxies', 'OPS Proxy管理', '', 'ops-center/proxies/index', 0, 'ops-center',
    'Proxy管理', 'ant-design:cluster-outlined', 0, 300
)
ON DUPLICATE KEY UPDATE
    updated_at = NOW(),
    title = 'Proxy管理',
    component = 'ops-center/proxies/index';

-- 2.1 Proxy管理权限按钮（使用唯一的名称避免冲突）
INSERT INTO sys_menus (id, created_at, updated_at, sort, tenant_id, menu_level, menu_type, path, name, title, icon, parent_id, service_name) VALUES
(326, NOW(), NOW(), 1, 1, 3, 2, '/proxy/list', 'OPS Proxy列表', '查询列表', '', 325, 'ops-center'),
(327, NOW(), NOW(), 2, 1, 3, 2, '/proxy/:id', 'OPS Proxy详情', '查看详情', '', 325, 'ops-center'),
(328, NOW(), NOW(), 3, 1, 3, 2, '/proxy/weight', 'OPS 更新Proxy权重', '更新权重', '', 325, 'ops-center'),
(329, NOW(), NOW(), 4, 1, 3, 2, '/proxy/:id/activate', 'OPS 激活Proxy', '激活Proxy', '', 325, 'ops-center'),
(330, NOW(), NOW(), 5, 1, 3, 2, '/proxy/:id/deactivate', 'OPS 停用Proxy', '停用Proxy', '', 325, 'ops-center'),
(331, NOW(), NOW(), 6, 1, 3, 2, '/proxy/delete', 'OPS 删除Proxy', '删除Proxy', '', 325, 'ops-center'),
(332, NOW(), NOW(), 7, 1, 3, 2, '/proxy/pick', 'OPS 选择Proxy', '选择Proxy', '', 325, 'ops-center'),
(333, NOW(), NOW(), 8, 1, 3, 2, '/proxy/metrics', 'OPS Proxy指标查询', '查询指标', '', 325, 'ops-center'),
(334, NOW(), NOW(), 9, 1, 3, 2, '/proxy/metrics/stats', 'OPS Proxy统计查询', '查询统计', '', 325, 'ops-center')
ON DUPLICATE KEY UPDATE
    updated_at = NOW();

-- ================================================
-- 3. Proxy指标监控页面 (ID: 335)
-- ================================================

INSERT INTO sys_menus (
    id, created_at, updated_at, sort, tenant_id, menu_level, menu_type,
    path, name, redirect, component, disabled, service_name,
    title, icon, hide_menu, parent_id
) VALUES (
    335, NOW(), NOW(), 7, 1, 2, 1,
    'proxies/metrics', 'OPS Proxy指标监控', '', 'ops-center/proxies/metrics', 0, 'ops-center',
    'Proxy指标监控', 'ant-design:line-chart-outlined', 1, 300
)
ON DUPLICATE KEY UPDATE
    updated_at = NOW(),
    component = 'ops-center/proxies/metrics';

-- ================================================
-- 4. Proxy选择测试页面 (ID: 336)
-- ================================================

INSERT INTO sys_menus (
    id, created_at, updated_at, sort, tenant_id, menu_level, menu_type,
    path, name, redirect, component, disabled, service_name,
    title, icon, hide_menu, parent_id
) VALUES (
    336, NOW(), NOW(), 8, 1, 2, 1,
    'proxies/pick', 'OPS Proxy选择测试', '', 'ops-center/proxies/pick', 0, 'ops-center',
    'Proxy选择测试', 'ant-design:experiment-outlined', 1, 300
)
ON DUPLICATE KEY UPDATE
    updated_at = NOW(),
    component = 'ops-center/proxies/pick';

-- ================================================
-- 5. Proxy详情页面 (ID: 337) - 动态路由
-- ================================================

INSERT INTO sys_menus (
    id, created_at, updated_at, sort, tenant_id, menu_level, menu_type,
    path, name, redirect, component, disabled, service_name,
    title, icon, hide_menu, hide_breadcrumb, parent_id
) VALUES (
    337, NOW(), NOW(), 9, 1, 2, 1,
    'proxies/:id', 'OPS Proxy详情页', '', 'ops-center/proxies/detail', 0, 'ops-center',
    'Proxy详情', 'ant-design:info-circle-outlined', 1, 1, 300
)
ON DUPLICATE KEY UPDATE
    updated_at = NOW(),
    component = 'ops-center/proxies/detail';

-- ================================================
-- 6. 为管理员角色分配新增菜单权限
-- ================================================

-- 查询管理员角色ID（假设code='admin'）
SET @admin_role_id = (SELECT id FROM sys_roles WHERE code = 'admin' AND tenant_id = 1 LIMIT 1);

-- 为管理员角色分配Dashboard和Worker权限
INSERT IGNORE INTO role_menus (role_id, menu_id)
SELECT @admin_role_id, id
FROM sys_menus
WHERE id BETWEEN 323 AND 337
  AND service_name = 'ops-center';

-- ================================================
-- 菜单汇总统计
-- ================================================
-- 父目录: 运维中心 (ID: 300) - 已存在
--
-- 新增菜单页面 (menu_type=1): 5个
--   - Dashboard概览 (ID: 323, sort: 5) - 显示在菜单
--   - Proxy管理 (ID: 325, sort: 6) - 显示在菜单
--   - Proxy指标监控 (ID: 335, sort: 7) - 隐藏在菜单（通过Proxy管理进入）
--   - Proxy选择测试 (ID: 336, sort: 8) - 隐藏在菜单（通过Proxy管理进入）
--   - Proxy详情 (ID: 337, sort: 9) - 隐藏在菜单（动态路由）
--
-- 新增权限按钮 (menu_type=2): 10个
--   - Dashboard权限: 1个 (ID: 324)
--   - Proxy管理权限: 9个 (ID: 326-334)
--
-- 总计: 15条菜单记录 (5个页面 + 10个按钮权限)
-- ================================================

-- ================================================
-- 验证SQL
-- ================================================
-- 查看新增菜单结构
SELECT
    id,
    name,
    title,
    path,
    menu_type,
    sort,
    hide_menu,
    service_name
FROM sys_menus
WHERE id BETWEEN 323 AND 337
ORDER BY id;

-- 查看运维中心完整菜单树
SELECT
    id,
    CONCAT(REPEAT('  ', menu_level - 1), title) as menu_tree,
    path,
    menu_type,
    sort,
    hide_menu
FROM sys_menus
WHERE id = 300 OR parent_id = 300
ORDER BY CASE WHEN id = 300 THEN 0 ELSE sort END, id;
