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

type ScriptCategory struct {
	ent.Schema
}

func (ScriptCategory) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TenantMixin{},
		mixins.DepartmentMixin{},
		mixins.StatusMixin{},
	}
}

func (ScriptCategory) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("分类名称").
			NotEmpty().
			MaxLen(100),

		field.String("code").
			Comment("分类标识代码").
			NotEmpty().
			MaxLen(100),

		field.String("description").
			Comment("分类描述").
			Optional().
			MaxLen(500),

		field.Uint64("parent_id").
			Comment("父分类ID（支持树形结构）").
			Optional(),

		field.Uint32("sort_order").
			Comment("排序顺序").
			Default(0),

		field.String("icon").
			Comment("分类图标").
			Optional().
			MaxLen(100),
	}
}

func (ScriptCategory) Edges() []ent.Edge {
	return []ent.Edge{
		// 父分类关系
		edge.To("children", ScriptCategory.Type).
			From("parent").
			Field("parent_id").
			Unique(),

		// 包含的脚本
		edge.To("scripts", Script.Type),
	}
}

func (ScriptCategory) Indexes() []ent.Index {
	return []ent.Index{
		// 租户内code唯一
		index.Fields("tenant_id", "code").
			Unique(),

		// 父分类查询
		index.Fields("tenant_id", "parent_id"),

		// 排序查询
		index.Fields("sort_order"),
	}
}

func (ScriptCategory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "ops_script_categories",
		},
	}
}
