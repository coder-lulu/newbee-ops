package agent

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent/agent"
	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentListLogic {
	return &GetAgentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAgentListLogic) GetAgentList(in *ops.AgentListReq) (*ops.AgentListResp, error) {
	var predicates []predicate.Agent
	if in.CreatedAt != nil {
		predicates = append(predicates, agent.CreatedAtGTE(time.UnixMilli(*in.CreatedAt)))
	}
	if in.UpdatedAt != nil {
		predicates = append(predicates, agent.UpdatedAtGTE(time.UnixMilli(*in.UpdatedAt)))
	}
	if in.Status != nil {
		predicates = append(predicates, agent.StatusEQ(uint8(*in.Status)))
	}
	if in.Name != nil {
		predicates = append(predicates, agent.NameContains(*in.Name))
	}
	if in.AgentId != nil {
		predicates = append(predicates, agent.AgentIDContains(*in.AgentId))
	}
	if in.Host != nil {
		predicates = append(predicates, agent.HostContains(*in.Host))
	}
	if in.Port != nil {
		predicates = append(predicates, agent.PortEQ(int(*in.Port)))
	}
	if in.ApiKey != nil {
		predicates = append(predicates, agent.APIKeyContains(*in.ApiKey))
	}
	if in.Description != nil {
		predicates = append(predicates, agent.DescriptionContains(*in.Description))
	}
	if in.AgentStatus != nil {
		predicates = append(predicates, agent.AgentStatusEQ(agent.AgentStatus(*in.AgentStatus)))
	}
	if in.LastHeartbeat != nil {
		predicates = append(predicates, agent.LastHeartbeatGTE(time.UnixMilli(*in.LastHeartbeat)))
	}
	if in.HeartbeatInterval != nil {
		predicates = append(predicates, agent.HeartbeatIntervalEQ(int(*in.HeartbeatInterval)))
	}
	if in.LastOnlineAt != nil {
		predicates = append(predicates, agent.LastOnlineAtGTE(time.UnixMilli(*in.LastOnlineAt)))
	}
	if in.MaxConcurrentTasks != nil {
		predicates = append(predicates, agent.MaxConcurrentTasksEQ(int(*in.MaxConcurrentTasks)))
	}
	if in.Region != nil {
		predicates = append(predicates, agent.RegionContains(*in.Region))
	}
	if in.LocalIp != nil {
		predicates = append(predicates, agent.LocalIPContains(*in.LocalIp))
	}
	if in.PublicIp != nil {
		predicates = append(predicates, agent.PublicIPContains(*in.PublicIp))
	}
	if in.ActiveSessions != nil {
		predicates = append(predicates, agent.ActiveSessionsEQ(int(*in.ActiveSessions)))
	}
	if in.CpuUsage != nil {
		predicates = append(predicates, agent.CPUUsageEQ(*in.CpuUsage))
	}
	if in.MemoryUsage != nil {
		predicates = append(predicates, agent.MemoryUsageEQ(*in.MemoryUsage))
	}
	if in.TotalRequests != nil {
		predicates = append(predicates, agent.TotalRequestsEQ(*in.TotalRequests))
	}
	if in.SuccessfulRequests != nil {
		predicates = append(predicates, agent.SuccessfulRequestsEQ(*in.SuccessfulRequests))
	}
	if in.FailedRequests != nil {
		predicates = append(predicates, agent.FailedRequestsEQ(*in.FailedRequests))
	}
	if in.Version != nil {
		predicates = append(predicates, agent.VersionContains(*in.Version))
	}
	if in.LastError != nil {
		predicates = append(predicates, agent.LastErrorContains(*in.LastError))
	}

	result, err := l.svcCtx.DB.Agent.Query().Where(predicates...).Page(l.ctx, in.Page, in.PageSize)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &ops.AgentListResp{}
	resp.Total = result.PageDetails.Total

	for _, v := range result.List {
		// Marshal JSON字段
		supportedProvidersJSON, _ := json.Marshal(v.SupportedProviders)
		capabilitiesJSON, _ := json.Marshal(v.Capabilities)
		tagsJSON, _ := json.Marshal(v.Tags)
		resourceLimitsJSON, _ := json.Marshal(v.ResourceLimits)
		networkSegmentsJSON, _ := json.Marshal(v.NetworkSegments)
		metadataJSON, _ := json.Marshal(v.Metadata)

		supportedProvidersStr := string(supportedProvidersJSON)
		capabilitiesStr := string(capabilitiesJSON)
		tagsStr := string(tagsJSON)
		resourceLimitsStr := string(resourceLimitsJSON)
		networkSegmentsStr := string(networkSegmentsJSON)
		metadataStr := string(metadataJSON)

		// 转换enum为string
		agentStatusStr := string(v.AgentStatus)

		resp.Data = append(resp.Data, &ops.AgentInfo{
			Id:                 &v.ID,
			CreatedAt:          pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:          pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			Status:             pointy.GetPointer(uint32(v.Status)),
			Name:               &v.Name,
			AgentId:            &v.AgentID,
			Host:               &v.Host,
			Port:               pointy.GetPointer(int64(v.Port)),
			ApiKey:             &v.APIKey,
			Description:        &v.Description,
			AgentStatus:        &agentStatusStr,
			LastHeartbeat:      pointy.GetUnixMilliPointer(v.LastHeartbeat.UnixMilli()),
			HeartbeatInterval:  pointy.GetPointer(int64(v.HeartbeatInterval)),
			LastOnlineAt:       pointy.GetUnixMilliPointer(v.LastOnlineAt.UnixMilli()),
			SupportedProviders: &supportedProvidersStr,
			Capabilities:       &capabilitiesStr,
			MaxConcurrentTasks: pointy.GetPointer(int64(v.MaxConcurrentTasks)),
			Tags:               &tagsStr,
			Region:             &v.Region,
			ResourceLimits:     &resourceLimitsStr,
			LocalIp:            &v.LocalIP,
			PublicIp:           &v.PublicIP,
			NetworkSegments:    &networkSegmentsStr,
			ActiveSessions:     pointy.GetPointer(int64(v.ActiveSessions)),
			CpuUsage:           &v.CPUUsage,
			MemoryUsage:        &v.MemoryUsage,
			TotalRequests:      &v.TotalRequests,
			SuccessfulRequests: &v.SuccessfulRequests,
			FailedRequests:     &v.FailedRequests,
			Version:            &v.Version,
			Metadata:           &metadataStr,
			LastError:          &v.LastError,
		})
	}

	return resp, nil
}
