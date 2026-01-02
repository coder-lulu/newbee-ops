package script

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScriptVersionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetScriptVersionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptVersionListLogic {
	return &GetScriptVersionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetScriptVersionListLogic) GetScriptVersionList(req *types.ScriptVersionListReq) (resp *types.ScriptVersionListResp, err error) {
	// 调用RPC获取脚本版本列表
	result, err := l.svcCtx.OpsClient.GetScriptVersionList(l.ctx, req.Page, req.PageSize, req.ScriptId)
	if err != nil {
		l.Errorw("Failed to get script version list from RPC",
			logx.Field("error", err),
			logx.Field("scriptId", req.ScriptId))
		return nil, err
	}

	// 转换数据格式
	items := make([]types.ScriptVersionItem, 0, len(result.Data))
	for _, v := range result.Data {
		item := types.ScriptVersionItem{}
		if v.Id != nil {
			item.Id = *v.Id
		}
		if v.ScriptId != nil {
			item.ScriptId = *v.ScriptId
		}
		if v.Version != nil {
			item.Version = *v.Version
		}
		if v.Content != nil {
			item.Content = *v.Content
		}
		if v.ScriptType != nil {
			item.ScriptType = *v.ScriptType
		}
		if v.Checksum != nil {
			item.ContentHash = *v.Checksum
		}
		if v.CreatedBy != nil {
			item.CreatedBy = *v.CreatedBy
		}
		if v.CreatedAt != nil {
			item.CreatedAt = *v.CreatedAt
		}
		if v.ChangeLog != nil {
			item.Remark = *v.ChangeLog
		}
		items = append(items, item)
	}

	// 构造响应
	resp = &types.ScriptVersionListResp{
		Code: 0,
		Msg:  "success",
		Data: types.ScriptVersionListData{
			Total: result.Total,
			Data:  items,
		},
	}

	return resp, nil
}
