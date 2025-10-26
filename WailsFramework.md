# Wails Framework

Wails is a framework for building desktop applications using Go and web technologies. It provides a lightweight alternative to Electron by using native rendering engines instead of embedding a browser. The framework wraps Go code and web frontends into a single binary, making it ideal for Go programmers who want to create desktop applications with HTML/JS/CSS interfaces without running a web server.

The framework supports multiple platforms (Windows 10/11, macOS 10.13+, Linux AMD64/ARM64) and offers native features like dialogs, menus, system themes, and events. It uses WebView2 on Windows, WebKit on macOS and Linux, and provides automatic TypeScript bindings generation for Go structs and methods. Wails includes a powerful CLI for project scaffolding, development with hot reload, and production builds.

## Installation and Project Setup

Installing Wails CLI and creating a new project

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verify installation
wails doctor

# Create a new project with React
wails init -n myproject -t react-ts

# Create a new project with Vue
wails init -n myproject -t vue-ts

# Create a new project with Svelte
wails init -n myproject -t svelte-ts

# Navigate to project
cd myproject

# Run in development mode with hot reload
wails dev

# Build for production
wails build

# Build with specific options
wails build -clean -upx -ldflags "-w -s"
```

## Basic Application Structure

Creating a simple Wails application with Go backend and web frontend

```go
package main

import (
    "context"
    "embed"
    "fmt"
    "github.com/wailsapp/wails/v2"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

type App struct {
    ctx context.Context
}

func NewApp() *App {
    return &App{}
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
}

func (a *App) Greet(name string) string {
    return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetData() map[string]interface{} {
    return map[string]interface{}{
        "message": "Data from Go backend",
        "count":   42,
    }
}

func main() {
    app := NewApp()

    err := wails.Run(&options.App{
        Title:  "MyApp",
        Width:  1024,
        Height: 768,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        BackgroundColour: options.NewRGB(27, 38, 54),
        OnStartup:        app.startup,
        Bind: []interface{}{
            app,
        },
    })

    if err != nil {
        println("Error:", err.Error())
    }
}
```
frontend
```
import {Greet} from "../wailsjs/go/main/App";

<button className="btn" onClick={greet}>Greet</button>
```
## Window Management

Controlling window properties and behavior from Go

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
    ctx context.Context
}

func (a *App) SetupWindow() {
    runtime.WindowSetTitle(a.ctx, "New Window Title")
    runtime.WindowSetSize(a.ctx, 800, 600)
    runtime.WindowCenter(a.ctx)
    runtime.WindowSetMinSize(a.ctx, 400, 300)
    runtime.WindowSetMaxSize(a.ctx, 1920, 1080)
}

func (a *App) ToggleFullscreen() {
    if runtime.WindowIsFullscreen(a.ctx) {
        runtime.WindowUnfullscreen(a.ctx)
    } else {
        runtime.WindowFullscreen(a.ctx)
    }
}

func (a *App) MaximizeWindow() {
    runtime.WindowMaximise(a.ctx)
}

func (a *App) MinimizeWindow() {
    runtime.WindowMinimise(a.ctx)
}

func (a *App) GetWindowSize() (int, int) {
    width, height := runtime.WindowGetSize(a.ctx)
    return width, height
}

func (a *App) GetWindowPosition() (int, int) {
    x, y := runtime.WindowGetPosition(a.ctx)
    return x, y
}

func (a *App) SetWindowAlwaysOnTop(alwaysOnTop bool) {
    runtime.WindowSetAlwaysOnTop(a.ctx, alwaysOnTop)
}

func (a *App) SetWindowTheme(theme string) {
    switch theme {
    case "dark":
        runtime.WindowSetDarkTheme(a.ctx)
    case "light":
        runtime.WindowSetLightTheme(a.ctx)
    default:
        runtime.WindowSetSystemDefaultTheme(a.ctx)
    }
}

func (a *App) HideWindow() {
    runtime.WindowHide(a.ctx)
}

func (a *App) ShowWindow() {
    runtime.WindowShow(a.ctx)
}

func (a *App) SetBackgroundColor(r, g, b uint8) {
    runtime.WindowSetBackgroundColour(a.ctx, r, g, b, 255)
}

func (a *App) ReloadWindow() {
    runtime.WindowReload(a.ctx)
}

func (a *App) PrintWindow() {
    runtime.WindowPrint(a.ctx)
}
```

## Dialog Operations

Native file dialogs and message boxes

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
    ctx context.Context
}

func (a *App) OpenFile() (string, error) {
    return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
        Title: "Select File",
        Filters: []runtime.FileFilter{
            {
                DisplayName: "Images (*.png;*.jpg)",
                Pattern:     "*.png;*.jpg;*.jpeg",
            },
            {
                DisplayName: "Text Files (*.txt)",
                Pattern:     "*.txt",
            },
        },
        DefaultDirectory:     "/home/user/documents",
        CanCreateDirectories: true,
        ShowHiddenFiles:      false,
    })
}

