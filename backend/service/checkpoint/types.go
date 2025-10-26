package checkpoint

import (
	"path/filepath"
	"time"
)

// Checkpoint 代表会话时间线中的一个检查点
type Checkpoint struct {
	ID                  string              `json:"id"`
	SessionID           string              `json:"sessionId"`
	ProjectID           string              `json:"projectId"`
	MessageIndex        int                 `json:"messageIndex"`
	Timestamp           time.Time           `json:"timestamp"`
	Description         *string             `json:"description"`
	ParentCheckpointID  *string             `json:"parentCheckpointId"`
	Metadata            CheckpointMetadata  `json:"metadata"`
}

// CheckpointMetadata 检查点元数据
type CheckpointMetadata struct {
	TotalTokens  int64  `json:"totalTokens"`
	ModelUsed    string `json:"modelUsed"`
	UserPrompt   string `json:"userPrompt"`
	FileChanges  int    `json:"fileChanges"`
	SnapshotSize int64  `json:"snapshotSize"`
}

// FileSnapshot 文件快照
type FileSnapshot struct {
	CheckpointID string `json:"checkpointId"`
	FilePath     string `json:"filePath"`
	Content      string `json:"content"`
	Hash         string `json:"hash"`
	IsDeleted    bool   `json:"isDeleted"`
	Permissions  *uint32 `json:"permissions"`
	Size         int64  `json:"size"`
}

// TimelineNode 时间线树节点
type TimelineNode struct {
	Checkpoint       Checkpoint     `json:"checkpoint"`
	Children         []TimelineNode `json:"children"`
	FileSnapshotIDs  []string       `json:"fileSnapshotIds"`
}

// SessionTimeline 会话的完整时间线
type SessionTimeline struct {
	SessionID             string             `json:"sessionId"`
	RootNode              *TimelineNode      `json:"rootNode"`
	CurrentCheckpointID   *string            `json:"currentCheckpointId"`
	AutoCheckpointEnabled bool               `json:"autoCheckpointEnabled"`
	CheckpointStrategy    CheckpointStrategy `json:"checkpointStrategy"`
	TotalCheckpoints      int                `json:"totalCheckpoints"`
}

// CheckpointStrategy 自动检查点创建策略
type CheckpointStrategy string

const (
	StrategyManual     CheckpointStrategy = "manual"
	StrategyPerPrompt  CheckpointStrategy = "per_prompt"
	StrategyPerToolUse CheckpointStrategy = "per_tool_use"
	StrategySmart      CheckpointStrategy = "smart"
)

// FileTracker 跟踪文件状态
type FileTracker struct {
	TrackedFiles map[string]*FileState
}

// FileState 跟踪的文件状态
type FileState struct {
	LastHash     string
	IsModified   bool
	LastModified time.Time
	Exists       bool
}

// CheckpointResult 检查点操作结果
type CheckpointResult struct {
	Checkpoint      Checkpoint `json:"checkpoint"`
	FilesProcessed  int        `json:"filesProcessed"`
	Warnings        []string   `json:"warnings"`
}

// CheckpointDiff 两个检查点之间的差异
type CheckpointDiff struct {
	FromCheckpointID string     `json:"fromCheckpointId"`
	ToCheckpointID   string     `json:"toCheckpointId"`
	ModifiedFiles    []FileDiff `json:"modifiedFiles"`
	AddedFiles       []string   `json:"addedFiles"`
	DeletedFiles     []string   `json:"deletedFiles"`
	TokenDelta       int64      `json:"tokenDelta"`
}

// FileDiff 单个文件的差异
type FileDiff struct {
	Path        string  `json:"path"`
	Additions   int     `json:"additions"`
	Deletions   int     `json:"deletions"`
	DiffContent *string `json:"diffContent"`
}

// CheckpointPaths 检查点存储路径
type CheckpointPaths struct {
	TimelineFile   string
	CheckpointsDir string
	FilesDir       string
	claudeDir      string
	projectID      string
	sessionID      string
}

// NewCheckpointPaths 创建检查点路径
func NewCheckpointPaths(claudeDir, projectID, sessionID string) *CheckpointPaths {
	baseDir := filepath.Join(claudeDir, "projects", projectID, ".timelines", sessionID)

	return &CheckpointPaths{
		TimelineFile:   filepath.Join(baseDir, "timeline.json"),
		CheckpointsDir: filepath.Join(baseDir, "checkpoints"),
		FilesDir:       filepath.Join(baseDir, "files"),
		claudeDir:      claudeDir,
		projectID:      projectID,
		sessionID:      sessionID,
	}
}

// CheckpointDir 获取检查点目录
func (p *CheckpointPaths) CheckpointDir(checkpointID string) string {
	return filepath.Join(p.CheckpointsDir, checkpointID)
}

// CheckpointMetadataFile 获取检查点元数据文件路径
func (p *CheckpointPaths) CheckpointMetadataFile(checkpointID string) string {
	return filepath.Join(p.CheckpointDir(checkpointID), "metadata.json")
}

// CheckpointMessagesFile 获取检查点消息文件路径
func (p *CheckpointPaths) CheckpointMessagesFile(checkpointID string) string {
	return filepath.Join(p.CheckpointDir(checkpointID), "messages.jsonl")
}

// FileSnapshotPath 获取文件快照路径（内容寻址存储）
func (p *CheckpointPaths) FileSnapshotPath(fileHash string) string {
	return filepath.Join(p.FilesDir, "content_pool", fileHash)
}

// FileReferencePath 获取文件引用路径
func (p *CheckpointPaths) FileReferencePath(checkpointID, safeFilename string) string {
	return filepath.Join(p.FilesDir, "refs", checkpointID, safeFilename+".json")
}

// NewSessionTimeline 创建新的会话时间线
func NewSessionTimeline(sessionID string) *SessionTimeline {
	return &SessionTimeline{
		SessionID:             sessionID,
		RootNode:              nil,
		CurrentCheckpointID:   nil,
		AutoCheckpointEnabled: false,
		CheckpointStrategy:    StrategySmart,
		TotalCheckpoints:      0,
	}
}

// FindCheckpoint 在时间线树中查找检查点
func (t *SessionTimeline) FindCheckpoint(checkpointID string) *TimelineNode {
	if t.RootNode == nil {
		return nil
	}
	return findInTree(t.RootNode, checkpointID)
}

// findInTree 递归查找节点
func findInTree(node *TimelineNode, checkpointID string) *TimelineNode {
	if node.Checkpoint.ID == checkpointID {
		return node
	}

	for i := range node.Children {
		if found := findInTree(&node.Children[i], checkpointID); found != nil {
			return found
		}
	}

	return nil
}
