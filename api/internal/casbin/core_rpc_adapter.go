package casbin

import (
	"context"

	dp "github.com/coder-lulu/newbee-common/v2/middleware/dataperm"
	coreclient "github.com/coder-lulu/newbee-core/rpc/coreclient"
	"github.com/coder-lulu/newbee-core/rpc/types/core"
)

// CoreRPCAdapter adapts core RPC client to dataperm.CoreRPCClient
// so DefaultCasbinProvider can call Core to check permissions and fetch roles.
type CoreRPCAdapter struct {
	coreRpc coreclient.Core
}

func NewCoreRPCAdapter(coreRpc coreclient.Core) *CoreRPCAdapter {
	return &CoreRPCAdapter{coreRpc: coreRpc}
}

// CheckPermissionWithRoles maps dataperm request to Core.CheckPermission RPC
// and converts Core response into dataperm.PermissionResult.
func (a *CoreRPCAdapter) CheckPermissionWithRoles(ctx context.Context, req *dp.PermissionCheckRequest) (*dp.PermissionResult, error) {
	enable := new(bool)
	*enable = true
	audit := new(bool)
	*audit = false
	rpcReq := &core.PermissionCheckReq{
		ServiceName: req.ServiceName,
		Subject:     req.Subject,
		Object:      req.Object,
		Action:      req.Action,
		EnableCache: enable,
		AuditLog:    audit,
	}
	resp, err := a.coreRpc.CheckPermission(ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	// Convert fields
	res := &dp.PermissionResult{
		Allowed:   resp.GetAllowed(),
		Reason:    resp.GetReason(),
		FromCache: resp.GetFromCache(),
	}
	if rules := resp.GetAppliedRules(); len(rules) > 0 {
		res.AppliedRules = append(res.AppliedRules, rules...)
	}
	return res, nil
}

// GetUserRolesWithCache fetches role codes for a user via Core RPC.
// Core does not expose a dedicated GetUserRolesWithCache RPC; we reuse
// GetUserById which returns role_codes in UserInfo.
func (a *CoreRPCAdapter) GetUserRolesWithCache(ctx context.Context, req *dp.UserRolesRequest) (*dp.UserRolesResponse, error) {
	// req.UserID 是用户UUID
	info, err := a.coreRpc.GetUserById(ctx, &core.UUIDReq{Id: req.UserID})
	if err != nil {
		return nil, err
	}
	return &dp.UserRolesResponse{Roles: append([]string{}, info.GetRoleCodes()...)}, nil
}
