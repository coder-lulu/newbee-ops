package rpcclient

import (
	"context"

	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"
	"github.com/zeromicro/go-zero/zrpc"
)

// OpsClient 统一RPC客户端接口（仅声明API侧用到的方法）
type OpsClient interface {
	GetProxyMetricsList(context.Context, *pb.ProxyMetricsListReq) (*pb.ProxyMetricsListResp, error)
	GetAgentGroupList(ctx context.Context, in *pb.AgentGroupListReq) (*pb.AgentGroupListResp, error)
	GetAgentGroupById(ctx context.Context, in *pb.IDReq) (*pb.AgentGroupInfo, error)
	GetProxyGroupList(ctx context.Context, in *pb.ProxyGroupListReq) (*pb.ProxyGroupListResp, error)
	GetProxyGroupById(ctx context.Context, in *pb.IDReq) (*pb.ProxyGroupInfo, error)

	// AccessProfile
	CreateAccessProfile(ctx context.Context, in *pb.AccessProfileInfo) (*pb.BaseIDResp, error)
	UpdateAccessProfile(ctx context.Context, in *pb.AccessProfileInfo) (*pb.BaseResp, error)
	GetAccessProfileList(ctx context.Context, in *pb.AccessProfileListReq) (*pb.AccessProfileListResp, error)
	DeleteAccessProfile(ctx context.Context, in *pb.IDsReq) (*pb.BaseResp, error)

	// Agent
	CreateAgent(ctx context.Context, in *pb.AgentInfo) (*pb.BaseIDResp, error)
	UpdateAgent(ctx context.Context, in *pb.AgentInfo) (*pb.BaseResp, error)
	GetAgentList(ctx context.Context, in *pb.AgentListReq) (*pb.AgentListResp, error)
	GetAgentById(ctx context.Context, in *pb.IDReq) (*pb.AgentInfo, error)
	DeleteAgent(ctx context.Context, in *pb.IDsReq) (*pb.BaseResp, error)
	SelectAgentForDiscovery(ctx context.Context, in *pb.AgentSelectionReq) (*pb.AgentSelectionResp, error)

	// Session
	CreateSession(ctx context.Context, in *pb.SessionInfo) (*pb.BaseIDResp, error)
	UpdateSession(ctx context.Context, in *pb.SessionInfo) (*pb.BaseResp, error)
	GetSessionList(ctx context.Context, in *pb.SessionListReq) (*pb.SessionListResp, error)
	GetSessionBySessionId(ctx context.Context, in *pb.SessionSIDReq) (*pb.SessionInfo, error)

	// Task
	CreateTask(ctx context.Context, in *pb.TaskInfo) (*pb.BaseIDResp, error)
	UpdateTask(ctx context.Context, in *pb.TaskInfo) (*pb.BaseResp, error)
	GetTaskList(ctx context.Context, in *pb.TaskListReq) (*pb.TaskListResp, error)
	GetTaskByTaskId(ctx context.Context, in *pb.TaskIdReq) (*pb.TaskInfo, error)
	ReportTaskResult(ctx context.Context, taskId, status, resultData, errorMessage string, executionTimeMs, endTime int64) (*pb.BaseResp, error)

	// WorkerGroup (deprecated, use Proxy instead)
	CreateWorkerGroup(ctx context.Context, in *pb.WorkerGroupInfo) (*pb.BaseIDResp, error)
	UpdateWorkerGroup(ctx context.Context, in *pb.WorkerGroupInfo) (*pb.BaseResp, error)
	GetWorkerGroupList(ctx context.Context, in *pb.WorkerGroupListReq) (*pb.WorkerGroupListResp, error)
	GetWorkerGroupById(ctx context.Context, in *pb.IDReq) (*pb.WorkerGroupInfo, error)
	GetWorkerGroupMembers(ctx context.Context, in *pb.IDReq) (*pb.WorkerGroupInfoWithMembers, error)
	DeleteWorkerGroup(ctx context.Context, in *pb.IDsReq) (*pb.BaseResp, error)
	AddWorkersToGroup(ctx context.Context, in *pb.WorkerGroupMemberReq) (*pb.BaseResp, error)
	RemoveWorkersFromGroup(ctx context.Context, in *pb.WorkerGroupMemberReq) (*pb.BaseResp, error)
	UpdateWorkerGroupMemberWeight(ctx context.Context, in *pb.WorkerGroupMemberInfo) (*pb.BaseResp, error)
	PickWorkerFromGroup(ctx context.Context, in *pb.PickWorkerFromGroupReq) (*pb.PickWorkerFromGroupResp, error)

	// Proxy
	CreateProxy(ctx context.Context, in *pb.ProxyInfo) (*pb.BaseIDResp, error)
	UpdateProxy(ctx context.Context, in *pb.ProxyInfo) (*pb.BaseResp, error)
	GetProxyList(ctx context.Context, in *pb.ProxyListReq) (*pb.ProxyListResp, error)
	GetProxyById(ctx context.Context, in *pb.IDReq) (*pb.ProxyInfo, error)
	GetProxyByProxyId(ctx context.Context, proxyId string) (*pb.ProxyInfo, error)
	DeleteProxy(ctx context.Context, in *pb.IDsReq) (*pb.BaseResp, error)

	// Script Category
	CreateScriptCategory(ctx context.Context, in *pb.ScriptCategoryInfo) (*pb.BaseIDResp, error)
	UpdateScriptCategory(ctx context.Context, in *pb.ScriptCategoryInfo) (*pb.BaseResp, error)
	GetScriptCategoryList(ctx context.Context, in *pb.ScriptCategoryListReq) (*pb.ScriptCategoryListResp, error)
	GetScriptCategoryById(ctx context.Context, in *pb.IDReq) (*pb.ScriptCategoryInfo, error)
	DeleteScriptCategory(ctx context.Context, in *pb.IDsReq) (*pb.BaseResp, error)

	// Script
	CreateScript(ctx context.Context, in *pb.ScriptInfo) (*pb.BaseIDResp, error)
	UpdateScript(ctx context.Context, in *pb.ScriptInfo) (*pb.BaseResp, error)
	GetScriptList(ctx context.Context, in *pb.ScriptListReq) (*pb.ScriptListResp, error)
	GetScriptById(ctx context.Context, in *pb.IDReq) (*pb.ScriptInfo, error)
	DeleteScript(ctx context.Context, in *pb.IDsReq) (*pb.BaseResp, error)

	// Script Version
	GetScriptVersionList(ctx context.Context, page, pageSize, scriptId uint64) (*pb.ScriptVersionListResp, error)
	GetScriptVersionById(ctx context.Context, id uint64) (*pb.ScriptVersionInfo, error)
}

