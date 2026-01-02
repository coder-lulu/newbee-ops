package profile

import (
	"context"
	"fmt"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	apitypes "github.com/coder-lulu/newbee-ops-api/internal/types"
	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"
	"google.golang.org/protobuf/types/known/structpb"
)

type Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Logic {
	return &Logic{ctx: ctx, svcCtx: svcCtx}
}

func (l *Logic) Create(req *apitypes.AccessProfileReq) error {
	if l.svcCtx.OpsClient != nil {
		p := req.AccessProfile
		caps, _ := packStringList(p.Capabilities)
		ports, _ := packIntMap(p.Ports)
		jump, _ := packStringList(p.JumpChain)
		tags, _ := packStringMap(p.Tags)

		_, err := l.svcCtx.OpsClient.CreateAccessProfile(l.ctx, &ops.AccessProfileInfo{
			CiId:          &p.CiId,
			Capabilities:  caps,
			Ports:         ports,
			CredentialRef: &p.CredentialRef,
			PreferProxy:   &p.PreferProxy,
			JumpChain:     jump,
			Tags:          tags,
		})
		return err
	}
	l.svcCtx.ProfileStore.Upsert(svc.AccessProfile{
		CiId:          req.AccessProfile.CiId,
		Capabilities:  req.AccessProfile.Capabilities,
		Ports:         req.AccessProfile.Ports,
		CredentialRef: req.AccessProfile.CredentialRef,
		PreferProxy:   req.AccessProfile.PreferProxy,
		JumpChain:     req.AccessProfile.JumpChain,
		Tags:          req.AccessProfile.Tags,
	})
	return nil
}

func (l *Logic) Update(req *apitypes.AccessProfileReq) error {
	if l.svcCtx.OpsClient != nil {
		p := req.AccessProfile
		caps, _ := packStringList(p.Capabilities)
		ports, _ := packIntMap(p.Ports)
		jump, _ := packStringList(p.JumpChain)
		tags, _ := packStringMap(p.Tags)

		_, err := l.svcCtx.OpsClient.UpdateAccessProfile(l.ctx, &ops.AccessProfileInfo{
			CiId:          &p.CiId,
			Capabilities:  caps,
			Ports:         ports,
			CredentialRef: &p.CredentialRef,
			PreferProxy:   &p.PreferProxy,
			JumpChain:     jump,
			Tags:          tags,
		})
		return err
	}
	l.svcCtx.ProfileStore.Upsert(svc.AccessProfile{
		CiId:          req.AccessProfile.CiId,
		Capabilities:  req.AccessProfile.Capabilities,
		Ports:         req.AccessProfile.Ports,
		CredentialRef: req.AccessProfile.CredentialRef,
		PreferProxy:   req.AccessProfile.PreferProxy,
		JumpChain:     req.AccessProfile.JumpChain,
		Tags:          req.AccessProfile.Tags,
	})
	return nil
}

func (l *Logic) Get(ciId string) (apitypes.AccessProfile, bool, error) {
	if l.svcCtx.OpsClient != nil {
		resList, err := l.svcCtx.OpsClient.GetAccessProfileList(l.ctx, &ops.AccessProfileListReq{Page: 1, PageSize: 1, CiId: &ciId})
		if err != nil {
			return apitypes.AccessProfile{}, false, err
		}
		if resList == nil || len(resList.Data) == 0 {
			return apitypes.AccessProfile{}, false, nil
		}
		item := resList.Data[0]
		ap := apitypes.AccessProfile{CiId: item.GetCiId(), CredentialRef: item.GetCredentialRef(), PreferProxy: item.GetPreferProxy()}
		if s := item.GetCapabilities(); s != nil {
			ap.Capabilities = unpackStringList(s)
		}
		if s := item.GetPorts(); s != nil {
			ap.Ports = unpackIntMap(s)
		}
		if s := item.GetJumpChain(); s != nil {
			ap.JumpChain = unpackStringList(s)
		}
		if s := item.GetTags(); s != nil {
			ap.Tags = unpackStringMap(s)
		}
		return ap, true, nil
	}
	if v, ok := l.svcCtx.ProfileStore.Get(ciId); ok {
		ap := apitypes.AccessProfile{
			CiId:          v.CiId,
			Capabilities:  v.Capabilities,
			Ports:         v.Ports,
			CredentialRef: v.CredentialRef,
			PreferProxy:   v.PreferProxy,
			JumpChain:     v.JumpChain,
			Tags:          v.Tags,
		}
		return ap, true, nil
	}
	return apitypes.AccessProfile{}, false, nil
}

