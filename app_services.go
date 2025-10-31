package main

import (
	"Codesk/backend/model"
)

// ==================== Usage Statistics ====================

// GetUsageStats 获取总体使用统计
func (a *App) GetUsageStats() (*model.UsageStats, error) {
	return a.usageService.GetUsageStats()
}

// GetUsageByDateRange 获取日期范围内的使用统计
func (a *App) GetUsageByDateRange(startDate, endDate string) (*model.UsageByDateRange, error) {
	return a.usageService.GetUsageByDateRange(startDate, endDate)
}

// GetSessionStats 获取会话统计
func (a *App) GetSessionStats(sessionID, projectPath string) (*model.SessionStats, error) {
	return a.usageService.GetSessionStats(sessionID, projectPath)
}

// ==================== MCP Management ====================

// MCPListServers 列出所有 MCP 服务器
func (a *App) MCPListServers() (map[string]model.MCPServer, error) {
	return a.mcpService.ListServers()
}

// MCPGetServer 获取指定的 MCP 服务器配置
func (a *App) MCPGetServer(name string) (*model.MCPServer, error) {
	return a.mcpService.GetServer(name)
}

// MCPAddServer 添加 MCP 服务器
func (a *App) MCPAddServer(name string, server model.MCPServer) error {
	return a.mcpService.AddServer(name, server)
}

// MCPRemoveServer 移除 MCP 服务器
func (a *App) MCPRemoveServer(name string) error {
	return a.mcpService.RemoveServer(name)
}

// MCPAddServerFromJSON 从 JSON 添加服务器
func (a *App) MCPAddServerFromJSON(name, jsonData string) error {
	return a.mcpService.AddServerFromJSON(name, jsonData)
}

// MCPReadProjectConfig 读取项目级 MCP 配置
func (a *App) MCPReadProjectConfig(projectPath string) (*model.MCPProjectConfig, error) {
	return a.mcpService.ReadProjectConfig(projectPath)
}

// MCPSaveProjectConfig 保存项目级 MCP 配置
func (a *App) MCPSaveProjectConfig(config *model.MCPProjectConfig) error {
	return a.mcpService.SaveProjectConfig(config)
}

// MCPResetProjectChoices 重置项目的 MCP 选择
func (a *App) MCPResetProjectChoices(projectPath string) error {
	return a.mcpService.ResetProjectChoices(projectPath)
}

// MCPGetServerStatus 获取服务器状态
func (a *App) MCPGetServerStatus(name string) (*model.MCPServerStatus, error) {
	return a.mcpService.GetServerStatus(name)
}

// MCPTestConnection 测试 MCP 服务器连接
func (a *App) MCPTestConnection(name string) (bool, error) {
	return a.mcpService.TestConnection(name)
}

// ==================== Proxy Settings ====================

// GetProxySettings 获取代理设置
func (a *App) GetProxySettings() (*model.ProxySettings, error) {
	return a.proxyService.GetProxySettings()
}

// SaveProxySettings 保存代理设置
func (a *App) SaveProxySettings(settings *model.ProxySettings) error {
	return a.proxyService.SaveProxySettings(settings)
}

// ==================== Slash Commands ====================

// ListSlashCommands 列出所有斜杠命令
func (a *App) ListSlashCommands() ([]*model.SlashCommand, error) {
	return a.slashService.ListCommands()
}

// GetSlashCommand 获取斜杠命令
func (a *App) GetSlashCommand(name string) (*model.SlashCommand, error) {
	return a.slashService.GetCommand(name)
}

// SaveSlashCommand 保存斜杠命令
func (a *App) SaveSlashCommand(cmd *model.SlashCommand) error {
	return a.slashService.SaveCommand(cmd)
}

// DeleteSlashCommand 删除斜杠命令
func (a *App) DeleteSlashCommand(name string) error {
	return a.slashService.DeleteCommand(name)
}

// ==================== Storage Management ====================

// StorageListTables 列出所有表
func (a *App) StorageListTables() ([]model.TableInfo, error) {
	return a.storageService.ListTables()
}

// StorageReadTable 读取表数据
func (a *App) StorageReadTable(tableName string, limit, offset int) ([]model.TableRow, error) {
	return a.storageService.ReadTable(tableName, limit, offset)
}

// StorageUpdateRow 更新表行
func (a *App) StorageUpdateRow(tableName string, id int64, data model.TableRow) error {
	return a.storageService.UpdateRow(tableName, id, data)
}

// StorageDeleteRow 删除表行
func (a *App) StorageDeleteRow(tableName string, id int64) error {
	return a.storageService.DeleteRow(tableName, id)
}

// StorageInsertRow 插入表行
func (a *App) StorageInsertRow(tableName string, data model.TableRow) error {
	return a.storageService.InsertRow(tableName, data)
}

// StorageExecuteSQL 执行 SQL 语句
func (a *App) StorageExecuteSQL(query string, args ...interface{}) (*model.SQLResult, error) {
	return a.storageService.ExecuteSQL(query, args...)
}

// StorageResetDatabase 重置数据库
func (a *App) StorageResetDatabase() error {
	return a.storageService.ResetDatabase()
}

// SetSetting 设置应用设置
func (a *App) SetSetting(key, value string) error {
	return a.storageService.SetSetting(key, value)
}
