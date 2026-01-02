-- Phase 2: 任务执行协同 - 数据库迁移脚本
-- 为 ops_tasks 表添加 Worker 绑定和执行时间跟踪字段

-- 1. 添加 worker_id 字段（绑定的Worker ID）
ALTER TABLE ops_tasks ADD COLUMN IF NOT EXISTS worker_id VARCHAR(100) DEFAULT NULL COMMENT '绑定的Worker ID';

-- 2. 添加任务分配时间
ALTER TABLE ops_tasks ADD COLUMN IF NOT EXISTS dispatched_at DATETIME DEFAULT NULL COMMENT '任务分配时间';

-- 3. 添加任务完成时间
ALTER TABLE ops_tasks ADD COLUMN IF NOT EXISTS completed_at DATETIME DEFAULT NULL COMMENT '任务完成时间';

-- 4. 添加执行耗时（毫秒）
ALTER TABLE ops_tasks ADD COLUMN IF NOT EXISTS execution_time_ms BIGINT DEFAULT NULL COMMENT '执行耗时(毫秒)';

-- 5. 添加完整结果数据字段（JSON格式）
ALTER TABLE ops_tasks ADD COLUMN IF NOT EXISTS result_data TEXT DEFAULT NULL COMMENT '完整结果数据(JSON)';

-- 6. 创建索引以优化查询性能
-- 按 worker_id 查询任务
CREATE INDEX IF NOT EXISTS idx_ops_tasks_worker_id ON ops_tasks(worker_id);

-- 按状态和租户查询（已存在的索引可能需要调整）
CREATE INDEX IF NOT EXISTS idx_ops_tasks_status_tenant ON ops_tasks(status_str, tenant_id);

-- 按分配时间查询（用于监控和超时检测）
CREATE INDEX IF NOT EXISTS idx_ops_tasks_dispatched_at ON ops_tasks(dispatched_at);

-- 执行结果验证
SELECT
    COUNT(*) as total_tasks,
    COUNT(worker_id) as tasks_with_worker,
    COUNT(dispatched_at) as dispatched_tasks,
    COUNT(completed_at) as completed_tasks
FROM ops_tasks;

-- 迁移完成提示
SELECT '✅ Phase 2 数据库迁移完成！' as status;
SELECT '新增字段: worker_id, dispatched_at, completed_at, execution_time_ms, result_data' as fields;
SELECT '新增索引: idx_ops_tasks_worker_id, idx_ops_tasks_status_tenant, idx_ops_tasks_dispatched_at' as indexes;
