package accessprofile

import (
    "context"
    "fmt"

    "github.com/coder-lulu/newbee-ops-rpc/internal/svc"
    "github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
    newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/coder-lulu/newbee-common/v2/i18n"
    "github.com/zeromicro/go-zero/core/logx"
)

type CreateAccessProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAccessProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAccessProfileLogic {
	return &CreateAccessProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAccessProfileLogic) CreateAccessProfile(in *newbee_ops_rpc.AccessProfileInfo) (*newbee_ops_rpc.BaseIDResp, error) {
	q := l.svcCtx.DB.AccessProfile.Create()
	if in.CiId != nil {
		q = q.SetCiID(in.GetCiId())
	}
	// 转换structpb.Struct到[]string (capabilities)
	if in.Capabilities != nil {
		capsAny := in.GetCapabilities().AsMap()
		caps := make([]string, 0)
		if list, ok := capsAny["list"].([]interface{}); ok {
			for _, v := range list {
				if str, ok := v.(string); ok {
					caps = append(caps, str)
				} else {
					l.Infow("skipping non-string capability value",
						logx.Field("type", fmt.Sprintf("%T", v)),
						logx.Field("value", v))
				}
			}
		}
		if len(caps) > 0 {
			q = q.SetCapabilities(caps)
		}
	}
	// 转换structpb.Struct到map[string]int (ports)
	if in.Ports != nil {
		portsAny := in.GetPorts().AsMap()
		ports := make(map[string]int)
		for k, v := range portsAny {
			var port int
			switch val := v.(type) {
			case float64:
				port = int(val)
			case int:
				port = val
			case int64:
				port = int(val)
			default:
				l.Infow("skipping non-numeric port value",
					logx.Field("key", k),
					logx.Field("type", fmt.Sprintf("%T", v)),
					logx.Field("value", v))
				continue
			}
			// 验证端口范围 0-65535
			if port < 0 || port > 65535 {
				l.Errorw("port number out of valid range (0-65535)",
					logx.Field("key", k),
					logx.Field("port", port))
				return nil, fmt.Errorf("invalid port number for %s: %d (must be 0-65535)", k, port)
			}
			ports[k] = port
		}
		if len(ports) > 0 {
			q = q.SetPorts(ports)
		}
	}
	if in.CredentialRef != nil {
		q = q.SetCredentialRef(in.GetCredentialRef())
	}
	if in.PreferProxy != nil {
		q = q.SetPreferProxy(in.GetPreferProxy())
	}
	// 转换structpb.Struct到[]string (jump_chain)
	if in.JumpChain != nil {
		jumpAny := in.GetJumpChain().AsMap()
		jump := make([]string, 0)
		if list, ok := jumpAny["list"].([]interface{}); ok {
			for _, v := range list {
				if str, ok := v.(string); ok {
					jump = append(jump, str)
				} else {
					l.Infow("skipping non-string jump_chain value",
						logx.Field("type", fmt.Sprintf("%T", v)),
						logx.Field("value", v))
				}
			}
		}
		if len(jump) > 0 {
			q = q.SetJumpChain(jump)
		}
	}
	// 转换structpb.Struct到map[string]string (tags)
	if in.Tags != nil {
		tagsAny := in.GetTags().AsMap()
		tags := make(map[string]string)
		for k, v := range tagsAny {
			if str, ok := v.(string); ok {
				tags[k] = str
			} else {
				l.Infow("skipping non-string tag value",
					logx.Field("key", k),
					logx.Field("type", fmt.Sprintf("%T", v)),
					logx.Field("value", v))
			}
		}
		if len(tags) > 0 {
			q = q.SetTags(tags)
		}
	}
	if in.Status != nil {
		q = q.SetStatus(uint8(in.GetStatus()))
	}
	res, err := q.Save(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}
	return &newbee_ops_rpc.BaseIDResp{Id: res.ID, Msg: i18n.CreateSuccess}, nil
}
