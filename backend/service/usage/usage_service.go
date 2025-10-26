package usage

import (
	"Codesk/backend/config"
	"Codesk/backend/model"
	"Codesk/backend/util"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UsageService 使用统计服务
type UsageService struct{}

// NewUsageService 创建使用统计服务实例
func NewUsageService() *UsageService {
	return &UsageService{}
}

// 定价常量（每百万 tokens）
const (
	// Claude 3.5 Sonnet
	Sonnet35InputPrice       = 3.0
	Sonnet35OutputPrice      = 15.0
	Sonnet35CacheWritePrice  = 3.75
	Sonnet35CacheReadPrice   = 0.30

	// Claude Opus
	OpusInputPrice      = 15.0
	OpusOutputPrice     = 75.0
	OpusCacheWritePrice = 18.75
	OpusCacheReadPrice  = 1.50
)

// GetUsageStats 获取总体使用统计
func (s *UsageService) GetUsageStats() (*model.UsageStats, error) {
	projectsDir := config.GetProjectsDir()

	if !util.DirExists(projectsDir) {
		return &model.UsageStats{}, nil
	}

	var totalSessions int
	var totalMessages int
	var totalTokens int64
	var totalCost float64

	// 遍历所有项目
	projects, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read projects directory: %w", err)
	}

	for _, project := range projects {
		if !project.IsDir() {
			continue
		}

		projectDir := filepath.Join(projectsDir, project.Name())
		sessions, err := os.ReadDir(projectDir)
		if err != nil {
			continue
		}

		for _, session := range sessions {
			if !strings.HasSuffix(session.Name(), ".jsonl") {
				continue
			}

			sessionPath := filepath.Join(projectDir, session.Name())
			tokens, messages, cost := s.parseSessionFile(sessionPath)

			totalSessions++
			totalMessages += messages
			totalTokens += tokens
			totalCost += cost
		}
	}

	avgTokens := float64(0)
	avgCost := float64(0)
	if totalSessions > 0 {
		avgTokens = float64(totalTokens) / float64(totalSessions)
		avgCost = totalCost / float64(totalSessions)
	}

	return &model.UsageStats{
		TotalSessions:  totalSessions,
		TotalMessages:  totalMessages,
		TotalTokens:    totalTokens,
		TotalCostUSD:   totalCost,
		AverageTokens:  avgTokens,
		AverageCostUSD: avgCost,
	}, nil
}

// GetUsageByDateRange 获取日期范围内的使用统计
func (s *UsageService) GetUsageByDateRange(startDate, endDate string) (*model.UsageByDateRange, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	projectsDir := config.GetProjectsDir()
	if !util.DirExists(projectsDir) {
		return &model.UsageByDateRange{
			StartDate: startDate,
			EndDate:   endDate,
		}, nil
	}

	dailyStatsMap := make(map[string]*model.DailyUsageStats)
	var totalSessions int
	var totalTokens int64
	var totalCost float64

	// 遍历所有项目和会话
	projects, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read projects directory: %w", err)
	}

	for _, project := range projects {
		if !project.IsDir() {
			continue
		}

		projectDir := filepath.Join(projectsDir, project.Name())
		sessions, err := os.ReadDir(projectDir)
		if err != nil {
			continue
		}

		for _, session := range sessions {
			if !strings.HasSuffix(session.Name(), ".jsonl") {
				continue
			}

			sessionPath := filepath.Join(projectDir, session.Name())
			info, err := os.Stat(sessionPath)
			if err != nil {
				continue
			}

			// 检查日期范围
			modTime := info.ModTime()
			if modTime.Before(start) || modTime.After(end.Add(24*time.Hour)) {
				continue
			}

			tokens, _, cost := s.parseSessionFile(sessionPath)
			dateKey := modTime.Format("2006-01-02")

			if _, exists := dailyStatsMap[dateKey]; !exists {
				dailyStatsMap[dateKey] = &model.DailyUsageStats{
					Date: dateKey,
				}
			}

			dailyStatsMap[dateKey].Sessions++
			dailyStatsMap[dateKey].Tokens += tokens
			dailyStatsMap[dateKey].CostUSD += cost

			totalSessions++
			totalTokens += tokens
			totalCost += cost
		}
	}

	// 转换为数组
	var dailyStats []model.DailyUsageStats
	for _, stats := range dailyStatsMap {
		dailyStats = append(dailyStats, *stats)
	}

	return &model.UsageByDateRange{
		StartDate:     startDate,
		EndDate:       endDate,
		TotalSessions: totalSessions,
		TotalTokens:   totalTokens,
		TotalCostUSD:  totalCost,
		DailyStats:    dailyStats,
	}, nil
}

