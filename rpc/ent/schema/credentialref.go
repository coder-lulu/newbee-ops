package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema/field"
    "github.com/coder-lulu/newbee-common/v2/orm/ent/mixins"
)

// CredentialRef 凭证引用（供会话/任务引用外部密管或本地凭证）
// 说明：仅定义 Schema，实际生成需按 CLAUDE.md 执行 make gen-ent
type CredentialRef struct{ ent.Schema }

func (CredentialRef) Mixin() []ent.Mixin {
    return []ent.Mixin{mixins.IDMixin{}, mixins.StatusMixin{}, mixins.TenantMixin{}}
}

func (CredentialRef) Fields() []ent.Field {
    return []ent.Field{
        field.String("provider").Annotations(entsql.Annotation{Size: 64}).Comment("凭证提供方，如 vault/local") ,
        field.String("ref").Annotations(entsql.Annotation{Size: 256}).Comment("外部引用ID或路径") ,
        field.String("scope").Annotations(entsql.Annotation{Size: 64}).Optional().Comment("作用域，如 global/tenant/user") ,
        field.String("created_by").Annotations(entsql.Annotation{Size: 64}).Optional().Comment("创建者") ,
        field.JSON("tags", map[string]string{}).Optional().Comment("扩展标签") ,
    }
}

