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

type ScriptExecutionDetail struct {
	ent.Schema
}

func (ScriptExecutionDetail) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TenantMixin{},
		mixins.StatusMixin{},
	}
}

func (ScriptExecutionDetail) Fields() []ent.Field {
	return []ent.Field{
		// 关联信息
		field.Uint64("execution_id").
			Comment("执行记录ID").
			Positive(),

		field.String("ci_id").
			Comment("CI ID").
			NotEmpty().
			MaxLen(64),

		field.String("ci_name").
			Comment("CI名称").
			Optional().
			MaxLen(200),

		field.String("ci_ip").
			Comment("CI IP地址").
			Optional().
			MaxLen(50),

		// 执行状态
		field.String("exec_status").
			Comment("执行状态：pending/running/success/failed/timeout/cancelled").
			Default("pending").
			MaxLen(20),

		// 时间统计
		field.Time("started_at").
			Comment("开始时间").
			Optional().
			Nillable(),

		field.Time("completed_at").
			Comment("完成时间").
			Optional().
			Nillable(),

		field.Uint64("duration_ms").
			Comment("执行时长（毫秒）").
			Default(0),

		// 执行结果
		field.Text("output").
			Comment("标准输出（截断至10KB）").
			Optional(),

		field.Text("error_output").
			Comment("错误输出").
			Optional(),

		field.Int("exit_code").
			Comment("退出代码").
			Optional().
			Nillable(),

		field.String("full_output_path").
			Comment("完整输出文件路径（S3/OSS）").
			Optional().
			MaxLen(500),

		// 关联信息
		field.String("task_id").
			Comment("关联的Task ID").
			Optional().
			MaxLen(64),

		field.String("worker_id").
			Comment("执行的Worker ID").
			Optional().
			MaxLen(64),

		field.String("proxy_id").
			Comment("使用的Proxy ID").
			Optional().
			MaxLen(64),

		// 其他信息
		field.Text("error_message").
			Comment("错误信息").
			Optional(),

		field.Uint32("retry_count").
			Comment("重试次数").
			Default(0),
	}
}

func (ScriptExecutionDetail) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("execution", ScriptExecution.Type).
			Ref("details").
			Field("execution_id").
			Unique().
			Required(),
	}
}

func (ScriptExecutionDetail) Indexes() []ent.Index {
	return []ent.Index{
		// 执行记录ID索引（查询某次执行的所有详情）
		index.Fields("execution_id"),

		// 租户ID + 执行记录ID索引
		index.Fields("tenant_id", "execution_id"),

		// CI ID索引（查询某个CI的执行历史）
		index.Fields("ci_id"),

		// 状态索引（查询失败/超时的任务）
		index.Fields("exec_status"),

		// Task ID索引（关联查询）
		index.Fields("task_id"),
	}
}

func (ScriptExecutionDetail) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "ops_script_execution_details",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
		},
		schema.Comment("脚本执行详情表"),
	}
}