func (a *App) OpenMultipleFiles() ([]string, error) {
    return runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
        Title: "Select Multiple Files",
        Filters: []runtime.FileFilter{
            {
                DisplayName: "All Files (*.*)",
                Pattern:     "*.*",
            },
        },
    })
}

func (a *App) OpenDirectory() (string, error) {
    return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
        Title:            "Select Directory",
        DefaultDirectory: "/home/user",
    })
}

func (a *App) SaveFile() (string, error) {
    return runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
        Title:           "Save File",
        DefaultFilename: "untitled.txt",
        Filters: []runtime.FileFilter{
            {
                DisplayName: "Text Files (*.txt)",
                Pattern:     "*.txt",
            },
        },
    })
}

func (a *App) ShowInfoDialog() (string, error) {
    return runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
        Type:    runtime.InfoDialog,
        Title:   "Information",
        Message: "This is an information message",
        Buttons: []string{"OK"},
    })
}

func (a *App) ShowWarningDialog() (string, error) {
    return runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
        Type:    runtime.WarningDialog,
        Title:   "Warning",
        Message: "This is a warning message",
        Buttons: []string{"OK", "Cancel"},
    })
}

func (a *App) ShowErrorDialog() (string, error) {
    return runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
        Type:    runtime.ErrorDialog,
        Title:   "Error",
        Message: "An error has occurred",
        Buttons: []string{"OK"},
    })
}

func (a *App) ShowQuestionDialog() (string, error) {
    return runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
        Type:          runtime.QuestionDialog,
        Title:         "Confirm Action",
        Message:       "Are you sure you want to continue?",
        Buttons:       []string{"Yes", "No", "Cancel"},
        DefaultButton: "Yes",
        CancelButton:  "Cancel",
    })
}
```

## Event System

Bidirectional event communication between Go and JavaScript

```go
package main

import (
    "context"
    "time"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
    ctx context.Context
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    a.setupEventListeners()
}

func (a *App) setupEventListeners() {
    runtime.EventsOn(a.ctx, "frontend:userAction", func(data ...interface{}) {
        if len(data) > 0 {
            runtime.LogInfo(a.ctx, "Received user action: "+data[0].(string))
        }
    })

    runtime.EventsOnce(a.ctx, "frontend:initialize", func(data ...interface{}) {
        runtime.LogInfo(a.ctx, "Frontend initialized")
        a.EmitWelcomeEvent()
    })

    runtime.EventsOnMultiple(a.ctx, "frontend:click", func(data ...interface{}) {
        runtime.LogInfo(a.ctx, "Click event received")
    }, 5)
}

func (a *App) EmitWelcomeEvent() {
    runtime.EventsEmit(a.ctx, "backend:welcome", map[string]interface{}{
        "message": "Welcome to Wails!",
        "time":    time.Now().Format(time.RFC3339),
    })
}

