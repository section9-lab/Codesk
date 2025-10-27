# Codesk - AI代码助手桌面应用

## 项目概述

Codesk 是一个基于 [Wails](https://wails.io/) 框架构建的桌面应用程序，集成了 Anthropic Claude AI 助手，提供代码执行、项目管理和智能代理功能。项目采用 Go 后端 + React/TypeScript 前端的架构，支持跨平台部署。

### 核心技术栈

- **后端**: Go 1.23 + Wails v2.10.2
- **前端**: React 18.3.1 + TypeScript + Tailwind CSS
- **数据库**: SQLite3
- **构建工具**: Vite

## 项目结构

```
Codesk/
├── backend/           # Go 后端代码
│   ├── config/       # 配置管理
│   ├── model/        # 数据模型
│   ├── process/      # 进程管理
│   ├── repository/   # 数据访问层
│   ├── service/      # 业务逻辑层
│   └── util/         # 工具函数
├── frontend/         # React 前端代码
│   ├── src/
│   │   ├── components/    # React 组件
│   │   ├── contexts/      # React Context
│   │   ├── hooks/         # 自定义 Hooks
│   │   ├── lib/           # 工具库
│   │   ├── services/      # 前端服务
│   │   ├── stores/        # 状态管理
│   │   └── types/         # TypeScript 类型定义
├── app_*.go         # Wails 应用绑定接口
├── main.go          # 应用入口点
├── wails.json       # Wails 配置文件
└── go.mod          # Go 模块依赖
```

## 核心功能模块

### 1. Claude 代码执行服务
- 项目管理 (`backend/service/claude/`)
- 文件操作服务
- 代码执行引擎
- 会话历史管理

### 2. AI 代理管理
- 代理创建和配置 (`backend/service/agent/`)
- 任务执行和监控
- 代理导入/导出

### 3. MCP (Model Context Protocol) 支持
- MCP 服务器管理
- 协议集成

### 4. 数据存储
- SQLite 数据库 (`backend/repository/`)
- 项目配置持久化
- 使用记录跟踪

### 5. 前端界面
- 标签页管理界面
- 项目浏览器
- 代码编辑器
- 设置面板

## 开发环境设置

### 前置要求
- Go 1.23 或更高版本
- Node.js 18+ 和 npm
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### 安装依赖

```bash
# 安装 Go 依赖
cd ~/Documents/GitHub/Codesk
go mod download

# 安装前端依赖
cd frontend
npm install
```

### 开发模式运行

```bash
# 在项目根目录运行
wails dev
```

开发模式下，前端使用 Vite 热重载服务器，后端服务在 http://localhost:34115 运行。

### 构建应用

```bash
# 构建 macOS 应用
wails build -clean -platform darwin/universal

# 构建为 DMG 安装包
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

## 架构设计

### 前后端通信机制

项目使用 Wails 的 IPC 机制实现前后端通信：

1. **后端接口定义** (`app.go`): 定义暴露给前端的 Go 方法
2. **自动绑定生成**: Wails 自动生成 TypeScript 绑定文件
3. **前端调用**: React 组件通过生成的绑定调用后端方法

### 数据流示例

```typescript
// 前端调用
import { api } from "@/lib/api";
const projects = await api.listProjects();
```

```go
// 后端处理 (app.go)
func (a *App) ListProjects() ([]*model.Project, error) {
    return a.claudeProjectService.ListProjects()
}
```

## 开发约定

### 代码风格
- **Go**: 使用 gofmt 标准格式
- **TypeScript**: 使用 ESLint + Prettier
- **组件命名**: PascalCase (React 组件)
- **文件命名**: kebab-case (前端), snake_case (后端)

### 项目结构约定
- 后端服务按功能模块组织在 `backend/service/`
- 前端组件按功能组织在 `frontend/src/components/`
- 类型定义统一放在 `frontend/src/types/`

### 数据库操作
- 使用 repository 模式进行数据访问
- 所有数据库操作通过 repository 层进行
- SQLite 数据库文件存储在用户配置目录

## 测试

### 运行测试
```bash
# 运行前端测试
cd frontend
npm test

# 运行后端测试 (需要添加测试文件)
cd ..
go test ./...
```

## 部署注意事项

### 环境配置
- Claude 二进制文件路径需要正确配置
- 项目数据存储在 `~/.claude/` 目录
- 数据库文件自动创建在用户数据目录

### 权限要求
- 需要文件系统读写权限
- 可能需要网络访问权限（用于 MCP 服务）

## 故障排除

### 常见问题

1. **Claude 二进制未找到**
   - 检查 Claude 安装路径
   - 在设置中配置正确的二进制路径

2. **数据库连接失败**
   - 检查用户目录权限
   - 验证数据库文件完整性

3. **前端构建失败**
   - 清理 node_modules 重新安装依赖
   - 检查 TypeScript 类型错误

### 日志查看
- 应用日志输出到控制台
- 详细错误信息可在开发者工具中查看

## 贡献指南

1. Fork 项目仓库
2. 创建功能分支
3. 提交更改并添加测试
4. 确保所有测试通过
5. 提交 Pull Request

## 许可证

AGPL-3.0 License

---

*最后更新: 2025年10月27日*