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
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

// Window and dialog methods for Wails IPC

// OpenFileDialog opens a file dialog for selecting files
func (a *App) OpenFileDialog(options map[string]interface{}) ([]string, error) {
	var dialogOptions []runtime.FileFilter
	var defaultDirectory string
	var multiple bool
	var directory bool

	// Parse options
	if filters, ok := options["filters"].([]interface{}); ok {
		for _, filter := range filters {
			if filterMap, ok := filter.(map[string]interface{}); ok {
				if displayName, ok := filterMap["displayName"].(string); ok {
					if pattern, ok := filterMap["pattern"].(string); ok {
						dialogOptions = append(dialogOptions, runtime.FileFilter{
							DisplayName: displayName,
							Pattern:     pattern,
						})
					}
				}
			}
		}
	}

	if dir, ok := options["defaultDirectory"].(string); ok {
		defaultDirectory = dir
	}

	if mult, ok := options["multiple"].(bool); ok {
		multiple = mult
	}

	if dir, ok := options["directory"].(bool); ok {
		directory = dir
	}

	if directory {
		// Open directory dialog
		selected, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
			Title:            "Select Directory",
			DefaultDirectory: defaultDirectory,
		})
		if err != nil {
			return nil, err
		}
		return []string{selected}, nil
	}

	if multiple {
		// Open multiple files dialog
		return runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
			Title:               "Select Files",
			DefaultDirectory:    defaultDirectory,
			Filters:             dialogOptions,
			ShowHiddenFiles:     false,
			CanCreateDirectories: true,
		})
	}

	// Open single file dialog
	selected, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:               "Select File",
		DefaultDirectory:    defaultDirectory,
		Filters:             dialogOptions,
		ShowHiddenFiles:     false,
		CanCreateDirectories: true,
	})
	if err != nil {
		return nil, err
	}
	return []string{selected}, nil
}

// SaveFileDialog opens a save file dialog
func (a *App) SaveFileDialog(options map[string]interface{}) (string, error) {
	var dialogOptions []runtime.FileFilter
	var defaultFilename string

	// Parse options
	if filters, ok := options["filters"].([]interface{}); ok {
		for _, filter := range filters {
			if filterMap, ok := filter.(map[string]interface{}); ok {
				if displayName, ok := filterMap["displayName"].(string); ok {
					if pattern, ok := filterMap["pattern"].(string); ok {
						dialogOptions = append(dialogOptions, runtime.FileFilter{
							DisplayName: displayName,
							Pattern:     pattern,
						})
					}
				}
			}
		}
	}

	if filename, ok := options["defaultFilename"].(string); ok {
		defaultFilename = filename
	}

	return runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save File",
		DefaultFilename: defaultFilename,
		Filters:         dialogOptions,
	})
}

// OpenExternal opens an external URL in the default browser
func (a *App) OpenExternal(url string) error {
	runtime.BrowserOpenURL(a.ctx, url)
	return nil
}

// GetClaudeBinaryPath gets the stored Claude binary path from settings
func (a *App) GetClaudeBinaryPath() (string, error) {
	settings, err := a.claudeProjectService.GetClaudeSettings()
	if err != nil {
		return "", err
	}

	if settings.Data == nil {
		return "", nil
	}

	if path, ok := settings.Data["claude_binary_path"].(string); ok {
		return path, nil
	}

	return "", nil
}

// SetClaudeBinaryPath sets the Claude binary path in settings
func (a *App) SetClaudeBinaryPath(path string) error {
	settings, err := a.claudeProjectService.GetClaudeSettings()
	if err != nil {
		return err
	}

	if settings.Data == nil {
		settings.Data = make(map[string]interface{})
	}

	settings.Data["claude_binary_path"] = path

	return a.claudeProjectService.SaveClaudeSettings(settings)
}