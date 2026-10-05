package ops

import (
	"context"
	"fmt"
	"github.com/coder-lulu/newbee-ops-api/internal/logic/profile"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) GetProfile(req *types.AccessProfileQueryReq) (resp *types.AccessProfileResp, err error) {
	if req.CiId == "" {
		return nil, fmt.Errorf("ciId is required")
	}
	item, found, err := profile.NewLogic(l.ctx, l.svcCtx).Get(req.CiId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("access profile not found")
	}
	return &types.AccessProfileResp{Msg: "success", Data: item}, nil
}
