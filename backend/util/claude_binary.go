package util

import (
	"Codesk/backend/model"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// FindClaudeBinary 查找 Claude Code 二进制文件路径
func FindClaudeBinary() (string, error) {
	// 首先尝试从 PATH 中查找
	path, err := exec.LookPath("claude")
	if err == nil {
		return path, nil
	}

	// 根据操作系统查找常见安装位置
	var searchPaths []string

	switch runtime.GOOS {
	case "darwin": // macOS
		homeDir, _ := os.UserHomeDir()
		searchPaths = []string{
			"/usr/local/bin/claude",
			"/opt/homebrew/bin/claude",
			filepath.Join(homeDir, ".local/bin/claude"),
			filepath.Join(homeDir, "bin/claude"),
		}
	case "linux":
		homeDir, _ := os.UserHomeDir()
		searchPaths = []string{
			"/usr/local/bin/claude",
			"/usr/bin/claude",
			filepath.Join(homeDir, ".local/bin/claude"),
			filepath.Join(homeDir, "bin/claude"),
		}
	case "windows":
		searchPaths = []string{
			"C:\\Program Files\\Claude\\claude.exe",
			"C:\\Program Files (x86)\\Claude\\claude.exe",
		}
		// 添加用户目录
		if homeDir, err := os.UserHomeDir(); err == nil {
			searchPaths = append(searchPaths,
				filepath.Join(homeDir, "AppData", "Local", "Claude", "claude.exe"),
			)
		}
	}

	// 搜索路径
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("claude binary not found")
}

// GetClaudeVersion 获取 Claude Code 版本
func GetClaudeVersion(claudePath string) (string, error) {
	cmd := exec.Command(claudePath, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get claude version: %w", err)
	}

	version := strings.TrimSpace(string(output))
	return version, nil
}

// CheckClaudeInstallation 检查 Claude Code 是否已安装
func CheckClaudeInstallation() (bool, string, error) {
	claudePath, err := FindClaudeBinary()
	if err != nil {
		return false, "", nil
	}

	version, err := GetClaudeVersion(claudePath)
	if err != nil {
		return true, "", err
	}

	return true, version, nil
}

// ListClaudeInstallations 列出所有可能的 Claude 安装位置
func ListClaudeInstallations() []model.ClaudeInstallation {
	var installations []model.ClaudeInstallation

	// 从 PATH 查找
	if path, err := exec.LookPath("claude"); err == nil {
		installations = append(installations, model.ClaudeInstallation{
			Path:             path,
			Source:           "PATH",
			InstallationType: "System",
		})
	}

	// 搜索常见位置
	var searchPaths []string
	homeDir, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "darwin":
		searchPaths = []string{
			"/usr/local/bin/claude",
			"/opt/homebrew/bin/claude",
			filepath.Join(homeDir, ".local/bin/claude"),
		}
	case "linux":
		searchPaths = []string{
			"/usr/local/bin/claude",
			"/usr/bin/claude",
			filepath.Join(homeDir, ".local/bin/claude"),
		}
	case "windows":
		searchPaths = []string{
			"C:\\Program Files\\Claude\\claude.exe",
			filepath.Join(homeDir, "AppData", "Local", "Claude", "claude.exe"),
		}
	}

	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			// 避免重复
			found := false
			for _, existing := range installations {
				if existing.Path == path {
					found = true
					break
				}
			}
			if !found {
				// 尝试获取版本信息
				version := ""
				if ver, err := GetClaudeVersion(path); err == nil {
					version = ver
				}

				// 确定安装类型
				installType := "Custom"
				source := "unknown"

				if strings.Contains(path, "homebrew") {
					installType = "System"
					source = "homebrew"
				} else if strings.Contains(path, "/usr/local/bin") || strings.Contains(path, "/usr/bin") {
					installType = "System"
					source = "system"
				} else if strings.Contains(path, ".local/bin") {
					installType = "Custom"
					source = "local"
				} else if strings.Contains(path, "Program Files") {
					installType = "System"
					source = "installer"
				} else if strings.Contains(path, "AppData") {
					installType = "Custom"
					source = "user"
				}

				installations = append(installations, model.ClaudeInstallation{
					Path:             path,
					Version:          version,
					Source:           source,
					InstallationType: installType,
				})
			}
		}
	}

	return installations
}

// GetClaudeDirectory 获取 Claude 配置目录 (~/.claude)
func GetClaudeDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	claudeDir := filepath.Join(homeDir, ".claude")
	if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
		return "", fmt.Errorf("claude directory not found: %s", claudeDir)
	}

	return claudeDir, nil
}

// GetClaudeProjectsDirectory 获取 Claude 项目目录
func GetClaudeProjectsDirectory() (string, error) {
	claudeDir, err := GetClaudeDirectory()
	if err != nil {
		return "", err
	}

	projectsDir := filepath.Join(claudeDir, "projects")
	return projectsDir, nil
}
