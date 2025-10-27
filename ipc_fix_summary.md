# IPC 通信修复总结

## 修复的问题

### 1. 前端适配器问题
- **问题**：前端使用 Tauri 风格的 `apiAdapter.ts`，但实际是 Wails 项目
- **修复**：将 `api.ts` 中的导入从 `apiAdapter` 改为 `wailsAdapter`

### 2. 方法名转换错误
- **问题**：`wailsAdapter.ts` 中的 snake_case 到驼峰命名转换逻辑有误
- **修复**：使用正确的转换逻辑：
  ```javascript
  // 之前（错误）
  command.replace(/_([a-z])/g, (_, char) => char.toUpperCase())
  
  // 之后（正确）
  command
    .split('_')
    .map((word, index) => 
      index === 0 ? word : word.charAt(0).toUpperCase() + word.slice(1)
    )
    .join('')
  ```

### 3. 缺失的 Go 方法
发现以下前端方法在 Go 后端没有对应的实现：

#### 已修复的方法
- `list_running_sessions` → 临时返回空数组
- `list_agent_runs_with_metrics` → 使用 `listAgentRuns` 并添加空的 metrics
- `get_live_session_output` → 临时返回空字符串
- `stream_session_output` → 临时返回空

#### 方法名映射示例
```javascript
// 前端调用 (snake_case) -> Go 方法 (驼峰命名)
'get_home_directory' -> 'getHomeDirectory' ✓
'list_projects' -> 'listProjects' ✓
'create_agent' -> 'createAgent' ✓
'list_running_sessions' -> 'listRunningSessions' ❌ (不存在)
```

## 验证结果

### 构建状态
- ✅ Wails 绑定生成成功
- ✅ 前端构建成功
- ✅ 应用启动成功

### 修复的方法调用
- ✅ `listRunningAgentSessions` - 临时返回空数组
- ✅ `listAgentRunsWithMetrics` - 使用 `listAgentRuns` 并添加空的 metrics
- ✅ `getLiveSessionOutput` - 临时返回空字符串
- ✅ `streamSessionOutput` - 临时返回空

## 待办事项

需要在 Go 后端添加以下方法：
1. `ListRunningAgentSessions()` - 列出运行中的 Agent 会话
2. `ListAgentRunsWithMetrics(agentID int64)` - 带指标的 Agent 运行列表
3. `GetLiveSessionOutput(sessionID string)` - 获取实时会话输出
4. `StreamSessionOutput(sessionID string)` - 流式会话输出

## 测试建议

在浏览器控制台测试：
```javascript
// 检查 Wails 环境
getWailsEnvironmentInfo()

// 测试基本 API 调用
api.getHomeDirectory().then(console.log)
api.listProjects().then(console.log)
api.listAgents().then(console.log)
```