package checkpoint

import (
	"Codesk/backend/util"
	"compress/zlib"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CheckpointStorage 管理检查点存储操作
type CheckpointStorage struct {
	claudeDir        string
	compressionLevel int
}

// NewCheckpointStorage 创建检查点存储实例
func NewCheckpointStorage(claudeDir string) *CheckpointStorage {
	return &CheckpointStorage{
		claudeDir:        claudeDir,
		compressionLevel: 6, // zlib 默认压缩级别
	}
}

// InitStorage 初始化会话的检查点存储
func (s *CheckpointStorage) InitStorage(projectID, sessionID string) error {
	paths := NewCheckpointPaths(s.claudeDir, projectID, sessionID)

	// 创建目录结构
	if err := os.MkdirAll(paths.CheckpointsDir, 0755); err != nil {
		return fmt.Errorf("failed to create checkpoints directory: %w", err)
	}

	if err := os.MkdirAll(paths.FilesDir, 0755); err != nil {
		return fmt.Errorf("failed to create files directory: %w", err)
	}

	// 初始化空时间线（如果不存在）
	if !util.FileExists(paths.TimelineFile) {
		timeline := NewSessionTimeline(sessionID)
		if err := s.SaveTimeline(paths.TimelineFile, timeline); err != nil {
			return err
		}
	}

	return nil
}

// SaveCheckpoint 保存检查点到磁盘
func (s *CheckpointStorage) SaveCheckpoint(
	projectID, sessionID string,
	checkpoint *Checkpoint,
	fileSnapshots []FileSnapshot,
	messages string,
) (*CheckpointResult, error) {
	paths := NewCheckpointPaths(s.claudeDir, projectID, sessionID)
	checkpointDir := paths.CheckpointDir(checkpoint.ID)

	// 创建检查点目录
	if err := os.MkdirAll(checkpointDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create checkpoint directory: %w", err)
	}

	// 保存检查点元数据
	metadataPath := paths.CheckpointMetadataFile(checkpoint.ID)
	metadataJSON, err := json.MarshalIndent(checkpoint, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to serialize checkpoint metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, metadataJSON, 0644); err != nil {
		return nil, fmt.Errorf("failed to write checkpoint metadata: %w", err)
	}

	// 保存消息（压缩）
	messagesPath := paths.CheckpointMessagesFile(checkpoint.ID)
	compressedMessages, err := s.compressData([]byte(messages))
	if err != nil {
		return nil, fmt.Errorf("failed to compress messages: %w", err)
	}

	if err := os.WriteFile(messagesPath, compressedMessages, 0644); err != nil {
		return nil, fmt.Errorf("failed to write compressed messages: %w", err)
	}

	// 保存文件快照
	var warnings []string
	filesProcessed := 0

	for _, snapshot := range fileSnapshots {
		if err := s.saveFileSnapshot(paths, &snapshot); err != nil {
			warnings = append(warnings, fmt.Sprintf("Failed to save %s: %v", snapshot.FilePath, err))
		} else {
			filesProcessed++
		}
	}

	// 更新时间线
	if err := s.updateTimelineWithCheckpoint(paths.TimelineFile, checkpoint, fileSnapshots); err != nil {
		return nil, fmt.Errorf("failed to update timeline: %w", err)
	}

	return &CheckpointResult{
		Checkpoint:     *checkpoint,
		FilesProcessed: filesProcessed,
		Warnings:       warnings,
	}, nil
}

// saveFileSnapshot 保存单个文件快照（内容寻址存储）
func (s *CheckpointStorage) saveFileSnapshot(paths *CheckpointPaths, snapshot *FileSnapshot) error {
	// 使用内容寻址存储：按哈希存储文件
	// 这可以防止跨检查点重复相同的文件内容
	contentPoolDir := filepath.Join(paths.FilesDir, "content_pool")
	if err := os.MkdirAll(contentPoolDir, 0755); err != nil {
		return fmt.Errorf("failed to create content pool directory: %w", err)
	}

	// 在内容池中存储实际内容
	contentFile := filepath.Join(contentPoolDir, snapshot.Hash)

	// 只有当内容不存在时才写入
	if !util.FileExists(contentFile) {
		// 压缩并保存文件内容
		compressedContent, err := s.compressData([]byte(snapshot.Content))
		if err != nil {
			return fmt.Errorf("failed to compress file content: %w", err)
		}

		if err := os.WriteFile(contentFile, compressedContent, 0644); err != nil {
			return fmt.Errorf("failed to write file content to pool: %w", err)
		}
	}

	// 在检查点特定目录中创建引用
	checkpointRefsDir := filepath.Join(paths.FilesDir, "refs", snapshot.CheckpointID)
	if err := os.MkdirAll(checkpointRefsDir, 0755); err != nil {
		return fmt.Errorf("failed to create checkpoint refs directory: %w", err)
	}

	// 保存文件元数据和内容引用
	refMetadata := map[string]interface{}{
		"path":        snapshot.FilePath,
		"hash":        snapshot.Hash,
		"is_deleted":  snapshot.IsDeleted,
		"permissions": snapshot.Permissions,
		"size":        snapshot.Size,
	}

	// 使用安全的文件名作为引用
	safeFilename := strings.ReplaceAll(snapshot.FilePath, "/", "_")
	safeFilename = strings.ReplaceAll(safeFilename, "\\", "_")
	refPath := filepath.Join(checkpointRefsDir, safeFilename+".json")

	refJSON, err := json.MarshalIndent(refMetadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize file reference: %w", err)
	}

	if err := os.WriteFile(refPath, refJSON, 0644); err != nil {
		return fmt.Errorf("failed to write file reference: %w", err)
	}

	return nil
}

// LoadCheckpoint 从磁盘加载检查点
func (s *CheckpointStorage) LoadCheckpoint(projectID, sessionID, checkpointID string) (*Checkpoint, error) {
	paths := NewCheckpointPaths(s.claudeDir, projectID, sessionID)
	metadataPath := paths.CheckpointMetadataFile(checkpointID)

	if !util.FileExists(metadataPath) {
		return nil, fmt.Errorf("checkpoint not found: %s", checkpointID)
	}

	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read checkpoint metadata: %w", err)
	}

	var checkpoint Checkpoint
	if err := json.Unmarshal(data, &checkpoint); err != nil {
		return nil, fmt.Errorf("failed to parse checkpoint metadata: %w", err)
	}

	return &checkpoint, nil
}

