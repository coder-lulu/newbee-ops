package types

// GetTaskStatusReq 获取任务状态请求
// 用于 GET /ops/task/status/:taskId
type GetTaskStatusReq struct {
	TaskId string `path:"taskId"` // 任务ID (path parameter)
}

// GetTaskResultReq 获取任务结果请求
// 用于 GET /ops/task/result/:taskId
type GetTaskResultReq struct {
	TaskId string `path:"taskId"` // 任务ID (path parameter)
}
