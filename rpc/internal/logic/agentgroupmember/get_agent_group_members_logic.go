package agentgroupmember

import (
	"context"
	"encoding/json"

	"github.com/coder-lulu/newbee-ops-rpc/ent/agentgroupmember"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentGroupMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentGroupMembersLogic {
	return &GetAgentGroupMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAgentGroupMembers 获取Agent组成员列表
func (l *GetAgentGroupMembersLogic) GetAgentGroupMembers(in *ops.IDReq) (*ops.AgentGroupInfoWithMembers, error) {
	// 获取Agent Group基本信息
	group, err := l.svcCtx.DB.AgentGroup.Get(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 获取所有成员（包含Agent详情）
	members, err := l.svcCtx.DB.AgentGroupMember.Query().
		Where(
			agentgroupmember.AgentGroupIDEQ(group.ID),
		).
		WithAgent().
		All(l.ctx)

	if err != nil {
		return nil, err
	}

	// 转换enum为string
	selectionStrategyStr := string(group.SelectionStrategy)

	// 构造响应
	result := &ops.AgentGroupInfoWithMembers{
		Id:                  &group.ID,
		CreatedAt:           pointy.GetPointer(group.CreatedAt.UnixMilli()),
		UpdatedAt:           pointy.GetPointer(group.UpdatedAt.UnixMilli()),
		Status:              pointy.GetPointer(uint32(group.Status)),
		TenantId:            &group.TenantID,
		Name:                &group.Name,
		Description:         &group.Description,
		SelectionStrategy:   &selectionStrategyStr,
		HealthCheckInterval: pointy.GetPointer(int64(group.HealthCheckInterval)),
		AutoFailover:        &group.AutoFailover,
		MaxRetryCount:       pointy.GetPointer(int64(group.MaxRetryCount)),
		Members:             make([]*ops.AgentGroupMemberWithAgent, 0, len(members)),
	}

	// 转换成员信息
	for _, member := range members {
		memberInfo := &ops.AgentGroupMemberWithAgent{
			Id:           &member.ID,
			AgentId:      &member.AgentID,
			AgentGroupId: &member.AgentGroupID,
			Priority:     pointy.GetPointer(int64(member.Priority)),
			Weight:       pointy.GetPointer(int64(member.Weight)),
		}

		if !member.JoinedAt.IsZero() {
			memberInfo.JoinedAt = pointy.GetPointer(member.JoinedAt.UnixMilli())
		}

		// 添加Agent详情
		if member.Edges.Agent != nil {
			agent := member.Edges.Agent

			// Marshal JSON fields
			supportedProvidersJSON, _ := json.Marshal(agent.SupportedProviders)
			capabilitiesJSON, _ := json.Marshal(agent.Capabilities)
			tagsJSON, _ := json.Marshal(agent.Tags)
			resourceLimitsJSON, _ := json.Marshal(agent.ResourceLimits)
			networkSegmentsJSON, _ := json.Marshal(agent.NetworkSegments)
			metadataJSON, _ := json.Marshal(agent.Metadata)

			supportedProvidersStr := string(supportedProvidersJSON)
			capabilitiesStr := string(capabilitiesJSON)
			tagsStr := string(tagsJSON)
			resourceLimitsStr := string(resourceLimitsJSON)
			networkSegmentsStr := string(networkSegmentsJSON)
			metadataStr := string(metadataJSON)

			// Convert enum to string
			agentStatusStr := string(agent.AgentStatus)

			agentInfo := &ops.AgentInfo{
				Id:                 &agent.ID,
				CreatedAt:          pointy.GetPointer(agent.CreatedAt.UnixMilli()),
				UpdatedAt:          pointy.GetPointer(agent.UpdatedAt.UnixMilli()),
				Status:             pointy.GetPointer(uint32(agent.Status)),
				Name:               &agent.Name,
				AgentId:            &agent.AgentID,
				Host:               &agent.Host,
				Port:               pointy.GetPointer(int64(agent.Port)),
				ApiKey:             &agent.APIKey,
				Description:        &agent.Description,
				AgentStatus:        &agentStatusStr,
				LastHeartbeat:      pointy.GetUnixMilliPointer(agent.LastHeartbeat.UnixMilli()),
				HeartbeatInterval:  pointy.GetPointer(int64(agent.HeartbeatInterval)),
				LastOnlineAt:       pointy.GetUnixMilliPointer(agent.LastOnlineAt.UnixMilli()),
				SupportedProviders: &supportedProvidersStr,
				Capabilities:       &capabilitiesStr,
				MaxConcurrentTasks: pointy.GetPointer(int64(agent.MaxConcurrentTasks)),
				Tags:               &tagsStr,
				Region:             &agent.Region,
				ResourceLimits:     &resourceLimitsStr,
				LocalIp:            &agent.LocalIP,
				PublicIp:           &agent.PublicIP,
				NetworkSegments:    &networkSegmentsStr,
				ActiveSessions:     pointy.GetPointer(int64(agent.ActiveSessions)),
				CpuUsage:           &agent.CPUUsage,
				MemoryUsage:        &agent.MemoryUsage,
				TotalRequests:      &agent.TotalRequests,
				SuccessfulRequests: &agent.SuccessfulRequests,
				FailedRequests:     &agent.FailedRequests,
				Version:            &agent.Version,
				Metadata:           &metadataStr,
				LastError:          &agent.LastError,
			}

			memberInfo.Agent = agentInfo
		}

		result.Members = append(result.Members, memberInfo)
	}

	return result, nil
}
