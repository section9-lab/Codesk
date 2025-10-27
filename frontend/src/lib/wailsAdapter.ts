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
  Quit as WailsQuit,
  BrowserOpenURL as WailsBrowserOpenURL
} from '../../wailsjs/runtime';

// 导入 Wails 自动生成的 App 方法
import {
  OpenFileDialog,
  SaveFileDialog,
  OpenExternal
} from '../../wailsjs/go/main/App';

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
  if (typeof window === 'undefined') {
    isWailsEnvironment = false;
    return false;
  }

  // 检查 Wails 特定的指示器
  const isWails = !!(window.runtime || window.go);
  
  console.log('[detectWailsEnvironment] isWails:', isWails, 'userAgent:', navigator.userAgent);
  
  isWailsEnvironment = isWails;
  return isWails;
}

/**
 * REST API 调用的响应包装器
 */


/**
 * 替换 Tauri 的 invoke 调用
 */
export async function wailsCall<T>(command: string, params?: any): Promise<T> {
  const isWails = detectWailsEnvironment();

  if (!isWails) {
    // 如果不是 Wails 环境，抛出错误
    throw new Error(`Wails environment not detected. Cannot call method: ${command}`);
  }

  // Wails 环境 - 调用 Go 后端方法
  console.log(`[Wails] Calling: ${command}`, params);
  
  if (window.go && window.go.main && window.go.main.App) {
    const app = window.go.main.App;
    // 将 snake_case 转换为驼峰命名
    const methodName = command
      .split('_')
      .map((word, index) => 
        index === 0 ? word : word.charAt(0).toUpperCase() + word.slice(1)
      )
      .join('');
    
    console.log(`[Wails] Looking for method: ${methodName}`);
    
    if (typeof app[methodName] === 'function') {
      try {
        const result = await app[methodName](params);
        console.log(`[Wails] Method ${methodName} executed successfully`);
        return result;
      } catch (error) {
        console.error(`[Wails] Method ${methodName} failed:`, error);
        throw error;
      }
    } else {
      throw new Error(`Method ${methodName} not found in Wails app`);
    }
  } else {
    throw new Error('Wails app not found - please ensure you are running in a Wails environment');
  }
}

/**
 * 事件监听 - 使用 Wails 自动生成的事件系统
 */
export function wailsListen(eventName: string, callback: (data: any) => void): () => void {
  if (detectWailsEnvironment()) {
    console.log(`[Wails] Listening to event: ${eventName}`);
    // Wails 事件监听使用自动生成的绑定
    // 注意：Wails 事件系统需要 Go 后端支持
    // 目前先返回一个空函数，后续需要实现事件系统
    console.warn(`[Wails] Event listening for ${eventName} not yet implemented with Wails IPC`);
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
      console.log('[Wails] Minimizing window');
      await WailsWindowMinimise();
    }
  },
  maximize: async () => {
    if (detectWailsEnvironment()) {
      console.log('[Wails] Maximizing window');
      await WailsWindowMaximise();
    }
  },
  unmaximize: async () => {
    if (detectWailsEnvironment()) {
      console.log('[Wails] Unmaximizing window');
      await WailsWindowUnmaximise();
    }
  },
  close: async () => {
    if (detectWailsEnvironment()) {
      console.log('[Wails] Closing window');
      await WailsQuit();
    }
  }
};

/**
 * 文件对话框 - 使用 Wails 自动生成的对话框方法
 */
export const wailsDialog = {
  open: async (options: any) => {
    if (detectWailsEnvironment()) {
      console.log('[Wails] Opening file dialog with options:', options);
      return await OpenFileDialog(options);
    }
    
    // Web 环境 - 使用原生 input
    console.log('[Web] Opening file dialog with native input');
    return new Promise((resolve) => {
      const input = document.createElement('input');
      input.type = 'file';
      input.multiple = options.multiple || false;
      if (options.directory) {
        input.webkitdirectory = true;
      }
      input.onchange = () => {
        resolve(input.files ? Array.from(input.files).map(f => (f as any).path || f.name) : null);
      };
      input.click();
    });
  },
  save: async (options: any) => {
    if (detectWailsEnvironment()) {
      console.log('[Wails] Opening save dialog with options:', options);
      return await SaveFileDialog(options);
    }
    
    // Web 环境 - 使用原生下载
    console.log('[Web] Save dialog not available in web mode');
    return null;
  }
};

/**
 * 外部链接打开 - 使用 Wails 自动生成的浏览器方法
 */
export const wailsShell = {
  openExternal: async (url: string) => {
    if (detectWailsEnvironment()) {
      console.log('[Wails] Opening external URL:', url);
      await OpenExternal(url);
    } else {
      // Web 环境 - 使用 window.open
      console.log('[Web] Opening external URL with window.open:', url);
      window.open(url, '_blank');
    }
  }
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
    console.error('Wails environment not detected. This application requires Wails to run.');
    throw new Error('Wails environment not detected. Please run this application in a Wails environment.');
  }
  
  if (!window.go || !window.go.main || !window.go.main.App) {
    console.error('Wails app not found. Please ensure Wails is properly initialized.');
    throw new Error('Wails app not found. Application cannot function without Wails backend.');
  }
  
  console.log('Wails environment initialized successfully');
}