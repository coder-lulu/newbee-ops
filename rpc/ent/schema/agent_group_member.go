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

// AgentGroupMember Agent组成员关联实体
type AgentGroupMember struct {
	ent.Schema
}

func (AgentGroupMember) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TenantMixin{},
		mixins.StatusMixin{},
	}
}

func (AgentGroupMember) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("agent_id").
			Comment("Agent ID"),

		field.Uint64("agent_group_id").
			Comment("Agent组ID"),

		field.Int("priority").
			Default(100).
			Comment("优先级（数字越小优先级越高）"),

		field.Int("weight").
			Default(1).
			Comment("权重（用于加权轮询）"),

		field.Time("joined_at").
			Optional().
			Nillable().
			Comment("加入时间"),
	}
}

func (AgentGroupMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("agent", Agent.Type).
			Unique().
			Required().
			Field("agent_id"),

		edge.To("agent_group", AgentGroup.Type).
			Unique().
			Required().
			Field("agent_group_id"),
	}
}

func (AgentGroupMember) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("agent_id", "agent_group_id").Unique(),
		index.Fields("agent_group_id"),
		index.Fields("tenant_id"),
	}
}

func (AgentGroupMember) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ops_agent_group_members"},
	}
}
