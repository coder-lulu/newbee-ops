package base

import (
	"context"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/coder-lulu/newbee-common/v2/errors"
	"github.com/coder-lulu/newbee-common/v2/msg/errormsg"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitDatabaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitDatabaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitDatabaseLogic {
	return &InitDatabaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InitDatabaseLogic) InitDatabase(in *ops.Empty) (*ops.BaseResp, error) {
	// This database may also contain Core, CMDB and IO tables. Only migrate
	// the Ops schema, retaining existing columns, indexes and business data.
	if err := l.svcCtx.DB.Schema.Create(l.ctx,
		schema.WithForeignKeys(false),
		schema.WithDropColumn(false),
		schema.WithDropIndex(false),
	); err != nil {
		l.Errorf("Failed to initialize Ops database: %v", err)
		return nil, errors.InternalWithCause("Ops database initialization failed", err)
	}
	if err := l.insertCoreData(); err != nil {
		return nil, errors.InternalWithCause("Ops catalog initialization failed; initialize Core first", err)
	}
	return &ops.BaseResp{Msg: errormsg.Success}, nil
}