// GrpcOpsClient 直接使用 protoc 生成的 gRPC 客户端（通过 zrpc 建连）
type GrpcOpsClient struct{ cli pb.OpsClient }

func NewGrpcOpsClient(c zrpc.Client) *GrpcOpsClient {
	return &GrpcOpsClient{cli: pb.NewOpsClient(c.Conn())}
}

// proxy methods
func (g *GrpcOpsClient) CreateAccessProfile(ctx context.Context, in *pb.AccessProfileInfo) (*pb.BaseIDResp, error) {
	return g.cli.CreateAccessProfile(ctx, in)
}
func (g *GrpcOpsClient) UpdateAccessProfile(ctx context.Context, in *pb.AccessProfileInfo) (*pb.BaseResp, error) {
	return g.cli.UpdateAccessProfile(ctx, in)
}
func (g *GrpcOpsClient) GetAccessProfileList(ctx context.Context, in *pb.AccessProfileListReq) (*pb.AccessProfileListResp, error) {
	return g.cli.GetAccessProfileList(ctx, in)
}
func (g *GrpcOpsClient) DeleteAccessProfile(ctx context.Context, in *pb.IDsReq) (*pb.BaseResp, error) {
	return g.cli.DeleteAccessProfile(ctx, in)
}

// Agent methods
func (g *GrpcOpsClient) CreateAgent(ctx context.Context, in *pb.AgentInfo) (*pb.BaseIDResp, error) {
	return g.cli.CreateAgent(ctx, in)
}
func (g *GrpcOpsClient) UpdateAgent(ctx context.Context, in *pb.AgentInfo) (*pb.BaseResp, error) {
	return g.cli.UpdateAgent(ctx, in)
}
func (g *GrpcOpsClient) GetAgentList(ctx context.Context, in *pb.AgentListReq) (*pb.AgentListResp, error) {
	return g.cli.GetAgentList(ctx, in)
}
func (g *GrpcOpsClient) GetAgentById(ctx context.Context, in *pb.IDReq) (*pb.AgentInfo, error) {
	return g.cli.GetAgentById(ctx, in)
}
func (g *GrpcOpsClient) DeleteAgent(ctx context.Context, in *pb.IDsReq) (*pb.BaseResp, error) {
	return g.cli.DeleteAgent(ctx, in)
}
func (g *GrpcOpsClient) SelectAgentForDiscovery(ctx context.Context, in *pb.AgentSelectionReq) (*pb.AgentSelectionResp, error) {
	return g.cli.SelectAgentForDiscovery(ctx, in)
}