// LoadFileSnapshots 加载检查点的文件快照
func (s *CheckpointStorage) LoadFileSnapshots(projectID, sessionID, checkpointID string) ([]FileSnapshot, error) {
	paths := NewCheckpointPaths(s.claudeDir, projectID, sessionID)
	refsDir := filepath.Join(paths.FilesDir, "refs", checkpointID)

	if !util.DirExists(refsDir) {
		return []FileSnapshot{}, nil
	}

	entries, err := os.ReadDir(refsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read refs directory: %w", err)
	}

	var snapshots []FileSnapshot
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		refPath := filepath.Join(refsDir, entry.Name())
		data, err := os.ReadFile(refPath)
		if err != nil {
			continue
		}

		var refMetadata map[string]interface{}
		if err := json.Unmarshal(data, &refMetadata); err != nil {
			continue
		}

		// 从内容池加载实际内容
		hash := refMetadata["hash"].(string)
		content, err := s.loadFileContent(paths, hash)
		if err != nil {
			continue
		}

		snapshot := FileSnapshot{
			CheckpointID: checkpointID,
			FilePath:     refMetadata["path"].(string),
			Content:      content,
			Hash:         hash,
			IsDeleted:    refMetadata["is_deleted"].(bool),
			Size:         int64(refMetadata["size"].(float64)),
		}

		if perms, ok := refMetadata["permissions"].(float64); ok {
			p := uint32(perms)
			snapshot.Permissions = &p
		}

		snapshots = append(snapshots, snapshot)
	}

	return snapshots, nil
}

