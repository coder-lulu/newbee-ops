package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/coder-lulu/newbee-common/v2/orm/ent/mixins"
)

type ScriptVersion struct {
	ent.Schema
}

func (ScriptVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TenantMixin{},
	}
}

func (ScriptVersion) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("script_id").
			Comment("脚本ID"),

		field.String("version").
			Comment("版本号").
			NotEmpty().
			MaxLen(50),

		// 内容快照
		field.Text("content").
			Comment("脚本内容快照"),

		field.JSON("parameters", map[string]interface{}{}).
			Comment("参数定义快照（JSON Schema格式）").
			Optional(),

		field.String("script_type").
			Comment("脚本类型快照").
			MaxLen(50),

		field.String("executor").
			Comment("执行器类型快照").
			MaxLen(50),

		// 变更信息
		field.Text("change_log").
			Comment("变更日志").
			Optional(),

		field.Uint64("created_by").
			Comment("创建者用户ID"),

		// Note: created_at and updated_at are provided by IDMixin

		// 防篡改
		field.String("checksum").
			Comment("内容校验和（SHA256）").
			MaxLen(64),
	}
}

func (ScriptVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("script", Script.Type).
			Ref("versions").
			Field("script_id").
			Unique().
			Required(),
	}
}

func (ScriptVersion) Indexes() []ent.Index {
	return []ent.Index{
		// 脚本版本查询
		index.Fields("script_id", "version"),

		// 租户隔离
		index.Fields("tenant_id", "script_id"),
	}
}

func (ScriptVersion) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "ops_script_versions",
		},
	}
}
