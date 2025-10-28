package checkpoint

import (
	"Codesk/backend/config"
	"fmt"
	"log"
	"sync"
)

// CheckpointService 检查点服务（统一接口）
type CheckpointService struct {
	managers map[string]*CheckpointManager // sessionID -> manager
	mu       sync.RWMutex
}

// NewCheckpointService 创建检查点服务实例
func NewCheckpointService() *CheckpointService {
	return &CheckpointService{
		managers: make(map[string]*CheckpointManager),
	}
}

// GetOrCreateManager 获取或创建检查点管理器
func (s *CheckpointService) GetOrCreateManager(projectID, sessionID, projectPath string) (*CheckpointManager, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if manager, exists := s.managers[sessionID]; exists {
		return manager, nil
	}

	claudeDir := config.GetClaudeDir()
	manager, err := NewCheckpointManager(projectID, sessionID, projectPath, claudeDir)
	if err != nil {
		return nil, err
	}

	s.managers[sessionID] = manager
	return manager, nil
}

// CreateCheckpoint 创建检查点
func (s *CheckpointService) CreateCheckpoint(projectID, sessionID, projectPath string, description *string, parentCheckpointID *string) (*CheckpointResult, error) {
	manager, err := s.GetOrCreateManager(projectID, sessionID, projectPath)
	if err != nil {
		return nil, err
	}

	return manager.CreateCheckpoint(description, parentCheckpointID)
}

// RestoreCheckpoint 恢复检查点
func (s *CheckpointService) RestoreCheckpoint(projectID, sessionID, projectPath, checkpointID string) error {
	manager, err := s.GetOrCreateManager(projectID, sessionID, projectPath)
	if err != nil {
		return err
	}

	return manager.RestoreCheckpoint(checkpointID)
}

// ListCheckpoints 列出检查点
func (s *CheckpointService) ListCheckpoints(projectID, sessionID, projectPath string) ([]*Checkpoint, error) {
	manager, err := s.GetOrCreateManager(projectID, sessionID, projectPath)
	if err != nil {
		return nil, err
	}

	return manager.ListCheckpoints()
}

// GetTimeline 获取会话时间线
func (s *CheckpointService) GetTimeline(projectID, sessionID, projectPath string) (*SessionTimeline, error) {
	manager, err := s.GetOrCreateManager(projectID, sessionID, projectPath)
	if err != nil {
		return nil, err
	}

	return manager.GetTimeline(), nil
}

// TrackMessage 跟踪消息
func (s *CheckpointService) TrackMessage(projectID, sessionID, projectPath, jsonlMessage string) error {
	manager, err := s.GetOrCreateManager(projectID, sessionID, projectPath)
	if err != nil {
		return err
	}

	return manager.TrackMessage(jsonlMessage)
}

// TrackFileModification 跟踪文件修改
func (s *CheckpointService) TrackFileModification(projectID, sessionID, projectPath, filePath string) error {
	manager, err := s.GetOrCreateManager(projectID, sessionID, projectPath)
	if err != nil {
		return err
	}

	return manager.TrackFileModification(filePath)
}

// GetCheckpointDiff 获取检查点差异
func (s *CheckpointService) GetCheckpointDiff(projectID, sessionID, fromCheckpointID, toCheckpointID string) (*CheckpointDiff, error) {
	claudeDir := config.GetClaudeDir()
	storage := NewCheckpointStorage(claudeDir)

	// 加载两个检查点
	fromCheckpoint, err := storage.LoadCheckpoint(projectID, sessionID, fromCheckpointID)
	if err != nil {
		return nil, fmt.Errorf("failed to load from checkpoint: %w", err)
	}

	toCheckpoint, err := storage.LoadCheckpoint(projectID, sessionID, toCheckpointID)
	if err != nil {
		return nil, fmt.Errorf("failed to load to checkpoint: %w", err)
	}

	// 加载文件快照
	fromSnapshots, err := storage.LoadFileSnapshots(projectID, sessionID, fromCheckpointID)
	if err != nil {
		return nil, fmt.Errorf("failed to load from snapshots: %w", err)
	}

	toSnapshots, err := storage.LoadFileSnapshots(projectID, sessionID, toCheckpointID)
	if err != nil {
		return nil, fmt.Errorf("failed to load to snapshots: %w", err)
	}

	// 计算差异
	diff := s.calculateDiff(fromSnapshots, toSnapshots)

	// 计算 token 差异
	tokenDelta := toCheckpoint.Metadata.TotalTokens - fromCheckpoint.Metadata.TotalTokens

	return &CheckpointDiff{
		FromCheckpointID: fromCheckpointID,
		ToCheckpointID:   toCheckpointID,
		ModifiedFiles:    diff.ModifiedFiles,
		AddedFiles:       diff.AddedFiles,
		DeletedFiles:     diff.DeletedFiles,
		TokenDelta:       tokenDelta,
	}, nil
}

