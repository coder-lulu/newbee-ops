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

type ScriptExecution struct {
	ent.Schema
}

func (ScriptExecution) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TenantMixin{},
		mixins.StatusMixin{},
	}
}

func (ScriptExecution) Fields() []ent.Field {
	return []ent.Field{
		// 关联信息
		field.Uint64("script_id").
			Comment("脚本ID").
			Positive(),

		field.String("execution_id").
			Comment("执行ID（UUID）").
			NotEmpty().
			MaxLen(64).
			Unique(),

		field.Uint64("executor_user_id").
			Comment("执行用户ID").
			Optional(),

		field.String("trigger_type").
			Comment("触发类型：manual/schedule/event/api").
			Default("manual").
			MaxLen(20),

		// 目标信息
		field.JSON("target_ci_ids", []string{}).
			Comment("目标CI ID列表").
			Optional(),

		field.Uint32("target_count").
			Comment("目标数量").
			Default(0),

		// 参数信息
		field.Text("parameters").
			Comment("渲染后的参数（JSON）").
			Optional(),

		field.Text("raw_parameters").
			Comment("用户输入的原始参数（JSON）").
			Optional(),

		// 执行状态
		field.String("execution_status").
			Comment("执行状态：pending/running/success/partial_success/failed/timeout/cancelled").
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

		// 结果统计
		field.Uint32("success_count").
			Comment("成功数量").
			Default(0),

		field.Uint32("failure_count").
			Comment("失败数量").
			Default(0),

		field.Uint32("timeout_count").
			Comment("超时数量").
			Default(0),

		// 关联Task系统
		field.String("task_id").
			Comment("关联的Task ID").
			Optional().
			MaxLen(64),

		// 其他信息
		field.Text("error_message").
			Comment("错误信息").
			Optional(),

		field.Text("execution_log").
			Comment("执行日志摘要").
			Optional(),
	}
}

func (ScriptExecution) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("details", ScriptExecutionDetail.Type).
			Comment("执行详情"),

		edge.From("script", Script.Type).
			Ref("executions").
			Field("script_id").
			Unique().
			Required(),
	}
}

func (ScriptExecution) Indexes() []ent.Index {
	return []ent.Index{
		// 租户ID + 执行ID 唯一索引
		index.Fields("tenant_id", "execution_id").
			Unique(),

		// 脚本ID + 开始时间索引（查询历史）
		index.Fields("script_id", "started_at"),

		// 执行状态索引（查询待处理/运行中的任务）
		index.Fields("execution_status"),

		// 租户ID + 执行状态索引
		index.Fields("tenant_id", "execution_status"),
	}
}

func (ScriptExecution) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "ops_script_executions",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
		},
		schema.Comment("脚本执行记录表"),
	}
}
