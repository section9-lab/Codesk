package model

import "time"

// Project 代表 ~/.claude/projects 目录中的一个项目
type Project struct {
	ID                string     `json:"id"`                  // 项目 ID（目录名）
	Path              string     `json:"path"`                // 原始项目路径
	Sessions          []string   `json:"sessions"`            // 会话 ID 列表
	CreatedAt         int64      `json:"created_at"`          // Unix 时间戳
	MostRecentSession *int64     `json:"most_recent_session"` // 最近会话时间戳
}

// Session 代表一个会话及其元数据
type Session struct {
	ID               string      `json:"id"`                 // 会话 ID (UUID)
	ProjectID        string      `json:"project_id"`         // 所属项目 ID
	ProjectPath      string      `json:"project_path"`       // 项目路径
	TodoData         interface{} `json:"todo_data"`          // 关联的 todo 数据
	CreatedAt        int64       `json:"created_at"`         // Unix 时间戳
	FirstMessage     *string     `json:"first_message"`      // 第一条用户消息
	MessageTimestamp *string     `json:"message_timestamp"`  // 第一条消息时间戳
}

// ClaudeSettings 代表 ~/.claude/settings.json 的设置
type ClaudeSettings struct {
	Data map[string]interface{} `json:"data"`
}

// ClaudeVersionStatus Claude Code 版本状态
type ClaudeVersionStatus struct {
	IsInstalled bool    `json:"is_installed"` // 是否已安装
	Version     *string `json:"version"`      // 版本字符串
	Output      string  `json:"output"`       // 命令输出
}

// SessionMessage JSONL 文件中的消息条目
type SessionMessage struct {
	Type      string           `json:"type"`
	Message   *MessageContent  `json:"message,omitempty"`
	Timestamp *string          `json:"timestamp,omitempty"`
}

// MessageContent 消息内容
type MessageContent struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// SessionHistory 会话历史记录
type SessionHistory struct {
	SessionID string           `json:"session_id"`
	Messages  []SessionMessage `json:"messages"`
}

// RunningSession 正在运行的会话信息
type RunningSession struct {
	SessionID   string    `json:"session_id"`
	ProjectPath string    `json:"project_path"`
	PID         uint32    `json:"pid"`
	StartedAt   time.Time `json:"started_at"`
}
