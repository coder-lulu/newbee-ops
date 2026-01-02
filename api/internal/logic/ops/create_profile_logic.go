package ops

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProfileLogic {
	return &CreateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProfileLogic) CreateProfile(req *types.AccessProfileReq) error {
	// todo: add your logic here and delete this line

	return nil
}
