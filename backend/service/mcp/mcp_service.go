package mcp

import (
	"Codesk/backend/config"
	"Codesk/backend/model"
	"Codesk/backend/util"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// MCPService MCP 服务
type MCPService struct{}

// NewMCPService 创建 MCP 服务实例
func NewMCPService() *MCPService {
	return &MCPService{}
}

// GetMCPConfigPath 获取 MCP 配置文件路径
func (s *MCPService) GetMCPConfigPath() string {
	return filepath.Join(config.GetClaudeDir(), "mcp.json")
}

// GetProjectMCPConfigPath 获取项目级 MCP 配置路径
func (s *MCPService) GetProjectMCPConfigPath(projectPath string) string {
	return filepath.Join(projectPath, ".claude", "mcp.json")
}

// ListServers 列出所有 MCP 服务器
func (s *MCPService) ListServers() (map[string]model.MCPServer, error) {
	configPath := s.GetMCPConfigPath()

	if !util.FileExists(configPath) {
		return make(map[string]model.MCPServer), nil
	}

	content, err := util.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read MCP config: %w", err)
	}

	var config model.MCPConfig
	if err := json.Unmarshal([]byte(content), &config); err != nil {
		return nil, fmt.Errorf("failed to parse MCP config: %w", err)
	}

	return config.MCPServers, nil
}

// GetServer 获取指定的 MCP 服务器配置
func (s *MCPService) GetServer(name string) (*model.MCPServer, error) {
	servers, err := s.ListServers()
	if err != nil {
		return nil, err
	}

	server, exists := servers[name]
	if !exists {
		return nil, fmt.Errorf("MCP server not found: %s", name)
	}

	return &server, nil
}

// AddServer 添加 MCP 服务器
func (s *MCPService) AddServer(name string, server model.MCPServer) error {
	servers, err := s.ListServers()
	if err != nil {
		return err
	}

	servers[name] = server

	return s.saveConfig(servers)
}

// RemoveServer 移除 MCP 服务器
func (s *MCPService) RemoveServer(name string) error {
	servers, err := s.ListServers()
	if err != nil {
		return err
	}

	if _, exists := servers[name]; !exists {
		return fmt.Errorf("MCP server not found: %s", name)
	}

	delete(servers, name)

	return s.saveConfig(servers)
}

// saveConfig 保存 MCP 配置
func (s *MCPService) saveConfig(servers map[string]model.MCPServer) error {
	config := model.MCPConfig{
		MCPServers: servers,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal MCP config: %w", err)
	}

	configPath := s.GetMCPConfigPath()

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := util.WriteFile(configPath, string(data)); err != nil {
		return fmt.Errorf("failed to write MCP config: %w", err)
	}

	return nil
}

// ReadProjectConfig 读取项目级 MCP 配置
func (s *MCPService) ReadProjectConfig(projectPath string) (*model.MCPProjectConfig, error) {
	configPath := s.GetProjectMCPConfigPath(projectPath)

	if !util.FileExists(configPath) {
		return &model.MCPProjectConfig{
			ProjectPath: projectPath,
			Servers:     make(map[string]model.MCPServer),
			Choices:     make(map[string]string),
		}, nil
	}

	content, err := util.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read project MCP config: %w", err)
	}

	var config model.MCPProjectConfig
	if err := json.Unmarshal([]byte(content), &config); err != nil {
		return nil, fmt.Errorf("failed to parse project MCP config: %w", err)
	}

	config.ProjectPath = projectPath

	return &config, nil
}

// SaveProjectConfig 保存项目级 MCP 配置
func (s *MCPService) SaveProjectConfig(config *model.MCPProjectConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project MCP config: %w", err)
	}

	configPath := s.GetProjectMCPConfigPath(config.ProjectPath)

	// 确保目录存在
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := util.WriteFile(configPath, string(data)); err != nil {
		return fmt.Errorf("failed to write project MCP config: %w", err)
	}

	return nil
}

// ResetProjectChoices 重置项目的 MCP 选择
func (s *MCPService) ResetProjectChoices(projectPath string) error {
	config, err := s.ReadProjectConfig(projectPath)
	if err != nil {
		return err
	}

	config.Choices = make(map[string]string)

	return s.SaveProjectConfig(config)
}

// AddServerFromJSON 从 JSON 添加服务器
func (s *MCPService) AddServerFromJSON(name, jsonData string) error {
	var server model.MCPServer
	if err := json.Unmarshal([]byte(jsonData), &server); err != nil {
		return fmt.Errorf("failed to parse server JSON: %w", err)
	}

	return s.AddServer(name, server)
}

// GetServerStatus 获取服务器状态（简化版本，实际需要进程管理）
func (s *MCPService) GetServerStatus(name string) (*model.MCPServerStatus, error) {
	server, err := s.GetServer(name)
	if err != nil {
		return nil, err
	}

	// 简化实现：只返回配置状态
	status := &model.MCPServerStatus{
		Name:      name,
		Status:    "unknown",
		Connected: false,
	}

	if server.Disabled {
		status.Status = "disabled"
	} else {
		status.Status = "configured"
	}

	return status, nil
}

// TestConnection 测试 MCP 服务器连接（简化版本）
func (s *MCPService) TestConnection(name string) (bool, error) {
	server, err := s.GetServer(name)
	if err != nil {
		return false, err
	}

	if server.Disabled {
		return false, fmt.Errorf("server is disabled")
	}

	// 简化实现：检查命令是否存在
	if server.Command == "" {
		return false, fmt.Errorf("server command is empty")
	}

	// 实际实现需要启动进程并测试连接
	// 这里只做基本验证
	return true, nil
}
