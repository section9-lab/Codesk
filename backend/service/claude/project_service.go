package claude

import (
	"bufio"
	"Codesk/backend/config"
	"Codesk/backend/model"
	"Codesk/backend/util"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectService Claude 项目管理服务
type ProjectService struct{}

// NewProjectService 创建项目服务实例
func NewProjectService() *ProjectService {
	return &ProjectService{}
}

// ListProjects 列出所有 Claude 项目
func (s *ProjectService) ListProjects() ([]*model.Project, error) {
	projectsDir := config.GetProjectsDir()

	if !util.DirExists(projectsDir) {
		return []*model.Project{}, nil
	}

	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read projects directory: %w", err)
	}

	var projects []*model.Project
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectDir := filepath.Join(projectsDir, entry.Name())
		project, err := s.loadProject(entry.Name(), projectDir)
		if err != nil {
			continue // 跳过无效项目
		}

		projects = append(projects, project)
	}

	return projects, nil
}

// loadProject 加载单个项目信息
func (s *ProjectService) loadProject(dirName, projectDir string) (*model.Project, error) {
	// 获取项目路径（从 JSONL 文件中读取）
	projectPath, err := s.getProjectPathFromSessions(projectDir)
	if err != nil {
		// 降级：使用目录名解码
		projectPath = s.decodeProjectPath(dirName)
	}

	// 获取会话列表
	sessions, err := s.listSessions(projectDir)
	if err != nil {
		sessions = []string{}
	}

	// 获取目录创建时间
	info, err := os.Stat(projectDir)
	if err != nil {
		return nil, err
	}

	// 获取最近会话时间
	var mostRecentSession *int64
	if len(sessions) > 0 {
		latestTime := s.getMostRecentSessionTime(projectDir, sessions)
		if latestTime > 0 {
			mostRecentSession = &latestTime
		}
	}

	return &model.Project{
		ID:                dirName,
		Path:              projectPath,
		Sessions:          sessions,
		CreatedAt:         info.ModTime().Unix(),
		MostRecentSession: mostRecentSession,
	}, nil
}

// getProjectPathFromSessions 从会话文件中读取项目路径
func (s *ProjectService) getProjectPathFromSessions(projectDir string) (string, error) {
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}

		jsonlPath := filepath.Join(projectDir, entry.Name())
		jsonlEntries, err := util.ReadJSONL(jsonlPath)
		if err != nil {
			continue
		}

		// 检查前 10 行，查找有效的 cwd
		for i, entry := range jsonlEntries {
			if i >= 10 {
				break
			}

			if cwd, ok := entry["cwd"].(string); ok && cwd != "" {
				return cwd, nil
			}
		}
	}

	return "", fmt.Errorf("could not determine project path from sessions")
}

// decodeProjectPath 解码项目目录名为路径（降级方案）
func (s *ProjectService) decodeProjectPath(encoded string) string {
	return strings.ReplaceAll(encoded, "-", "/")
}

// listSessions 列出项目的所有会话
func (s *ProjectService) listSessions(projectDir string) ([]string, error) {
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return nil, err
	}

	var sessions []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".jsonl") {
			sessionID := strings.TrimSuffix(name, ".jsonl")
			sessions = append(sessions, sessionID)
		}
	}

	return sessions, nil
}

// getMostRecentSessionTime 获取最近会话的时间戳
func (s *ProjectService) getMostRecentSessionTime(projectDir string, sessions []string) int64 {
	var latestTime int64

	for _, sessionID := range sessions {
		sessionPath := filepath.Join(projectDir, sessionID+".jsonl")
		info, err := os.Stat(sessionPath)
		if err != nil {
			continue
		}

		modTime := info.ModTime().Unix()
		if modTime > latestTime {
			latestTime = modTime
		}
	}

	return latestTime
}

