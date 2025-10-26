package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// JSONLEntry JSONL 文件中的一行
type JSONLEntry map[string]interface{}

// ReadJSONL 读取 JSONL 文件
func ReadJSONL(path string) ([]JSONLEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSONL file: %w", err)
	}
	defer file.Close()

	var entries []JSONLEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var entry JSONLEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// 跳过无效行
			continue
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read JSONL file: %w", err)
	}

	return entries, nil
}

// WriteJSONL 写入 JSONL 文件
func WriteJSONL(path string, entries []JSONLEntry) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create JSONL file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, entry := range entries {
		data, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("failed to marshal entry: %w", err)
		}

		if _, err := writer.Write(data); err != nil {
			return fmt.Errorf("failed to write entry: %w", err)
		}

		if _, err := writer.WriteString("\n"); err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}
	}

	return nil
}

// AppendJSONL 追加到 JSONL 文件
func AppendJSONL(path string, entry JSONLEntry) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open JSONL file: %w", err)
	}
	defer file.Close()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write entry: %w", err)
	}

	if _, err := file.WriteString("\n"); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	return nil
}

// ParseJSONLMetrics 从 JSONL 解析指标
func ParseJSONLMetrics(entries []JSONLEntry) (totalTokens int64, messageCount int64, costUSD float64) {
	for _, entry := range entries {
		// 解析 token 使用
		if usage, ok := entry["usage"].(map[string]interface{}); ok {
			if inputTokens, ok := usage["input_tokens"].(float64); ok {
				totalTokens += int64(inputTokens)
			}
			if outputTokens, ok := usage["output_tokens"].(float64); ok {
				totalTokens += int64(outputTokens)
			}
		}

		// 计算消息数
		if msgType, ok := entry["type"].(string); ok {
			if msgType == "message" {
				messageCount++
			}
		}
	}

	// 简单的成本估算（根据模型调整）
	// Claude 3.5 Sonnet: $3/M input, $15/M output
	// 这里简化为平均 $9/M tokens
	costUSD = float64(totalTokens) * 9.0 / 1000000.0

	return
}

// GetFirstUserMessage 获取第一条用户消息
func GetFirstUserMessage(entries []JSONLEntry) (string, string) {
	for _, entry := range entries {
		if msgType, ok := entry["type"].(string); ok && msgType == "message" {
			if message, ok := entry["message"].(map[string]interface{}); ok {
				if role, ok := message["role"].(string); ok && role == "user" {
					if content, ok := message["content"].(string); ok {
						timestamp := ""
						if ts, ok := entry["timestamp"].(string); ok {
							timestamp = ts
						}
						return content, timestamp
					}
				}
			}
		}
	}
	return "", ""
}

// StreamJSONL 流式读取 JSONL 文件
func StreamJSONL(path string, callback func(JSONLEntry) error) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open JSONL file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var entry JSONLEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		if err := callback(entry); err != nil {
			return err
		}
	}

	return scanner.Err()
}