// loadFileContent 从内容池加载文件内容
func (s *CheckpointStorage) loadFileContent(paths *CheckpointPaths, hash string) (string, error) {
	contentFile := paths.FileSnapshotPath(hash)

	if !util.FileExists(contentFile) {
		return "", fmt.Errorf("content file not found: %s", hash)
	}

	compressedData, err := os.ReadFile(contentFile)
	if err != nil {
		return "", fmt.Errorf("failed to read content file: %w", err)
	}

	decompressedData, err := s.decompressData(compressedData)
	if err != nil {
		return "", fmt.Errorf("failed to decompress content: %w", err)
	}

	return string(decompressedData), nil
}

// SaveTimeline 保存时间线
func (s *CheckpointStorage) SaveTimeline(timelineFile string, timeline *SessionTimeline) error {
	data, err := json.MarshalIndent(timeline, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize timeline: %w", err)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(timelineFile), 0755); err != nil {
		return fmt.Errorf("failed to create timeline directory: %w", err)
	}

	if err := os.WriteFile(timelineFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write timeline: %w", err)
	}

	return nil
}

// LoadTimeline 加载时间线
func (s *CheckpointStorage) LoadTimeline(timelineFile string) (*SessionTimeline, error) {
	if !util.FileExists(timelineFile) {
		return nil, fmt.Errorf("timeline file not found")
	}

	data, err := os.ReadFile(timelineFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read timeline: %w", err)
	}

	var timeline SessionTimeline
	if err := json.Unmarshal(data, &timeline); err != nil {
		return nil, fmt.Errorf("failed to parse timeline: %w", err)
	}

	return &timeline, nil
}

// updateTimelineWithCheckpoint 更新时间线添加新检查点
func (s *CheckpointStorage) updateTimelineWithCheckpoint(timelineFile string, checkpoint *Checkpoint, snapshots []FileSnapshot) error {
	timeline, err := s.LoadTimeline(timelineFile)
	if err != nil {
		// 如果时间线不存在，创建新的
		timeline = NewSessionTimeline(checkpoint.SessionID)
	}

	// 创建新节点
	snapshotIDs := make([]string, len(snapshots))
	for i, snapshot := range snapshots {
		snapshotIDs[i] = snapshot.Hash
	}

	newNode := TimelineNode{
		Checkpoint:      *checkpoint,
		Children:        []TimelineNode{},
		FileSnapshotIDs: snapshotIDs,
	}

	// 添加到时间线
	if timeline.RootNode == nil {
		// 第一个检查点
		timeline.RootNode = &newNode
	} else if checkpoint.ParentCheckpointID != nil {
		// 找到父节点并添加为子节点
		if parent := timeline.FindCheckpoint(*checkpoint.ParentCheckpointID); parent != nil {
			parent.Children = append(parent.Children, newNode)
		} else {
			// 如果找不到父节点，添加为根节点的子节点
			timeline.RootNode.Children = append(timeline.RootNode.Children, newNode)
		}
	} else {
		// 添加为根节点的子节点
		timeline.RootNode.Children = append(timeline.RootNode.Children, newNode)
	}

	timeline.CurrentCheckpointID = &checkpoint.ID
	timeline.TotalCheckpoints++

	return s.SaveTimeline(timelineFile, timeline)
}

// CalculateFileHash 计算文件哈希
func (s *CheckpointStorage) CalculateFileHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

// compressData 压缩数据
func (s *CheckpointStorage) compressData(data []byte) ([]byte, error) {
	var buf strings.Builder
	writer := zlib.NewWriter(&buf)

	if _, err := writer.Write(data); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}

// decompressData 解压数据
func (s *CheckpointStorage) decompressData(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}
