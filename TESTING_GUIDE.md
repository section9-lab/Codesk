# Codesk 后端测试指南

## 测试环境准备

### 1. 安装依赖

```bash
# 确保 Go 版本 >= 1.23
go version

# 安装依赖
go mod tidy
```

### 2. 数据库初始化

数据库会在应用启动时自动初始化，位置：`~/.codesk/codesk.db`

### 3. Claude Code 安装

确保系统已安装 Claude Code CLI：

```bash
# 检查 Claude 是否安装
claude --version

# 如果未安装，请访问 https://claude.ai/download
```

## 功能测试清单

### ✅ Phase 1: 基础设施

#### 数据库
- [ ] 数据库文件创建成功
- [ ] 所有表创建成功（agents, agent_runs, app_settings, slash_commands, usage_stats, checkpoints）
- [ ] 索引创建成功

#### 配置管理
- [ ] 配置目录创建 (~/.codesk)
- [ ] Claude 目录识别 (~/.claude)
- [ ] 路径管理正确

#### 工具函数
- [ ] 文件读写正常
- [ ] JSONL 解析正常
- [ ] 哈希计算正确
- [ ] Claude 二进制查找成功

### ✅ Phase 2: 核心服务

#### Claude 项目管理
```bash
# 测试命令（在 wails dev 运行后，通过前端调用）
ListProjects()
GetProjectSessions("/path/to/project")
GetClaudeSettings()
LoadSessionHistory("session-id", "/path/to/project")
```

**预期结果**:
- [ ] 能列出 ~/.claude/projects 下的所有项目
- [ ] 能获取项目的会话列表
- [ ] 能读取 ~/.claude/settings.json
- [ ] 能加载会话的 JSONL 历史

#### Claude 文件操作
```bash
ListDirectoryContents("/path/to/project")
SearchFiles("/path/to/project", "*.go", 10)
FindClaudeMdFiles("/path/to/project")
```

**预期结果**:
- [ ] 能列出目录内容
- [ ] 能搜索文件
- [ ] 能找到 CLAUDE.md 文件

#### Claude 执行
```bash
ExecuteClaude("/path/to/project", "Hello world", "sonnet")
GetClaudeSessionStatus("session-id", "/path/to/project")
CancelClaude("session-id")
```

**预期结果**:
- [ ] 能启动 Claude Code 进程
- [ ] 能获取会话状态
- [ ] 能取消执行

#### Agent 管理
```bash
CreateAgent({name: "Test Agent", icon: "🤖", ...})
ListAgents()
ExecuteAgent(1, "Test task", "/path/to/project")
```

**预期结果**:
- [ ] 能创建 Agent
- [ ] 能列出所有 Agents
- [ ] 能执行 Agent

### ✅ Phase 3: 高级功能

#### Checkpoint 系统
```bash
CreateCheckpoint("project-id", "session-id", "/path/to/project", "Test checkpoint", nil)
ListCheckpoints("project-id", "session-id", "/path/to/project")
RestoreCheckpoint("project-id", "session-id", "/path/to/project", "checkpoint-id")
```

**预期结果**:
- [ ] 能创建检查点
- [ ] 文件快照保存成功
- [ ] 能恢复检查点
- [ ] 时间线树结构正确

#### MCP 管理
```bash
MCPListServers()
MCPAddServer("test-server", {...})
MCPGetServerStatus("test-server")
```

**预期结果**:
- [ ] 能读取 MCP 配置
- [ ] 能添加/删除服务器
- [ ] 能获取服务器状态

#### 使用统计
```bash
GetUsageStats()
GetUsageByDateRange("2024-01-01", "2024-12-31")
GetSessionStats("session-id", "/path/to/project")
```

**预期结果**:
- [ ] 能统计总体使用情况
- [ ] 能按日期范围统计
- [ ] 能计算成本

### ✅ Phase 4: Wails 集成

#### IPC 接口
- [ ] 所有方法都能从前端调用
- [ ] 参数序列化正常
- [ ] 返回值反序列化正常
- [ ] 错误正确传递到前端

#### 初始化
- [ ] 应用启动时配置初始化成功
- [ ] 数据库自动创建
- [ ] 代理设置自动加载

## 性能测试

### 1. 数据库性能

```bash
# 测试大量数据插入
for i in {1..1000}; do
  CreateAgent(...)
done

# 测试查询性能
time ListAgents()
```

**性能目标**:
- [ ] 1000 条记录插入 < 5 秒
- [ ] 查询 1000 条记录 < 100ms

### 2. 文件操作性能

```bash
# 测试大目录扫描
time ListDirectoryContents("/large/directory")

# 测试文件搜索
time SearchFiles("/large/directory", "*.go", 100)
```

