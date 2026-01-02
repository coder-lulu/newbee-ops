-- Phase 2: 任务执行协同 - 数据库迁移脚本（修复版）
-- 为 ops_tasks 表添加 Worker 绑定和执行时间跟踪字段

-- 先检查表结构
SHOW COLUMNS FROM ops_tasks;

-- 添加字段（如果字段已存在会报错，但不影响后续执行）
ALTER TABLE ops_tasks ADD COLUMN worker_id VARCHAR(100) DEFAULT NULL COMMENT '绑定的Worker ID';
ALTER TABLE ops_tasks ADD COLUMN dispatched_at DATETIME DEFAULT NULL COMMENT '任务分配时间';
ALTER TABLE ops_tasks ADD COLUMN completed_at DATETIME DEFAULT NULL COMMENT '任务完成时间';
ALTER TABLE ops_tasks ADD COLUMN execution_time_ms BIGINT DEFAULT NULL COMMENT '执行耗时(毫秒)';
ALTER TABLE ops_tasks ADD COLUMN result_data TEXT DEFAULT NULL COMMENT '完整结果数据(JSON)';

-- 创建索引（如果索引已存在会报错，但不影响后续执行）
CREATE INDEX idx_ops_tasks_worker_id ON ops_tasks(worker_id);
CREATE INDEX idx_ops_tasks_status_tenant ON ops_tasks(status_str, tenant_id);
CREATE INDEX idx_ops_tasks_dispatched_at ON ops_tasks(dispatched_at);

-- 执行结果验证
SELECT
    COUNT(*) as total_tasks,
    COUNT(worker_id) as tasks_with_worker,
    COUNT(dispatched_at) as dispatched_tasks,
    COUNT(completed_at) as completed_tasks
FROM ops_tasks;

-- 查看最终表结构
SHOW COLUMNS FROM ops_tasks;