func (g *GrpcOpsClient) CreateSession(ctx context.Context, in *pb.SessionInfo) (*pb.BaseIDResp, error) {
	return g.cli.CreateSession(ctx, in)
}
func (g *GrpcOpsClient) UpdateSession(ctx context.Context, in *pb.SessionInfo) (*pb.BaseResp, error) {
	return g.cli.UpdateSession(ctx, in)
}
func (g *GrpcOpsClient) GetSessionList(ctx context.Context, in *pb.SessionListReq) (*pb.SessionListResp, error) {
	return g.cli.GetSessionList(ctx, in)
}
func (g *GrpcOpsClient) GetSessionBySessionId(ctx context.Context, in *pb.SessionSIDReq) (*pb.SessionInfo, error) {
	return g.cli.GetSessionBySessionId(ctx, in)
}

func (g *GrpcOpsClient) CreateTask(ctx context.Context, in *pb.TaskInfo) (*pb.BaseIDResp, error) {
	return g.cli.CreateTask(ctx, in)
}
func (g *GrpcOpsClient) UpdateTask(ctx context.Context, in *pb.TaskInfo) (*pb.BaseResp, error) {
	return g.cli.UpdateTask(ctx, in)
}
func (g *GrpcOpsClient) GetTaskList(ctx context.Context, in *pb.TaskListReq) (*pb.TaskListResp, error) {
	return g.cli.GetTaskList(ctx, in)
}
func (g *GrpcOpsClient) GetTaskByTaskId(ctx context.Context, in *pb.TaskIdReq) (*pb.TaskInfo, error) {
	return g.cli.GetTaskByTaskId(ctx, in)
}

// ReportTaskResult reports task execution result from Proxy (Phase 2)
func (g *GrpcOpsClient) ReportTaskResult(ctx context.Context, taskId, status, resultData, errorMessage string, executionTimeMs, endTime int64) (*pb.BaseResp, error) {
	req := &pb.TaskResultReq{
		TaskId:          taskId,
		Status:          status,
		ResultData:      &resultData,
		ErrorMessage:    &errorMessage,
		ExecutionTimeMs: &executionTimeMs,
		EndTime:         &endTime,
	}
	return g.cli.ReportTaskResult(ctx, req)
}

