package proxy

import proxyservice "github.com/coder-lulu/newbee-ops-api/internal/services/proxy"

type ListLogic struct { r *proxyservice.Registry }
func NewListLogic(r *proxyservice.Registry) *ListLogic { return &ListLogic{r: r} }
func (l *ListLogic) List() any { if l.r == nil { return []any{} }; return l.r.List() }
