# Codesk - Claude Code 桌面应用

## 项目概述

Codesk 是一个基于 Wails 框架构建的桌面应用程序，专门用于管理和运行 Claude Code 会话。它是一个现代化的 Go + React/TypeScript 应用程序，提供了直观的图形界面来管理 Claude 项目、会话和设置。

### 核心技术栈

- **后端**: Go (1.23) 使用 Wails v2 框架
- **前端**: React 18 + TypeScript + Tailwind CSS
- **数据库**: SQLite3 (用于数据持久化)
- **构建工具**: Vite (前端构建)
- **UI 组件**: Radix UI + 自定义组件

### 项目架构

```
Codesk/
├── app.go                    # 主应用结构，Wails 绑定接口
├── main.go                   # 应用入口点
├── backend/                  # 后端业务逻辑
│   ├── config/              # 配置管理
│   ├── model/               # 数据模型定义
│   ├── process/             # 进程管理
│   ├── repository/          # 数据访问层
│   └── service/             # 业务服务层
│       ├── agent/           # AI 代理服务
│       ├── checkpoint/      # 检查点管理
│       ├── claude/          # Claude 集成服务
│       ├── mcp/             # 模型控制协议服务
│       └── usage/           # 使用统计服务
├── frontend/                 # 前端应用
│   ├── src/
│   │   ├── components/      # React 组件
│   │   ├── contexts/        # React 上下文
│   │   ├── hooks/           # 自定义 Hooks
│   │   ├── lib/             # 工具库和 API
│   │   ├── stores/          # 状态管理 (Zustand)
│   │   └── types/           # TypeScript 类型定义
│   └── wailsjs/             # Wails 自动生成的绑定
└── build/                   # 构建配置和资源
```

## 构建和运行

### 开发环境

#### 前置要求
- Go 1.23 或更高版本
- Node.js 和 npm
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

#### 开发模式运行
```bash
# 进入项目目录
cd ~/Documents/GitHub/Codesk

# 启动开发服务器（支持热重载）
wails dev
```

开发模式下：
- 前端开发服务器运行在 http://localhost:34115
- 支持前后端代码的热重载
- 可访问浏览器开发者工具调试

#### 构建生产版本
```bash
# 构建可执行文件（macOS）
wails build -clean -platform darwin/universal

# 构建后的文件位于 build/bin/
```

#### 创建 DMG 安装包
```bash
create-dmg \
  --volname "Codesk" \
  --window-pos 200 120 \
  --window-size 500 400 \
  --icon-size 100 \
  --icon "Codesk.app" 100 150 \
  --app-drop-link 300 155 \
  --hide-extension "Codesk.app" \
  "$HOME/Downloads/Codesk_v0.0.1_$(date +%Y%m%d_%H%M).dmg" \
  "build/bin/"
```

### 前后端联动机制

Codesk 使用 Wails 框架实现 Go 后端与 React 前端的无缝集成：

1. **后端服务层** (`backend/service/`): 实现具体的业务逻辑
2. **App 接口层** (`app.go`): 暴露给前端的方法接口
3. **Wails 绑定**: 自动生成 TypeScript 类型定义和 JavaScript 绑定
4. **前端调用**: 通过生成的 `wailsjs/go/main/App.js` 调用后端方法

#### 数据流向示例
```
React 组件 → wailsjs 绑定 → App 结构体方法 → 后端服务 → 返回结果
```

## 核心功能

### Claude 项目管理
- 自动扫描 `~/.claude/projects` 目录
- 项目会话管理
- CLAUDE.md 文件编辑
- 会话历史查看

### AI 代理系统
- 自定义 AI 代理配置
- 代理执行和监控
- 代理模板管理

### 模型控制协议 (MCP)
- MCP 服务器管理
- 协议配置和监控

### 使用统计
- 会话使用统计
- 性能指标监控
- 数据分析仪表板

## 开发约定

### 代码风格
- **Go**: 遵循 Go 标准代码风格，使用 `gofmt`
- **TypeScript**: 使用 ESLint 和 Prettier 配置
- **React**: 函数式组件 + Hooks，使用 TypeScript 严格模式

### 项目结构
- 清晰的关注点分离（后端/前端）
- 模块化服务架构
- 统一的错误处理
- 类型安全的接口定义

### 数据库
- 使用 SQLite3 进行数据持久化
- 数据库文件位置：`~/.config/Codesk/database.db`
- 通过 repository 层进行数据访问

## 依赖管理

### 后端依赖 (go.mod)
- `github.com/wailsapp/wails/v2` - Wails 框架
- `github.com/mattn/go-sqlite3` - SQLite3 驱动
- `github.com/google/uuid` - UUID 生成

### 前端依赖 (package.json)
- `react` + `react-dom` - React 框架
- `@radix-ui/*` - 无障碍 UI 组件
- `tailwindcss` - CSS 框架
- `lucide-react` - 图标库
- `zustand` - 状态管理
- `framer-motion` - 动画库

## 配置管理

### 应用配置
- 配置文件：`~/.config/Codesk/config.json`
- Claude 设置：`~/.claude/settings.json`
- 代理配置：通过服务动态管理

### 环境配置
- 开发环境：`wails dev` 模式
- 生产环境：构建后的可执行文件
- 平台特定配置在 `wails.json` 中定义

## 故障排除

### 常见问题

1. **Claude Code 未找到**
   - 应用会显示配置对话框
   - 需要设置正确的 Claude Code 二进制路径

2. **数据库初始化失败**
   - 检查 `~/.config/Codesk/` 目录权限
   - 确保 SQLite3 驱动正常工作

3. **前端构建失败**
   - 运行 `npm install` 重新安装依赖
   - 检查 Node.js 版本兼容性

### 日志和调试
- 开发模式下查看控制台输出
- 生产模式下日志写入系统日志
- 使用浏览器开发者工具调试前端

## 扩展开发

### 添加新功能
1. 在后端 `service/` 目录创建新服务
2. 在 `app.go` 中暴露接口方法
3. 运行 `wails dev` 生成前端绑定
4. 在前端组件中调用新方法

### 自定义组件
- 组件位于 `frontend/src/components/`
- 使用 Tailwind CSS 进行样式设计
- 遵循 Radix UI 的无障碍设计原则

## 部署和维护

### 版本管理
- 使用语义化版本控制
- 通过 GitHub Actions 自动化构建
- 定期更新依赖包

### 用户数据
- 用户数据存储在 `~/.claude/` 和 `~/.config/Codesk/`
- 定期备份重要数据
- 提供数据迁移工具

---

*此文档基于项目当前状态生成，随着项目发展可能需要更新。*