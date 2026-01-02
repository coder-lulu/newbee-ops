package script

import (
	"context"

	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteScriptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除脚本
func NewDeleteScriptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteScriptLogic {
	return &DeleteScriptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteScriptLogic) DeleteScript(req *types.IDsReq) (resp *types.BaseResp, err error) {
	// 调用RPC服务
	rpcReq := &ops.IDsReq{
		Ids: req.Ids,
	}

	result, err := l.svcCtx.OpsClient.DeleteScript(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorw("Failed to delete script via RPC",
			logx.Field("error", err),
			logx.Field("ids", req.Ids))
		return nil, err
	}

	// 返回结果
	return &types.BaseResp{
		Code: 0,
		Msg:  result.Msg,
	}, nil
}
