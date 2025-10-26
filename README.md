# README

## About

This is the official Wails React-TS template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

### 1. build bin
```sh
cd ~/Documents/GitHub/Codesk
wails build -clean -platform darwin/universal
```
### 2. package dmg

```sh
create-dmg \
  --volname "Codesk" \
  --window-pos 200 120 \
  --window-size 400 400 \
  --icon-size 100 \
  --icon "Codesk.app" 200 190 \
  --app-drop-link 600 185 \
  --hide-extension "Codesk.app" \
  "/Codesk_v0.0.1.dmg" \
  "build/bin/"
```

## 前后端联动机制

本项目基于 Wails 框架，实现了 Go 后端与 React 前端的无缝集成。

### 架构概览

```
Frontend (React/TypeScript)
    ↓ IPC 调用
Wails 自动生成的绑定层
    ↓ 
App 结构体 (app.go)
    ↓ 依赖注入
Backend Services (backend/service/)
```

### 联动流程

#### 1. 后端服务层 (Backend Services)
```go
// backend/service/greet_service.go
type GreetService struct {}

func (s *GreetService) Greet(name string) string {
    return fmt.Sprintf("Hello %s, It's show time!", name)
}

// backend/service/time_service.go  
type TimeService struct {}

func (s *TimeService) GetCurrentTime() string {
    return time.Now().Format("2006-01-02 15:04:05")
}
```

#### 2. App 接口层 (app.go)
```go
type App struct {
    ctx          context.Context
    greetService *service.GreetService
    timeService  *service.TimeService
}

// 暴露给前端的方法
func (a *App) Greet(name string) string {
    return a.greetService.Greet(name)
}

func (a *App) GetCurrentTime() string {
    return a.timeService.GetCurrentTime()
}
```

#### 3. Wails 绑定 (main.go)
```go
err := wails.Run(&options.App{
    // ...其他配置
    Bind: []interface{}{
        app,  // 将 App 结构体绑定到前端
    },
})
```

#### 4. 自动生成的 TypeScript 绑定
Wails 自动扫描 App 结构体的导出方法，生成：

**frontend/wailsjs/go/main/App.js**
```javascript
export function Greet(arg1) {
  return window['go']['main']['App']['Greet'](arg1);
}

export function GetCurrentTime() {
  return window['go']['main']['App']['GetCurrentTime']();
}
```

**frontend/wailsjs/go/main/App.d.ts**
```typescript
export function Greet(arg1:string):Promise<string>;
export function GetCurrentTime():Promise<string>;
```

#### 5. 前端调用 (App.tsx)
```typescript
import {Greet, GetCurrentTime} from "../wailsjs/go/main/App";

function greet() {
    Greet(name).then(updateResultText);
}

function getCurrentTime() {
    GetCurrentTime().then(setTimeText);
}
```

### 关键特性

- **零配置 IPC**：无需手动设置 HTTP 服务器或 WebSocket
- **类型安全**：自动生成 TypeScript 定义，编译时类型检查
- **异步调用**：所有后端方法返回 Promise，支持异步操作
- **热重载**：开发模式下支持前后端代码热重载
- **单一二进制**：打包后生成单个可执行文件

### 数据流向

```
用户操作 → React 组件状态更新 → 调用 Wails 绑定方法 
    ↓
window.go.main.App.Method() [JavaScript Bridge]
    ↓  
Go App 结构体方法
    ↓
Backend Service 业务逻辑
    ↓
返回结果 → Promise 回调 → 更新 UI 状态
```

### 开发工作流

1. **后端开发**：在 `backend/service/` 中实现业务逻辑
2. **接口暴露**：在 `app.go` 中添加导出方法
3. **自动绑定**：运行 `wails dev`，Wails 自动生成前端绑定
4. **前端调用**：在 React 组件中导入并调用生成的方法

这种架构让你可以像调用本地函数一样调用后端方法，同时保持了前后端代码的清晰分离。

Ref:
> https://github.com/winfunc/opcode
