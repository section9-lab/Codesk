/**
 * Wails 适配器 - 使用 Wails 自动生成的 IPC 绑定
 *
 * 这个模块完全使用 Wails 自动生成的 IPC 调用，移除对 @wailsio/runtime 的依赖
 * 所有窗口、对话框、事件操作都通过自动生成的 Go 方法调用
 */

// 导入 Wails 自动生成的 IPC 绑定
import {
  WindowMinimise as WailsWindowMinimise,
  WindowMaximise as WailsWindowMaximise,
  WindowUnmaximise as WailsWindowUnmaximise,
  WindowGetPosition as WailsWindowGetPosition,
  WindowSetPosition as WailsWindowSetPosition,
  Quit as WailsQuit,
  BrowserOpenURL as WailsBrowserOpenURL,
} from "../../wailsjs/runtime";

// 导入 Wails 自动生成的 App 方法
import {
  OpenFileDialog,
  SaveFileDialog,
  OpenExternal,
} from "../../wailsjs/go/main/App";

// Type for unlisten function
export type UnlistenFn = () => void;

// 扩展 Window 接口以支持 Wails
declare global {
  interface Window {
    runtime?: any;
    go?: {
      main: {
        App: any;
      };
    };
  }
}

// 环境检测
let isWailsEnvironment: boolean | null = null;

/**
 * 检测是否在 Wails 环境中运行
 */
function detectWailsEnvironment(): boolean {
  if (isWailsEnvironment !== null) {
    return isWailsEnvironment;
  }

  // 首先检查是否在浏览器环境中
  if (typeof window === "undefined") {
    isWailsEnvironment = false;
    return false;
  }

  // 检查 Wails 特定的指示器
  const hasRuntime = !!window.runtime;
  const hasGo = !!(window.go && window.go.main && window.go.main.App);

  // 检查是否在 Electron 应用中（Wails 使用 Electron）
  const isElectron = navigator.userAgent.toLowerCase().indexOf("electron") > -1;

  // 综合判断
  const isWails = hasRuntime || hasGo || isElectron;

  console.log("[detectWailsEnvironment] Detection results:", {
    hasRuntime,
    hasGo,
    isElectron,
    userAgent: navigator.userAgent,
    url: window.location.href,
    finalResult: isWails,
  });

  isWailsEnvironment = isWails;
  return isWails;
}

// 导入所有 Wails 自动生成的 App 方法
import * as AppMethods from "../../wailsjs/go/main/App";

/**
 * Wails 方法映射表
 * 将 snake_case 命令映射到对应的生成方法
 */
