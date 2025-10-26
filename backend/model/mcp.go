package model

// MCPServer MCP 服务器配置
type MCPServer struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	Env         map[string]string `json:"env"`
	Disabled    bool              `json:"disabled"`
	AutoApprove []string          `json:"auto_approve"`
}

// MCPConfig MCP 配置文件
type MCPConfig struct {
	MCPServers map[string]MCPServer `json:"mcpServers"`
}

// MCPServerStatus MCP 服务器状态
type MCPServerStatus struct {
	Name      string `json:"name"`
	Status    string `json:"status"` // running, stopped, error
	PID       *int   `json:"pid"`
	Error     string `json:"error"`
	Connected bool   `json:"connected"`
}

// MCPProjectConfig 项目级别的 MCP 配置
type MCPProjectConfig struct {
	ProjectPath string              `json:"project_path"`
	Servers     map[string]MCPServer `json:"servers"`
	Choices     map[string]string   `json:"choices"` // 用户选择的服务器
}
