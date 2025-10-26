package claude

import (
	"Codesk/backend/util"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileService Claude 文件管理服务
type FileService struct{}

// NewFileService 创建文件服务实例
func NewFileService() *FileService {
	return &FileService{}
}

// FileEntry 文件或目录条目
type FileEntry struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	IsDirectory bool    `json:"is_directory"`
	Size        int64   `json:"size"`
	Extension   *string `json:"extension"`
}

// ClaudeMdFile Claude.md 文件信息
type ClaudeMdFile struct {
	RelativePath string `json:"relative_path"`
	AbsolutePath string `json:"absolute_path"`
	Size         int64  `json:"size"`
	Modified     int64  `json:"modified"`
}

// ListDirectoryContents 列出目录内容
func (s *FileService) ListDirectoryContents(path string) ([]FileEntry, error) {
	if !util.DirExists(path) {
		return nil, fmt.Errorf("directory not found: %s", path)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var fileEntries []FileEntry
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())
		fileEntry := FileEntry{
			Name:        entry.Name(),
			Path:        fullPath,
			IsDirectory: entry.IsDir(),
			Size:        info.Size(),
		}

		if !entry.IsDir() {
			ext := filepath.Ext(entry.Name())
			if ext != "" {
				fileEntry.Extension = &ext
			}
		}

		fileEntries = append(fileEntries, fileEntry)
	}

	return fileEntries, nil
}

// SearchFiles 搜索文件
func (s *FileService) SearchFiles(rootPath, pattern string, maxResults int) ([]string, error) {
	if !util.DirExists(rootPath) {
		return nil, fmt.Errorf("directory not found: %s", rootPath)
	}

	matches, err := util.SearchFiles(rootPath, pattern)
	if err != nil {
		return nil, err
	}

	// 限制结果数量
	if maxResults > 0 && len(matches) > maxResults {
		matches = matches[:maxResults]
	}

	return matches, nil
}

// GetRecentlyModifiedFiles 获取最近修改的文件
func (s *FileService) GetRecentlyModifiedFiles(rootPath string, limit int) ([]string, error) {
	if !util.DirExists(rootPath) {
		return nil, fmt.Errorf("directory not found: %s", rootPath)
	}

	return util.GetRecentlyModifiedFiles(rootPath, limit)
}

// FindClaudeMdFiles 查找所有 CLAUDE.md 文件
func (s *FileService) FindClaudeMdFiles(rootPath string) ([]ClaudeMdFile, error) {
	if !util.DirExists(rootPath) {
		return nil, fmt.Errorf("directory not found: %s", rootPath)
	}

	var claudeFiles []ClaudeMdFile

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过错误
		}

		if info.IsDir() {
			return nil
		}

		fileName := strings.ToUpper(filepath.Base(path))
		if fileName == "CLAUDE.MD" {
			relPath, _ := filepath.Rel(rootPath, path)

			claudeFiles = append(claudeFiles, ClaudeMdFile{
				RelativePath: relPath,
				AbsolutePath: path,
				Size:         info.Size(),
				Modified:     info.ModTime().Unix(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return claudeFiles, nil
}

// ReadClaudeMdFile 读取 CLAUDE.md 文件
func (s *FileService) ReadClaudeMdFile(path string) (string, error) {
	if !util.FileExists(path) {
		return "", fmt.Errorf("file not found: %s", path)
	}

	content, err := util.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return content, nil
}

// SaveClaudeMdFile 保存 CLAUDE.md 文件
func (s *FileService) SaveClaudeMdFile(path, content string) error {
	if err := util.WriteFile(path, content); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// ReadFile 读取任意文件
func (s *FileService) ReadFile(path string) (string, error) {
	if !util.FileExists(path) {
		return "", fmt.Errorf("file not found: %s", path)
	}

	return util.ReadFile(path)
}

// WriteFile 写入文件
func (s *FileService) WriteFile(path, content string) error {
	return util.WriteFile(path, content)
}