// GetProjectSessions 获取项目的会话列表（带详细信息）
func (s *ProjectService) GetProjectSessions(projectPath string) ([]*model.Session, error) {
	// 编码项目路径
	encodedPath := strings.ReplaceAll(projectPath, "/", "-")
	projectDir := filepath.Join(config.GetProjectsDir(), encodedPath)

	if !util.DirExists(projectDir) {
		return []*model.Session{}, nil
	}

	sessionIDs, err := s.listSessions(projectDir)
	if err != nil {
		return nil, err
	}

	var sessions []*model.Session
	for _, sessionID := range sessionIDs {
		session, err := s.loadSession(sessionID, projectPath, projectDir)
		if err != nil {
			continue
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// loadSession 加载会话详细信息
func (s *ProjectService) loadSession(sessionID, projectPath, projectDir string) (*model.Session, error) {
	sessionPath := filepath.Join(projectDir, sessionID+".jsonl")

	info, err := os.Stat(sessionPath)
	if err != nil {
		return nil, err
	}

	// 读取第一条用户消息
	firstMessage, timestamp := s.extractFirstUserMessage(sessionPath)

	session := &model.Session{
		ID:               sessionID,
		ProjectID:        filepath.Base(projectDir),
		ProjectPath:      projectPath,
		CreatedAt:        info.ModTime().Unix(),
		FirstMessage:     firstMessage,
		MessageTimestamp: timestamp,
	}

	return session, nil
}

// extractFirstUserMessage 提取第一条用户消息
func (s *ProjectService) extractFirstUserMessage(jsonlPath string) (*string, *string) {
	entries, err := util.ReadJSONL(jsonlPath)
	if err != nil {
		return nil, nil
	}

	for _, entry := range entries {
		if msgType, ok := entry["type"].(string); ok && msgType == "message" {
			if message, ok := entry["message"].(map[string]interface{}); ok {
				if role, ok := message["role"].(string); ok && role == "user" {
					if content, ok := message["content"].(string); ok {
						// 跳过系统消息
						if strings.Contains(content, "Caveat: The messages below were generated") {
							continue
						}
						if strings.HasPrefix(content, "<command-name>") {
							continue
						}

						timestamp := ""
						if ts, ok := entry["timestamp"].(string); ok {
							timestamp = ts
						}

						return &content, &timestamp
					}
				}
			}
		}
	}

	return nil, nil
}

// GetClaudeSettings 获取 Claude 设置
func (s *ProjectService) GetClaudeSettings() (*model.ClaudeSettings, error) {
	settingsPath := config.GetSettingsPath()

	if !util.FileExists(settingsPath) {
		return &model.ClaudeSettings{Data: make(map[string]interface{})}, nil
	}

	content, err := util.ReadFile(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return nil, fmt.Errorf("failed to parse settings: %w", err)
	}

	return &model.ClaudeSettings{Data: data}, nil
}

// SaveClaudeSettings 保存 Claude 设置
func (s *ProjectService) SaveClaudeSettings(settings *model.ClaudeSettings) error {
	settingsPath := config.GetSettingsPath()

	data, err := json.MarshalIndent(settings.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := util.WriteFile(settingsPath, string(data)); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	return nil
}

// LoadSessionHistory 加载会话历史
func (s *ProjectService) LoadSessionHistory(sessionID, projectPath string) (*model.SessionHistory, error) {
	encodedPath := strings.ReplaceAll(projectPath, "/", "-")
	sessionPath := filepath.Join(config.GetProjectsDir(), encodedPath, sessionID+".jsonl")

	if !util.FileExists(sessionPath) {
		return nil, fmt.Errorf("session file not found")
	}

	// 读取原始 JSONL 行
	file, err := os.Open(sessionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open session file: %w", err)
	}
	defer file.Close()

	var messages []map[string]interface{}
	scanner := bufio.NewScanner(file)

	// Increase the maximum token size to handle very long lines (up to 10MB)
	buf := make([]byte, 0, 10*1024*1024) // 10MB buffer
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// 跳过无效行，但记录日志
			continue
		}

		messages = append(messages, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	// 直接返回原始消息，保持完整的 JSON 数据结构
	return &model.SessionHistory{
		SessionID: sessionID,
		Messages:  messages,
	}, nil
}
