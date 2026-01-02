package agent

import (
	"context"
	"encoding/json"

	"github.com/coder-lulu/newbee-ops-rpc/ent/agent"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/msg/errormsg"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentLogic {
	return &CreateAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAgentLogic) CreateAgent(in *ops.AgentInfo) (*ops.BaseIDResp, error) {
	// 解析JSON字段
	var supportedProviders []string
	if in.SupportedProviders != nil && *in.SupportedProviders != "" {
		if err := json.Unmarshal([]byte(*in.SupportedProviders), &supportedProviders); err != nil {
			logx.Errorw("Failed to unmarshal supported_providers", logx.Field("error", err))
		}
	}

	var capabilities []string
	if in.Capabilities != nil && *in.Capabilities != "" {
		if err := json.Unmarshal([]byte(*in.Capabilities), &capabilities); err != nil {
			logx.Errorw("Failed to unmarshal capabilities", logx.Field("error", err))
		}
	}

	var tags []string
	if in.Tags != nil && *in.Tags != "" {
		if err := json.Unmarshal([]byte(*in.Tags), &tags); err != nil {
			logx.Errorw("Failed to unmarshal tags", logx.Field("error", err))
		}
	}

	var resourceLimits map[string]interface{}
	if in.ResourceLimits != nil && *in.ResourceLimits != "" {
		if err := json.Unmarshal([]byte(*in.ResourceLimits), &resourceLimits); err != nil {
			logx.Errorw("Failed to unmarshal resource_limits", logx.Field("error", err))
		}
	}

	var networkSegments []string
	if in.NetworkSegments != nil && *in.NetworkSegments != "" {
		if err := json.Unmarshal([]byte(*in.NetworkSegments), &networkSegments); err != nil {
			logx.Errorw("Failed to unmarshal network_segments", logx.Field("error", err))
		}
	}

	var metadata map[string]interface{}
	if in.Metadata != nil && *in.Metadata != "" {
		if err := json.Unmarshal([]byte(*in.Metadata), &metadata); err != nil {
			logx.Errorw("Failed to unmarshal metadata", logx.Field("error", err))
		}
	}

	// 创建Agent
	query := l.svcCtx.DB.Agent.Create().
		SetNotNilName(in.Name).
		SetNotNilAgentID(in.AgentId).
		SetNotNilHost(in.Host).
		SetNotNilAPIKey(in.ApiKey).
		SetNotNilDescription(in.Description).
		SetNotNilLastHeartbeat(pointy.GetTimeMilliPointer(in.LastHeartbeat)).
		SetNotNilLastOnlineAt(pointy.GetTimeMilliPointer(in.LastOnlineAt)).
		SetNotNilRegion(in.Region).
		SetNotNilLocalIP(in.LocalIp).
		SetNotNilPublicIP(in.PublicIp).
		SetNotNilCPUUsage(in.CpuUsage).
		SetNotNilMemoryUsage(in.MemoryUsage).
		SetNotNilTotalRequests(in.TotalRequests).
		SetNotNilSuccessfulRequests(in.SuccessfulRequests).
		SetNotNilFailedRequests(in.FailedRequests).
		SetNotNilVersion(in.Version).
		SetNotNilLastError(in.LastError)

	// 设置agent_status枚举
	if in.AgentStatus != nil {
		query.SetAgentStatus(agent.AgentStatus(*in.AgentStatus))
	}

	// 设置JSON字段
	if len(supportedProviders) > 0 {
		query.SetSupportedProviders(supportedProviders)
	}
	if len(capabilities) > 0 {
		query.SetCapabilities(capabilities)
	}
	if len(tags) > 0 {
		query.SetTags(tags)
	}
	if len(resourceLimits) > 0 {
		query.SetResourceLimits(resourceLimits)
	}
	if len(networkSegments) > 0 {
		query.SetNetworkSegments(networkSegments)
	}
	if len(metadata) > 0 {
		query.SetMetadata(metadata)
	}

	// 设置其他整型字段
	if in.Status != nil {
		query.SetNotNilStatus(pointy.GetPointer(uint8(*in.Status)))
	}
	if in.Port != nil {
		query.SetNotNilPort(pointy.GetPointer(int(*in.Port)))
	}
	if in.HeartbeatInterval != nil {
		query.SetNotNilHeartbeatInterval(pointy.GetPointer(int(*in.HeartbeatInterval)))
	}
	if in.MaxConcurrentTasks != nil {
		query.SetNotNilMaxConcurrentTasks(pointy.GetPointer(int(*in.MaxConcurrentTasks)))
	}
	if in.ActiveSessions != nil {
		query.SetNotNilActiveSessions(pointy.GetPointer(int(*in.ActiveSessions)))
	}

	result, err := query.Save(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &ops.BaseIDResp{Id: result.ID, Msg: errormsg.CreateSuccess}, nil
}