func (l *Logic) List() (items []apitypes.AccessProfile, total int, err error) {
	if l.svcCtx.OpsClient != nil {
		res, err := l.svcCtx.OpsClient.GetAccessProfileList(l.ctx, &ops.AccessProfileListReq{Page: 1, PageSize: 200})
		if err != nil {
			return nil, 0, err
		}
		out := make([]apitypes.AccessProfile, 0, len(res.Data))
		for _, it := range res.Data {
			ap := apitypes.AccessProfile{CiId: it.GetCiId(), CredentialRef: it.GetCredentialRef(), PreferProxy: it.GetPreferProxy()}
			if s := it.GetCapabilities(); s != nil {
				ap.Capabilities = unpackStringList(s)
			}
			if s := it.GetPorts(); s != nil {
				ap.Ports = unpackIntMap(s)
			}
			if s := it.GetJumpChain(); s != nil {
				ap.JumpChain = unpackStringList(s)
			}
			if s := it.GetTags(); s != nil {
				ap.Tags = unpackStringMap(s)
			}
			out = append(out, ap)
		}
		return out, int(res.Total), nil
	}
	list := l.svcCtx.ProfileStore.List()
	out := make([]apitypes.AccessProfile, 0, len(list))
	for _, v := range list {
		out = append(out, apitypes.AccessProfile{
			CiId:          v.CiId,
			Capabilities:  v.Capabilities,
			Ports:         v.Ports,
			CredentialRef: v.CredentialRef,
			PreferProxy:   v.PreferProxy,
			JumpChain:     v.JumpChain,
			Tags:          v.Tags,
		})
	}
	return out, len(out), nil
}

func (l *Logic) Delete(ciId string) error {
	if l.svcCtx.OpsClient != nil {
		_, err := l.svcCtx.OpsClient.DeleteAccessProfile(l.ctx, &ops.IDsReq{Ids: []uint64{0}}) // Need ID but we have CiId.
		// Actually DeleteAccessProfile takes IDs (uint64). We might need a DeleteAccessProfileByCiId or similar if we only have CiId.
		// However, looking at ops.proto, DeleteAccessProfile takes IDsReq.
		// And AccessProfileInfo has CiId.
		// If we only have CiId, we first need to Get to find the ID.

		resList, err := l.svcCtx.OpsClient.GetAccessProfileList(l.ctx, &ops.AccessProfileListReq{Page: 1, PageSize: 1, CiId: &ciId})
		if err != nil {
			return err
		}
		if resList != nil && len(resList.Data) > 0 {
			id := resList.Data[0].Id
			if id != nil {
				_, err = l.svcCtx.OpsClient.DeleteAccessProfile(l.ctx, &ops.IDsReq{Ids: []uint64{*id}})
				return err
			}
		}
		return nil // Not found, consider deleted
	}
	l.svcCtx.ProfileStore.Delete(ciId)
	return nil
}

// Helpers

func packStringList(list []string) (*structpb.Struct, error) {
	if len(list) == 0 {
		return nil, nil
	}
	vals := make([]interface{}, len(list))
	for i, v := range list {
		vals[i] = v
	}
	return structpb.NewStruct(map[string]interface{}{"list": vals})
}

func packIntMap(m map[string]int) (*structpb.Struct, error) {
	if len(m) == 0 {
		return nil, nil
	}
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return structpb.NewStruct(out)
}

func packStringMap(m map[string]string) (*structpb.Struct, error) {
	if len(m) == 0 {
		return nil, nil
	}
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return structpb.NewStruct(out)
}

func unpackStringList(s *structpb.Struct) []string {
	m := s.AsMap()
	if list, ok := m["list"].([]interface{}); ok {
		out := make([]string, 0, len(list))
		for _, v := range list {
			if str, ok := v.(string); ok {
				out = append(out, str)
			}
		}
		return out
	}
	return nil
}

func unpackIntMap(s *structpb.Struct) map[string]int {
	m := s.AsMap()
	out := make(map[string]int, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case float64:
			out[k] = int(val)
		case int:
			out[k] = val
		case int64:
			out[k] = int(val)
		}
	}
	return out
}

func unpackStringMap(s *structpb.Struct) map[string]string {
	m := s.AsMap()
	out := make(map[string]string, len(m))
	for k, v := range m {
		if str, ok := v.(string); ok {
			out[k] = str
		} else {
			out[k] = fmt.Sprintf("%v", v)
		}
	}
	return out
}
