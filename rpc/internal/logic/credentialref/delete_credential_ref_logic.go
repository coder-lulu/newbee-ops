package credentialref

import (
    "context"

    "github.com/coder-lulu/newbee-ops-rpc/ent/credentialref"
    "github.com/coder-lulu/newbee-ops-rpc/internal/svc"
    "github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
    newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/coder-lulu/newbee-common/v2/i18n"
    "github.com/zeromicro/go-zero/core/logx"
)

type DeleteCredentialRefLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCredentialRefLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCredentialRefLogic {
	return &DeleteCredentialRefLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCredentialRefLogic) DeleteCredentialRef(in *newbee_ops_rpc.IDsReq) (*newbee_ops_rpc.BaseResp, error) {
	_, err := l.svcCtx.DB.CredentialRef.Delete().Where(credentialref.IDIn(in.Ids...)).Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &newbee_ops_rpc.BaseResp{Msg: i18n.DeleteSuccess}, nil
}
