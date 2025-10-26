package main

import "Codesk/backend/model"

// ==================== Agent Management ====================

// CreateAgent 创建新 Agent
func (a *App) CreateAgent(agent *model.Agent) error {
	return a.agentService.CreateAgent(agent)
}

// UpdateAgent 更新 Agent
func (a *App) UpdateAgent(agent *model.Agent) error {
	return a.agentService.UpdateAgent(agent)
}

// DeleteAgent 删除 Agent
func (a *App) DeleteAgent(id int64) error {
	return a.agentService.DeleteAgent(id)
}

// GetAgent 获取 Agent
func (a *App) GetAgent(id int64) (*model.Agent, error) {
	return a.agentService.GetAgent(id)
}

// ListAgents 列出所有 Agents
func (a *App) ListAgents() ([]*model.Agent, error) {
	return a.agentService.ListAgents()
}

// ExecuteAgent 执行 Agent
func (a *App) ExecuteAgent(agentID int64, task, projectPath string) (*model.AgentRun, error) {
	return a.agentService.ExecuteAgent(agentID, task, projectPath)
}

// ListAgentRuns 列出 Agent 的运行记录
func (a *App) ListAgentRuns(agentID int64) ([]*model.AgentRun, error) {
	return a.agentService.ListAgentRuns(agentID)
}

// GetAgentRun 获取运行记录
func (a *App) GetAgentRun(runID int64) (*model.AgentRun, error) {
	return a.agentService.GetAgentRun(runID)
}

// KillAgentSession 终止 Agent 会话
func (a *App) KillAgentSession(sessionID string) error {
	return a.agentService.KillAgentSession(sessionID)
}

// GetAgentSessionOutput 获取 Agent 会话输出
func (a *App) GetAgentSessionOutput(sessionID, projectPath string) (string, error) {
	return a.agentService.GetSessionOutput(sessionID, projectPath)
}

// GetAgentSessionStatus 获取 Agent 会话状态
func (a *App) GetAgentSessionStatus(sessionID, projectPath string) (string, error) {
	return a.agentService.GetSessionStatus(sessionID, projectPath)
}

// ExportAgent 导出 Agent
func (a *App) ExportAgent(agentID int64) (*model.AgentExport, error) {
	return a.agentService.ExportAgent(agentID)
}

// ExportAgentToJSON 导出 Agent 为 JSON
func (a *App) ExportAgentToJSON(agentID int64) (string, error) {
	return a.agentService.ExportAgentToJSON(agentID)
}

// ImportAgent 导入 Agent
func (a *App) ImportAgent(exportData *model.AgentExport) (*model.Agent, error) {
	return a.agentService.ImportAgent(exportData)
}

// ImportAgentFromJSON 从 JSON 导入 Agent
func (a *App) ImportAgentFromJSON(jsonData string) (*model.Agent, error) {
	return a.agentService.ImportAgentFromJSON(jsonData)
}
