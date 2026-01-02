package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema/field"
    "github.com/coder-lulu/newbee-common/v2/orm/ent/mixins"
)

// Session 实体（持久化会话）
// 按 CLAUDE.md 要求：包含 TenantMixin 与审计字段；高变字段采用 JSON。
type Session struct{ ent.Schema }

func (Session) Mixin() []ent.Mixin {
    return []ent.Mixin{
        mixins.IDMixin{},
        mixins.StatusMixin{},
        mixins.TenantMixin{},
    }
}

func (Session) Fields() []ent.Field {
    return []ent.Field{
        field.String("session_id").Unique().Annotations(entsql.Annotation{Size: 128}).Comment("会话ID"),
        field.String("user_id").Optional().Annotations(entsql.Annotation{Size: 64}).Comment("用户ID"),
        field.String("ci_id").Annotations(entsql.Annotation{Size: 128}).Comment("CMDB CI ID"),
        field.String("protocol").Annotations(entsql.Annotation{Size: 16}).Comment("协议: ssh|telnet|rdp|vnc"),
        field.String("proxy_id").Optional().Annotations(entsql.Annotation{Size: 128}).Comment("选路Proxy ID"),
        field.String("endpoint").Optional().Annotations(entsql.Annotation{Size: 256}).Comment("接入端点"),

        // Phase 3: Worker绑定字段
        field.String("worker_id").Optional().Annotations(entsql.Annotation{Size: 100}).Comment("绑定的Worker ID"),
        field.String("worker_ip").Optional().Annotations(entsql.Annotation{Size: 100}).Comment("Worker IP地址"),
        field.Int("worker_port").Optional().Comment("Worker端口"),

        field.String("status_str").Default("active").Annotations(entsql.Annotation{Size: 16}).Comment("active|closed"),
        field.Int64("expires_at").Optional().Comment("过期时间(Unix秒)"),
        field.Int64("closed_at").Optional().Comment("关闭时间(Unix秒)"),
        field.JSON("tags", map[string]string{}).Optional().Comment("标签/扩展"),
    }
}
