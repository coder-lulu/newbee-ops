package agent

import (
	"context"
	"fmt"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Agent列表
func NewGetAgentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentListLogic {
	return &GetAgentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAgentListLogic) GetAgentList(req *types.AgentListReq) (resp *types.AgentListResp, err error) {
	// 检查OpsClient是否已配置
	if l.svcCtx.OpsClient == nil {
		logx.Error("OpsClient is not configured")
		return nil, fmt.Errorf("ops rpc client is not configured")
	}

	// 构建RPC请求
	rpcReq := &ops.AgentListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// 可选参数
	if req.Status > 0 {
		rpcReq.Status = &req.Status
	}
	if req.AgentStatus != "" {
		rpcReq.AgentStatus = &req.AgentStatus
	}
	if req.Region != "" {
		rpcReq.Region = &req.Region
	}
	if req.Name != "" {
		rpcReq.Name = &req.Name
	}
	if req.Tags != "" {
		rpcReq.Tags = &req.Tags
	}

	// 调用RPC服务
	rpcResp, err := l.svcCtx.OpsClient.GetAgentList(l.ctx, rpcReq)
	if err != nil {
		logx.Errorw("Failed to call GetAgentList RPC",
			logx.Field("error", err),
			logx.Field("request", req))
		return nil, fmt.Errorf("failed to get agent list: %w", err)
	}

	// 转换响应
	resp = &types.AgentListResp{
		Total: rpcResp.Total,
		Data:  make([]types.AgentItem, 0, len(rpcResp.Data)),
	}

	for _, item := range rpcResp.Data {
		agentItem := types.AgentItem{
			Id:          getUint64Value(item.Id),
			Name:        getStringValue(item.Name),
			AgentId:     getStringValue(item.AgentId),
			Host:        getStringValue(item.Host),
			Port:        getInt64Value(item.Port),
			AgentStatus: getStringValue(item.AgentStatus),
			Region:      getStringValue(item.Region),
			Tags:        getStringValue(item.Tags),
			Version:     getStringValue(item.Version),
		}

		if item.LastHeartbeat != nil {
			agentItem.LastHeartbeat = *item.LastHeartbeat
		}
		if item.CpuUsage != nil {
			agentItem.CpuUsage = *item.CpuUsage
		}
		if item.MemoryUsage != nil {
			agentItem.MemoryUsage = *item.MemoryUsage
		}

		resp.Data = append(resp.Data, agentItem)
	}

	return resp, nil
}

// Helper functions to safely extract pointer values
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func getUint64Value(ptr *uint64) uint64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func getInt64Value(ptr *int64) int64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}
