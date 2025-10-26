package model

import "time"

// Agent 代表一个 CC Agent
type Agent struct {
	ID              *int64     `json:"id"`
	Name            string     `json:"name"`
	Icon            string     `json:"icon"`
	SystemPrompt    string     `json:"system_prompt"`
	DefaultTask     *string    `json:"default_task"`
	Model           string     `json:"model"`
	EnableFileRead  bool       `json:"enable_file_read"`
	EnableFileWrite bool       `json:"enable_file_write"`
	EnableNetwork   bool       `json:"enable_network"`
	Hooks           *string    `json:"hooks"` // JSON string
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// AgentRun 代表一次 Agent 执行
type AgentRun struct {
	ID               *int64     `json:"id"`
	AgentID          int64      `json:"agent_id"`
	AgentName        string     `json:"agent_name"`
	AgentIcon        string     `json:"agent_icon"`
	Task             string     `json:"task"`
	Model            string     `json:"model"`
	ProjectPath      string     `json:"project_path"`
	SessionID        string     `json:"session_id"` // UUID
	Status           string     `json:"status"`     // pending, running, completed, failed, cancelled
	PID              *uint32    `json:"pid"`
	ProcessStartedAt *time.Time `json:"process_started_at"`
	CreatedAt        time.Time  `json:"created_at"`
	CompletedAt      *time.Time `json:"completed_at"`
}

// AgentRunMetrics 从 JSONL 计算的运行时指标
type AgentRunMetrics struct {
	DurationMs   *int64   `json:"duration_ms"`
	TotalTokens  *int64   `json:"total_tokens"`
	CostUSD      *float64 `json:"cost_usd"`
	MessageCount *int64   `json:"message_count"`
}

// AgentRunWithMetrics Agent 运行记录和实时指标
type AgentRunWithMetrics struct {
	AgentRun
	Metrics *AgentRunMetrics `json:"metrics"`
	Output  *string          `json:"output"` // 实时 JSONL 内容
}

// AgentExport Agent 导出格式
type AgentExport struct {
	Version    uint32     `json:"version"`
	ExportedAt time.Time  `json:"exported_at"`
	Agent      AgentData  `json:"agent"`
}

// AgentData 导出的 Agent 数据
type AgentData struct {
	Name         string  `json:"name"`
	Icon         string  `json:"icon"`
	SystemPrompt string  `json:"system_prompt"`
	DefaultTask  *string `json:"default_task"`
	Model        string  `json:"model"`
	Hooks        *string `json:"hooks"`
}