const wailsMethodMap: Record<string, keyof typeof AppMethods | string> = {
  // Claude 相关方法
  get_home_directory: "GetHomeDirectory",
  list_projects: "ListProjects",
  create_project: "CreateProject",
  get_project_sessions: "GetProjectSessions",
  get_claude_settings: "GetClaudeSettings",
  save_claude_settings: "SaveClaudeSettings",
  open_new_session: "OpenNewSession",
  get_system_prompt: "not_implemented", // 需要在 Go 后端添加
  save_system_prompt: "not_implemented", // 需要在 Go 后端添加
  check_claude_version: "CheckClaudeVersion",
  find_claude_md_files: "FindClaudeMdFiles",
  read_claude_md_file: "ReadClaudeMdFile",
  save_claude_md_file: "SaveClaudeMdFile",
  execute_claude_code: "ExecuteClaude",
  continue_claude_code: "ContinueClaude",
  resume_claude_code: "ResumeClaude",
  cancel_claude_execution: "CancelClaude",
  list_running_claude_sessions: "ListRunningClaudeSessions",
  get_claude_session_output: "GetClaudeSessionOutput",
  get_claude_session_status: "GetClaudeSessionStatus",
  list_directory_contents: "ListDirectoryContents",
  search_files: "SearchFiles",

  // Agent 相关方法
  list_agents: "ListAgents",
  create_agent: "CreateAgent",
  update_agent: "UpdateAgent",
  delete_agent: "DeleteAgent",
  get_agent: "GetAgent",
  export_agent: "ExportAgent",
  export_agent_to_json: "ExportAgentToJSON",
  import_agent: "ImportAgent",
  import_agent_from_json: "ImportAgentFromJSON",
  execute_agent: "ExecuteAgent",
  list_agent_runs: "ListAgentRuns",
  get_agent_run: "GetAgentRun",
  get_agent_session_output: "GetAgentSessionOutput",
  get_agent_session_status: "GetAgentSessionStatus",
  kill_agent_session: "KillAgentSession",
  cleanup_finished_processes: "CleanupFinishedProcesses",
  load_session_history: "LoadSessionHistory",

  // Checkpoint 相关方法
  create_checkpoint: "CreateCheckpoint",
  restore_checkpoint: "RestoreCheckpoint",
  list_checkpoints: "ListCheckpoints",
  fork_from_checkpoint: "ForkFromCheckpoint",
  get_session_timeline: "GetSessionTimeline",
  update_checkpoint_settings: "UpdateCheckpointSettings",
  get_checkpoint_diff: "GetCheckpointDiff",
  track_checkpoint_message: "TrackCheckpointMessage",
  cleanup_old_checkpoints: "CleanupOldCheckpoints",
  get_checkpoint_settings: "GetCheckpointSettings",

  // Usage 统计方法
  get_usage_stats: "GetUsageStats",
  get_usage_by_date_range: "GetUsageByDateRange",
  get_session_stats: "GetSessionStats",
  get_usage_details: "not_implemented", // 需要在 Go 后端添加

  // MCP 相关方法
  mcp_add: "MCPAddServer",
  mcp_list: "MCPListServers",
  mcp_get: "MCPGetServer",
  mcp_remove: "MCPRemoveServer",
  mcp_add_json: "MCPAddServerFromJSON",
  mcp_add_from_claude_desktop: "not_implemented", // 需要在 Go 后端添加
  mcp_serve: "not_implemented", // 需要在 Go 后端添加
  mcp_test_connection: "MCPTestConnection",
  mcp_reset_project_choices: "MCPResetProjectChoices",
  mcp_get_server_status: "MCPGetServerStatus",
  mcp_read_project_config: "MCPReadProjectConfig",
  mcp_save_project_config: "MCPSaveProjectConfig",

  // Storage 相关方法
  storage_list_tables: "StorageListTables",
  storage_read_table: "StorageReadTable",
  storage_update_row: "StorageUpdateRow",
  storage_delete_row: "StorageDeleteRow",
  storage_insert_row: "StorageInsertRow",
  storage_execute_sql: "StorageExecuteSQL",
  storage_reset_database: "StorageResetDatabase",
  set_setting: "SetSetting",

  // Slash Commands 相关方法
  slash_commands_list: "ListSlashCommands",
  slash_command_get: "GetSlashCommand",
  slash_command_save: "SaveSlashCommand",
  slash_command_delete: "DeleteSlashCommand",

  // 其他方法
  get_claude_binary_path: "GetClaudeBinaryPath", // ✅ 已实现
  set_claude_binary_path: "SetClaudeBinaryPath", // ✅ 已实现
  list_claude_installations: "ListClaudeInstallations",
  get_recently_modified_files: "GetRecentlyModifiedFiles",
  get_proxy_settings: "GetProxySettings",
  save_proxy_settings: "SaveProxySettings",
  validate_hook_command: "not_implemented", // 需要在 Go 后端添加
  get_hooks_config: "not_implemented", // 需要在 Go 后端添加
  update_hooks_config: "not_implemented", // 需要在 Go 后端添加
  track_session_messages: "not_implemented", // 需要在 Go 后端添加
  track_file_modification: "TrackFileModification",
  clear_checkpoint_manager: "ClearCheckpointManager", // ✅ 已实现
  check_auto_checkpoint: "not_implemented", // 需要在 Go 后端添加
  fetch_github_agents: "not_implemented", // 需要在 Go 后端添加
  fetch_github_agent_content: "not_implemented", // 需要在 Go 后端添加
  import_agent_from_github: "not_implemented", // 需要在 Go 后端添加
  get_agent_run_with_real_time_metrics: "not_implemented", // 需要在 Go 后端添加
  list_agent_runs_with_metrics: "not_implemented", // 需要在 Go 后端添加
  list_running_agent_sessions: "not_implemented", // 需要在 Go 后端添加
  get_session_output: "not_implemented", // 需要在 Go 后端添加
  get_live_session_output: "not_implemented", // 需要在 Go 后端添加
  stream_session_output: "not_implemented", // 需要在 Go 后端添加
  get_session_status: "not_implemented", // 需要在 Go 后端添加
};

/**
 * 使用 Wails 自动生成的绑定
 */
export async function wailsCall<T>(command: string, params?: any): Promise<T> {
  const isWails = detectWailsEnvironment();

  if (!isWails) {
    // 如果不是 Wails 环境，抛出错误
    throw new Error(
      `Wails environment not detected. Cannot call method: ${command}`,
    );
  }

  console.log(`[Wails] Calling: ${command}`, params);

  // 查找对应的方法名
  const methodName = wailsMethodMap[command];
  if (!methodName) {
    throw new Error(`Method mapping not found for command: ${command}`);
  }

  // 检查方法是否已实现
  if (methodName === "not_implemented") {
    throw new Error(
      `Method ${command} is not implemented in the Go backend yet`,
    );
  }

  // 获取对应的方法
  const method = (AppMethods as any)[methodName];
  if (typeof method !== "function") {
    throw new Error(`Method ${methodName} not found in Wails App bindings`);
  }

  try {
    // 根据参数数量调用方法
    let result: T;
    if (params && typeof params === "object") {
      // 如果是对象参数，按属性顺序传递
      const paramValues = Object.values(params);
      result = await method(...paramValues);
    } else if (params !== undefined) {
      // 单个参数
      result = await method(params);
    } else {
      // 无参数
      result = await method();
    }

    console.log(`[Wails] Method ${methodName} executed successfully`);
    return result;
  } catch (error) {
    console.error(`[Wails] Method ${methodName} failed:`, error);
    throw error;
  }
}

