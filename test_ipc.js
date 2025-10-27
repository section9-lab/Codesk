// 简单的 IPC 测试脚本
// 在浏览器控制台中运行这个脚本来测试 IPC 通信

async function testListProjects() {
  try {
    console.log('🧪 Testing IPC communication...');
    console.log('📱 Environment:', navigator.userAgent);
    console.log('🔗 URL:', window.location.href);

    // 检查 Wails 环境检测
    console.log('🔍 Checking Wails environment...');

    // 尝试调用 ListProjects
    console.log('📡 Calling ListProjects...');

    // 导入 wailsCall 方法
    const { wailsCall } = await import('./src/lib/wailsAdapter.js');

    const projects = await wailsCall('list_projects');
    console.log('✅ Success! Projects:', projects);

    // 测试其他方法
    console.log('📡 Testing GetHomeDirectory...');
    const homeDir = await wailsCall('get_home_directory');
    console.log('✅ Home directory:', homeDir);

    console.log('🎉 All IPC tests passed!');

  } catch (error) {
    console.error('❌ IPC test failed:', error);
    console.error('Stack trace:', error.stack);
  }
}

// 运行测试
testListProjects();

console.log('🚀 IPC test script loaded. Run testListProjects() to start testing.');