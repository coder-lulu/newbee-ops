-- ================================================
-- 修复 OPS 服务菜单组件路径
-- 生成时间: 2025-12-17
-- 说明: 将旧的 ops/ 路径更新为 ops-center/ 路径
-- ================================================

-- 1. 更新会话管理页面路径
UPDATE sys_menus
SET component = 'ops-center/sessions/index',
    updated_at = NOW()
WHERE id = 301
  AND component = 'ops/session/index';

-- 2. 更新代理管理页面路径
UPDATE sys_menus
SET component = 'ops-center/proxy/index',
    updated_at = NOW()
WHERE id = 307
  AND component = 'ops/proxy/index';

-- 3. 更新访问配置页面路径
UPDATE sys_menus
SET component = 'ops-center/profiles/index',
    updated_at = NOW()
WHERE id = 312
  AND component = 'ops/profile/index';

-- 4. 更新任务编排页面路径
UPDATE sys_menus
SET component = 'ops-center/tasks/index',
    updated_at = NOW()
WHERE id = 318
  AND component = 'ops/task/index';

-- ================================================
-- 验证更新结果
-- ================================================
SELECT
    id,
    name,
    component,
    title,
    service_name
FROM sys_menus
WHERE id IN (301, 307, 312, 318)
ORDER BY id;
