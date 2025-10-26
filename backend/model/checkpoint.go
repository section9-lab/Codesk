package model

import "time"

// Checkpoint 代表一个会话检查点
type Checkpoint struct {
	ID          string    `json:"id"`           // 检查点 ID
	SessionID   string    `json:"session_id"`   // 会话 ID
	ProjectPath string    `json:"project_path"` // 项目路径
	Name        string    `json:"name"`         // 检查点名称
	Description string    `json:"description"`  // 描述
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	MessageID   *string   `json:"message_id"`   // 关联的消息 ID
	FileStates  []FileState `json:"file_states"` // 文件状态快照
}

// FileState 文件状态快照
type FileState struct {
	Path         string `json:"path"`          // 文件路径
	Content      string `json:"content"`       // 文件内容
	Hash         string `json:"hash"`          // 内容哈希
	LastModified int64  `json:"last_modified"` // 最后修改时间
}

// CheckpointSettings 检查点设置
type CheckpointSettings struct {
	AutoCheckpoint        bool `json:"auto_checkpoint"`         // 是否自动创建检查点
	CheckpointInterval    int  `json:"checkpoint_interval"`     // 检查点间隔（消息数）
	MaxCheckpoints        int  `json:"max_checkpoints"`         // 最大检查点数
	EnableAutoCleanup     bool `json:"enable_auto_cleanup"`     // 是否自动清理
	CleanupRetentionDays  int  `json:"cleanup_retention_days"`  // 保留天数
}

// CheckpointDiff 检查点差异
type CheckpointDiff struct {
	CheckpointID string     `json:"checkpoint_id"`
	Files        []FileDiff `json:"files"`
}

// FileDiff 文件差异
type FileDiff struct {
	Path      string `json:"path"`
	Status    string `json:"status"` // added, modified, deleted
	OldHash   string `json:"old_hash"`
	NewHash   string `json:"new_hash"`
	DiffLines string `json:"diff_lines"` // 差异内容
}

// TimelineEntry 时间线条目
type TimelineEntry struct {
	Type        string    `json:"type"`         // message, checkpoint, file_change
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
	Data        interface{} `json:"data"`
}

// CheckpointStateStats 检查点状态统计
type CheckpointStateStats struct {
	TotalCheckpoints   int       `json:"total_checkpoints"`
	TotalSessions      int       `json:"total_sessions"`
	OldestCheckpoint   *time.Time `json:"oldest_checkpoint"`
	NewestCheckpoint   *time.Time `json:"newest_checkpoint"`
	TotalStorageBytes  int64     `json:"total_storage_bytes"`
}
