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

// Proxy Proxy运维操作执行节点
type Proxy struct {
	ent.Schema
}

func (Proxy) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TenantMixin{},
		mixins.StatusMixin{},
	}
}

func (Proxy) Fields() []ent.Field {
	return []ent.Field{
		// ========== 基础信息 ==========
		field.String("worker_id").
			MaxLen(100).
			NotEmpty().
			Unique().
			Comment("Proxy唯一标识（由Proxy生成，格式：proxy-{uuid}）"),

		field.String("name").
			MaxLen(100).
			NotEmpty().
			Comment("Proxy友好名称"),

		field.String("ip").
			MaxLen(100).
			NotEmpty().
			Comment("Proxy IP地址"),

		field.Int("port").
			Default(8889).
			Comment("Proxy服务端口"),

		field.String("version").
			MaxLen(50).
			Optional().
			Comment("Proxy版本号"),

		// ========== 地理位置 ==========
		field.String("region").
			MaxLen(100).
			Default("").
			Comment("地理区域（如：cn-beijing, us-west-1）"),

		field.String("zone").
			MaxLen(100).
			Default("").
			Comment("可用区（如：az-1, az-2）"),

		// ========== 能力与标签 ==========
		field.JSON("capabilities", []string{}).
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL: "json",
			}).
			Comment("支持的协议能力 [ssh, telnet, rdp, ipmi, snmp]"),

		field.JSON("tags", []string{}).
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL: "json",
			}).
			Comment("标签列表，用于分类和筛选"),

		field.JSON("endpoints", map[string]string{}).
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL: "json",
			}).
			Comment("服务端点映射 {\"http\":\"http://...\", \"ws\":\"ws://...\"}"),

		// ========== 状态信息 ==========
		field.Enum("worker_status").
			Values("online", "degraded", "offline").
			Default("offline").
			Comment("Proxy运行状态"),

		field.Time("last_heartbeat").
			Optional().
			Nillable().
			Comment("最后心跳时间"),

		field.Time("register_time").
			Optional().
			Nillable().
			Comment("首次注册时间"),

		// ========== 负载均衡配置 ==========
		field.Int("weight").
			Default(100).
			Range(1, 1000).
			Comment("负载均衡权重（1-1000），越大分配越多"),

		field.Int("priority").
			Default(0).
			Comment("优先级，数字越大优先级越高"),

		// ========== 资源指标（实时） ==========
		field.Float("cpu_usage").
			Default(0).
			Comment("CPU使用率（0-100）"),

		field.Float("memory_usage").
			Default(0).
			Comment("内存使用率（0-100）"),

		field.Float("disk_usage").
			Default(0).
			Comment("磁盘使用率（0-100）"),

		field.Uint64("network_in").
			Default(0).
			Comment("网络入流量（字节）"),

		field.Uint64("network_out").
			Default(0).
			Comment("网络出流量（字节）"),

		// ========== 业务指标（累计） ==========
		field.Int("active_sessions").
			Default(0).
			Comment("当前活跃会话数"),

		field.Uint64("total_requests").
			Default(0).
			Comment("总请求次数"),

		field.Uint64("success_count").
			Default(0).
			Comment("成功请求次数"),

		field.Uint64("failure_count").
			Default(0).
			Comment("失败请求次数"),

		field.Int("max_sessions").
			Default(1000).
			Comment("最大并发会话数"),

		// ========== 健康检查 ==========
		field.Int("health_check_failures").
			Default(0).
			Comment("连续健康检查失败次数"),

		field.Time("last_health_check").
			Optional().
			Nillable().
			Comment("最后健康检查时间"),

		field.String("health_check_url").
			MaxLen(500).
			Default("").
			Comment("健康检查URL（为空则不主动探测）"),

		// ========== 网络信息 ==========
		field.String("local_ip").
			MaxLen(100).
			Default("").
			Comment("本地IP地址"),

		field.String("public_ip").
			MaxLen(100).
			Default("").
			Comment("公网IP地址"),

		field.JSON("network_segments", []string{}).
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL: "json",
			}).
			Comment("管理的网段列表 [\"192.168.1.0/24\", \"10.0.0.0/8\"]"),

		// ========== 扩展元数据 ==========
		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL: "json",
			}).
			Comment("扩展元数据（用户自定义字段）"),

		field.Text("last_error").
			Optional().
			Comment("最后错误信息"),
	}
}

func (Proxy) Indexes() []ent.Index {
	return []ent.Index{
		// 租户+状态查询（最常用）
		index.Fields("tenant_id", "worker_status", "status"),

		// 租户+区域+状态查询（地理位置选择）
		index.Fields("tenant_id", "region", "worker_status"),

		// Proxy ID查找
		index.Fields("worker_id"),

		// 心跳过期清理
		index.Fields("last_heartbeat"),

		// 租户查询
		index.Fields("tenant_id"),
	}
}

func (Proxy) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ops_workers"},
	}
}
