package main

import (
	"fmt"
	"context"
	"Codesk/backend/config"
	"Codesk/backend/repository"
	"Codesk/backend/service/agent"
	"Codesk/backend/service/checkpoint"
	"Codesk/backend/service/claude"
	"Codesk/backend/service/mcp"
	"Codesk/backend/service/proxy"
	"Codesk/backend/service/slash"
	"Codesk/backend/service/storage"
	"Codesk/backend/service/usage"
)

// App struct
type App struct {
	ctx context.Context

	// Services
	claudeProjectService   *claude.ProjectService
	claudeFileService      *claude.FileService
	claudeExecutionService *claude.ExecutionService
	agentService           *agent.AgentService
	checkpointService      *checkpoint.CheckpointService
	mcpService             *mcp.MCPService
	usageService           *usage.UsageService
	proxyService           *proxy.ProxyService
	slashService           *slash.SlashService
	storageService         *storage.StorageService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		claudeProjectService:   claude.NewProjectService(),
		claudeFileService:      claude.NewFileService(),
		claudeExecutionService: claude.NewExecutionService(),
		agentService:           agent.NewAgentService(),
		checkpointService:      checkpoint.NewCheckpointService(),
		mcpService:             mcp.NewMCPService(),
		usageService:           usage.NewUsageService(),
		proxyService:           proxy.NewProxyService(),
		slashService:           slash.NewSlashService(),
		storageService:         storage.NewStorageService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize configuration
	fmt.Println("Initialize configuration ...")
	config.InitConfig()

	// Initialize database
	fmt.Println("Initialize database ...")
	dbPath := config.GetDatabasePath()
	if err := repository.InitDatabase(dbPath); err != nil {
		panic("Failed to initialize database: " + err.Error())
	}

	// Load and apply proxy settings
	proxySettings, err := a.proxyService.GetProxySettings()
	if err == nil {
		a.proxyService.ApplyProxySettings(proxySettings)
	}
	fmt.Println("startup ok!")
}