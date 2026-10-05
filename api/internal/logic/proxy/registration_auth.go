package proxy

import (
	"context"
	"crypto/subtle"
	"strconv"

	"github.com/coder-lulu/newbee-common/v2/middleware/keys"
	"github.com/coder-lulu/newbee-common/v2/orm/ent/hooks"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"google.golang.org/grpc/metadata"
)

// registrationContext binds an authenticated machine to its configured tenant.
// Request-supplied tenancy and system flags never authorize machine access.
func registrationContext(ctx context.Context, s *svc.ServiceContext, psk string) (context.Context, *types.BaseResp) {
	cfg := s.Config.Ops.Registration
	if cfg.PSK == "" || psk == "" || subtle.ConstantTimeCompare([]byte(psk), []byte(cfg.PSK)) != 1 {
		return nil, &types.BaseResp{Code: 401, Msg: "Invalid PSK"}
	}
	if cfg.TenantID == 0 {
		return nil, &types.BaseResp{Code: 503, Msg: "Registration tenant is not configured"}
	}
	if s.OpsClient == nil {
		return nil, &types.BaseResp{Code: 503, Msg: "Ops service is unavailable"}
	}
	tenant := strconv.FormatUint(cfg.TenantID, 10)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		md = md.Copy()
		md.Delete(keys.SystemContextKey.String())
		md.Set(keys.TenantIDKey.String(), tenant)
		md.Set(keys.OriginalTenantIDKey.String(), tenant)
		ctx = metadata.NewIncomingContext(ctx, md)
	}
	ctx = hooks.SetTenantIDToContext(ctx, cfg.TenantID)
	ctx = keys.NewContextManager().SetOriginalTenantID(ctx, tenant)
	return ctx, nil
}
