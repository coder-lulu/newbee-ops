package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/coder-lulu/newbee-common/v2/orm/ent/mixins"
)

// AgentGroup Agent组实体
type AgentGroup struct {
	ent.Schema
}

func (AgentGroup) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TenantMixin{},
		mixins.StatusMixin{},
	}
}

func (AgentGroup) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(100).
			NotEmpty().
			Comment("Agent组名称"),

		field.String("description").
			MaxLen(500).
			Optional().
			Comment("描述信息"),

		field.Enum("selection_strategy").
			Values("round_robin", "least_loaded", "random", "priority", "geo_nearest").
			Default("round_robin").
			Comment("Agent选择策略"),

		field.Int("health_check_interval").
			Default(60).
			Comment("健康检查间隔（秒）"),

		field.Bool("auto_failover").
			Default(true).
			Comment("是否启用自动故障转移"),

		field.Int("max_retry_count").
			Default(3).
			Comment("Agent失败时最大重试次数"),
	}
}

func (AgentGroup) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("status"),
	}
}

func (AgentGroup) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ops_agent_groups"},
	}
}
