Wails IPC适配重写计划：

  前端Tauri IPC调用分析总结

  通过全面分析，我发现了以下Tauri API的使用情况：

  1. Tauri IPC调用类型分析

  A. 核心IPC调用 (`invoke`)
   - 位置：src/lib/apiAdapter.ts、src/components/ProxySettings.tsx、src/components/Agents.tsx、src/components/CCAgents.tsx、src/components/AgentsModal.tsx
   - 功能：前后端通信，调用Rust命令

  B. 事件监听 (`listen`)
   - 位置：src/components/SessionOutputViewer.tsx、src/components/AgentExecution.tsx、src/components/AgentRunOutputViewer.tsx、src/components/ClaudeCodeSession.tsx、src/c
     omponents/claude-code-session/useClaudeMessages.ts
   - 功能：监听后端事件，实现实时数据流

  C. 窗口管理 (`getCurrentWindow`)
   - 位置：src/components/CustomTitlebar.tsx
   - 功能：窗口最小化/最大化/关闭操作

  D. 文件对话框 (`plugin-dialog`)
   - 位置：src/components/TabContent.tsx、src/components/CCAgents.tsx、src/components/Agents.tsx、src/components/ClaudeCodeSession.refactored.tsx、src/components/AgentsMo
     dal.tsx
   - 功能：文件选择、保存对话框

  E. 文件路径转换 (`convertFileSrc`)
   - 位置：src/components/ImagePreview.tsx
   - 功能：将本地文件路径转换为可访问的URL

  F. 外部链接打开 (`plugin-opener`, `plugin-shell`)
   - 位置：src/components/NFOCredits.tsx、src/components/GitHubAgentBrowser.tsx、src/components/ToolWidgets.tsx
   - 功能：打开外部URL和文件

 ╭────────────────────────────────────────────────────────────────────────────────╮
 │ ✅ todo_write 更新待办事项列表（0个待处理，1个进行中，3个已完成）                       │
 │    ·已更新待办事项列表                                                             │
 │      ⎿ ✅ 分析前端代码中所有使用Tauri IPC调用的地方                                  │
 │        ✅ 识别所有@tauri-apps/api和@tauri-apps/plugin的导入和使用                   │
 │        ✅ 制定Wails IPC重写计划                                                   │
 │        ☐ 提供具体的代码迁移方案                                                    │
 ╰────────────────────────────────────────────────────────────────────────────────╯
