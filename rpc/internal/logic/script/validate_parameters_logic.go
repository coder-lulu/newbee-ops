package script

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateParametersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateParametersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateParametersLogic {
	return &ValidateParametersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ValidateParametersLogic) ValidateParameters(in *ops.ValidateParametersReq) (*ops.ValidateParametersResp, error) {
	// 确定JSON Schema来源
	var schemaJSON string

	if in.ScriptId != nil && *in.ScriptId > 0 {
		// 从数据库获取脚本的参数Schema
		script, err := l.svcCtx.DB.Script.Get(l.ctx, *in.ScriptId)
		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}

		if len(script.Parameters) == 0 {
			// 脚本没有定义参数Schema，视为验证通过
			return &ops.ValidateParametersResp{
				Valid:   true,
				Errors:  []string{},
				Message: pointy.GetPointer("No parameter schema defined, validation skipped"),
			}, nil
		}

// 		schemaJSON = script.Parameters
	} else if in.SchemaJson != nil && *in.SchemaJson != "" {
		// 使用请求中提供的Schema
		schemaJSON = *in.SchemaJson
	} else {
		// 既没有script_id也没有schema_json
		return &ops.ValidateParametersResp{
			Valid:   false,
			Errors:  []string{"either script_id or schema_json must be provided"},
			Message: pointy.GetPointer("Missing schema source"),
		}, nil
	}

	// 解析参数JSON字符串为map
	var parameters map[string]interface{}
	if err := json.Unmarshal([]byte(in.ParametersJson), &parameters); err != nil {
		return &ops.ValidateParametersResp{
			Valid:   false,
			Errors:  []string{fmt.Sprintf("invalid parameters JSON: %v", err)},
			Message: pointy.GetPointer("Parameters JSON parse error"),
		}, nil
	}

	// 使用ParameterValidator验证
	valid, errors, err := l.svcCtx.ParameterValidator.ValidateWithDetail(schemaJSON, parameters)
	if err != nil {
		// 验证过程出错（通常是schema格式错误）
		l.Errorf("Parameter validation error: %v", err)
		return &ops.ValidateParametersResp{
			Valid:   false,
			Errors:  []string{fmt.Sprintf("validation error: %v", err)},
			Message: pointy.GetPointer("Schema validation error"),
		}, nil
	}

	// 返回验证结果
	if valid {
		return &ops.ValidateParametersResp{
			Valid:   true,
			Errors:  []string{},
			Message: pointy.GetPointer("Parameters validation passed"),
		}, nil
	}

	return &ops.ValidateParametersResp{
		Valid:   false,
		Errors:  errors,
		Message: pointy.GetPointer("Parameters validation failed"),
	}, nil
}
