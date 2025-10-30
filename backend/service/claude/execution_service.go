package claude

import (
	"bufio"
	"Codesk/backend/config"
	"Codesk/backend/model"
	"Codesk/backend/process"
	"Codesk/backend/util"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// ExecutionService Claude 执行服务
type ExecutionService struct{}

// NewExecutionService 创建执行服务实例
func NewExecutionService() *ExecutionService {
	return &ExecutionService{}
}

// ExecuteOptions 执行选项
type ExecuteOptions struct {
	ProjectPath string
	Task        string
	Model       string
	Continue    bool
	Resume      bool
	SessionID   string
}

// ExecuteResult 执行结果
type ExecuteResult struct {
	SessionID   string `json:"session_id"`
	ProjectPath string `json:"project_path"`
	PID         int    `json:"pid"`
	Message     string `json:"message"`
}

// Execute 执行 Claude Code
func (s *ExecutionService) Execute(opts ExecuteOptions) (*ExecuteResult, error) {
	// 查找 Claude 二进制
	claudePath, err := util.FindClaudeBinary()
	if err != nil {
		return nil, fmt.Errorf("claude binary not found: %w", err)
	}

	// 生成或使用现有 session ID
	sessionID := opts.SessionID
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	// 构建命令
	ctx, cancel := context.WithCancel(context.Background())
	cmd := s.buildCommand(ctx, claudePath, opts, sessionID)

	// 设置工作目录
	cmd.Dir = opts.ProjectPath

	// 创建管道来捕获输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// 启动进程
	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to start claude: %w", err)
	}

	// 注册进程
	process.GetManager().Register(sessionID, opts.ProjectPath, cmd, cancel)

	// 启动输出读取协程
	go s.readOutput(stdout, stderr, sessionID)

	// 异步等待进程完成
	go s.waitForCompletion(cmd, sessionID, cancel)

	return &ExecuteResult{
		SessionID:   sessionID,
		ProjectPath: opts.ProjectPath,
		PID:         cmd.Process.Pid,
		Message:     "Claude Code started successfully",
	}, nil
}

// buildCommand 构建执行命令
func (s *ExecutionService) buildCommand(ctx context.Context, claudePath string, opts ExecuteOptions, sessionID string) *exec.Cmd {
	args := []string{}

	// 添加模型参数
	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}

	// 添加会话 ID
	args = append(args, "--session-id", sessionID)

	// 添加任务或继续/恢复标志
	if opts.Continue {
		args = append(args, "--continue")
	} else if opts.Resume {
		args = append(args, "--resume")
		// Resume模式也需要包含prompt
		if opts.Task != "" {
			args = append(args, opts.Task)
		}
	} else if opts.Task != "" {
		args = append(args, opts.Task)
	}

	cmd := exec.CommandContext(ctx, claudePath, args...)

	// 设置环境变量
	cmd.Env = s.buildEnvironment()

	return cmd
}

// buildEnvironment 构建环境变量
func (s *ExecutionService) buildEnvironment() []string {
	env := os.Environ()

	// 确保包含必要的环境变量
	requiredVars := []string{
		"PATH", "HOME", "USER", "SHELL", "LANG",
		"NODE_PATH", "NVM_DIR", "NVM_BIN",
		"HOMEBREW_PREFIX", "HOMEBREW_CELLAR",
	}

	envMap := make(map[string]string)
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	// 确保所有必需变量都存在
	var result []string
	for _, key := range requiredVars {
		if val, exists := envMap[key]; exists {
			result = append(result, fmt.Sprintf("%s=%s", key, val))
		}
	}

	// 添加其他环境变量
	for key, val := range envMap {
		found := false
		for _, req := range requiredVars {
			if key == req {
				found = true
				break
			}
		}
		if !found {
			result = append(result, fmt.Sprintf("%s=%s", key, val))
		}
	}

	return result
}