✦ Wails IPC适配重写详细方案

  阶段1：创建Wails适配层

  1.1 创建 `src/lib/wailsAdapter.ts`
     1 // 替换 src/lib/apiAdapter.ts 中的Tauri专用代码
     2 import { EventsEmit, EventsOn, EventsOff, WindowMinimise, WindowMaximise, WindowUnmaximise, WindowClose } from '@wailsjs/runtime';
     3 
     4 // Wails环境检测
     5 let isWailsEnvironment: boolean | null = null;
     6 
     7 function detectWailsEnvironment(): boolean {
     8   if (isWailsEnvironment !== null) {
     9     return isWailsEnvironment;
    10   }
    11 
    12   // 检查Wails运行时
    13   const isWails = !!(window.runtime || window.go);
    14   isWailsEnvironment = isWails;
    15   return isWails;
    16 }
    17 
    18 // 替换Tauri的invoke调用
    19 export async function wailsCall<T>(command: string, params?: any): Promise<T> {
    20   const isWails = detectWailsEnvironment();
    21 
    22   if (isWails) {
    23     // Wails环境 - 调用Go后端方法
    24     if (window.go && window.go.main && window.go.main.App) {
    25       const app = window.go.main.App;
    26       const methodName = command.replace(/_/g, ''); // 转换为驼峰命名
    27 
    28       if (typeof app[methodName] === 'function') {
    29         return await app[methodName](params);
    30       }
    31     }
    32     throw new Error(`Wails method ${command} not found`);
    33   }
    34 
    35   // Web环境 - 保持原有REST API逻辑
    36   return apiCall<T>(command, params);
    37 }
    38 
    39 // 替换Tauri事件监听
    40 export function wailsListen(eventName: string, callback: (data: any) => void): () => void {
    41   if (detectWailsEnvironment()) {
    42     const cleanup = EventsOn(eventName, callback);
    43     return () => EventsOff(eventName, cleanup);
    44   }
    45 
    46   // Web环境 - 使用自定义事件
    47   const handler = (e: CustomEvent) => callback(e.detail);
    48   window.addEventListener(eventName, handler as EventListener);
    49   return () => window.removeEventListener(eventName, handler as EventListener);
    50 }
    51 
    52 // 窗口操作
    53 export const wailsWindow = {
    54   minimize: async () => {
    55     if (detectWailsEnvironment() && WindowMinimise) {
    56       await WindowMinimise();
    57     }
    58   },
    59   maximize: async () => {
    60     if (detectWailsEnvironment() && WindowMaximise) {
    61       await WindowMaximise();
    62     }
    63   },
    64   unmaximize: async () => {
    65     if (detectWailsEnvironment() && WindowUnmaximise) {
    66       await WindowUnmaximise();
    67     }
    68   },
    69   close: async () => {
    70     if (detectWailsEnvironment() && WindowClose) {
    71       await WindowClose();
    72     }
    73   }
    74 };
    75 
    76 // 文件对话框
    77 export const wailsDialog = {
    78   open: async (options: any) => {
    79     if (detectWailsEnvironment() && window.go && window.go.main && window.go.main.App) {
    80       return await window.go.main.App.openFileDialog(options);
    81     }
    82     // Web环境 - 使用原生input
    83     return new Promise((resolve) => {
    84       const input = document.createElement('input');
    85       input.type = 'file';
    86       input.multiple = options.multiple || false;
    87       if (options.directory) {
    88         input.webkitdirectory = true;
    89       }
    90       input.onchange = () => {
    91         resolve(input.files ? Array.from(input.files).map(f => f.path || f.name) : null);
    92       };
    93       input.click();
    94     });
    95   },
    96   save: async (options: any) => {
    97     if (detectWailsEnvironment() && window.go && window.go.main && window.go.main.App) {
    98       return await window.go.main.App.saveFileDialog(options);
    99     }
   100     // Web环境 - 使用原生下载
   101     return null;
   102   }
   103 };
   104 
   105 // 外部链接打开
   106 export const wailsShell = {
   107   openExternal: async (url: string) => {
   108     if (detectWailsEnvironment() && window.go && window.go.main && window.go.main.App) {
   109       return await window.go.main.App.openExternal(url);
   110     }
   111     // Web环境 - 使用window.open
   112     window.open(url, '_blank');
   113   }
   114 };

  阶段2：逐步替换Tauri API调用

  2.1 替换所有 `invoke` 调用
   - 在 src/lib/api.ts 中，将所有 apiCall 调用改为 wailsCall
   - 更新所有组件中的直接 invoke 调用

  2.2 替换事件监听
   - 将 listen 调用替换为 wailsListen
   - 更新事件处理逻辑

  2.3 替换窗口操作
   - 在 CustomTitlebar.tsx 中替换 getCurrentWindow 相关调用

  2.4 替换文件对话框
   - 在 TabContent.tsx、CCAgents.tsx 等文件中替换文件对话框调用

  阶段3：具体文件修改清单

  3.1 核心文件
   - src/lib/apiAdapter.ts → 替换为 wailsAdapter.ts
   - src/lib/api.ts → 更新所有调用为 wailsCall

  3.2 组件文件
   - src/components/CustomTitlebar.tsx → 替换窗口操作
   - src/components/ProxySettings.tsx → 替换 invoke 调用
   - src/components/CCAgents.tsx → 替换文件对话框和 invoke
   - src/components/Agents.tsx → 替换文件对话框和 invoke
   - src/components/AgentsModal.tsx → 替换文件对话框和 invoke
   - src/components/TabContent.tsx → 替换文件对话框
   - src/components/ImagePreview.tsx → 替换文件路径转换
   - src/components/NFOCredits.tsx → 替换外部链接打开
   - src/components/GitHubAgentBrowser.tsx → 替换外部链接打开
   - src/components/ToolWidgets.tsx → 替换外部链接打开

  3.3 事件监听组件
   - src/components/SessionOutputViewer.tsx
   - src/components/AgentExecution.tsx
   - src/components/AgentRunOutputViewer.tsx
   - src/components/ClaudeCodeSession.tsx
   - src/components/claude-code-session/useClaudeMessages.ts

  阶段4：Web环境兼容性

  保持现有的 apiAdapter.ts 中的Web环境逻辑，确保在浏览器中仍然可以工作。



IPC 两种通讯方式：
`frontend/wailsjs/go/main/App.d.ts` 确实是 **调用 Go 方法** 的接口；
而 `@wailsio/runtime` 则是 **控制运行时环境 + 事件系统 + 窗口行为** 的接口。

我们来详细拆解一下 👇

---

## 🧭 一、Wails 有两种通信机制

