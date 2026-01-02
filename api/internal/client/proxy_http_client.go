package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// ProxyHTTPClient Proxy HTTP客户端（ops-center调用proxy的HTTP接口）
type ProxyHTTPClient struct {
	httpClient *http.Client
	logger     logx.Logger
}

// NewProxyHTTPClient 创建Proxy HTTP客户端
func NewProxyHTTPClient() *ProxyHTTPClient {
	return &ProxyHTTPClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logx.WithContext(context.Background()),
	}
}

// =========================
// 任务执行接口
// =========================

// CommandExecuteRequest 命令执行请求
type CommandExecuteRequest struct {
	Target       string            `json:"target"`
	Port         int32             `json:"port"`
	Protocol     string            `json:"protocol"`      // ssh/telnet
	Username     string            `json:"username"`
	Password     string            `json:"password"`
	PrivateKey   string            `json:"private_key"`
	Command      string            `json:"command"`
	WorkingDir   string            `json:"working_dir"`
	Timeout      int32             `json:"timeout"`
	UseSudo      bool              `json:"use_sudo"`
	SudoPassword string            `json:"sudo_password"`
	Environment  map[string]string `json:"environment"`
}

// ScriptExecuteRequest 脚本执行请求
type ScriptExecuteRequest struct {
	Target         string            `json:"target"`
	Port           int32             `json:"port"`
	Protocol       string            `json:"protocol"`
	Username       string            `json:"username"`
	Password       string            `json:"password"`
	PrivateKey     string            `json:"private_key"`
	ScriptContent  string            `json:"script_content"`
	ScriptType     string            `json:"script_type"`
	WorkingDir     string            `json:"working_dir"`
	Timeout        int32             `json:"timeout"`
	UseSudo        bool              `json:"use_sudo"`
	SudoPassword   string            `json:"sudo_password"`
	Environment    map[string]string `json:"environment"`
	RemoteFileMode string            `json:"remote_file_mode"`
	RemoteFilePath string            `json:"remote_file_path"`
	CleanupAfter   bool              `json:"cleanup_after"`
}

