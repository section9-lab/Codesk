package main

import "Codesk/backend/service/checkpoint"

// ==================== Checkpoint Management ====================

// CreateCheckpoint 创建检查点
func (a *App) CreateCheckpoint(projectID, sessionID, projectPath string, description *string, parentCheckpointID *string) (*checkpoint.CheckpointResult, error) {
	return a.checkpointService.CreateCheckpoint(projectID, sessionID, projectPath, description, parentCheckpointID)
}

// RestoreCheckpoint 恢复检查点
func (a *App) RestoreCheckpoint(projectID, sessionID, projectPath, checkpointID string) error {
	return a.checkpointService.RestoreCheckpoint(projectID, sessionID, projectPath, checkpointID)
}

// ListCheckpoints 列出检查点
func (a *App) ListCheckpoints(projectID, sessionID, projectPath string) ([]*checkpoint.Checkpoint, error) {
	return a.checkpointService.ListCheckpoints(projectID, sessionID, projectPath)
}

// GetSessionTimeline 获取会话时间线
func (a *App) GetSessionTimeline(projectID, sessionID, projectPath string) (*checkpoint.SessionTimeline, error) {
	return a.checkpointService.GetTimeline(projectID, sessionID, projectPath)
}

// TrackCheckpointMessage 跟踪消息
func (a *App) TrackCheckpointMessage(projectID, sessionID, projectPath, jsonlMessage string) error {
	return a.checkpointService.TrackMessage(projectID, sessionID, projectPath, jsonlMessage)
}

// TrackFileModification 跟踪文件修改
func (a *App) TrackFileModification(projectID, sessionID, projectPath, filePath string) error {
	return a.checkpointService.TrackFileModification(projectID, sessionID, projectPath, filePath)
}

// GetCheckpointDiff 获取检查点差异
func (a *App) GetCheckpointDiff(projectID, sessionID, fromCheckpointID, toCheckpointID string) (*checkpoint.CheckpointDiff, error) {
	return a.checkpointService.GetCheckpointDiff(projectID, sessionID, fromCheckpointID, toCheckpointID)
}

// ForkFromCheckpoint 从检查点创建分支
func (a *App) ForkFromCheckpoint(projectID, sessionID, projectPath, checkpointID string, description *string) (*checkpoint.CheckpointResult, error) {
	return a.checkpointService.ForkFromCheckpoint(projectID, sessionID, projectPath, checkpointID, description)
}

// GetCheckpointSettings 获取检查点设置
func (a *App) GetCheckpointSettings(projectID, sessionID, projectPath string) (*checkpoint.CheckpointSettings, error) {
	return a.checkpointService.GetCheckpointSettings(projectID, sessionID, projectPath)
}

// UpdateCheckpointSettings 更新检查点设置
func (a *App) UpdateCheckpointSettings(projectID, sessionID, projectPath string, settings *checkpoint.CheckpointSettings) error {
	return a.checkpointService.UpdateCheckpointSettings(projectID, sessionID, projectPath, settings)
}

// CleanupOldCheckpoints 清理旧检查点
func (a *App) CleanupOldCheckpoints(retentionDays int) (int, error) {
	return a.checkpointService.CleanupOldCheckpoints(retentionDays)
}