/**
 * 事件监听 - 使用 Wails 自动生成的事件系统
 */
export function wailsListen(
  eventName: string,
  callback: (data: any) => void,
): () => void {
  if (detectWailsEnvironment()) {
    console.log(`[Wails] Listening to event: ${eventName}`);
    // Wails 事件监听使用自动生成的绑定
    // 注意：Wails 事件系统需要 Go 后端支持
    // 目前先返回一个空函数，后续需要实现事件系统
    console.warn(
      `[Wails] Event listening for ${eventName} not yet implemented with Wails IPC`,
    );
    return () => {
      console.log(`[Wails] Unsubscribing from event: ${eventName}`);
    };
  }

  // Web 环境 - 使用自定义事件
  console.log(`[Web] Listening to custom event: ${eventName}`);
  const handler = (e: CustomEvent) => callback(e.detail);
  window.addEventListener(eventName, handler as EventListener);
  return () => {
    console.log(`[Web] Removing listener for custom event: ${eventName}`);
    window.removeEventListener(eventName, handler as EventListener);
  };
}

/**
 * 窗口操作 - 使用 Wails 自动生成的窗口方法
 */
export const wailsWindow = {
  minimize: async () => {
    if (detectWailsEnvironment()) {
      console.log("[Wails] Minimizing window");
      await WailsWindowMinimise();
    }
  },
  maximize: async () => {
    if (detectWailsEnvironment()) {
      console.log("[Wails] Maximizing window");
      await WailsWindowMaximise();
    }
  },
  unmaximize: async () => {
    if (detectWailsEnvironment()) {
      console.log("[Wails] Unmaximizing window");
      await WailsWindowUnmaximise();
    }
  },
  close: async () => {
    if (detectWailsEnvironment()) {
      console.log("[Wails] Closing window");
      await WailsQuit();
    }
  },
  getPosition: async () => {
    if (detectWailsEnvironment()) {
      console.log("[Wails] Getting window position");
      return await WailsWindowGetPosition();
    }
    return { x: 0, y: 0 };
  },
  setPosition: async (x: number, y: number) => {
    if (detectWailsEnvironment()) {
      console.log("[Wails] Setting window position:", { x, y });
      await WailsWindowSetPosition(x, y);
    }
  },
};

/**
 * 文件对话框 - 使用 Wails 自动生成的对话框方法
 */
export const wailsDialog = {
  open: async (options: any) => {
    if (detectWailsEnvironment()) {
      console.log("[Wails] Opening file dialog with options:", options);
      return await OpenFileDialog(options);
    }

    // Web 环境 - 使用原生 input
    console.log("[Web] Opening file dialog with native input");
    return new Promise((resolve) => {
      const input = document.createElement("input");
      input.type = "file";
      input.multiple = options.multiple || false;
      if (options.directory) {
        input.webkitdirectory = true;
      }
      input.onchange = () => {
        resolve(
          input.files
            ? Array.from(input.files).map((f) => (f as any).path || f.name)
            : null,
        );
      };
      input.click();
    });
  },
  save: async (options: any) => {
    if (detectWailsEnvironment()) {
      console.log("[Wails] Opening save dialog with options:", options);
      return await SaveFileDialog(options);
    }

    // Web 环境 - 使用原生下载
    console.log("[Web] Save dialog not available in web mode");
    return null;
  },
};

/**
 * 外部链接打开 - 使用 Wails 自动生成的浏览器方法
 */
export const wailsShell = {
  openExternal: async (url: string) => {
    if (detectWailsEnvironment()) {
      console.log("[Wails] Opening external URL:", url);
      await OpenExternal(url);
    } else {
      // Web 环境 - 使用 window.open
      console.log("[Web] Opening external URL with window.open:", url);
      window.open(url, "_blank");
    }
  },
};

/**
 * 获取环境信息用于调试
 */
export function getWailsEnvironmentInfo() {
  return {
    isWails: detectWailsEnvironment(),
    userAgent: navigator.userAgent,
    location: window.location.href,
    hasRuntime: !!window.runtime,
    hasGoApp: !!(window.go && window.go.main && window.go.main.App),
  };
}

/**
 * 初始化 Wails 模式
 * 检查 Wails 环境是否正常
 */
export function initializeWailsMode() {
  const isWails = detectWailsEnvironment();

  if (!isWails) {
    console.error(
      "Wails environment not detected. This application requires Wails to run.",
    );
    throw new Error(
      "Wails environment not detected. Please run this application in a Wails environment.",
    );
  }

  if (!window.go || !window.go.main || !window.go.main.App) {
    console.error(
      "Wails app not found. Please ensure Wails is properly initialized.",
    );
    throw new Error(
      "Wails app not found. Application cannot function without Wails backend.",
    );
  }

  console.log("Wails environment initialized successfully");
}