// GetSessionStats 获取会话统计
func (s *UsageService) GetSessionStats(sessionID, projectPath string) (*model.SessionStats, error) {
	encodedPath := strings.ReplaceAll(projectPath, "/", "-")
	sessionPath := filepath.Join(config.GetProjectsDir(), encodedPath, sessionID+".jsonl")

	if !util.FileExists(sessionPath) {
		return nil, fmt.Errorf("session file not found")
	}

	entries, err := util.ReadJSONL(sessionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read session: %w", err)
	}

	var messageCount int
	var inputTokens, outputTokens, totalTokens int64
	var cost float64
	var startTime, endTime *time.Time

	for _, entry := range entries {
		// 解析时间戳
		if ts, ok := entry["timestamp"].(string); ok {
			t, err := time.Parse(time.RFC3339, ts)
			if err == nil {
				if startTime == nil || t.Before(*startTime) {
					startTime = &t
				}
				if endTime == nil || t.After(*endTime) {
					endTime = &t
				}
			}
		}

		// 解析消息
		if msgType, ok := entry["type"].(string); ok && msgType == "message" {
			messageCount++

			// 解析 usage
			if message, ok := entry["message"].(map[string]interface{}); ok {
				if usage, ok := message["usage"].(map[string]interface{}); ok {
					if input, ok := usage["input_tokens"].(float64); ok {
						inputTokens += int64(input)
					}
					if output, ok := usage["output_tokens"].(float64); ok {
						outputTokens += int64(output)
					}
				}

				// 获取模型
				model := "sonnet"
				if m, ok := message["model"].(string); ok {
					model = m
				}

				// 计算成本
				if usage, ok := message["usage"].(map[string]interface{}); ok {
					cost += s.calculateCost(model, usage)
				}
			}
		}
	}

	totalTokens = inputTokens + outputTokens

	var duration int64
	if startTime != nil && endTime != nil {
		duration = endTime.Sub(*startTime).Milliseconds()
	}

	stats := &model.SessionStats{
		SessionID:    sessionID,
		MessageCount: messageCount,
		TotalTokens:  totalTokens,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		CostUSD:      cost,
		Duration:     duration,
	}

	if startTime != nil {
		stats.StartTime = *startTime
	}
	if endTime != nil {
		stats.EndTime = endTime
	}

	return stats, nil
}

// parseSessionFile 解析会话文件
func (s *UsageService) parseSessionFile(sessionPath string) (tokens int64, messages int, cost float64) {
	entries, err := util.ReadJSONL(sessionPath)
	if err != nil {
		return 0, 0, 0
	}

	tokens, messages64, cost := util.ParseJSONLMetrics(entries)
	messages = int(messages64)

	return
}

// calculateCost 计算成本
func (s *UsageService) calculateCost(model string, usage map[string]interface{}) float64 {
	inputTokens := s.getTokenValue(usage, "input_tokens")
	outputTokens := s.getTokenValue(usage, "output_tokens")
	cacheCreation := s.getTokenValue(usage, "cache_creation_input_tokens")
	cacheRead := s.getTokenValue(usage, "cache_read_input_tokens")

	var inputPrice, outputPrice, cacheWritePrice, cacheReadPrice float64

	// 根据模型选择定价
	if strings.Contains(strings.ToLower(model), "opus") {
		inputPrice = OpusInputPrice
		outputPrice = OpusOutputPrice
		cacheWritePrice = OpusCacheWritePrice
		cacheReadPrice = OpusCacheReadPrice
	} else {
		// 默认使用 Sonnet 定价
		inputPrice = Sonnet35InputPrice
		outputPrice = Sonnet35OutputPrice
		cacheWritePrice = Sonnet35CacheWritePrice
		cacheReadPrice = Sonnet35CacheReadPrice
	}

	// 计算成本（价格是每百万 tokens）
	cost := (float64(inputTokens) * inputPrice / 1000000.0) +
		(float64(outputTokens) * outputPrice / 1000000.0) +
		(float64(cacheCreation) * cacheWritePrice / 1000000.0) +
		(float64(cacheRead) * cacheReadPrice / 1000000.0)

	return cost
}

// getTokenValue 从 usage 中获取 token 值
func (s *UsageService) getTokenValue(usage map[string]interface{}, key string) int64 {
	if val, ok := usage[key].(float64); ok {
		return int64(val)
	}
	return 0
}