func (a *App) EmitProgressUpdate(progress int) {
    runtime.EventsEmit(a.ctx, "backend:progress", progress)
}

func (a *App) StartLongRunningTask() {
    go func() {
        for i := 0; i <= 100; i += 10 {
            time.Sleep(500 * time.Millisecond)
            runtime.EventsEmit(a.ctx, "backend:taskProgress", map[string]interface{}{
                "progress": i,
                "status":   "Processing...",
            })
        }
        runtime.EventsEmit(a.ctx, "backend:taskComplete", "Task finished successfully")
    }()
}

func (a *App) UnsubscribeEvents() {
    runtime.EventsOff(a.ctx, "frontend:userAction")
}

func (a *App) UnsubscribeAllEvents() {
    runtime.EventsOffAll(a.ctx)
}
```

JavaScript frontend event handling:

```javascript
import { EventsOn, EventsEmit, EventsOff } from './wailsjs/runtime/runtime';

// Listen for events from Go
EventsOn('backend:welcome', (data) => {
    console.log('Welcome message:', data.message);
    console.log('Time:', data.time);
});

EventsOn('backend:progress', (progress) => {
    console.log('Progress:', progress + '%');
});

EventsOn('backend:taskProgress', (data) => {
    updateProgressBar(data.progress);
    updateStatusText(data.status);
});

EventsOn('backend:taskComplete', (message) => {
    console.log('Task complete:', message);
});

// Emit events to Go
function handleUserAction(action) {
    EventsEmit('frontend:userAction', action);
}

function handleClick() {
    EventsEmit('frontend:click', { x: 100, y: 200 });
}

// Unsubscribe
function cleanup() {
    EventsOff('backend:welcome');
}
```

## Application Configuration

Advanced application options and lifecycle hooks

```go
package main

import (
    "context"
    "embed"
    "github.com/wailsapp/wails/v2"
    "github.com/wailsapp/wails/v2/pkg/logger"
    "github.com/wailsapp/wails/v2/pkg/menu"
    "github.com/wailsapp/wails/v2/pkg/menu/keys"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
    "github.com/wailsapp/wails/v2/pkg/options/linux"
    "github.com/wailsapp/wails/v2/pkg/options/mac"
    "github.com/wailsapp/wails/v2/pkg/options/windows"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

type App struct {
    ctx context.Context
}

func NewApp() *App {
    return &App{}
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    runtime.LogInfo(ctx, "Application started")
}

func (a *App) domReady(ctx context.Context) {
    runtime.LogInfo(ctx, "DOM is ready")
    runtime.WindowCenter(ctx)
}

func (a *App) beforeClose(ctx context.Context) bool {
    dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
        Type:    runtime.QuestionDialog,
        Title:   "Quit?",
        Message: "Are you sure you want to quit?",
        Buttons: []string{"Yes", "No"},
    })
    if err != nil {
        return false
    }
    return dialog != "Yes"
}

func (a *App) shutdown(ctx context.Context) {
    runtime.LogInfo(ctx, "Application is shutting down")
}

func createMenu() *menu.Menu {
    appMenu := menu.NewMenu()

    fileMenu := appMenu.AddSubmenu("File")
    fileMenu.AddText("New", keys.CmdOrCtrl("n"), func(_ *menu.CallbackData) {
        runtime.LogInfo(nil, "New file clicked")
    })
    fileMenu.AddText("Open", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
        runtime.LogInfo(nil, "Open file clicked")
    })
    fileMenu.AddSeparator()
    fileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
        runtime.Quit(nil)
    })

    editMenu := appMenu.AddSubmenu("Edit")
    editMenu.AddText("Cut", keys.CmdOrCtrl("x"), nil)
    editMenu.AddText("Copy", keys.CmdOrCtrl("c"), nil)
    editMenu.AddText("Paste", keys.CmdOrCtrl("v"), nil)

    viewMenu := appMenu.AddSubmenu("View")
    viewMenu.AddCheckbox("Show Sidebar", true, nil, func(_ *menu.CallbackData) {
        runtime.LogInfo(nil, "Toggle sidebar")
    })

    return appMenu
}