**性能目标**:
- [ ] 扫描 1000 个文件 < 1 秒
- [ ] 搜索 1000 个文件 < 500ms

### 3. Checkpoint 性能

```bash
# 测试大项目检查点创建
time CreateCheckpoint(...) # 100+ 文件的项目
```

**性能目标**:
- [ ] 100 个文件的检查点创建 < 3 秒
- [ ] 压缩率 > 50%

## 集成测试场景

### 场景 1: 完整的 Agent 工作流

1. 创建 Agent
2. 执行 Agent
3. 监控执行状态
4. 获取输出
5. 查看运行历史

### 场景 2: Checkpoint 工作流

1. 开始会话
2. 跟踪文件修改
3. 创建检查点
4. 继续修改
5. 创建第二个检查点
6. 恢复到第一个检查点
7. 从检查点创建分支

### 场景 3: 使用统计工作流

1. 执行多个会话
2. 查看总体统计
3. 按日期范围查询
4. 查看单个会话详情

## 错误处理测试

### 1. 文件不存在
```bash
ReadClaudeMdFile("/nonexistent/file")
# 预期: 返回错误 "file not found"
```

### 2. 数据库错误
```bash
GetAgent(999999)
# 预期: 返回错误 "agent not found"
```

### 3. 进程错误
```bash
ExecuteClaude("/invalid/path", "task", "model")
# 预期: 返回错误 "failed to start claude"
```

### 4. 权限错误
```bash
SaveClaudeMdFile("/root/file", "content")
# 预期: 返回错误 "permission denied"
```

## 边界条件测试

### 1. 空数据
- [ ] 空项目列表
- [ ] 空会话历史
- [ ] 空检查点列表

### 2. 大数据
- [ ] 1000+ 个 Agents
- [ ] 100+ 个检查点
- [ ] 10MB+ 的 JSONL 文件

### 3. 特殊字符
- [ ] 文件名包含空格
- [ ] 文件名包含特殊字符
- [ ] 路径包含中文

## 内存泄漏测试

```bash
# 长时间运行测试
for i in {1..1000}; do
  ExecuteClaude(...)
  sleep 1
  CancelClaude(...)
done

# 监控内存使用
# 预期: 内存使用稳定，无持续增长
```

## 并发测试

```bash
# 同时执行多个操作
parallel ::: \
  "ExecuteClaude(...)" \
  "CreateCheckpoint(...)" \
  "GetUsageStats()" \
  "ListAgents()"

# 预期: 无数据竞争，无死锁
```

## 跨平台测试

### macOS
- [ ] 所有功能正常
- [ ] Claude 二进制查找成功
- [ ] 文件权限正确

### Linux
- [ ] 所有功能正常
- [ ] 路径处理正确

### Windows
- [ ] 所有功能正常
- [ ] 路径分隔符处理正确

## 测试报告模板

```markdown
## 测试日期: YYYY-MM-DD
## 测试人员: [Name]
## 测试环境:
- OS: [macOS/Linux/Windows]
- Go Version: [version]
- Wails Version: [version]

## 测试结果:

### 功能测试
- Phase 1: ✅/❌
- Phase 2: ✅/❌
- Phase 3: ✅/❌
- Phase 4: ✅/❌

### 性能测试
- 数据库性能: ✅/❌
- 文件操作性能: ✅/❌
- Checkpoint 性能: ✅/❌

### 问题列表
1. [问题描述]
2. [问题描述]

### 建议
1. [优化建议]
2. [优化建议]
```

## 自动化测试（未来）

```go
// 示例单元测试
func TestAgentService_CreateAgent(t *testing.T) {
    service := agent.NewAgentService()
    
    agent := &model.Agent{
        Name: "Test Agent",
        Icon: "🤖",
        SystemPrompt: "Test prompt",
        Model: "sonnet",
    }
    
    err := service.CreateAgent(agent)
    assert.NoError(t, err)
    assert.NotNil(t, agent.ID)
}
```

## 测试工具

### 1. 数据库查看
```bash
sqlite3 ~/.codesk/codesk.db
.tables
.schema agents
SELECT * FROM agents;
```

### 2. 日志查看
```bash
# Wails 开发模式会输出日志到控制台
wails dev
```

### 3. 性能分析
```bash
# Go pprof
go tool pprof http://localhost:6060/debug/pprof/profile
```

## 测试完成标准

- [ ] 所有功能测试通过
- [ ] 性能测试达标
- [ ] 无内存泄漏
- [ ] 无数据竞争
- [ ] 跨平台兼容
- [ ] 错误处理完善
- [ ] 文档完整
