package main

import (
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
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/progrium/darwinkit/macos/appkit"
	"github.com/progrium/darwinkit/macos/foundation"
	"github.com/progrium/darwinkit/objc"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Global variables for model selection management
var (
	globalSelectedModel = "Claude"
	globalContext       context.Context
	globalStatusItem    *appkit.StatusItem
	modelChangeEvent    = make(chan string, 10)   // Buffered channel for model changes
	statusBarClickEvent = make(chan struct{}, 10) // Buffered channel for status bar clicks

	// Menu item references for updating selection state
	globalMenuItems map[string]appkit.MenuItem
)

// updateModelSelection safely updates the selected model
func updateModelSelection(modelName string) {
	globalSelectedModel = modelName

	// Update menu item states
	updateMenuStates(modelName)

	// Send model name to channel instead of direct event emission
	select {
	case modelChangeEvent <- modelName:
		log.Printf("✓ 发送模型变更事件: %s", modelName)
	default:
		log.Printf("⚠ 模型变更事件队列已满")
	}

	log.Printf("✓ 更新模型选择: %s", modelName)
}

// updateMenuStates updates the visual state of menu items to show which model is selected
func updateMenuStates(selectedModel string) {
	if globalMenuItems == nil {
		return
	}

	// Update menu states on main queue for thread safety
	mainQueue := foundation.OperationQueue_MainQueue()
	mainQueue.AddOperationWithBlock(func() {
		for modelName, menuItem := range globalMenuItems {
			if modelName == selectedModel {
				// Set selected state (checked with checkmark)
				menuItem.SetState(1) // NSControlStateValueOn
				log.Printf("✓ 设置菜单项选中状态: %s", modelName)
			} else {
				// Set unselected state
				menuItem.SetState(0) // NSControlStateValueOff
				log.Printf("✓ 取消菜单项选中状态: %s", modelName)
			}
		}
	})
}

// startModelChangeListener starts a goroutine to safely handle model change events
func startModelChangeListener(ctx context.Context) {
	go func() {
		for modelName := range modelChangeEvent {
			// This runs in a separate goroutine, safe from CGO callback issues
			if ctx != nil {
				// Send event to frontend
				wailsruntime.EventsEmit(ctx, "modelChanged", map[string]interface{}{
					"model": modelName,
				})

				log.Printf("✓ 已处理模型变更: %s", modelName)
			}
		}
	}()
}

// startStatusBarClickListener starts a goroutine to safely handle status bar click events
func startStatusBarClickListener(ctx context.Context) {
	go func() {
		for range statusBarClickEvent {
			// This runs in a separate goroutine, safe from CGO callback issues
			if ctx != nil {
				// Send event to frontend safely
				wailsruntime.EventsEmit(ctx, "statusBarClicked", nil)
				log.Printf("✓ 已处理状态栏点击")
			}
		}
	}()
}

// App struct
type App struct {
	statusBarItem *appkit.StatusItem
	statusBarMenu *appkit.Menu

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

// OnShutdown 应用关闭时清理状态栏图标
func (a *App) OnShutdown(ctx context.Context) {
	if a.statusBarItem != nil {
		appkit.StatusBar_SystemStatusBar().RemoveStatusItem(a.statusBarItem)
	}
}

// AppDelegate 处理 macOS 原生回调
type AppDelegate struct {
	ctx context.Context
}

func (a *AppDelegate) onStatusBarClick(sender objc.Object) {
	// Safely send status bar click event via channel
	select {
	case statusBarClickEvent <- struct{}{}:
		log.Printf("✓ 发送状态栏点击事件")
	default:
		log.Printf("⚠ 状态栏点击事件队列已满")
	}
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
	// Store context globally for use in callbacks
	globalContext = ctx
	a.ctx = ctx

	// Start the model change event listener
	startModelChangeListener(ctx)

	// Start the status bar click event listener
	startStatusBarClickListener(ctx)

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

	// Initialize status bar on main thread using operation queue
	runtime.LockOSThread()
	mainQueue := foundation.OperationQueue_MainQueue()
	mainQueue.AddOperationWithBlock(func() {
		a.initStatusBar(ctx)
	})

	fmt.Println("startup ok!")
}

// initStatusBar initializes the status bar on the main thread
func (a *App) initStatusBar(ctx context.Context) {
	log.Printf("=== 开始初始化状态栏 ===")

	// 1. 获取 macOS 系统状态栏
	statusBar := appkit.StatusBar_SystemStatusBar()
	log.Printf("✓ 获取系统状态栏成功")

	// 2. 创建状态栏项（使用固定长度以确保可见性）
	statusItem := statusBar.StatusItemWithLength(15)
	log.Printf("✓ 创建状态栏项成功，固定长度: 15")

	// 3. 确保状态栏项可见
	statusItem.SetVisible(true)
	log.Printf("✓ 设置状态栏项可见")

	// 4. 使用文本作为主要标识（更可靠）
	button := statusItem.Button()
	log.Printf("✓ 设置按钮标题: 🔧")

	// 5. 添加简单图标作为备选
	// localIcon := appkit.NewImageFromFile(core.NSString(iconPath))
	icon := appkit.Image_ImageWithSystemSymbolNameAccessibilityDescription("gear", "Settings icon")
	if icon.IsValid() {
		icon.SetTemplate(true) // 使用模板模式适配深色/浅色模式
		button.SetImage(icon)
		log.Printf("✓ 设置按钮图标成功")
	} else {
		button.SetTitle("🔧")
		log.Printf("⚠ 创建系统图标失败，仅使用文本")
	}

	// 6. 创建菜单
	menu := appkit.NewMenuWithTitle("Codesk")
	a.statusBarMenu = &menu
	log.Printf("✓ 创建菜单成功")

	// Initialize global menu items map
	globalMenuItems = make(map[string]appkit.MenuItem)
	log.Printf("✓ 初始化全局菜单项映射")

	// 显示/隐藏应用
	showHideItem := appkit.NewMenuItemWithAction("Select Model", "", func(sender objc.Object) {
		log.Printf("🔔 显示/隐藏 Codesk 被点击")
	})
	menu.AddItem(showHideItem)
	log.Printf("✓ 添加显示/隐藏菜单项")

	// 分隔线
	menu.AddItem(appkit.MenuItem_SeparatorItem())
	log.Printf("✓ 添加分隔线")

	// ===== 具体选项 ======
	// 1.OpenAI
	openAiItem := appkit.NewMenuItemWithAction("OpenAI", "", func(sender objc.Object) {
		log.Printf("🔔 OpenAI被点击")
		updateModelSelection("OpenAI")
	})
	menu.AddItem(openAiItem)
	globalMenuItems["OpenAI"] = openAiItem
	log.Printf("✓ 添加OpenAI")

	// 2.Claude
	claudeItem := appkit.NewMenuItemWithAction("Claude", "", func(sender objc.Object) {
		log.Printf("🔔 Claude被点击")
		updateModelSelection("Claude")
	})
	menu.AddItem(claudeItem)
	globalMenuItems["Claude"] = claudeItem
	log.Printf("✓ 添加Claude")

	// 3.Gemini
	geminiItem := appkit.NewMenuItemWithAction("Gemini", "", func(sender objc.Object) {
		log.Printf("🔔 Gemini被点击")
		updateModelSelection("Gemini")
	})
	menu.AddItem(geminiItem)
	globalMenuItems["Gemini"] = geminiItem
	log.Printf("✓ 添加Gemini")

	// 4.DeepSeek
	deepSeekItem := appkit.NewMenuItemWithAction("DeepSeek", "", func(sender objc.Object) {
		log.Printf("🔔 DeepSeek被点击")
		updateModelSelection("DeepSeek")
	})
	menu.AddItem(deepSeekItem)
	globalMenuItems["DeepSeek"] = deepSeekItem
	log.Printf("✓ 添加DeepSeek")

	// 5.Qwen
	qwenItem := appkit.NewMenuItemWithAction("Qwen", "", func(sender objc.Object) {
		log.Printf("🔔 Qwen被点击")
		updateModelSelection("Qwen")
	})
	menu.AddItem(qwenItem)
	globalMenuItems["Qwen"] = qwenItem
	log.Printf("✓ 添加Qwen")

	// 6.GLM
	glmItem := appkit.NewMenuItemWithAction("GLM", "", func(sender objc.Object) {
		log.Printf("🔔 GLM被点击")
		updateModelSelection("GLM")
	})
	menu.AddItem(glmItem)
	globalMenuItems["GLM"] = glmItem
	log.Printf("✓ 添加GLM")

	// 7.Kimi
	kimiItem := appkit.NewMenuItemWithAction("Kimi", "", func(sender objc.Object) {
		log.Printf("🔔 Kimi被点击")
		updateModelSelection("Kimi")
	})
	menu.AddItem(kimiItem)
	globalMenuItems["Kimi"] = kimiItem
	log.Printf("✓ 添加Kimi")
	// ====== 选项结束  =======

	// 设置菜单
	statusItem.SetMenu(menu)
	log.Printf("✓ 设置菜单到状态栏")

	// 保存状态栏项的引用并 retain
	a.statusBarItem = &statusItem
	globalStatusItem = &statusItem
	objc.Retain(&statusItem)
	log.Printf("✓ 保存并 retain 状态栏项")

	// Set initial status bar title
	// updateStatusBarTitle(globalSelectedModel)
	// log.Printf("✓ 设置初始状态栏标题: %s", getDisplayText(globalSelectedModel))

	// Set initial menu states
	updateMenuStates(globalSelectedModel)
	log.Printf("✓ 设置初始菜单选中状态: %s", globalSelectedModel)

	// 9. 调试信息
	log.Printf("📊 状态栏信息:")
	log.Printf("   - 长度: %v", statusItem.Length())
	log.Printf("   - 可见: %v", statusItem.IsVisible())
	log.Printf("   - 图标有效: %v", icon.IsValid())
	log.Printf("   - 菜单项数量: %v", len(menu.ItemArray()))
	log.Printf("🎉 状态栏初始化完成")
}

// Window and dialog methods for Wails IPC

// OpenFileDialog opens a file dialog for selecting files
func (a *App) OpenFileDialog(options map[string]interface{}) ([]string, error) {
	var dialogOptions []wailsruntime.FileFilter
	var defaultDirectory string
	var multiple bool
	var directory bool

	// Parse options
	if filters, ok := options["filters"].([]interface{}); ok {
		for _, filter := range filters {
			if filterMap, ok := filter.(map[string]interface{}); ok {
				if displayName, ok := filterMap["displayName"].(string); ok {
					if pattern, ok := filterMap["pattern"].(string); ok {
						dialogOptions = append(dialogOptions, wailsruntime.FileFilter{
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
		selected, err := wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{
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
		return wailsruntime.OpenMultipleFilesDialog(a.ctx, wailsruntime.OpenDialogOptions{
			Title:                "Select Files",
			DefaultDirectory:     defaultDirectory,
			Filters:              dialogOptions,
			ShowHiddenFiles:      false,
			CanCreateDirectories: true,
		})
	}

	// Open single file dialog
	selected, err := wailsruntime.OpenFileDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title:                "Select File",
		DefaultDirectory:     defaultDirectory,
		Filters:              dialogOptions,
		ShowHiddenFiles:      false,
		CanCreateDirectories: true,
	})
	if err != nil {
		return nil, err
	}
	return []string{selected}, nil
}

// SaveFileDialog opens a save file dialog
func (a *App) SaveFileDialog(options map[string]interface{}) (string, error) {
	var dialogOptions []wailsruntime.FileFilter
	var defaultFilename string

	// Parse options
	if filters, ok := options["filters"].([]interface{}); ok {
		for _, filter := range filters {
			if filterMap, ok := filter.(map[string]interface{}); ok {
				if displayName, ok := filterMap["displayName"].(string); ok {
					if pattern, ok := filterMap["pattern"].(string); ok {
						dialogOptions = append(dialogOptions, wailsruntime.FileFilter{
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

	return wailsruntime.SaveFileDialog(a.ctx, wailsruntime.SaveDialogOptions{
		Title:           "Save File",
		DefaultFilename: defaultFilename,
		Filters:         dialogOptions,
	})
}

// OpenExternal opens an external URL in the default browser
func (a *App) OpenExternal(url string) error {
	wailsruntime.BrowserOpenURL(a.ctx, url)
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