// WorkerGroup methods
func (g *GrpcOpsClient) CreateWorkerGroup(ctx context.Context, in *pb.WorkerGroupInfo) (*pb.BaseIDResp, error) {
	return g.cli.CreateWorkerGroup(ctx, in)
}
func (g *GrpcOpsClient) UpdateWorkerGroup(ctx context.Context, in *pb.WorkerGroupInfo) (*pb.BaseResp, error) {
	return g.cli.UpdateWorkerGroup(ctx, in)
}
func (g *GrpcOpsClient) GetWorkerGroupList(ctx context.Context, in *pb.WorkerGroupListReq) (*pb.WorkerGroupListResp, error) {
	return g.cli.GetWorkerGroupList(ctx, in)
}
func (g *GrpcOpsClient) GetWorkerGroupById(ctx context.Context, in *pb.IDReq) (*pb.WorkerGroupInfo, error) {
	return g.cli.GetWorkerGroupById(ctx, in)
}
func (g *GrpcOpsClient) GetWorkerGroupMembers(ctx context.Context, in *pb.IDReq) (*pb.WorkerGroupInfoWithMembers, error) {
	return g.cli.GetWorkerGroupMembers(ctx, in)
}
func (g *GrpcOpsClient) DeleteWorkerGroup(ctx context.Context, in *pb.IDsReq) (*pb.BaseResp, error) {
	return g.cli.DeleteWorkerGroup(ctx, in)
}
func (g *GrpcOpsClient) AddWorkersToGroup(ctx context.Context, in *pb.WorkerGroupMemberReq) (*pb.BaseResp, error) {
	return g.cli.AddWorkersToGroup(ctx, in)
}
func (g *GrpcOpsClient) RemoveWorkersFromGroup(ctx context.Context, in *pb.WorkerGroupMemberReq) (*pb.BaseResp, error) {
	return g.cli.RemoveWorkersFromGroup(ctx, in)
}
func (g *GrpcOpsClient) UpdateWorkerGroupMemberWeight(ctx context.Context, in *pb.WorkerGroupMemberInfo) (*pb.BaseResp, error) {
	return g.cli.UpdateWorkerGroupMemberWeight(ctx, in)
}
func (g *GrpcOpsClient) PickWorkerFromGroup(ctx context.Context, in *pb.PickWorkerFromGroupReq) (*pb.PickWorkerFromGroupResp, error) {
	return g.cli.PickWorkerFromGroup(ctx, in)
}

// Proxy methods
func (g *GrpcOpsClient) CreateProxy(ctx context.Context, in *pb.ProxyInfo) (*pb.BaseIDResp, error) {
	return g.cli.CreateProxy(ctx, in)
}

func (g *GrpcOpsClient) UpdateProxy(ctx context.Context, in *pb.ProxyInfo) (*pb.BaseResp, error) {
	return g.cli.UpdateProxy(ctx, in)
}

func (g *GrpcOpsClient) GetProxyList(ctx context.Context, in *pb.ProxyListReq) (*pb.ProxyListResp, error) {
	return g.cli.GetProxyList(ctx, in)
}

func (g *GrpcOpsClient) GetProxyById(ctx context.Context, in *pb.IDReq) (*pb.ProxyInfo, error) {
	return g.cli.GetProxyById(ctx, in)
}

func (g *GrpcOpsClient) GetProxyByProxyId(ctx context.Context, proxyId string) (*pb.ProxyInfo, error) {
	// Query proxy list filtered by proxy_id
	resp, err := g.cli.GetProxyList(ctx, &pb.ProxyListReq{
		Page:     1,
		PageSize: 1,
		ProxyId:  &proxyId,
	})
	if err != nil {
		return nil, err
	}
	if resp.Total == 0 || len(resp.Data) == 0 {
		return nil, nil
	}
	return resp.Data[0], nil
}

func (g *GrpcOpsClient) DeleteProxy(ctx context.Context, in *pb.IDsReq) (*pb.BaseResp, error) {
	return g.cli.DeleteProxy(ctx, in)
}