// waitForCompletion 等待进程完成
func (s *ExecutionService) waitForCompletion(cmd *exec.Cmd, sessionID string, cancel context.CancelFunc) {
	defer cancel()

	err := cmd.Wait()

	if err != nil {
		process.GetManager().UpdateStatus(sessionID, "failed")
	} else {
		process.GetManager().UpdateStatus(sessionID, "completed")
	}

	// 延迟注销进程（给前端时间读取状态）
	// time.Sleep(5 * time.Second)
	// process.GetManager().Unregister(sessionID)
}

// Cancel 取消执行
func (s *ExecutionService) Cancel(sessionID string) error {
	return process.GetManager().Kill(sessionID)
}

// GetSessionOutput 获取会话输出
func (s *ExecutionService) GetSessionOutput(sessionID, projectPath string) (string, error) {
	reader := process.NewSessionOutputReader(sessionID, projectPath, config.GetClaudeDir())
	return reader.ReadOutput()
}

// GetSessionStatus 获取会话状态
func (s *ExecutionService) GetSessionStatus(sessionID, projectPath string) (string, error) {
	reader := process.NewSessionOutputReader(sessionID, projectPath, config.GetClaudeDir())
	return reader.GetStatus()
}

// ListRunningSessions 列出所有运行中的会话
func (s *ExecutionService) ListRunningSessions() []map[string]interface{} {
	processes := process.GetManager().List()

	var sessions []map[string]interface{}
	for _, proc := range processes {
		sessions = append(sessions, map[string]interface{}{
			"session_id":   proc.SessionID,
			"project_path": proc.ProjectPath,
			"pid":          proc.PID,
			"started_at":   proc.StartedAt,
			"status":       proc.Status,
		})
	}

	return sessions
}

// CleanupFinishedProcesses 清理已完成的进程
func (s *ExecutionService) CleanupFinishedProcesses() int {
	return process.GetManager().CleanupFinished()
}

// CheckClaudeVersion 检查 Claude Code 版本
func (s *ExecutionService) CheckClaudeVersion() (bool, string, error) {
	isInstalled, version, err := util.CheckClaudeInstallation()
	return isInstalled, version, err
}

// ListClaudeInstallations 列出所有 Claude 安装位置
func (s *ExecutionService) ListClaudeInstallations() []model.ClaudeInstallation {
	return util.ListClaudeInstallations()
}

// GetHomeDirectory 获取用户主目录
func (s *ExecutionService) GetHomeDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return homeDir, nil
}

// CreateProject 创建新项目（确保目录存在）
func (s *ExecutionService) CreateProject(projectPath string) error {
	if !util.DirExists(projectPath) {
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}
	}

	// 创建 CLAUDE.md 文件（如果不存在）
	claudeMdPath := filepath.Join(projectPath, "CLAUDE.md")
	if !util.FileExists(claudeMdPath) {
		defaultContent := "# Project Instructions\n\nAdd your project-specific instructions here.\n"
		if err := util.WriteFile(claudeMdPath, defaultContent); err != nil {
			return fmt.Errorf("failed to create CLAUDE.md: %w", err)
		}
	}

	return nil
}

// OpenNewSession 打开新会话（生成 session ID）
func (s *ExecutionService) OpenNewSession() string {
	return uuid.New().String()
}

// readOutput 读取进程输出并保存到session文件
func (s *ExecutionService) readOutput(stdout, stderr io.ReadCloser, sessionID string) {
	defer stdout.Close()
	defer stderr.Close()

	// 创建session输出文件
	sessionFile := filepath.Join(config.GetClaudeDir(), "session-env", sessionID+".jsonl")
	if err := os.MkdirAll(filepath.Dir(sessionFile), 0755); err != nil {
		fmt.Printf("Failed to create session directory: %v\n", err)
		return
	}

	file, err := os.OpenFile(sessionFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Failed to create session file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// 读取stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			// 写入session文件
			fmt.Fprintf(writer, "%s\n", line)
			writer.Flush()

			// 这里可以添加事件发送逻辑
			// TODO: 发送claude-output事件
		}
	}()

	// 读取stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			// 写入session文件
			fmt.Fprintf(writer, "%s\n", line)
			writer.Flush()

			// 这里可以添加错误事件发送逻辑
			// TODO: 发送claude-error事件
		}
	}()
}
