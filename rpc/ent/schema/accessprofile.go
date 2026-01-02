package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema/field"
    "github.com/coder-lulu/newbee-common/v2/orm/ent/mixins"
)

// AccessProfile 按 CLAUDE.md 要求：包含 TenantMixin
type AccessProfile struct{ ent.Schema }

func (AccessProfile) Mixin() []ent.Mixin {
    return []ent.Mixin{
        mixins.IDMixin{},
        mixins.StatusMixin{},
        mixins.TenantMixin{},
    }
}

func (AccessProfile) Fields() []ent.Field {
    return []ent.Field{
        field.String("ci_id").Unique().Comment("CMDB CI ID").Annotations(entsql.Annotation{Size: 128}),
        field.JSON("capabilities", []string{}).Optional().Comment("协议能力列表"),
        field.JSON("ports", map[string]int{}).Optional().Comment("协议端口映射"),
        field.String("credential_ref").Optional().Comment("凭证引用"),
        field.String("prefer_proxy").Optional().Comment("首选Proxy ID"),
        field.JSON("jump_chain", []string{}).Optional().Comment("跳板链"),
        field.JSON("tags", map[string]string{}).Optional().Comment("标签"),
    }
}