func main() {
    app := NewApp()

    err := wails.Run(&options.App{
        Title:             "Advanced Wails App",
        Width:             1280,
        Height:            720,
        MinWidth:          800,
        MinHeight:         600,
        MaxWidth:          1920,
        MaxHeight:         1080,
        DisableResize:     false,
        Fullscreen:        false,
        Frameless:         false,
        StartHidden:       false,
        HideWindowOnClose: false,
        AlwaysOnTop:       false,
        BackgroundColour:  options.NewRGBA(33, 37, 43, 255),

        AssetServer: &assetserver.Options{
            Assets: assets,
        },

        Menu:               createMenu(),
        Logger:             nil,
        LogLevel:           logger.DEBUG,
        LogLevelProduction: logger.ERROR,

        OnStartup:     app.startup,
        OnDomReady:    app.domReady,
        OnBeforeClose: app.beforeClose,
        OnShutdown:    app.shutdown,

        Bind: []interface{}{
            app,
        },

        WindowStartState: options.Normal,

        CSSDragProperty: "--wails-draggable",
        CSSDragValue:    "drag",

        EnableDefaultContextMenu:         false,
        EnableFraudulentWebsiteDetection: false,

        Windows: &windows.Options{
            WebviewIsTransparent: false,
            WindowIsTranslucent:  false,
            DisableWindowIcon:    false,
        },

        Mac: &mac.Options{
            TitleBar: &mac.TitleBar{
                TitlebarAppearsTransparent: true,
                HideTitle:                  false,
                HideTitleBar:               false,
                FullSizeContent:            false,
                UseToolbar:                 false,
                HideToolbarSeparator:       true,
            },
            WebviewIsTransparent: true,
            WindowIsTranslucent:  true,
            About: &mac.AboutInfo{
                Title:   "My App",
                Message: "© 2024 My Company",
            },
        },

        Linux: &linux.Options{
            Icon: []byte{},
        },

        DragAndDrop: &options.DragAndDrop{
            EnableFileDrop:     true,
            DisableWebViewDrop: false,
            CSSDropProperty:    "--wails-drop-target",
            CSSDropValue:       "drop",
        },
    })

    if err != nil {
        println("Error:", err.Error())
    }
}
```

## Logging System

Structured logging from Go backend

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v2/pkg/logger"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
    ctx context.Context
}

func (a *App) DemonstrateLogging() {
    runtime.LogTrace(a.ctx, "This is a trace message")
    runtime.LogDebug(a.ctx, "This is a debug message")
    runtime.LogInfo(a.ctx, "This is an info message")
    runtime.LogWarning(a.ctx, "This is a warning message")
    runtime.LogError(a.ctx, "This is an error message")

    runtime.LogTracef(a.ctx, "User %s performed action at %d", "john", 1234567890)
    runtime.LogDebugf(a.ctx, "Debug info: value=%v", map[string]int{"count": 42})
    runtime.LogInfof(a.ctx, "Processing file: %s", "data.json")
    runtime.LogWarningf(a.ctx, "Memory usage: %d%%", 85)
    runtime.LogErrorf(a.ctx, "Failed to connect to %s:%d", "localhost", 8080)
}

func (a *App) SetLogLevel(level string) {
    var logLevel logger.LogLevel
    switch level {
    case "trace":
        logLevel = logger.TRACE
    case "debug":
        logLevel = logger.DEBUG
    case "info":
        logLevel = logger.INFO
    case "warning":
        logLevel = logger.WARNING
    case "error":
        logLevel = logger.ERROR
    default:
        logLevel = logger.INFO
    }
    runtime.LogSetLogLevel(a.ctx, logLevel)
}
```

