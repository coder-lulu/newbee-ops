package casbin

import (
	"context"
	"fmt"

	commontypes "github.com/coder-lulu/newbee-common/v2/casbin/types"
	"github.com/coder-lulu/newbee-core/rpc/types/core"
	coreclient "github.com/coder-lulu/newbee-core/rpc/coreclient"
)

// RpcCasbinRuleQuerier Ops-Center API服务端的Casbin规则查询器
// 通过RPC调用从Core RPC服务查询规则
type RpcCasbinRuleQuerier struct {
	coreRpc coreclient.Core
}

// NewRpcCasbinRuleQuerier 创建RPC查询器
func NewRpcCasbinRuleQuerier(coreRpc coreclient.Core) *RpcCasbinRuleQuerier {
	return &RpcCasbinRuleQuerier{
		coreRpc: coreRpc,
	}
}

// QueryCasbinRules 实现CasbinRuleQuerier接口
// 通过RPC从Core服务查询Casbin规则
func (q *RpcCasbinRuleQuerier) QueryCasbinRules(ctx context.Context, tenantID uint64) ([]commontypes.CasbinRuleEntity, error) {
	// 调用RPC获取Casbin规则列表
	resp, err := q.coreRpc.GetCasbinRuleList(ctx, &core.CasbinRuleListReq{
		Page:     1,
		PageSize: 10000,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query casbin rules via RPC: %w", err)
	}

	// 转换RPC响应为通用接口
	result := make([]commontypes.CasbinRuleEntity, 0, len(resp.Data))
	for _, rule := range resp.Data {
		result = append(result, &RpcCasbinRuleWrapper{rule: rule})
	}

	return result, nil
}

// RpcCasbinRuleWrapper 包装RPC响应数据以实现CasbinRuleEntity接口
type RpcCasbinRuleWrapper struct {
	rule *core.CasbinRuleInfo
}

func (w *RpcCasbinRuleWrapper) GetID() uint64 {
	if w.rule.Id != nil {
		return *w.rule.Id
	}
	return 0
}

func (w *RpcCasbinRuleWrapper) GetPtype() string {
	return w.rule.Ptype
}

func (w *RpcCasbinRuleWrapper) GetV0() string {
	if w.rule.V0 != nil {
		return *w.rule.V0
	}
	return ""
}

func (w *RpcCasbinRuleWrapper) GetV1() string {
	if w.rule.V1 != nil {
		return *w.rule.V1
	}
	return ""
}

func (w *RpcCasbinRuleWrapper) GetV2() string {
	if w.rule.V2 != nil {
		return *w.rule.V2
	}
	return ""
}

func (w *RpcCasbinRuleWrapper) GetV3() string {
	if w.rule.V3 != nil {
		return *w.rule.V3
	}
	return ""
}

func (w *RpcCasbinRuleWrapper) GetV4() string {
	if w.rule.V4 != nil {
		return *w.rule.V4
	}
	return ""
}

func (w *RpcCasbinRuleWrapper) GetV5() string {
	if w.rule.V5 != nil {
		return *w.rule.V5
	}
	return ""
}

func (w *RpcCasbinRuleWrapper) GetTenantID() uint64 {
	if w.rule.TenantId != nil {
		return *w.rule.TenantId
	}
	return 0
}

func (w *RpcCasbinRuleWrapper) GetStatus() uint8 {
	if w.rule.Status != nil {
		return uint8(*w.rule.Status)
	}
	return 0
}

func (w *RpcCasbinRuleWrapper) GetRequireApproval() bool {
	if w.rule.RequireApproval != nil {
		return *w.rule.RequireApproval
	}
	return false
}

func (w *RpcCasbinRuleWrapper) GetApprovalStatus() string {
	if w.rule.ApprovalStatus != nil {
		return *w.rule.ApprovalStatus
	}
	return "pending"
}

func (w *RpcCasbinRuleWrapper) HasEffectiveFrom() bool {
	return w.rule.EffectiveFrom != nil && *w.rule.EffectiveFrom > 0
}

func (w *RpcCasbinRuleWrapper) HasEffectiveTo() bool {
	return w.rule.EffectiveTo != nil && *w.rule.EffectiveTo > 0
}
