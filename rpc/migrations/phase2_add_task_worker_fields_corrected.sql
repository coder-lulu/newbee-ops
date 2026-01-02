-- Phase 2: 任务执行协同 - 数据库迁移脚本（正确版本）
-- 为 tasks 表添加 Worker 绑定和执行时间跟踪字段

-- 先检查表结构
SHOW COLUMNS FROM tasks;

-- 添加字段
ALTER TABLE tasks ADD COLUMN worker_id VARCHAR(100) DEFAULT NULL COMMENT '绑定的Worker ID';
ALTER TABLE tasks ADD COLUMN dispatched_at DATETIME DEFAULT NULL COMMENT '任务分配时间';
ALTER TABLE tasks ADD COLUMN completed_at DATETIME DEFAULT NULL COMMENT '任务完成时间';
ALTER TABLE tasks ADD COLUMN execution_time_ms BIGINT DEFAULT NULL COMMENT '执行耗时(毫秒)';
ALTER TABLE tasks ADD COLUMN result_data TEXT DEFAULT NULL COMMENT '完整结果数据(JSON)';

-- 创建索引
CREATE INDEX idx_tasks_worker_id ON tasks(worker_id);
CREATE INDEX idx_tasks_status_tenant ON tasks(status_str, tenant_id);
CREATE INDEX idx_tasks_dispatched_at ON tasks(dispatched_at);

-- 执行结果验证
SELECT
    COUNT(*) as total_tasks,
    COUNT(worker_id) as tasks_with_worker,
    COUNT(dispatched_at) as dispatched_tasks,
    COUNT(completed_at) as completed_tasks
FROM tasks;

-- 查看最终表结构
SHOW COLUMNS FROM tasks;

SELECT '✅ Phase 2 数据库迁移完成！' as status;