## Screen Information

Retrieving screen and display information

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
    ctx context.Context
}

type ScreenInfo struct {
    IsCurrent    bool   `json:"isCurrent"`
    IsPrimary    bool   `json:"isPrimary"`
    Width        int    `json:"width"`
    Height       int    `json:"height"`
    PhysicalWidth int   `json:"physicalWidth"`
    PhysicalHeight int  `json:"physicalHeight"`
}

func (a *App) GetScreens() ([]ScreenInfo, error) {
    screens, err := runtime.ScreenGetAll(a.ctx)
    if err != nil {
        return nil, err
    }

    result := make([]ScreenInfo, len(screens))
    for i, screen := range screens {
        result[i] = ScreenInfo{
            IsCurrent:      screen.IsCurrent,
            IsPrimary:      screen.IsPrimary,
            Width:          screen.Size.Width,
            Height:         screen.Size.Height,
            PhysicalWidth:  screen.PhysicalSize.Width,
            PhysicalHeight: screen.PhysicalSize.Height,
        }
    }

    return result, nil
}

func (a *App) GetPrimaryScreen() (*ScreenInfo, error) {
    screens, err := a.GetScreens()
    if err != nil {
        return nil, err
    }

    for _, screen := range screens {
        if screen.IsPrimary {
            return &screen, nil
        }
    }

    return nil, nil
}
```

## Clipboard Operations

Reading and writing clipboard text

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
    ctx context.Context
}

func (a *App) CopyToClipboard(text string) error {
    return runtime.ClipboardSetText(a.ctx, text)
}

func (a *App) PasteFromClipboard() (string, error) {
    return runtime.ClipboardGetText(a.ctx)
}

func (a *App) AppendToClipboard(text string) error {
    current, err := runtime.ClipboardGetText(a.ctx)
    if err != nil {
        return err
    }
    return runtime.ClipboardSetText(a.ctx, current+text)
}
```

## Browser Integration

Opening URLs in default system browser

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
    ctx context.Context
}

func (a *App) OpenURL(url string) {
    runtime.BrowserOpenURL(a.ctx, url)
}

func (a *App) OpenDocumentation() {
    runtime.BrowserOpenURL(a.ctx, "https://wails.io/docs")
}

func (a *App) OpenGitHub() {
    runtime.BrowserOpenURL(a.ctx, "https://github.com/wailsapp/wails")
}
```

## Environment Information

Retrieving build and platform information

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
    ctx context.Context
}

func (a *App) GetEnvironment() runtime.EnvironmentInfo {
    return runtime.Environment(a.ctx)
}

func (a *App) CheckPlatform() string {
    env := runtime.Environment(a.ctx)
    return env.Platform
}

func (a *App) IsDevelopment() bool {
    env := runtime.Environment(a.ctx)
    return env.BuildType == "dev"
}

func (a *App) IsProduction() bool {
    env := runtime.Environment(a.ctx)
    return env.BuildType == "production"
}
```

## Summary

Wails provides a comprehensive framework for building native desktop applications using Go and modern web technologies. The main use cases include creating lightweight alternatives to Electron apps, building desktop GUIs for Go CLI tools, developing cross-platform applications with native features, and wrapping existing Go services with modern web interfaces. The framework excels at applications requiring native system integration, file operations, and desktop-specific functionality while maintaining the flexibility of web-based UIs.

Integration patterns typically involve embedding the frontend assets into the Go binary, binding Go methods for JavaScript to call, using the bidirectional event system for reactive updates, and leveraging the runtime package for native features like dialogs and window management. The framework supports standard web development workflows with hot reload during development and produces single-binary executables for distribution. Common patterns include using TypeScript with auto-generated bindings, implementing pub/sub patterns with the event system, and structuring Go backends with clear separation between business logic and UI methods exposed to the frontend.
