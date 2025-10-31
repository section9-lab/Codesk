#!/bin/bash

# Codesk 自动打包脚本
# 功能：1. 构建应用 2. 打包为 DMG

set -e  # 遇到错误立即退出

# 获取当前脚本所在目录（项目根目录）
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$SCRIPT_DIR"

# 配置变量
APP_NAME="Codesk"
VERSION="v0.0.1"
BUILD_DIR="$PROJECT_DIR/build/bin"
TIMESTAMP=$(date +%Y%m%d_%H%M)
DMG_NAME="${APP_NAME}_${VERSION}_${TIMESTAMP}.dmg"
DMG_PATH="$HOME/Downloads/$DMG_NAME"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log() {
    echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

# 检查是否在项目目录中
check_project_dir() {
    log "检查项目目录..."
    info "当前目录: $PROJECT_DIR"
    
    # 检查必要的项目文件
    if [ ! -f "$PROJECT_DIR/wails.json" ]; then
        error "未找到 wails.json 文件，请确保在正确的项目目录中运行脚本"
    fi
    
    if [ ! -f "$PROJECT_DIR/go.mod" ]; then
        error "未找到 go.mod 文件，请确保在正确的项目目录中运行脚本"
    fi
    
    log "✓ 项目目录验证通过"
}

# 检查依赖
check_dependencies() {
    log "检查依赖工具..."
    
    if ! command -v wails &> /dev/null; then
        error "wails 命令未找到，请先安装 Wails"
    fi
    
    if ! command -v create-dmg &> /dev/null; then
        error "create-dmg 未安装，请运行: brew install create-dmg"
    fi
    
    log "✓ 所有依赖工具就绪"
}

# 步骤1: 构建应用
build_application() {
    log "步骤1: 构建应用程序..."
    
    cd "$PROJECT_DIR" || error "无法进入项目目录: $PROJECT_DIR"
    
    # 清理并构建通用二进制
    log "执行构建命令: wails build -clean -platform darwin/universal"
    
    if wails build -clean -platform darwin/universal; then
        log "✓ 应用构建成功"
    else
        error "应用构建失败"
    fi
    
    # 验证构建结果
    if [ -f "$BUILD_DIR/${APP_NAME}.app/Contents/MacOS/${APP_NAME}" ]; then
        log "✓ 可执行文件验证通过"
        
        # 检查架构
        local arch_info=$(file "$BUILD_DIR/${APP_NAME}.app/Contents/MacOS/${APP_NAME}")
        if echo "$arch_info" | grep -q "universal binary"; then
            log "✓ 确认是通用二进制 (支持 Intel 和 Apple Silicon)"
        else
            warn "可能不是通用二进制"
        fi
    else
        error "构建产物验证失败"
    fi
}

# 步骤2: 打包 DMG
package_dmg() {
    log "步骤2: 打包为 DMG 文件..."
    info "输出路径: $DMG_PATH"
    
    cd "$PROJECT_DIR" || error "无法进入项目目录"
    
    # 检查源文件是否存在
    if [ ! -d "$BUILD_DIR" ]; then
        error "构建目录不存在: $BUILD_DIR"
    fi
    
    if [ ! -f "$BUILD_DIR/${APP_NAME}.app/Contents/MacOS/${APP_NAME}" ]; then
        error "应用文件未找到，请先执行构建步骤"
    fi
    
    # 删除已存在的 DMG 文件
    if [ -f "$DMG_PATH" ]; then
        warn "删除已存在的 DMG 文件: $(basename "$DMG_PATH")"
        rm -f "$DMG_PATH"
    fi
    
    # 执行 DMG 打包
    log "执行 DMG 打包命令..."
    
    create-dmg \
        --volname "$APP_NAME" \
        --window-pos 200 120 \
        --window-size 500 400 \
        --icon-size 100 \
        --icon "${APP_NAME}.app" 100 150 \
        --app-drop-link 300 155 \
        --hide-extension "${APP_NAME}.app" \
        "$DMG_PATH" \
        "$BUILD_DIR/"
    
    # 检查打包结果
    if [ $? -eq 0 ] && [ -f "$DMG_PATH" ]; then
        log "✓ DMG 打包成功"
        
        # 显示文件信息
        local dmg_size=$(du -h "$DMG_PATH" | cut -f1)
        log "文件大小: $dmg_size"
        log "完整路径: $DMG_PATH"
    else
        error "DMG 打包失败"
    fi
}

# 验证 DMG 文件
verify_dmg() {
    log "验证 DMG 文件..."
    
    # 基本验证
    if [ ! -f "$DMG_PATH" ]; then
        error "DMG 文件不存在: $DMG_PATH"
    fi
    
    # 测试挂载
    local mount_point=$(mktemp -d)
    log "测试挂载 DMG..."
    
    if hdiutil attach "$DMG_PATH" -mountpoint "$mount_point" -nobrowse -quiet; then
        log "✓ DMG 挂载测试通过"
        
        # 检查内容
        if [ -d "$mount_point/${APP_NAME}.app" ]; then
            log "✓ 应用包验证通过"
        else
            warn "应用包在 DMG 中未找到"
        fi
        
        # 卸载
        hdiutil detach "$mount_point" -quiet
        rmdir "$mount_point"
    else
        warn "DMG 挂载测试失败，但文件可能仍可用"
    fi
}

# 显示摘要信息
show_summary() {
    echo
    info "🎉 打包完成!"
    echo "=========================================="
    info "应用名称: $APP_NAME"
    info "版本号: $VERSION"
    info "时间戳: $TIMESTAMP"
    info "DMG 文件: $(basename "$DMG_PATH")"
    info "文件位置: $DMG_PATH"
    info "文件大小: $(du -h "$DMG_PATH" | cut -f1)"
    info "项目目录: $PROJECT_DIR"
    echo "=========================================="
}

# 交互选项
interactive_options() {
    echo
    read -p "是否打开 DMG 文件进行安装测试? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        open "$DMG_PATH"
        log "已打开 DMG 文件"
    fi
}

# 主函数
main() {
    echo
    log "🚀 开始打包 $APP_NAME"
    info "版本: $VERSION"
    info "时间: $(date)"
    echo
    
    # 执行步骤
    check_project_dir
    check_dependencies
    build_application
    package_dmg
    verify_dmg
    show_summary
    interactive_options
    
    echo
    log "✅ 所有步骤完成!"
}

# 显示帮助
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help    显示帮助信息"
    echo "  -v, --version 显示版本信息"
    echo ""
    echo "功能:"
    echo "  自动构建 Codesk 应用并打包为 DMG 文件"
    echo "  输出文件: ~/Downloads/Codesk_v0.0.1_时间戳.dmg"
    echo ""
    echo "说明:"
    echo "  脚本会自动检测当前目录，请确保在项目根目录下运行"
}

# 参数处理
case "$1" in
    -h|--help)
        show_help
        exit 0
        ;;
    -v|--version)
        echo "Codesk 打包脚本 v1.0"
        exit 0
        ;;
    *)
        main
        ;;
esac