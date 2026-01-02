-- ================================================
-- OPS-CENTER Database Migration: Worker → Proxy
-- 生成时间: 2025-12-28
-- 服务名称: ops-center
-- 目的: 将 worker 相关表重命名为 proxy
-- ================================================

-- 说明：
-- 1. 此脚本将现有的 worker 表重命名为 proxy 表
-- 2. 重命名相关字段（worker_id → proxy_id, worker_status → proxy_status）
-- 3. 更新索引名称以匹配新的表名和字段名
-- 4. 保留所有现有数据

-- ================================================
-- 方案选择
-- ================================================
-- 方案 A：重命名表和字段（保留数据）- 推荐用于生产环境
-- 方案 B：删除旧表并重建（丢失数据）- 仅用于开发环境

-- ================================================
-- 方案 A：重命名表和字段（保留数据）
-- ================================================

-- 1. 重命名主表
-- ================================================

-- 1.1 重命名 ops_workers 表
ALTER TABLE ops_workers RENAME TO ops_proxies;

-- 1.2 重命名字段
ALTER TABLE ops_proxies CHANGE COLUMN worker_id proxy_id VARCHAR(100);
ALTER TABLE ops_proxies CHANGE COLUMN worker_status proxy_status ENUM('online','degraded','offline');

-- 1.3 更新索引
-- 删除旧索引
ALTER TABLE ops_proxies DROP INDEX IF EXISTS ops_workers_tenant_id_worker_status_status;
ALTER TABLE ops_proxies DROP INDEX IF EXISTS ops_workers_tenant_id_region_worker_status;
ALTER TABLE ops_proxies DROP INDEX IF EXISTS ops_workers_worker_id;

-- 添加新索引
ALTER TABLE ops_proxies ADD INDEX idx_tenant_proxy_status (tenant_id, proxy_status, status);
ALTER TABLE ops_proxies ADD INDEX idx_tenant_region_proxy_status (tenant_id, region, proxy_status);
ALTER TABLE ops_proxies ADD INDEX idx_proxy_id (proxy_id);

-- 2. 重命名 ProxyGroup 表
-- ================================================

-- 2.1 重命名 ops_worker_groups 表
ALTER TABLE ops_worker_groups RENAME TO ops_proxy_groups;

-- 2.2 更新字段（如果有需要，根据schema检查）
-- ProxyGroup schema 中将 min_healthy_workers 改为了 min_healthy_proxies
ALTER TABLE ops_proxy_groups CHANGE COLUMN min_healthy_workers min_healthy_proxies INT DEFAULT 1;

-- 3. 重命名 ProxyGroupMember 表
-- ================================================

-- 3.1 重命名 ops_worker_group_members 表
ALTER TABLE ops_worker_group_members RENAME TO ops_proxy_group_members;

-- 3.2 重命名字段
ALTER TABLE ops_proxy_group_members CHANGE COLUMN worker_id proxy_id BIGINT UNSIGNED;
ALTER TABLE ops_proxy_group_members CHANGE COLUMN worker_group_id proxy_group_id BIGINT UNSIGNED;

-- 3.3 更新索引
-- 删除旧索引
ALTER TABLE ops_proxy_group_members DROP INDEX IF EXISTS ops_worker_group_members_worker_id_worker_group_id;
ALTER TABLE ops_proxy_group_members DROP INDEX IF EXISTS ops_worker_group_members_worker_group_id_priority_weight;
ALTER TABLE ops_proxy_group_members DROP INDEX IF EXISTS ops_worker_group_members_worker_id;

-- 添加新索引
ALTER TABLE ops_proxy_group_members ADD UNIQUE INDEX idx_proxy_id_proxy_group_id (proxy_id, proxy_group_id);
ALTER TABLE ops_proxy_group_members ADD INDEX idx_proxy_group_id_priority_weight (proxy_group_id, priority, weight);
ALTER TABLE ops_proxy_group_members ADD INDEX idx_proxy_id (proxy_id);

-- 4. 重命名 ProxyMetrics 表
-- ================================================

-- 4.1 重命名 ops_worker_metrics 表
ALTER TABLE ops_worker_metrics RENAME TO ops_proxy_metrics;

-- 4.2 重命名字段
ALTER TABLE ops_proxy_metrics CHANGE COLUMN worker_id proxy_id VARCHAR(100);
ALTER TABLE ops_proxy_metrics CHANGE COLUMN worker_status proxy_status VARCHAR(50);

-- 4.3 更新索引
-- 删除旧索引
ALTER TABLE ops_proxy_metrics DROP INDEX IF EXISTS ops_worker_metrics_worker_id_timestamp;
ALTER TABLE ops_proxy_metrics DROP INDEX IF EXISTS ops_worker_metrics_tenant_id_timestamp;

-- 添加新索引
ALTER TABLE ops_proxy_metrics ADD INDEX idx_proxy_id_timestamp (proxy_id, timestamp);
ALTER TABLE ops_proxy_metrics ADD INDEX idx_tenant_id_timestamp (tenant_id, timestamp);

-- ================================================
-- 验证迁移结果
-- ================================================

-- 查看表是否存在
SHOW TABLES LIKE 'ops_proxy%';

-- 查看 ops_proxies 表结构
DESCRIBE ops_proxies;

-- 查看 ops_proxy_groups 表结构
DESCRIBE ops_proxy_groups;

-- 查看 ops_proxy_group_members 表结构
DESCRIBE ops_proxy_group_members;

-- 查看 ops_proxy_metrics 表结构
DESCRIBE ops_proxy_metrics;

-- 检查数据是否完整
SELECT COUNT(*) AS proxy_count FROM ops_proxies;
SELECT COUNT(*) AS proxy_group_count FROM ops_proxy_groups;
SELECT COUNT(*) AS proxy_group_member_count FROM ops_proxy_group_members;
SELECT COUNT(*) AS proxy_metrics_count FROM ops_proxy_metrics;

-- ================================================
-- 回滚脚本（如果迁移失败）
-- ================================================
-- IMPORTANT: 仅在迁移失败时使用！
--
-- -- 重命名回 worker
-- ALTER TABLE ops_proxies RENAME TO ops_workers;
-- ALTER TABLE ops_workers CHANGE COLUMN proxy_id worker_id VARCHAR(100);
-- ALTER TABLE ops_workers CHANGE COLUMN proxy_status worker_status ENUM('online','degraded','offline');
--
-- ALTER TABLE ops_proxy_groups RENAME TO ops_worker_groups;
-- ALTER TABLE ops_worker_groups CHANGE COLUMN min_healthy_proxies min_healthy_workers INT DEFAULT 1;
--
-- ALTER TABLE ops_proxy_group_members RENAME TO ops_worker_group_members;
-- ALTER TABLE ops_worker_group_members CHANGE COLUMN proxy_id worker_id BIGINT UNSIGNED;
-- ALTER TABLE ops_worker_group_members CHANGE COLUMN proxy_group_id worker_group_id BIGINT UNSIGNED;
--
-- ALTER TABLE ops_proxy_metrics RENAME TO ops_worker_metrics;
-- ALTER TABLE ops_worker_metrics CHANGE COLUMN proxy_id worker_id VARCHAR(100);
-- ALTER TABLE ops_worker_metrics CHANGE COLUMN proxy_status worker_status VARCHAR(50);

-- ================================================
-- 迁移完成
-- ================================================
