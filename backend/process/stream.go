package process

import (
	"Codesk/backend/util"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// OutputBuffer 输出缓冲区
type OutputBuffer struct {
	sessionID string
	lines     []string
	mu        sync.RWMutex
}

// NewOutputBuffer 创建输出缓冲区
func NewOutputBuffer(sessionID string) *OutputBuffer {
	return &OutputBuffer{
		sessionID: sessionID,
		lines:     make([]string, 0),
	}
}

// Append 追加输出行
func (b *OutputBuffer) Append(line string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.lines = append(b.lines, line)
}

// GetAll 获取所有输出
func (b *OutputBuffer) GetAll() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make([]string, len(b.lines))
	copy(result, b.lines)
	return result
}

// GetLast 获取最后 n 行
func (b *OutputBuffer) GetLast(n int) []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if n >= len(b.lines) {
		result := make([]string, len(b.lines))
		copy(result, b.lines)
		return result
	}

	result := make([]string, n)
	copy(result, b.lines[len(b.lines)-n:])
	return result
}

// Clear 清空缓冲区
func (b *OutputBuffer) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.lines = make([]string, 0)
}

// OutputBufferManager 输出缓冲区管理器
type OutputBufferManager struct {
	buffers map[string]*OutputBuffer
	mu      sync.RWMutex
}

var (
	bufferManager     *OutputBufferManager
	bufferManagerOnce sync.Once
)

// GetBufferManager 获取全局缓冲区管理器
func GetBufferManager() *OutputBufferManager {
	bufferManagerOnce.Do(func() {
		bufferManager = &OutputBufferManager{
			buffers: make(map[string]*OutputBuffer),
		}
	})
	return bufferManager
}

// GetOrCreate 获取或创建缓冲区
func (m *OutputBufferManager) GetOrCreate(sessionID string) *OutputBuffer {
	m.mu.Lock()
	defer m.mu.Unlock()

	if buffer, exists := m.buffers[sessionID]; exists {
		return buffer
	}

	buffer := NewOutputBuffer(sessionID)
	m.buffers[sessionID] = buffer
	return buffer
}

// Get 获取缓冲区
func (m *OutputBufferManager) Get(sessionID string) (*OutputBuffer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	buffer, exists := m.buffers[sessionID]
	return buffer, exists
}

// Remove 移除缓冲区
func (m *OutputBufferManager) Remove(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.buffers, sessionID)
}

// SessionOutputReader 会话输出读取器
type SessionOutputReader struct {
	sessionID   string
	projectPath string
	claudeDir   string
}

// NewSessionOutputReader 创建会话输出读取器
func NewSessionOutputReader(sessionID, projectPath, claudeDir string) *SessionOutputReader {
	return &SessionOutputReader{
		sessionID:   sessionID,
		projectPath: projectPath,
		claudeDir:   claudeDir,
	}
}

// GetSessionFilePath 获取会话文件路径
func (r *SessionOutputReader) GetSessionFilePath() string {
	encodedPath := strings.ReplaceAll(r.projectPath, "/", "-")
	return filepath.Join(r.claudeDir, "projects", encodedPath, r.sessionID+".jsonl")
}

// ReadOutput 读取会话输出
func (r *SessionOutputReader) ReadOutput() (string, error) {
	filePath := r.GetSessionFilePath()

	if !util.FileExists(filePath) {
		return "", fmt.Errorf("session file not found: %s", filePath)
	}

	content, err := util.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read session file: %w", err)
	}

	return content, nil
}

// ReadOutputLines 读取会话输出行
func (r *SessionOutputReader) ReadOutputLines() ([]string, error) {
	content, err := r.ReadOutput()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(content, "\n")
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}

	return result, nil
}

// StreamOutput 流式读取输出（用于实时监控）
func (r *SessionOutputReader) StreamOutput(callback func(line string) error) error {
	filePath := r.GetSessionFilePath()

	if !util.FileExists(filePath) {
		return fmt.Errorf("session file not found: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open session file: %w", err)
	}
	defer file.Close()

	return util.StreamJSONL(filePath, func(entry util.JSONLEntry) error {
		data, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		return callback(string(data))
	})
}

// GetMetrics 获取会话指标
func (r *SessionOutputReader) GetMetrics() (totalTokens int64, messageCount int64, costUSD float64, err error) {
	entries, err := util.ReadJSONL(r.GetSessionFilePath())
	if err != nil {
		return 0, 0, 0, err
	}

	totalTokens, messageCount, costUSD = util.ParseJSONLMetrics(entries)
	return
}

// GetStatus 获取会话状态
func (r *SessionOutputReader) GetStatus() (string, error) {
	// 检查进程管理器中的状态
	if info, exists := GetManager().Get(r.sessionID); exists {
		return info.Status, nil
	}

	// 检查文件是否存在
	filePath := r.GetSessionFilePath()
	if !util.FileExists(filePath) {
		return "not_found", nil
	}

	// 如果文件存在但进程不在运行，则认为已完成
	return "completed", nil
}
