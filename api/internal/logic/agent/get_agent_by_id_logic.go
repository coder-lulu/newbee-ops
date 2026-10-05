package agent

import (
	"context"
	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Agent详情
func NewGetAgentByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentByIdLogic {
	return &GetAgentByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAgentByIdLogic) GetAgentById(req *types.IDReq) (resp *types.AgentDetail, err error) {
	item, err := l.svcCtx.OpsClient.GetAgentById(l.ctx, &pb.IDReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.AgentDetail{Id: item.GetId(),
		CreatedAt:          item.GetCreatedAt(),
		UpdatedAt:          item.GetUpdatedAt(),
		Status:             item.GetStatus(),
		Name:               item.GetName(),
		AgentId:            item.GetAgentId(),
		Host:               item.GetHost(),
		Port:               item.GetPort(),
		Description:        item.GetDescription(),
		AgentStatus:        item.GetAgentStatus(),
		LastHeartbeat:      item.GetLastHeartbeat(),
		HeartbeatInterval:  item.GetHeartbeatInterval(),
		LastOnlineAt:       item.GetLastOnlineAt(),
		SupportedProviders: item.GetSupportedProviders(),
		Capabilities:       item.GetCapabilities(),
		MaxConcurrentTasks: item.GetMaxConcurrentTasks(),
		Tags:               item.GetTags(),
		Region:             item.GetRegion(),
		LocalIp:            item.GetLocalIp(),
		PublicIp:           item.GetPublicIp(),
		NetworkSegments:    item.GetNetworkSegments(),
		ActiveSessions:     item.GetActiveSessions(),
		CpuUsage:           item.GetCpuUsage(),
		MemoryUsage:        item.GetMemoryUsage(),
		TotalRequests:      item.GetTotalRequests(),
		SuccessfulRequests: item.GetSuccessfulRequests(),
		FailedRequests:     item.GetFailedRequests(),
		Version:            item.GetVersion()}, nil
}
