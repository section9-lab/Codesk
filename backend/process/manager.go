package process

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

// ProcessInfo 进程信息
type ProcessInfo struct {
	SessionID   string
	ProjectPath string
	PID         int
	StartedAt   time.Time
	Cmd         *exec.Cmd
	Cancel      context.CancelFunc
	Status      string // running, completed, failed, cancelled
}

// ProcessManager 进程管理器
type ProcessManager struct {
	processes map[string]*ProcessInfo
	mu        sync.RWMutex
}

var (
	globalManager *ProcessManager
	once          sync.Once
)

// GetManager 获取全局进程管理器实例
func GetManager() *ProcessManager {
	once.Do(func() {
		globalManager = &ProcessManager{
			processes: make(map[string]*ProcessInfo),
		}
	})
	return globalManager
}

// Register 注册进程
func (m *ProcessManager) Register(sessionID, projectPath string, cmd *exec.Cmd, cancel context.CancelFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.processes[sessionID] = &ProcessInfo{
		SessionID:   sessionID,
		ProjectPath: projectPath,
		PID:         cmd.Process.Pid,
		StartedAt:   time.Now(),
		Cmd:         cmd,
		Cancel:      cancel,
		Status:      "running",
	}
}

// Unregister 注销进程
func (m *ProcessManager) Unregister(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.processes, sessionID)
}

// Get 获取进程信息
func (m *ProcessManager) Get(sessionID string) (*ProcessInfo, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	info, exists := m.processes[sessionID]
	return info, exists
}

// List 列出所有运行中的进程
func (m *ProcessManager) List() []*ProcessInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []*ProcessInfo
	for _, info := range m.processes {
		list = append(list, info)
	}

	return list
}

// Kill 终止进程
func (m *ProcessManager) Kill(sessionID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	info, exists := m.processes[sessionID]
	if !exists {
		return fmt.Errorf("process not found: %s", sessionID)
	}

	// 调用取消函数
	if info.Cancel != nil {
		info.Cancel()
	}

	// 终止进程
	if info.Cmd != nil && info.Cmd.Process != nil {
		if err := info.Cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
	}

	info.Status = "cancelled"
	delete(m.processes, sessionID)

	return nil
}

// UpdateStatus 更新进程状态
func (m *ProcessManager) UpdateStatus(sessionID, status string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if info, exists := m.processes[sessionID]; exists {
		info.Status = status
	}
}

// CleanupFinished 清理已完成的进程
func (m *ProcessManager) CleanupFinished() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	count := 0
	for sessionID, info := range m.processes {
		if info.Status == "completed" || info.Status == "failed" || info.Status == "cancelled" {
			delete(m.processes, sessionID)
			count++
		}
	}

	return count
}

// OutputHandler 输出处理器接口
type OutputHandler func(line string) error

// StreamOutput 流式读取进程输出
func StreamOutput(scanner *bufio.Scanner, handler OutputHandler) error {
	for scanner.Scan() {
		line := scanner.Text()
		if err := handler(line); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// CombineOutputs 合并标准输出和错误输出
func CombineOutputs(cmd *exec.Cmd, handler OutputHandler) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	var wg sync.WaitGroup
	var streamErr error

	// 读取 stdout
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		if err := StreamOutput(scanner, handler); err != nil {
			streamErr = err
		}
	}()

	// 读取 stderr
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		if err := StreamOutput(scanner, handler); err != nil {
			streamErr = err
		}
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return streamErr
}
