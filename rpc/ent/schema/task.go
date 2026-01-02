package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema/field"
    "github.com/coder-lulu/newbee-common/v2/orm/ent/mixins"
)

// Task 持久化任务模型
type Task struct{ ent.Schema }

func (Task) Mixin() []ent.Mixin {
    return []ent.Mixin{
        mixins.IDMixin{},
        mixins.StatusMixin{},
        mixins.TenantMixin{},
    }
}

func (Task) Fields() []ent.Field {
    return []ent.Field{
        field.String("task_id").Unique().Annotations(entsql.Annotation{Size: 128}).Comment("任务ID"),
        field.String("name").Optional().Annotations(entsql.Annotation{Size: 128}).Comment("任务名称"),
        field.String("creator_id").Optional().Annotations(entsql.Annotation{Size: 64}).Comment("创建人ID"),

        // 目标与执行器
        field.JSON("ci_ids", []string{}).Optional().Comment("目标CI列表"),
        field.String("executor").Default("agent").Annotations(entsql.Annotation{Size: 16}).Comment("执行器: agent|ssh|telnet"),
        field.String("worker_id").Optional().Annotations(entsql.Annotation{Size: 100}).Comment("绑定的Worker ID"),

        // 命令参数
        field.String("command_content").Optional().Annotations(entsql.Annotation{Size: 1024}).Comment("命令/脚本内容"),
        field.Int("command_timeout").Optional().Comment("超时(秒)"),
        field.String("command_workdir").Optional().Annotations(entsql.Annotation{Size: 256}).Comment("工作目录"),
        field.JSON("command_env", map[string]string{}).Optional().Comment("环境变量"),

        // 状态与结果
        field.String("status_str").Default("pending").Annotations(entsql.Annotation{Size: 16}).Comment("pending|running|success|failed"),
        field.Text("result_output").Optional().Comment("输出结果(截断/采样)"),
        field.Text("result_data").Optional().Comment("完整结果数据(JSON)"),
        field.Text("error_msg").Optional().Comment("错误信息"),
        field.Int("exit_code").Optional().Comment("退出码"),
        field.String("remote_task_id").Optional().Annotations(entsql.Annotation{Size: 128}).Comment("远端任务ID(nb-agent)"),
        field.Int64("start_time").Optional().Comment("开始时间(Unix秒)"),
        field.Int64("end_time").Optional().Comment("结束时间(Unix秒)"),
        field.Time("dispatched_at").Optional().Comment("任务分配时间"),
        field.Time("completed_at").Optional().Comment("任务完成时间"),
        field.Int64("execution_time_ms").Optional().Comment("执行耗时(毫秒)"),

        field.JSON("tags", map[string]string{}).Optional().Comment("标签/扩展"),
    }
}
