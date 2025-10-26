package main

import (
	"Codesk/backend/model"
	"Codesk/backend/service/claude"
)

// ==================== Claude Project Management ====================

// ListProjects 列出所有 Claude 项目
func (a *App) ListProjects() ([]*model.Project, error) {
	return a.claudeProjectService.ListProjects()
}

// GetProjectSessions 获取项目的会话列表
func (a *App) GetProjectSessions(projectPath string) ([]*model.Session, error) {
	return a.claudeProjectService.GetProjectSessions(projectPath)
}

// GetClaudeSettings 获取 Claude 设置
func (a *App) GetClaudeSettings() (*model.ClaudeSettings, error) {
	return a.claudeProjectService.GetClaudeSettings()
}

// SaveClaudeSettings 保存 Claude 设置
func (a *App) SaveClaudeSettings(settings *model.ClaudeSettings) error {
	return a.claudeProjectService.SaveClaudeSettings(settings)
}

// LoadSessionHistory 加载会话历史
func (a *App) LoadSessionHistory(sessionID, projectPath string) (*model.SessionHistory, error) {
	return a.claudeProjectService.LoadSessionHistory(sessionID, projectPath)
}

// ==================== Claude File Operations ====================

// ListDirectoryContents 列出目录内容
func (a *App) ListDirectoryContents(path string) ([]claude.FileEntry, error) {
	return a.claudeFileService.ListDirectoryContents(path)
}

// SearchFiles 搜索文件
func (a *App) SearchFiles(rootPath, pattern string, maxResults int) ([]string, error) {
	return a.claudeFileService.SearchFiles(rootPath, pattern, maxResults)
}

// GetRecentlyModifiedFiles 获取最近修改的文件
func (a *App) GetRecentlyModifiedFiles(rootPath string, limit int) ([]string, error) {
	return a.claudeFileService.GetRecentlyModifiedFiles(rootPath, limit)
}

// FindClaudeMdFiles 查找所有 CLAUDE.md 文件
func (a *App) FindClaudeMdFiles(rootPath string) ([]claude.ClaudeMdFile, error) {
	return a.claudeFileService.FindClaudeMdFiles(rootPath)
}

// ReadClaudeMdFile 读取 CLAUDE.md 文件
func (a *App) ReadClaudeMdFile(path string) (string, error) {
	return a.claudeFileService.ReadClaudeMdFile(path)
}

// SaveClaudeMdFile 保存 CLAUDE.md 文件
func (a *App) SaveClaudeMdFile(path, content string) error {
	return a.claudeFileService.SaveClaudeMdFile(path, content)
}

// ==================== Claude Execution ====================

// ExecuteClaude 执行 Claude Code
func (a *App) ExecuteClaude(projectPath, task, model string) (*claude.ExecuteResult, error) {
	opts := claude.ExecuteOptions{
		ProjectPath: projectPath,
		Task:        task,
		Model:       model,
	}
	return a.claudeExecutionService.Execute(opts)
}

// ContinueClaude 继续 Claude 执行
func (a *App) ContinueClaude(projectPath, sessionID, model string) (*claude.ExecuteResult, error) {
	opts := claude.ExecuteOptions{
		ProjectPath: projectPath,
		Model:       model,
		SessionID:   sessionID,
		Continue:    true,
	}
	return a.claudeExecutionService.Execute(opts)
}

// ResumeClaude 恢复 Claude 执行
func (a *App) ResumeClaude(projectPath, sessionID, model string) (*claude.ExecuteResult, error) {
	opts := claude.ExecuteOptions{
		ProjectPath: projectPath,
		Model:       model,
		SessionID:   sessionID,
		Resume:      true,
	}
	return a.claudeExecutionService.Execute(opts)
}

// CancelClaude 取消 Claude 执行
func (a *App) CancelClaude(sessionID string) error {
	return a.claudeExecutionService.Cancel(sessionID)
}

// GetClaudeSessionOutput 获取会话输出
func (a *App) GetClaudeSessionOutput(sessionID, projectPath string) (string, error) {
	return a.claudeExecutionService.GetSessionOutput(sessionID, projectPath)
}

// GetClaudeSessionStatus 获取会话状态
func (a *App) GetClaudeSessionStatus(sessionID, projectPath string) (string, error) {
	return a.claudeExecutionService.GetSessionStatus(sessionID, projectPath)
}

// ListRunningClaudeSessions 列出运行中的会话
func (a *App) ListRunningClaudeSessions() []map[string]interface{} {
	return a.claudeExecutionService.ListRunningSessions()
}

// CleanupFinishedProcesses 清理已完成的进程
func (a *App) CleanupFinishedProcesses() int {
	return a.claudeExecutionService.CleanupFinishedProcesses()
}

// CheckClaudeVersion 检查 Claude Code 版本
func (a *App) CheckClaudeVersion() (bool, string, error) {
	return a.claudeExecutionService.CheckClaudeVersion()
}

// ListClaudeInstallations 列出 Claude 安装位置
func (a *App) ListClaudeInstallations() []string {
	return a.claudeExecutionService.ListClaudeInstallations()
}

// GetHomeDirectory 获取用户主目录
func (a *App) GetHomeDirectory() (string, error) {
	return a.claudeExecutionService.GetHomeDirectory()
}

// CreateProject 创建新项目
func (a *App) CreateProject(projectPath string) error {
	return a.claudeExecutionService.CreateProject(projectPath)
}

// OpenNewSession 打开新会话
func (a *App) OpenNewSession() string {
	return a.claudeExecutionService.OpenNewSession()
}
