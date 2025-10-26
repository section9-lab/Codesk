package agent

import (
	"Codesk/backend/model"
	"Codesk/backend/repository"
	"Codesk/backend/service/claude"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AgentService Agent 管理服务
type AgentService struct {
	repo           *repository.AgentRepository
	executionSvc   *claude.ExecutionService
}

// NewAgentService 创建 Agent 服务实例
func NewAgentService() *AgentService {
	return &AgentService{
		repo:         repository.NewAgentRepository(),
		executionSvc: claude.NewExecutionService(),
	}
}

// CreateAgent 创建新 Agent
func (s *AgentService) CreateAgent(agent *model.Agent) error {
	// 设置默认值
	if agent.Icon == "" {
		agent.Icon = "🤖"
	}
	if agent.Model == "" {
		agent.Model = "sonnet"
	}

	return s.repo.Create(agent)
}

// UpdateAgent 更新 Agent
func (s *AgentService) UpdateAgent(agent *model.Agent) error {
	// 验证 Agent 是否存在
	if agent.ID == nil {
		return fmt.Errorf("agent ID is required")
	}

	existing, err := s.repo.GetByID(*agent.ID)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	// 保留创建时间
	agent.CreatedAt = existing.CreatedAt

	return s.repo.Update(agent)
}

// DeleteAgent 删除 Agent
func (s *AgentService) DeleteAgent(id int64) error {
	return s.repo.Delete(id)
}

// GetAgent 获取 Agent
func (s *AgentService) GetAgent(id int64) (*model.Agent, error) {
	return s.repo.GetByID(id)
}

// ListAgents 列出所有 Agents
func (s *AgentService) ListAgents() ([]*model.Agent, error) {
	return s.repo.List()
}

// ExecuteAgent 执行 Agent
func (s *AgentService) ExecuteAgent(agentID int64, task, projectPath string) (*model.AgentRun, error) {
	// 获取 Agent 信息
	agent, err := s.repo.GetByID(agentID)
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	// 使用默认任务（如果未提供）
	if task == "" && agent.DefaultTask != nil {
		task = *agent.DefaultTask
	}

	if task == "" {
		return nil, fmt.Errorf("task is required")
	}

	// 生成 session ID
	sessionID := uuid.New().String()

	// 创建运行记录
	now := time.Now()
	run := &model.AgentRun{
		AgentID:          agentID,
		AgentName:        agent.Name,
		AgentIcon:        agent.Icon,
		Task:             task,
		Model:            agent.Model,
		ProjectPath:      projectPath,
		SessionID:        sessionID,
		Status:           "pending",
		ProcessStartedAt: &now,
	}

	if err := s.repo.CreateRun(run); err != nil {
		return nil, fmt.Errorf("failed to create run record: %w", err)
	}

	// 执行 Claude Code
	opts := claude.ExecuteOptions{
		ProjectPath: projectPath,
		Task:        s.buildTaskPrompt(agent, task),
		Model:       agent.Model,
		SessionID:   sessionID,
	}

	result, err := s.executionSvc.Execute(opts)
	if err != nil {
		// 更新状态为失败
		s.repo.UpdateRunStatus(*run.ID, "failed", nil)
		return nil, fmt.Errorf("failed to execute agent: %w", err)
	}

	// 更新运行记录
	pid := uint32(result.PID)
	run.PID = &pid
	run.Status = "running"

	return run, nil
}

// buildTaskPrompt 构建任务提示词
func (s *AgentService) buildTaskPrompt(agent *model.Agent, task string) string {
	// 组合系统提示词和任务
	prompt := agent.SystemPrompt + "\n\n"
	prompt += "Task: " + task

	return prompt
}

// ListAgentRuns 列出 Agent 的运行记录
func (s *AgentService) ListAgentRuns(agentID int64) ([]*model.AgentRun, error) {
	return s.repo.ListRuns(agentID)
}

// GetAgentRun 获取运行记录
func (s *AgentService) GetAgentRun(runID int64) (*model.AgentRun, error) {
	return s.repo.GetRunByID(runID)
}

// KillAgentSession 终止 Agent 会话
func (s *AgentService) KillAgentSession(sessionID string) error {
	return s.executionSvc.Cancel(sessionID)
}

// GetSessionOutput 获取会话输出
func (s *AgentService) GetSessionOutput(sessionID, projectPath string) (string, error) {
	return s.executionSvc.GetSessionOutput(sessionID, projectPath)
}

// GetSessionStatus 获取会话状态
func (s *AgentService) GetSessionStatus(sessionID, projectPath string) (string, error) {
	return s.executionSvc.GetSessionStatus(sessionID, projectPath)
}

// ExportAgent 导出 Agent
func (s *AgentService) ExportAgent(agentID int64) (*model.AgentExport, error) {
	agent, err := s.repo.GetByID(agentID)
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	export := &model.AgentExport{
		Version:    1,
		ExportedAt: time.Now(),
		Agent: model.AgentData{
			Name:         agent.Name,
			Icon:         agent.Icon,
			SystemPrompt: agent.SystemPrompt,
			DefaultTask:  agent.DefaultTask,
			Model:        agent.Model,
			Hooks:        agent.Hooks,
		},
	}

	return export, nil
}

// ExportAgentToJSON 导出 Agent 为 JSON 字符串
func (s *AgentService) ExportAgentToJSON(agentID int64) (string, error) {
	export, err := s.ExportAgent(agentID)
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal export: %w", err)
	}

	return string(data), nil
}

// ImportAgent 导入 Agent
func (s *AgentService) ImportAgent(exportData *model.AgentExport) (*model.Agent, error) {
	agent := &model.Agent{
		Name:         exportData.Agent.Name,
		Icon:         exportData.Agent.Icon,
		SystemPrompt: exportData.Agent.SystemPrompt,
		DefaultTask:  exportData.Agent.DefaultTask,
		Model:        exportData.Agent.Model,
		Hooks:        exportData.Agent.Hooks,
		EnableFileRead:  true,
		EnableFileWrite: true,
		EnableNetwork:   false,
	}

	if err := s.repo.Create(agent); err != nil {
		return nil, fmt.Errorf("failed to import agent: %w", err)
	}

	return agent, nil
}

// ImportAgentFromJSON 从 JSON 字符串导入 Agent
func (s *AgentService) ImportAgentFromJSON(jsonData string) (*model.Agent, error) {
	var exportData model.AgentExport
	if err := json.Unmarshal([]byte(jsonData), &exportData); err != nil {
		return nil, fmt.Errorf("failed to parse import data: %w", err)
	}

	return s.ImportAgent(&exportData)
}
