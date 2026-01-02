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

type Script struct {
	ent.Schema
}

func (Script) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TenantMixin{},
		mixins.DepartmentMixin{},
		mixins.StatusMixin{},
	}
}

func (Script) Fields() []ent.Field {
	return []ent.Field{
		// 基础信息
		field.String("name").
			Comment("脚本名称").
			NotEmpty().
			MaxLen(100),

		field.String("code").
			Comment("脚本唯一标识代码").
			NotEmpty().
			MaxLen(100),

		field.String("description").
			Comment("脚本描述").
			Optional().
			MaxLen(500),

		field.Uint64("category_id").
			Comment("分类ID").
			Optional(),

		field.JSON("tags", []string{}).
			Comment("标签列表").
			Optional(),

		// 脚本内容
		field.String("script_type").
			Comment("脚本类型：shell, bash, python, perl, sql, javascript, http, snmp").
			NotEmpty().
			MaxLen(50),

		field.String("executor").
			Comment("执行器类型：agent, ssh, telnet, rdp, http, snmp").
			NotEmpty().
			MaxLen(50),

		field.Text("content").
			Comment("脚本内容"),

		// 模板支持
		field.Bool("is_template").
			Comment("是否为模板脚本").
			Default(false),

		field.String("template_engine").
			Comment("模板引擎：go-template, jinja2").
			Optional().
			MaxLen(50),

		field.JSON("parameters", map[string]interface{}{}).
			Comment("参数定义（JSON Schema格式）").
			Optional(),

		// 版本管理
		field.String("version").
			Comment("版本号").
			Default("1.0.0").
			MaxLen(50),

		field.Uint64("base_version_id").
			Comment("基础版本ID（用于版本链）").
			Optional(),

		field.Bool("is_latest").
			Comment("是否为最新版本").
			Default(true),

		// CMDB集成
		field.JSON("target_ci_types", []string{}).
			Comment("适用的CI类型列表：server, vm, container, network_device").
			Optional(),

		field.JSON("target_os_types", []string{}).
			Comment("适用的OS类型列表：linux, windows, aix, solaris").
			Optional(),

		field.JSON("target_selector", map[string]interface{}{}).
			Comment("CMDB目标选择器配置（JSON格式）").
			Optional(),

		// 凭证配置
		field.String("credential_ref").
			Comment("凭证引用（指向AccessProfile或凭证库）").
			Optional().
			MaxLen(200),

		field.JSON("required_capabilities", []string{}).
			Comment("执行所需的能力列表：sudo, docker, kubectl").
			Optional(),

		// 执行配置
		field.Uint32("default_timeout").
			Comment("默认超时时间（秒）").
			Default(300),

		field.String("default_workdir").
			Comment("默认工作目录").
			Optional().
			MaxLen(500),

		field.JSON("default_env", map[string]string{}).
			Comment("默认环境变量").
			Optional(),

		// 安全控制
		field.Bool("require_confirmation").
			Comment("执行前是否需要确认").
			Default(false),

		field.String("risk_level").
			Comment("风险级别：low, medium, high, critical").
			Default("low").
			MaxLen(20),

		// 调度配置
		field.Bool("schedulable").
			Comment("是否支持调度").
			Default(true),

		field.String("default_schedule").
			Comment("默认调度表达式（Cron格式）").
			Optional().
			MaxLen(100),

		// 统计信息
		field.Uint64("execution_count").
			Comment("总执行次数").
			Default(0),

		field.Uint64("success_count").
			Comment("成功执行次数").
			Default(0),

		field.Uint64("failure_count").
			Comment("失败执行次数").
			Default(0),

		field.Int64("last_executed_at").
			Comment("最后执行时间（Unix时间戳）").
			Optional(),
	}
}

func (Script) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", ScriptCategory.Type).
			Ref("scripts").
			Field("category_id").
			Unique(),

		edge.To("versions", ScriptVersion.Type),

		// Phase 3: 执行记录
		edge.To("executions", ScriptExecution.Type),

		// Note: Edge to ScriptSchedule will be added in Phase 5
	}
}

func (Script) Indexes() []ent.Index {
	return []ent.Index{
		// 租户内code唯一
		index.Fields("tenant_id", "code").
			Unique(),

		// 分类查询
		index.Fields("tenant_id", "category_id"),

		// 类型查询
		index.Fields("script_type"),

		// 最新版本查询
		index.Fields("is_latest"),
	}
}

func (Script) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "ops_scripts",
		},
	}
}
