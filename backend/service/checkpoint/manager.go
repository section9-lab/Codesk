package checkpoint

import (
	"Codesk/backend/util"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// CheckpointManager 管理会话的检查点操作
type CheckpointManager struct {
	projectID       string
	sessionID       string
	projectPath     string
	fileTracker     *FileTracker
	storage         *CheckpointStorage
	timeline        *SessionTimeline
	currentMessages []string
	mu              sync.RWMutex
}

// NewCheckpointManager 创建新的检查点管理器
func NewCheckpointManager(projectID, sessionID, projectPath, claudeDir string) (*CheckpointManager, error) {
	storage := NewCheckpointStorage(claudeDir)

	// 初始化存储
	if err := storage.InitStorage(projectID, sessionID); err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// 加载或创建时间线
	paths := NewCheckpointPaths(claudeDir, projectID, sessionID)
	timeline, err := storage.LoadTimeline(paths.TimelineFile)
	if err != nil {
		timeline = NewSessionTimeline(sessionID)
	}

	fileTracker := &FileTracker{
		TrackedFiles: make(map[string]*FileState),
	}

	return &CheckpointManager{
		projectID:       projectID,
		sessionID:       sessionID,
		projectPath:     projectPath,
		fileTracker:     fileTracker,
		storage:         storage,
		timeline:        timeline,
		currentMessages: []string{},
	}, nil
}

// TrackMessage 跟踪会话中的新消息
func (m *CheckpointManager) TrackMessage(jsonlMessage string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.currentMessages = append(m.currentMessages, jsonlMessage)

	// 解析消息检查工具使用
	var msg map[string]interface{}
	if err := json.Unmarshal([]byte(jsonlMessage), &msg); err != nil {
		return nil // 忽略解析错误
	}

	if message, ok := msg["message"].(map[string]interface{}); ok {
		if content, ok := message["content"].([]interface{}); ok {
			for _, item := range content {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if itemType, ok := itemMap["type"].(string); ok && itemType == "tool_use" {
						if toolName, ok := itemMap["name"].(string); ok {
							if input, ok := itemMap["input"].(map[string]interface{}); ok {
								m.trackToolOperation(toolName, input)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

// trackToolOperation 跟踪工具操作的文件修改
func (m *CheckpointManager) trackToolOperation(tool string, input map[string]interface{}) {
	toolLower := strings.ToLower(tool)

	switch toolLower {
	case "edit", "write", "multiedit":
		if filePath, ok := input["file_path"].(string); ok {
			m.TrackFileModification(filePath)
		}
	case "bash":
		if command, ok := input["command"].(string); ok {
			m.trackBashSideEffects(command)
		}
	}
}

// TrackFileModification 跟踪文件修改
func (m *CheckpointManager) TrackFileModification(filePath string) error {
	fullPath := filepath.Join(m.projectPath, filePath)

	// 读取当前文件状态
	var hash string
	var exists bool
	var modified time.Time

	if util.FileExists(fullPath) {
		content, err := util.ReadFile(fullPath)
		if err == nil {
			hash = m.storage.CalculateFileHash(content)
			exists = true

			info, err := os.Stat(fullPath)
			if err == nil {
				modified = info.ModTime()
			} else {
				modified = time.Now()
			}
		}
	} else {
		hash = ""
		exists = false
		modified = time.Now()
	}

	// 检查文件是否真的改变了
	isModified := true
	if existingState, ok := m.fileTracker.TrackedFiles[filePath]; ok {
		isModified = existingState.LastHash != hash ||
			existingState.Exists != exists ||
			existingState.IsModified
	}

	m.fileTracker.TrackedFiles[filePath] = &FileState{
		LastHash:     hash,
		IsModified:   isModified,
		LastModified: modified,
		Exists:       exists,
	}

	return nil
}

// trackBashSideEffects 跟踪 bash 命令的潜在文件更改
func (m *CheckpointManager) trackBashSideEffects(command string) {
	// 常见的文件修改命令
	fileCommands := []string{
		"echo", "cat", "cp", "mv", "rm", "touch", "sed", "awk",
		"npm", "yarn", "pnpm", "bun", "cargo", "make", "gcc", "g++",
	}

	// 简单启发式：如果命令包含文件修改操作
	for _, cmd := range fileCommands {
		if strings.Contains(command, cmd) {
			// 标记所有跟踪的文件为可能已修改
			for _, state := range m.fileTracker.TrackedFiles {
				state.IsModified = true
			}
			break
		}
	}
}

// CreateCheckpoint 创建检查点
func (m *CheckpointManager) CreateCheckpoint(description *string, parentCheckpointID *string) (*CheckpointResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	messageIndex := len(m.currentMessages) - 1
	if messageIndex < 0 {
		messageIndex = 0
	}

	// 从最后的用户消息中提取元数据
	userPrompt, modelUsed, totalTokens := m.extractCheckpointMetadata()

	// 确保项目中的每个文件都被跟踪
	if err := m.scanProjectFiles(); err != nil {
		return nil, fmt.Errorf("failed to scan project files: %w", err)
	}

	// 收集修改的文件
	var fileSnapshots []FileSnapshot
	var snapshotSize int64

	for filePath, state := range m.fileTracker.TrackedFiles {
		if !state.Exists {
			// 文件已删除
			fileSnapshots = append(fileSnapshots, FileSnapshot{
				CheckpointID: "", // 稍后设置
				FilePath:     filePath,
				Content:      "",
				Hash:         state.LastHash,
				IsDeleted:    true,
				Size:         0,
			})
			continue
		}

		fullPath := filepath.Join(m.projectPath, filePath)
		content, err := util.ReadFile(fullPath)
		if err != nil {
			continue
		}

		hash := m.storage.CalculateFileHash(content)
		info, _ := os.Stat(fullPath)

		var perms *uint32
		if info != nil {
			p := uint32(info.Mode().Perm())
			perms = &p
		}

		snapshot := FileSnapshot{
			CheckpointID: "", // 稍后设置
			FilePath:     filePath,
			Content:      content,
			Hash:         hash,
			IsDeleted:    false,
			Permissions:  perms,
			Size:         int64(len(content)),
		}

		fileSnapshots = append(fileSnapshots, snapshot)
		snapshotSize += snapshot.Size

		// 更新跟踪状态
		state.LastHash = hash
		state.IsModified = false
	}

	// 创建检查点
	checkpointID := uuid.New().String()
	checkpoint := &Checkpoint{
		ID:                 checkpointID,
		SessionID:          m.sessionID,
		ProjectID:          m.projectID,
		MessageIndex:       messageIndex,
		Timestamp:          time.Now(),
		Description:        description,
		ParentCheckpointID: parentCheckpointID,
		Metadata: CheckpointMetadata{
			TotalTokens:  totalTokens,
			ModelUsed:    modelUsed,
			UserPrompt:   userPrompt,
			FileChanges:  len(fileSnapshots),
			SnapshotSize: snapshotSize,
		},
	}

	// 设置快照的检查点 ID
	for i := range fileSnapshots {
		fileSnapshots[i].CheckpointID = checkpointID
	}

	// 合并所有消息
	messages := strings.Join(m.currentMessages, "\n")

	// 保存检查点
	result, err := m.storage.SaveCheckpoint(m.projectID, m.sessionID, checkpoint, fileSnapshots, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to save checkpoint: %w", err)
	}

	// 更新时间线
	m.timeline.CurrentCheckpointID = &checkpointID
	m.timeline.TotalCheckpoints++

	return result, nil
}

// extractCheckpointMetadata 从消息中提取元数据
func (m *CheckpointManager) extractCheckpointMetadata() (userPrompt string, modelUsed string, totalTokens int64) {
	userPrompt = ""
	modelUsed = "sonnet"
	totalTokens = 0

	// 从最后的消息中提取
	for i := len(m.currentMessages) - 1; i >= 0; i-- {
		var msg map[string]interface{}
		if err := json.Unmarshal([]byte(m.currentMessages[i]), &msg); err != nil {
			continue
		}

		if message, ok := msg["message"].(map[string]interface{}); ok {
			// 提取模型
			if model, ok := message["model"].(string); ok {
				modelUsed = model
			}

			// 提取 tokens
			if usage, ok := message["usage"].(map[string]interface{}); ok {
				if inputTokens, ok := usage["input_tokens"].(float64); ok {
					totalTokens += int64(inputTokens)
				}
				if outputTokens, ok := usage["output_tokens"].(float64); ok {
					totalTokens += int64(outputTokens)
				}
			}

			// 提取用户提示
			if role, ok := message["role"].(string); ok && role == "user" {
				if content, ok := message["content"].(string); ok {
					userPrompt = content
					break
				}
			}
		}
	}

	return
}

// scanProjectFiles 扫描项目文件
func (m *CheckpointManager) scanProjectFiles() error {
	return filepath.Walk(m.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过错误
		}

		if info.IsDir() {
			// 跳过隐藏目录和常见的忽略目录
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "target" {
				return filepath.SkipDir
			}
			return nil
		}

		// 跳过隐藏文件和大文件
		if strings.HasPrefix(info.Name(), ".") || info.Size() > 1024*1024 {
			return nil
		}

		relPath, err := filepath.Rel(m.projectPath, path)
		if err != nil {
			return nil
		}

		// 如果文件还未被跟踪，添加它
		if _, exists := m.fileTracker.TrackedFiles[relPath]; !exists {
			m.TrackFileModification(relPath)
		}

		return nil
	})
}

// RestoreCheckpoint 恢复检查点
func (m *CheckpointManager) RestoreCheckpoint(checkpointID string) error {
	// 加载检查点
	checkpoint, err := m.storage.LoadCheckpoint(m.projectID, m.sessionID, checkpointID)
	if err != nil {
		return fmt.Errorf("failed to load checkpoint: %w", err)
	}

	// 加载文件快照
	snapshots, err := m.storage.LoadFileSnapshots(m.projectID, m.sessionID, checkpointID)
	if err != nil {
		return fmt.Errorf("failed to load file snapshots: %w", err)
	}

	// 恢复文件
	for _, snapshot := range snapshots {
		fullPath := filepath.Join(m.projectPath, snapshot.FilePath)

		if snapshot.IsDeleted {
			// 删除文件
			if util.FileExists(fullPath) {
				if err := os.Remove(fullPath); err != nil {
					return fmt.Errorf("failed to delete file %s: %w", snapshot.FilePath, err)
				}
			}
		} else {
			// 恢复文件内容
			if err := util.WriteFile(fullPath, snapshot.Content); err != nil {
				return fmt.Errorf("failed to restore file %s: %w", snapshot.FilePath, err)
			}

			// 恢复权限
			if snapshot.Permissions != nil {
				if err := os.Chmod(fullPath, os.FileMode(*snapshot.Permissions)); err != nil {
					// 忽略权限错误
				}
			}
		}
	}

	// 更新当前检查点
	m.timeline.CurrentCheckpointID = &checkpoint.ID

	return nil
}

// ListCheckpoints 列出所有检查点
func (m *CheckpointManager) ListCheckpoints() ([]*Checkpoint, error) {
	if m.timeline.RootNode == nil {
		return []*Checkpoint{}, nil
	}

	var checkpoints []*Checkpoint
	m.collectCheckpoints(m.timeline.RootNode, &checkpoints)

	return checkpoints, nil
}

// collectCheckpoints 递归收集检查点
func (m *CheckpointManager) collectCheckpoints(node *TimelineNode, checkpoints *[]*Checkpoint) {
	*checkpoints = append(*checkpoints, &node.Checkpoint)

	for i := range node.Children {
		m.collectCheckpoints(&node.Children[i], checkpoints)
	}
}

// GetTimeline 获取时间线
func (m *CheckpointManager) GetTimeline() *SessionTimeline {
	return m.timeline
}
