package config

import (
	"os"
	"path/filepath"
)

// Config 应用配置
type Config struct {
	DatabasePath string
	ClaudeDir    string
	AppDataDir   string
}

var globalConfig *Config

// InitConfig 初始化配置
func InitConfig() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// 应用数据目录
	appDataDir := filepath.Join(homeDir, ".codesk")
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		return nil, err
	}

	// Claude 目录
	claudeDir := filepath.Join(homeDir, ".claude")

	globalConfig = &Config{
		DatabasePath: filepath.Join(appDataDir, "codesk.db"),
		ClaudeDir:    claudeDir,
		AppDataDir:   appDataDir,
	}

	return globalConfig, nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if globalConfig == nil {
		InitConfig()
	}
	return globalConfig
}

// GetDatabasePath 获取数据库路径
func GetDatabasePath() string {
	return GetConfig().DatabasePath
}

// GetClaudeDir 获取 Claude 目录
func GetClaudeDir() string {
	return GetConfig().ClaudeDir
}

// GetAppDataDir 获取应用数据目录
func GetAppDataDir() string {
	return GetConfig().AppDataDir
}

// GetProjectsDir 获取 Claude 项目目录
func GetProjectsDir() string {
	return filepath.Join(GetClaudeDir(), "projects")
}

// GetSettingsPath 获取 Claude 设置文件路径
func GetSettingsPath() string {
	return filepath.Join(GetClaudeDir(), "settings.json")
}
