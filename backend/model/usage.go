package model

import "time"

// UsageStats 使用统计
type UsageStats struct {
	TotalSessions    int     `json:"total_sessions"`
	TotalMessages    int     `json:"total_messages"`
	TotalTokens      int64   `json:"total_tokens"`
	TotalCostUSD     float64 `json:"total_cost_usd"`
	AverageTokens    float64 `json:"average_tokens"`
	AverageCostUSD   float64 `json:"average_cost_usd"`
}

// UsageByDateRange 按日期范围的使用统计
type UsageByDateRange struct {
	StartDate     string           `json:"start_date"`
	EndDate       string           `json:"end_date"`
	TotalSessions int              `json:"total_sessions"`
	TotalTokens   int64            `json:"total_tokens"`
	TotalCostUSD  float64          `json:"total_cost_usd"`
	DailyStats    []DailyUsageStats `json:"daily_stats"`
}

// DailyUsageStats 每日使用统计
type DailyUsageStats struct {
	Date     string  `json:"date"`
	Sessions int     `json:"sessions"`
	Tokens   int64   `json:"tokens"`
	CostUSD  float64 `json:"cost_usd"`
}

// UsageDetail 详细使用记录
type UsageDetail struct {
	SessionID   string    `json:"session_id"`
	ProjectPath string    `json:"project_path"`
	Model       string    `json:"model"`
	Tokens      int64     `json:"tokens"`
	CostUSD     float64   `json:"cost_usd"`
	StartTime   time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Duration    *int64    `json:"duration"` // 毫秒
}

// SessionStats 会话统计
type SessionStats struct {
	SessionID      string    `json:"session_id"`
	MessageCount   int       `json:"message_count"`
	TotalTokens    int64     `json:"total_tokens"`
	InputTokens    int64     `json:"input_tokens"`
	OutputTokens   int64     `json:"output_tokens"`
	CostUSD        float64   `json:"cost_usd"`
	Duration       int64     `json:"duration"` // 毫秒
	StartTime      time.Time `json:"start_time"`
	EndTime        *time.Time `json:"end_time"`
}

// TokenUsage Token 使用详情
type TokenUsage struct {
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
	TotalTokens  int64 `json:"total_tokens"`
}

// CostCalculation 成本计算
type CostCalculation struct {
	Model        string  `json:"model"`
	InputTokens  int64   `json:"input_tokens"`
	OutputTokens int64   `json:"output_tokens"`
	InputCost    float64 `json:"input_cost"`
	OutputCost   float64 `json:"output_cost"`
	TotalCost    float64 `json:"total_cost"`
}