// calculateDiff 计算文件差异
func (s *CheckpointService) calculateDiff(fromSnapshots, toSnapshots []FileSnapshot) *CheckpointDiff {
	fromMap := make(map[string]FileSnapshot)
	for _, snapshot := range fromSnapshots {
		fromMap[snapshot.FilePath] = snapshot
	}

	toMap := make(map[string]FileSnapshot)
	for _, snapshot := range toSnapshots {
		toMap[snapshot.FilePath] = snapshot
	}

	var modifiedFiles []FileDiff
	var addedFiles []string
	var deletedFiles []string

	// 检查修改和删除
	for path, fromSnapshot := range fromMap {
		if toSnapshot, exists := toMap[path]; exists {
			if fromSnapshot.Hash != toSnapshot.Hash {
				// 文件已修改
				modifiedFiles = append(modifiedFiles, FileDiff{
					Path:      path,
					Additions: 0, // 简化实现，不计算具体行数
					Deletions: 0,
				})
			}
		} else {
			// 文件已删除
			deletedFiles = append(deletedFiles, path)
		}
	}

	// 检查新增
	for path := range toMap {
		if _, exists := fromMap[path]; !exists {
			addedFiles = append(addedFiles, path)
		}
	}

	return &CheckpointDiff{
		ModifiedFiles: modifiedFiles,
		AddedFiles:    addedFiles,
		DeletedFiles:  deletedFiles,
	}
}

// ForkFromCheckpoint 从检查点创建分支
func (s *CheckpointService) ForkFromCheckpoint(projectID, sessionID, projectPath, checkpointID string, description *string) (*CheckpointResult, error) {
	// 先恢复检查点
	if err := s.RestoreCheckpoint(projectID, sessionID, projectPath, checkpointID); err != nil {
		return nil, fmt.Errorf("failed to restore checkpoint for fork: %w", err)
	}

	// 创建新检查点，指定父检查点
	return s.CreateCheckpoint(projectID, sessionID, projectPath, description, &checkpointID)
}

// CleanupOldCheckpoints 清理旧检查点
func (s *CheckpointService) CleanupOldCheckpoints(retentionDays int) (int, error) {
	// 简化实现：遍历所有管理器并清理
	// 实际应该遍历文件系统
	return 0, nil
}

// GetCheckpointSettings 获取检查点设置
func (s *CheckpointService) GetCheckpointSettings(projectID, sessionID, projectPath string) (*CheckpointSettings, error) {
	manager, err := s.GetOrCreateManager(projectID, sessionID, projectPath)
	if err != nil {
		return nil, err
	}

	timeline := manager.GetTimeline()

	return &CheckpointSettings{
		AutoCheckpointEnabled: timeline.AutoCheckpointEnabled,
		CheckpointStrategy:    string(timeline.CheckpointStrategy),
	}, nil
}

// UpdateCheckpointSettings 更新检查点设置
func (s *CheckpointService) UpdateCheckpointSettings(projectID, sessionID, projectPath string, settings *CheckpointSettings) error {
	manager, err := s.GetOrCreateManager(projectID, sessionID, projectPath)
	if err != nil {
		return err
	}

	timeline := manager.GetTimeline()
	timeline.AutoCheckpointEnabled = settings.AutoCheckpointEnabled
	timeline.CheckpointStrategy = CheckpointStrategy(settings.CheckpointStrategy)

	return nil
}

// ClearCheckpointManager 清理指定会话的检查点管理器，释放相关资源
// 对应 Rust 版本的 clear_checkpoint_manager 函数
func (s *CheckpointService) ClearCheckpointManager(sessionID string) error {
	log.Printf("Clearing checkpoint manager for session: %s", sessionID)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.managers[sessionID]; exists {
		delete(s.managers, sessionID)
		log.Printf("Successfully cleared checkpoint manager for session: %s", sessionID)
	} else {
		log.Printf("No checkpoint manager found for session: %s", sessionID)
	}

	return nil
}

// CheckpointSettings 检查点设置
type CheckpointSettings struct {
	AutoCheckpointEnabled bool   `json:"autoCheckpointEnabled"`
	CheckpointStrategy    string `json:"checkpointStrategy"`
}
