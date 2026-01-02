package agentclient

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Client struct {
    endpoints []string
    http      *http.Client
    token     string
}

func New(endpoints []string, timeoutMS int, token string) *Client {
    if timeoutMS <= 0 { timeoutMS = 5000 }
    return &Client{endpoints: endpoints, http: &http.Client{Timeout: time.Duration(timeoutMS) * time.Millisecond}, token: token}
}

// ExecuteCommand 调用 nb-agent 执行命令（示例路径：/api/v1/tasks 或 /agent/task/execute）
// 具体路径随后与 nb-agent 对齐；此处优先尝试 /api/v1/tasks，失败时尝试 /agent/task/execute。
func (c *Client) ExecuteCommand(payload any) (string, error) {
    if len(c.endpoints) == 0 { return "", fmt.Errorf("no agent endpoints configured") }
    body, _ := json.Marshal(payload)
    var lastErr error
    for _, base := range c.endpoints {
        // 1) try /api/v1/tasks
        if id, err := c.postForTaskID(base+"/api/v1/tasks", body); err == nil { return id, nil } else { lastErr = err }
        // 2) try /agent/task/execute
        if id, err := c.postForTaskID(base+"/agent/task/execute", body); err == nil { return id, nil } else { lastErr = err }
    }
    if lastErr == nil { lastErr = fmt.Errorf("agent execute failed") }
    return "", lastErr
}

func (c *Client) postForTaskID(url string, body []byte) (string, error) {
    req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
    if err != nil { return "", err }
    req.Header.Set("Content-Type", "application/json")
    if c.token != "" { req.Header.Set("Authorization", "Bearer "+c.token) }
    resp, err := c.http.Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return "", fmt.Errorf("status %d", resp.StatusCode)
    }
    // 兼容不同返回结构，尝试读取常见字段
    var x struct{ TaskId string `json:"taskId"`; ID string `json:"id"` }
    _ = json.NewDecoder(resp.Body).Decode(&x)
    if x.TaskId != "" { return x.TaskId, nil }
    if x.ID != "" { return x.ID, nil }
    return "", nil
}

