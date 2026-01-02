package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema/field"
    "github.com/coder-lulu/newbee-common/v2/orm/ent/mixins"
)

// AuditEvent 审计事件索引（基础元数据，不存储大对象）
type AuditEvent struct{ ent.Schema }

func (AuditEvent) Mixin() []ent.Mixin {
    return []ent.Mixin{mixins.IDMixin{}, mixins.StatusMixin{}, mixins.TenantMixin{}}
}

func (AuditEvent) Fields() []ent.Field {
    return []ent.Field{
        field.String("event_id").Unique().Annotations(entsql.Annotation{Size: 128}).Comment("事件ID") ,
        field.String("type").Annotations(entsql.Annotation{Size: 32}).Comment("事件类型：session|task|profile|proxy") ,
        field.String("subject_user").Annotations(entsql.Annotation{Size: 64}).Optional().Comment("用户ID") ,
        field.String("subject_ci").Annotations(entsql.Annotation{Size: 128}).Optional().Comment("CI ID") ,
        field.String("subject_proxy").Annotations(entsql.Annotation{Size: 128}).Optional().Comment("Proxy ID") ,
        field.Int64("occur_at").Optional().Comment("发生时间(Unix秒)") ,
        field.JSON("meta", map[string]string{}).Optional().Comment("附加元数据") ,
    }
}

