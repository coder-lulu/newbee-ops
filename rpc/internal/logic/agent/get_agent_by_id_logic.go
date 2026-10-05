package agent

import (
	"context"
	"encoding/json"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentByIdLogic {
	return &GetAgentByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAgentByIdLogic) GetAgentById(in *ops.IDReq) (*ops.AgentInfo, error) {
	result, err := l.svcCtx.DB.Agent.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	// Marshal JSON字段
	supportedProvidersJSON, _ := json.Marshal(result.SupportedProviders)
	capabilitiesJSON, _ := json.Marshal(result.Capabilities)
	tagsJSON, _ := json.Marshal(result.Tags)
	resourceLimitsJSON, _ := json.Marshal(result.ResourceLimits)
	networkSegmentsJSON, _ := json.Marshal(result.NetworkSegments)
	metadataJSON, _ := json.Marshal(result.Metadata)

	// 转换为string指针
	supportedProvidersStr := string(supportedProvidersJSON)
	capabilitiesStr := string(capabilitiesJSON)
	tagsStr := string(tagsJSON)
	resourceLimitsStr := string(resourceLimitsJSON)
	networkSegmentsStr := string(networkSegmentsJSON)
	metadataStr := string(metadataJSON)

	// 转换enum为string
	agentStatusStr := string(result.AgentStatus)

	return &ops.AgentInfo{
		Id:                 &result.ID,
		CreatedAt:          pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:          pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		Status:             pointy.GetPointer(uint32(result.Status)),
		Name:               &result.Name,
		AgentId:            &result.AgentID,
		Host:               &result.Host,
		Port:               pointy.GetPointer(int64(result.Port)),
		ApiKey:             &result.APIKey,
		Description:        &result.Description,
		AgentStatus:        &agentStatusStr,
		LastHeartbeat:      optionalUnixMilli(result.LastHeartbeat),
		HeartbeatInterval:  pointy.GetPointer(int64(result.HeartbeatInterval)),
		LastOnlineAt:       optionalUnixMilli(result.LastOnlineAt),
		SupportedProviders: &supportedProvidersStr,
		Capabilities:       &capabilitiesStr,
		MaxConcurrentTasks: pointy.GetPointer(int64(result.MaxConcurrentTasks)),
		Tags:               &tagsStr,
		Region:             &result.Region,
		ResourceLimits:     &resourceLimitsStr,
		LocalIp:            &result.LocalIP,
		PublicIp:           &result.PublicIP,
		NetworkSegments:    &networkSegmentsStr,
		ActiveSessions:     pointy.GetPointer(int64(result.ActiveSessions)),
		CpuUsage:           &result.CPUUsage,
		MemoryUsage:        &result.MemoryUsage,
		TotalRequests:      &result.TotalRequests,
		SuccessfulRequests: &result.SuccessfulRequests,
		FailedRequests:     &result.FailedRequests,
		Version:            &result.Version,
		Metadata:           &metadataStr,
		LastError:          &result.LastError,
	}, nil
}