// Script Category methods
func (g *GrpcOpsClient) CreateScriptCategory(ctx context.Context, in *pb.ScriptCategoryInfo) (*pb.BaseIDResp, error) {
	return g.cli.CreateScriptCategory(ctx, in)
}

func (g *GrpcOpsClient) UpdateScriptCategory(ctx context.Context, in *pb.ScriptCategoryInfo) (*pb.BaseResp, error) {
	return g.cli.UpdateScriptCategory(ctx, in)
}

func (g *GrpcOpsClient) GetScriptCategoryList(ctx context.Context, in *pb.ScriptCategoryListReq) (*pb.ScriptCategoryListResp, error) {
	return g.cli.GetScriptCategoryList(ctx, in)
}

func (g *GrpcOpsClient) GetScriptCategoryById(ctx context.Context, in *pb.IDReq) (*pb.ScriptCategoryInfo, error) {
	return g.cli.GetScriptCategoryById(ctx, in)
}

func (g *GrpcOpsClient) DeleteScriptCategory(ctx context.Context, in *pb.IDsReq) (*pb.BaseResp, error) {
	return g.cli.DeleteScriptCategory(ctx, in)
}

// Script methods
func (g *GrpcOpsClient) CreateScript(ctx context.Context, in *pb.ScriptInfo) (*pb.BaseIDResp, error) {
	return g.cli.CreateScript(ctx, in)
}

func (g *GrpcOpsClient) UpdateScript(ctx context.Context, in *pb.ScriptInfo) (*pb.BaseResp, error) {
	return g.cli.UpdateScript(ctx, in)
}

func (g *GrpcOpsClient) GetScriptList(ctx context.Context, in *pb.ScriptListReq) (*pb.ScriptListResp, error) {
	return g.cli.GetScriptList(ctx, in)
}

func (g *GrpcOpsClient) GetScriptById(ctx context.Context, in *pb.IDReq) (*pb.ScriptInfo, error) {
	return g.cli.GetScriptById(ctx, in)
}

func (g *GrpcOpsClient) DeleteScript(ctx context.Context, in *pb.IDsReq) (*pb.BaseResp, error) {
	return g.cli.DeleteScript(ctx, in)
}

// Script Version methods
func (g *GrpcOpsClient) GetScriptVersionList(ctx context.Context, page, pageSize, scriptId uint64) (*pb.ScriptVersionListResp, error) {
	return g.cli.GetScriptVersionList(ctx, &pb.ScriptVersionListReq{
		Page:     page,
		PageSize: pageSize,
		ScriptId: &scriptId,
	})
}

func (g *GrpcOpsClient) GetScriptVersionById(ctx context.Context, id uint64) (*pb.ScriptVersionInfo, error) {
	return g.cli.GetScriptVersionById(ctx, &pb.IDReq{Id: id})
}

// pointy helper
func pointy[T any](v T) *T { return &v }

func (g *GrpcOpsClient) GetAgentGroupList(ctx context.Context, in *pb.AgentGroupListReq) (*pb.AgentGroupListResp, error) {
	return g.cli.GetAgentGroupList(ctx, in)
}

func (g *GrpcOpsClient) GetAgentGroupById(ctx context.Context, in *pb.IDReq) (*pb.AgentGroupInfo, error) {
	return g.cli.GetAgentGroupById(ctx, in)
}

func (g *GrpcOpsClient) GetProxyGroupList(ctx context.Context, in *pb.ProxyGroupListReq) (*pb.ProxyGroupListResp, error) {
	return g.cli.GetProxyGroupList(ctx, in)
}

func (g *GrpcOpsClient) GetProxyGroupById(ctx context.Context, in *pb.IDReq) (*pb.ProxyGroupInfo, error) {
	return g.cli.GetProxyGroupById(ctx, in)
}

func (g *GrpcOpsClient) GetProxyMetricsList(ctx context.Context, in *pb.ProxyMetricsListReq) (*pb.ProxyMetricsListResp, error) {
	return g.cli.GetProxyMetricsList(ctx, in)
}