| 通信方向               | 主要用途             | 使用方式                              | 示例                                 |
| ------------------ | ---------------- | --------------------------------- | ---------------------------------- |
| ✅ Go ⇄ JS 方法调用     | 调用 Go 暴露的函数      | 使用 `frontend/wailsjs/go/.../*.ts` | `window.backend.App.Greet("Jack")` |
| ⚙️ 运行时交互（事件、窗口、环境） | 控制窗口、监听事件、获取系统信息 | 使用 `@wailsio/runtime`  | `WindowSetTitle()`, `EventsEmit()` |

---

## 🧩 二、`frontend/wailsjs/go/...` 是什么？

在你运行：

```bash
wails generate module
# 或 wails dev / build 时自动生成
```

Wails 会在前端项目中生成一套 **类型定义和包装代码**，一般位于：

```
frontend/wailsjs/go/
  ├── main/App.d.ts
  ├── main/App.js
  ├── models.ts
```

这相当于 **Go 函数在 JS 端的自动绑定层**。

### 🧠 作用

让你在前端能直接调用 Go 方法，就像调用 JS 函数一样：

```js
import { Greet } from '../../wailsjs/go/main/App'

Greet("Jack").then(result => console.log(result))
```

或使用全局别名（若你配置过）：

```js
window.backend.App.Greet("Jack")
```

这些函数是由 Wails 根据你 Go 代码里的 `App` 结构体自动生成的桥接层。

---

## 🧠 三、那 `@wailsio/runtime` 干什么用？

它不是针对某个 Go 结构体的，而是提供**全局运行时功能**。

包括：

| 类别        | 函数                                                                  | 用途              |
| --------- | ------------------------------------------------------------------- | --------------- |
| 🪟 窗口控制   | `WindowSetTitle(title)`、`WindowMinimise()`、`WindowToggleMaximise()` | 控制窗口            |
| 🔄 事件系统   | `EventsOn(name, callback)`、`EventsEmit(name, data)`                 | Go ⇄ 前端事件通知     |
| ⚙️ 环境信息   | `Environment()`                                                     | 获取当前平台、构建信息     |
| 🧾 日志     | `LogInfo(msg)`、`LogError(msg)`                                      | 输出日志到 Wails 控制台 |
| 📡 后端通信辅助 | `BrowserOpenURL(url)`                                               | 打开外部链接等系统功能     |

它对应的是 Go 端的：

```go
"github.com/wailsapp/wails/v2/pkg/runtime"
```

---

## 🧩 四、两者的区别可以这么理解：

| 对比项  | `wailsjs/go/...`  | `@wailsio/runtime`        |
| ---- | ----------------- | ------------------------- |
| 来源   | 自动生成的桥接代码         | 官方 NPM 包                  |
| 通信类型 | 前端调用 Go 方法        | 前后端事件通信 / 控制运行时      |
| 调用目标 | 你的 Go 业务逻辑函数      | Wails 内部运行时（窗口、事件等）|
| 自动生成 | ✅ 是（每次 build 时更新） | ❌ 否（你自己引入）           |
| 示例   | `Greet("Jack")`   | `WindowSetTitle("Hello")` |

---

## 🧱 五、两者配合使用的典型例子

**Go 后端：**

```go
func (a *App) Startup(ctx context.Context) {
    a.ctx = ctx
    runtime.EventsEmit(ctx, "ready", "Go 启动完成！")
}

func (a *App) Greet(name string) string {
    return "你好, " + name
}
```

**前端：**

```js
import { EventsOn } from "@wailsio/runtime"
import { Greet } from "../wailsjs/go/main/App"

EventsOn("ready", (msg) => {
  console.log("来自 Go 的事件:", msg)
})

Greet("Jack").then(res => console.log(res))
```

🧩 解释：

* `Greet` → 调用 Go 方法（业务函数）
* `EventsOn` → 监听运行时事件（例如 Go 的启动完成通知）

---

## 📦 六、为什么 Wails 要拆成两层？

因为它想**分离职责**：

* `wailsjs/go/*` 依赖具体的 Go 模块和函数签名，会随你的业务代码变化；
* `@wailsio/runtime` 是固定的运行时 API，不受项目结构变化影响。

所以：

* `frontend/wailsjs/...` 每次改 Go 代码后都要重新生成；
* `@wailsio/runtime` 永远不用重建，是稳定的 SDK。

---

## ✅ 总结一句话

> * 想调用你自己写的 Go 函数 → 用 `frontend/wailsjs/go/...`
> * 想控制窗口、监听事件、发消息 → 用 `@wailsio/runtime`
业务中大部分都是用`frontend/wailsjs/go/...`


