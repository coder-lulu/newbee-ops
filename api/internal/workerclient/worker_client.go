package workerclient

import (
	"context"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent"
	"github.com/coder-lulu/newbee-ops-rpc/ent/proxy"
	"github.com/coder-lulu/newbee-ops-rpc/ent/proxymetrics"
	"github.com/zeromicro/go-zero/core/logx"
)

// WorkerClient Worker客户端（用于API服务访问Worker数据）
type WorkerClient struct {
	db     *ent.Client
	logger logx.Logger
}

// NewWorkerClient 创建Worker客户端
func NewWorkerClient(db *ent.Client) *WorkerClient {
	return &WorkerClient{
		db:     db,
		logger: logx.WithContext(context.Background()),
	}
}

// PSKValidation PSK验证
func (c *WorkerClient) ValidatePSK(psk string, configPSK string) bool {
	return psk != "" && psk == configPSK
}

// GetWorkerByWorkerID 根据Worker ID获取Worker
func (c *WorkerClient) GetWorkerByWorkerID(ctx context.Context, workerID string, tenantID uint64) (*ent.Proxy, error) {
	return c.db.Proxy.Query().
		Where(
			proxy.WorkerIDEQ(workerID),
			proxy.TenantIDEQ(tenantID),
		).
		First(ctx)
}

// GetWorkerByID 根据ID获取Worker
func (c *WorkerClient) GetWorkerByID(ctx context.Context, id uint64) (*ent.Proxy, error) {
	return c.db.Proxy.Get(ctx, id)
}

// ListWorkers 获取Worker列表
func (c *WorkerClient) ListWorkers(ctx context.Context, tenantID uint64, page, pageSize uint64, filters map[string]interface{}) ([]*ent.Proxy, uint64, error) {
	query := c.db.Proxy.Query().
		Where(proxy.TenantIDEQ(tenantID))

	// 应用过滤条件
	if status, ok := filters["status"].(uint32); ok {
		query = query.Where(proxy.StatusEQ(uint8(status)))
	}
	if workerStatus, ok := filters["worker_status"].(string); ok && workerStatus != "" {
		query = query.Where(proxy.WorkerStatusEQ(proxy.WorkerStatus(workerStatus)))
	}
	if region, ok := filters["region"].(string); ok && region != "" {
		query = query.Where(proxy.RegionEQ(region))
	}
	if zone, ok := filters["zone"].(string); ok && zone != "" {
		query = query.Where(proxy.ZoneEQ(zone))
	}
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where(proxy.NameContains(name))
	}

	// 计算总数
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	workers, err := query.
		Order(ent.Desc(proxy.FieldID)).
		Offset(int(offset)).
		Limit(int(pageSize)).
		All(ctx)

	return workers, uint64(total), err
}

// UpdateWorkerWeight 更新Worker权重
func (c *WorkerClient) UpdateWorkerWeight(ctx context.Context, id uint64, weight int) error {
	return c.db.Proxy.UpdateOneID(id).
		SetWeight(weight).
		Exec(ctx)
}

// ActivateWorker 激活Worker
func (c *WorkerClient) ActivateWorker(ctx context.Context, id uint64) error {
	return c.db.Proxy.UpdateOneID(id).
		SetStatus(1).
		SetWorkerStatus(proxy.WorkerStatusOnline).
		Exec(ctx)
}

// DeactivateWorker 停用Worker
func (c *WorkerClient) DeactivateWorker(ctx context.Context, id uint64) error {
	return c.db.Proxy.UpdateOneID(id).
		SetStatus(0). // 0=禁用，1=启用（符合StatusMixin约定）
		SetWorkerStatus(proxy.WorkerStatusOffline).
		Exec(ctx)
}

// DeleteWorkers 删除Workers
func (c *WorkerClient) DeleteWorkers(ctx context.Context, ids []uint64) error {
	_, err := c.db.Proxy.Delete().
		Where(proxy.IDIn(ids...)).
		Exec(ctx)
	return err
}

// GetWorkerMetrics 获取Worker指标
func (c *WorkerClient) GetWorkerMetrics(ctx context.Context, workerID string, tenantID uint64, startTime, endTime *time.Time, limit int) ([]*ent.ProxyMetrics, error) {
	query := c.db.ProxyMetrics.Query().
		Where(
			proxymetrics.WorkerIDEQ(workerID),
			proxymetrics.TenantIDEQ(tenantID),
		)

	if startTime != nil {
		query = query.Where(proxymetrics.TimestampGTE(*startTime))
	}
	if endTime != nil {
		query = query.Where(proxymetrics.TimestampLTE(*endTime))
	}

	return query.
		Order(ent.Desc(proxymetrics.FieldTimestamp)).
		Limit(limit).
		All(ctx)
}
