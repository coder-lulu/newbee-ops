package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/coder-lulu/newbee-common/v2/orm/ent/mixins"
)

// ProxyGroup Proxy分组实体 - 将Proxy组织成逻辑分组，实现负载均衡和故障转移
type ProxyGroup struct {
	ent.Schema
}

func (ProxyGroup) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TenantMixin{},
		mixins.DepartmentMixin{},
		mixins.StatusMixin{},
	}
}

func (ProxyGroup) Fields() []ent.Field {
	return []ent.Field{
		// ========== 基础信息 ==========
		field.String("name").
			MaxLen(100).
			NotEmpty().
			Comment("Proxy分组名称"),

		field.String("description").
			MaxLen(500).
			Optional().
			Comment("描述信息"),

		// ========== 负载均衡策略 ==========
		field.Enum("selection_strategy").
			Values(
				"round_robin",        // 轮询
				"least_connections",  // 最少连接数
				"weighted",           // 加权轮询
				"random",             // 随机
				"consistent_hash",    // 一致性哈希（基于session_id）
			).
			Default("least_connections").
			Comment("Proxy选择策略"),

		// ========== 健康检查配置 ==========
		field.Int("health_check_interval").
			Default(60).
			Range(10, 3600).
			Comment("健康检查间隔（秒），10-3600秒"),

		field.Bool("auto_failover").
			Default(true).
			Comment("是否启用自动故障转移"),

		field.Int("max_retry_count").
			Default(3).
			Range(0, 10).
			Comment("Proxy失败时最大重试次数（0-10次）"),

		field.Int("min_healthy_workers").
			Default(1).
			Range(1, 100).
			Comment("最少健康Proxy数量，低于此值告警"),

		// ========== 超时配置 ==========
		field.Int("connection_timeout").
			Default(30).
			Range(5, 300).
			Comment("连接超时时间（秒），5-300秒"),

		field.Int("request_timeout").
			Default(60).
			Range(10, 600).
			Comment("请求超时时间（秒），10-600秒"),

		// ========== 统计信息（缓存字段，定期更新） ==========
		field.Int("total_members").
			Default(0).
			Comment("总成员数（缓存字段）"),

		field.Int("online_members").
			Default(0).
			Comment("在线成员数（缓存字段）"),

		field.Int("total_weight").
			Default(0).
			Comment("总权重（缓存字段，用于加权算法）"),
	}
}

func (ProxyGroup) Indexes() []ent.Index {
	return []ent.Index{
		// 租户+状态查询（最常用）
		index.Fields("tenant_id", "status"),

		// 租户查询
		index.Fields("tenant_id"),

		// 名称模糊搜索
		index.Fields("name"),

		// 状态查询
		index.Fields("status"),
	}
}

func (ProxyGroup) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ops_worker_groups"},
	}
}
