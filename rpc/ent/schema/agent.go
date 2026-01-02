package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/coder-lulu/newbee-common/v2/orm/ent/mixins"
)

// Agent 代理节点实体
type Agent struct {
	ent.Schema
}

func (Agent) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TenantMixin{},
		mixins.StatusMixin{},
	}
}

func (Agent) Fields() []ent.Field {
	return []ent.Field{
		// 基础信息
		field.String("name").
			MaxLen(100).
			NotEmpty().
			Comment("Agent名称"),

		field.String("agent_id").
			MaxLen(100).
			NotEmpty().
			Unique().
			Comment("Agent唯一标识（由Agent生成）"),

		field.String("host").
			MaxLen(255).
			NotEmpty().
			Comment("Agent主机地址或IP"),

		field.Int("port").
			Default(8080).
			Comment("Agent服务端口"),

		field.String("api_key").
			MaxLen(255).
			Sensitive().
			Comment("API认证密钥"),

		field.String("description").
			MaxLen(500).
			Optional().
			Comment("描述信息"),

		// 状态字段
		field.Enum("agent_status").
			Values("online", "offline", "error", "maintenance").
			Default("offline").
			Comment("Agent运行状态"),

		field.Time("last_heartbeat").
			Optional().
			Nillable().
			Comment("最后心跳时间"),

		field.Int("heartbeat_interval").
			Default(30).
			Comment("心跳间隔（秒）"),

		field.Time("last_online_at").
			Optional().
			Nillable().
			Comment("最后在线时间"),

		// 能力字段
		field.JSON("supported_providers", []string{}).
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL: "json",
			}).
			Comment("支持的Provider类型列表"),

		field.JSON("capabilities", []string{}).
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL: "json",
			}).
			Comment("支持的协议能力 [ssh, telnet, rdp, ipmi, snmp]"),

		field.Int("max_concurrent_tasks").
			Default(10).
			Comment("最大并发任务数"),

		// 分类和标签
		field.JSON("tags", []string{}).
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL: "json",
			}).
			Comment("标签列表，用于分类和筛选"),

		field.String("region").
			MaxLen(100).
			Optional().
			Comment("所属区域"),

		// 资源限制
		field.JSON("resource_limits", map[string]interface{}{}).
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL: "json",
			}).
			Comment("资源限制配置"),

		// 网络信息
		field.String("local_ip").
			MaxLen(100).
			Optional().
			Comment("本地IP地址"),

		field.String("public_ip").
			MaxLen(100).
			Optional().
			Comment("公网IP地址"),

		field.JSON("network_segments", []string{}).
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL: "json",
			}).
			Comment("管理的网段列表"),

		// 监控指标
		field.Int("active_sessions").
			Default(0).
			Comment("当前活跃会话数"),

		field.Float("cpu_usage").
			Default(0).
			Comment("CPU使用率（百分比）"),

		field.Float("memory_usage").
			Default(0).
			Comment("内存使用率（百分比）"),

		field.Int64("total_requests").
			Default(0).
			Comment("总请求次数"),

		field.Int64("successful_requests").
			Default(0).
			Comment("成功请求次数"),

		field.Int64("failed_requests").
			Default(0).
			Comment("失败请求次数"),

		// Agent版本信息
		field.String("version").
			MaxLen(50).
			Optional().
			Comment("Agent版本号"),

		// 扩展元数据
		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL: "json",
			}).
			Comment("扩展元数据"),

		field.Text("last_error").
			Optional().
			Comment("最后错误信息"),
	}
}

func (Agent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("agent_id"),
		index.Fields("agent_status"),
		index.Fields("last_heartbeat"),
	}
}

func (Agent) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ops_agents"},
	}
}
