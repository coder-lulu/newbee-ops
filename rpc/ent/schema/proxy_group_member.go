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

// ProxyGroupMember Proxy分组成员关联实体
type ProxyGroupMember struct {
	ent.Schema
}

func (ProxyGroupMember) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TenantMixin{},
		mixins.DepartmentMixin{},
		mixins.StatusMixin{},
	}
}

func (ProxyGroupMember) Fields() []ent.Field {
	return []ent.Field{
		// ========== 关联关系 ==========
		field.Uint64("worker_id").
			Comment("Proxy ID"),

		field.Uint64("worker_group_id").
			Comment("Proxy分组ID"),

		// ========== 负载均衡参数 ==========
		field.Int("weight").
			Default(100).
			Range(1, 1000).
			Comment("权重（1-1000），用于加权轮询策略"),

		field.Int("priority").
			Default(0).
			Range(-100, 100).
			Comment("优先级（-100到100），数字越大优先级越高"),

		// ========== 状态信息 ==========
		field.Time("joined_at").
			Optional().
			Nillable().
			Comment("加入分组的时间"),

		field.Time("last_selected_at").
			Optional().
			Nillable().
			Comment("最后被选中的时间（用于轮询算法）"),

		field.Int("select_count").
			Default(0).
			Comment("被选中次数（统计字段）"),

		// ========== 健康状态 ==========
		field.Bool("is_backup").
			Default(false).
			Comment("是否为备用节点（仅在主节点全部故障时使用）"),

		field.Int("consecutive_failures").
			Default(0).
			Comment("连续失败次数（用于自动摘除）"),

		field.Time("last_failure_at").
			Optional().
			Nillable().
			Comment("最后失败时间"),
	}
}

func (ProxyGroupMember) Edges() []ent.Edge {
	return []ent.Edge{
		// 关联到Proxy表
		edge.To("proxy", Proxy.Type).
			Unique().
			Required().
			Field("worker_id"),

		// 关联到ProxyGroup表
		edge.To("proxy_group", ProxyGroup.Type).
			Unique().
			Required().
			Field("worker_group_id"),
	}
}

func (ProxyGroupMember) Indexes() []ent.Index {
	return []ent.Index{
		// 唯一约束：同一个Proxy不能重复加入同一个分组
		index.Fields("worker_id", "worker_group_id").Unique(),

		// 分组查询（最常用）- 按优先级和权重排序
		index.Fields("worker_group_id", "priority", "weight"),

		// Proxy查询（查询某个Proxy属于哪些分组）
		index.Fields("worker_id"),

		// 租户查询
		index.Fields("tenant_id"),

		// 状态查询
		index.Fields("status"),
	}
}

func (ProxyGroupMember) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ops_worker_group_members"},
	}
}
