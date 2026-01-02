package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/coder-lulu/newbee-common/v2/orm/ent/mixins"
)

// ProxyMetrics Proxy指标快照（时序数据）
type ProxyMetrics struct {
	ent.Schema
}

func (ProxyMetrics) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TenantMixin{},
	}
}

func (ProxyMetrics) Fields() []ent.Field {
	return []ent.Field{
		field.String("worker_id").
			MaxLen(100).
			NotEmpty().
			Comment("Proxy唯一标识"),

		field.Time("timestamp").
			Comment("指标采集时间戳"),

		// 资源指标
		field.Float("cpu_usage").
			Default(0).
			Comment("CPU使用率（%）"),

		field.Float("memory_usage").
			Default(0).
			Comment("内存使用率（%）"),

		field.Float("disk_usage").
			Default(0).
			Comment("磁盘使用率（%）"),

		field.Uint64("network_in_delta").
			Default(0).
			Comment("网络入流量增量（字节）"),

		field.Uint64("network_out_delta").
			Default(0).
			Comment("网络出流量增量（字节）"),

		// 业务指标
		field.Int("active_sessions").
			Default(0).
			Comment("活跃会话数（快照）"),

		field.Uint64("request_count_delta").
			Default(0).
			Comment("请求数增量"),

		field.Uint64("success_count_delta").
			Default(0).
			Comment("成功数增量"),

		field.Uint64("failure_count_delta").
			Default(0).
			Comment("失败数增量"),

		field.Float("avg_latency_ms").
			Default(0).
			Comment("平均延迟（毫秒）"),

		// 状态
		field.String("worker_status").
			Default("offline").
			Comment("Proxy状态快照"),
	}
}

func (ProxyMetrics) Indexes() []ent.Index {
	return []ent.Index{
		// 时序查询（最常用）
		index.Fields("worker_id", "timestamp"),

		// 租户查询
		index.Fields("tenant_id", "timestamp"),

		// TTL清理（按时间删除旧数据）
		index.Fields("timestamp"),
	}
}

func (ProxyMetrics) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ops_worker_metrics"},
	}
}