// FileTransferRequest 文件传输请求
type FileTransferRequest struct {
	Target     string `json:"target"`
	Port       int32  `json:"port"`
	Protocol   string `json:"protocol"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
	SourcePath string `json:"source_path"`
	TargetPath string `json:"target_path"`
	Direction  string `json:"direction"` // upload/download
	Timeout    int32  `json:"timeout"`
}

// TaskResponse 任务响应
type TaskResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	TaskID  string `json:"task_id"`
	Error   string `json:"error,omitempty"`
}

// TaskStatusResponse 任务状态响应
type TaskStatusResponse struct {
	TaskID     string `json:"task_id"`
	TaskType   string `json:"task_type"`
	Status     string `json:"status"`
	StartTime  string `json:"start_time"`
	UpdateTime string `json:"update_time"`
	Duration   int64  `json:"duration"`
}

// TaskResultResponse 任务结果响应
type TaskResultResponse struct {
	TaskID         string                 `json:"task_id"`
	TaskType       string                 `json:"task_type"`
	Status         string                 `json:"status"`
	StartTime      string                 `json:"start_time"`
	UpdateTime     string                 `json:"update_time"`
	Duration       int64                  `json:"duration"`
	ResultStatus   string                 `json:"result_status"`
	ErrorMessage   string                 `json:"error_message"`
	HasResult      bool                   `json:"has_result"`
	DetailedResult map[string]interface{} `json:"detailed_result"`
	ParsedResult   map[string]interface{} `json:"parsed_result"`
}

// ExecuteCommand 执行命令
func (c *ProxyHTTPClient) ExecuteCommand(ctx context.Context, proxyEndpoint string, req *CommandExecuteRequest) (*TaskResponse, error) {
	url := fmt.Sprintf("%s/api/task/command", proxyEndpoint)
	return c.postJSON(ctx, url, req)
}

// ExecuteScript 执行脚本
func (c *ProxyHTTPClient) ExecuteScript(ctx context.Context, proxyEndpoint string, req *ScriptExecuteRequest) (*TaskResponse, error) {
	url := fmt.Sprintf("%s/api/task/script", proxyEndpoint)
	return c.postJSON(ctx, url, req)
}

// ExecuteFileTransfer 执行文件传输
func (c *ProxyHTTPClient) ExecuteFileTransfer(ctx context.Context, proxyEndpoint string, req *FileTransferRequest) (*TaskResponse, error) {
	url := fmt.Sprintf("%s/api/task/file", proxyEndpoint)
	return c.postJSON(ctx, url, req)
}

// GetTaskStatus 获取任务状态
func (c *ProxyHTTPClient) GetTaskStatus(ctx context.Context, proxyEndpoint string, taskID string) (*TaskStatusResponse, error) {
	url := fmt.Sprintf("%s/api/task/status/%s", proxyEndpoint, taskID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("proxy returned status %d: %s", resp.StatusCode, string(body))
	}

	var result TaskStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &result, nil
}

// GetTaskResult 获取任务结果
func (c *ProxyHTTPClient) GetTaskResult(ctx context.Context, proxyEndpoint string, taskID string) (*TaskResultResponse, error) {
	url := fmt.Sprintf("%s/api/task/result/%s", proxyEndpoint, taskID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("proxy returned status %d: %s", resp.StatusCode, string(body))
	}

	var result TaskResultResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &result, nil
}

// CancelTask 取消任务
func (c *ProxyHTTPClient) CancelTask(ctx context.Context, proxyEndpoint string, taskID string) error {
	url := fmt.Sprintf("%s/api/task/cancel", proxyEndpoint)

	reqBody := map[string]string{"task_id": taskID}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("proxy returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// =========================
// WebSocket会话URL生成
// =========================

// GenerateSSHWebSocketURL 生成SSH WebSocket URL
func (c *ProxyHTTPClient) GenerateSSHWebSocketURL(proxyEndpoint string, target string, port int, username string, password string) string {
	// 将 http:// 替换为 ws://
	wsEndpoint := replaceHTTPWithWS(proxyEndpoint)
	return fmt.Sprintf("%s/ws/ssh?target=%s&port=%d&username=%s&password=%s",
		wsEndpoint, target, port, username, password)
}

// GenerateTelnetWebSocketURL 生成Telnet WebSocket URL
func (c *ProxyHTTPClient) GenerateTelnetWebSocketURL(proxyEndpoint string, target string, port int, username string, password string) string {
	wsEndpoint := replaceHTTPWithWS(proxyEndpoint)
	return fmt.Sprintf("%s/ws/telnet?target=%s&port=%d&username=%s&password=%s",
		wsEndpoint, target, port, username, password)
}

// GenerateRDPWebSocketURL 生成RDP WebSocket URL (Guacamole协议)
func (c *ProxyHTTPClient) GenerateRDPWebSocketURL(proxyEndpoint string, target string, port int, username string, password string) string {
	wsEndpoint := replaceHTTPWithWS(proxyEndpoint)
	return fmt.Sprintf("%s/api/rdp/websocket?target=%s&port=%d&username=%s&password=%s",
		wsEndpoint, target, port, username, password)
}

// GenerateVNCWebSocketURL 生成VNC WebSocket URL (Guacamole协议)
func (c *ProxyHTTPClient) GenerateVNCWebSocketURL(proxyEndpoint string, target string, port int, password string) string {
	wsEndpoint := replaceHTTPWithWS(proxyEndpoint)
	return fmt.Sprintf("%s/api/vnc/websocket?target=%s&port=%d&password=%s",
		wsEndpoint, target, port, password)
}

// =========================
// 内部辅助方法
// =========================

// postJSON 发送JSON POST请求
func (c *ProxyHTTPClient) postJSON(ctx context.Context, url string, reqBody interface{}) (*TaskResponse, error) {
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	var result TaskResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w, body: %s", err, string(body))
	}

	if !result.Success {
		return &result, fmt.Errorf("proxy task failed: %s (error: %s)", result.Message, result.Error)
	}

	c.logger.Infof("Proxy task submitted successfully - TaskID: %s, URL: %s", result.TaskID, url)

	return &result, nil
}

// replaceHTTPWithWS 将HTTP协议替换为WebSocket协议
func replaceHTTPWithWS(endpoint string) string {
	if len(endpoint) > 7 && endpoint[:7] == "http://" {
		return "ws://" + endpoint[7:]
	}
	if len(endpoint) > 8 && endpoint[:8] == "https://" {
		return "wss://" + endpoint[8:]
	}
	return endpoint
}
