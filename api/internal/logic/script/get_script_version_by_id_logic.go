package script

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScriptVersionByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetScriptVersionByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptVersionByIdLogic {
	return &GetScriptVersionByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetScriptVersionByIdLogic) GetScriptVersionById(req *types.IDReq) (resp *types.ScriptVersionDetailResp, err error) {
	// 调用RPC获取脚本版本详情
	result, err := l.svcCtx.OpsClient.GetScriptVersionById(l.ctx, req.Id)
	if err != nil {
		l.Errorw("Failed to get script version by id from RPC",
			logx.Field("error", err),
			logx.Field("id", req.Id))
		return nil, err
	}

	// 构造响应
	item := types.ScriptVersionItem{}
	if result.Id != nil {
		item.Id = *result.Id
	}
	if result.ScriptId != nil {
		item.ScriptId = *result.ScriptId
	}
	if result.Version != nil {
		item.Version = *result.Version
	}
	if result.Content != nil {
		item.Content = *result.Content
	}
	if result.ScriptType != nil {
		item.ScriptType = *result.ScriptType
	}
	if result.Checksum != nil {
		item.ContentHash = *result.Checksum
	}
	if result.CreatedBy != nil {
		item.CreatedBy = *result.CreatedBy
	}
	if result.CreatedAt != nil {
		item.CreatedAt = *result.CreatedAt
	}
	if result.ChangeLog != nil {
		item.Remark = *result.ChangeLog
	}

	resp = &types.ScriptVersionDetailResp{
		Code: 0,
		Msg:  "success",
		Data: item,
	}

	return resp, nil
}
