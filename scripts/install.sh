#!/bin/bash

# AI Shell - 安装脚本
# 用于安装AI Shell到系统

set -e  # 出错时退出

# 脚本配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BINARY_NAME="aishell"
DEFAULT_INSTALL_DIR="/usr/local/bin"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 显示帮助信息
show_help() {
    cat << EOF
AI Shell 安装脚本

用法: $0 [选项]

选项:
    -h, --help              显示此帮助信息
    -d, --dir DIR          指定安装目录 (默认: $DEFAULT_INSTALL_DIR)
    --user                 安装到用户目录 (~/.local/bin)
    --build                安装前重新构建
    --clean-install        清理安装（先卸载再安装）
    --dry-run             模拟安装过程（不实际安装）
    -f, --force            强制安装（覆盖现有文件）

示例:
    $0                     # 标准安装到 $DEFAULT_INSTALL_DIR
    $0 --user              # 安装到用户目录
    $0 -d /opt/bin         # 安装到指定目录
    $0 --build             # 重新构建后安装
    $0 --clean-install     # 清理安装

注意:
    - 安装到系统目录可能需要sudo权限
    - 用户目录安装不需要sudo权限
    - 安装后请确保安装目录在PATH环境变量中

EOF
}

# 检查权限
check_permissions() {
    local install_dir="$1"
    
    if [ ! -d "$install_dir" ]; then
        print_info "创建安装目录: $install_dir"
        if ! mkdir -p "$install_dir" 2>/dev/null; then
            print_error "无法创建安装目录，请检查权限或使用sudo"
            return 1
        fi
    fi
    
    if [ ! -w "$install_dir" ]; then
        print_warning "没有写入权限到 $install_dir"
        if [ "$EUID" -ne 0 ] && [[ "$install_dir" == /usr/* || "$install_dir" == /opt/* ]]; then
            print_info "提示：系统目录安装可能需要sudo权限"
            print_info "请使用: sudo $0 $*"
            return 1
        fi
        return 1
    fi
    
    return 0
}

# 构建应用
build_app() {
    print_info "开始构建应用..."
    
    if [ -f "$SCRIPT_DIR/build.sh" ]; then
        if bash "$SCRIPT_DIR/build.sh" --release; then
            print_success "构建完成"
            return 0
        else
            print_error "构建失败"
            return 1
        fi
    else
        print_error "未找到构建脚本: $SCRIPT_DIR/build.sh"
        return 1
    fi
}

# 检查二进制文件
check_binary() {
    local binary_path="$PROJECT_ROOT/build/$BINARY_NAME"
    
    if [ ! -f "$binary_path" ]; then
        print_warning "未找到编译后的二进制文件: $binary_path"
        print_info "将自动构建应用..."
        if ! build_app; then
            return 1
        fi
    fi
    
    if [ ! -x "$binary_path" ]; then
        print_error "二进制文件没有执行权限: $binary_path"
        return 1
    fi
    
    print_success "二进制文件检查通过: $binary_path"
    return 0
}

# 卸载现有版本
uninstall_existing() {
    local install_dir="$1"
    local install_path="$install_dir/$BINARY_NAME"
    
    if [ -f "$install_path" ]; then
        print_info "发现现有安装版本: $install_path"
        
        # 显示当前版本信息
        if "$install_path" --version &>/dev/null; then
            local current_version=$("$install_path" --version 2>/dev/null || echo "unknown")
            print_info "当前版本: $current_version"
        fi
        
        print_info "正在卸载现有版本..."
        if rm "$install_path"; then
            print_success "现有版本已卸载"
        else
            print_error "卸载失败"
            return 1
        fi
    else
        print_info "未发现现有安装版本"
    fi
    
    return 0
}

# 安装二进制文件
install_binary() {
    local install_dir="$1"
    local force="$2"
    local binary_path="$PROJECT_ROOT/build/$BINARY_NAME"
    local install_path="$install_dir/$BINARY_NAME"
    
    # 检查目标文件是否存在
    if [ -f "$install_path" ] && [ "$force" != "true" ]; then
        print_warning "目标文件已存在: $install_path"
        print_info "使用 --force 参数强制覆盖，或使用 --clean-install 清理安装"
        return 1
    fi
    
    print_info "安装二进制文件..."
    print_info "源文件: $binary_path"
    print_info "目标路径: $install_path"
    
    # 复制文件
    if cp "$binary_path" "$install_path"; then
        print_success "文件复制完成"
    else
        print_error "文件复制失败"
        return 1
    fi
    
    # 设置执行权限
    if chmod +x "$install_path"; then
        print_success "执行权限设置完成"
    else
        print_error "设置执行权限失败"
        return 1
    fi
    
    return 0
}

# 创建配置目录
create_config_dir() {
    local config_dir="$HOME/.config/aishell"
    
    if [ ! -d "$config_dir" ]; then
        print_info "创建配置目录: $config_dir"
        if mkdir -p "$config_dir"; then
            print_success "配置目录创建完成"
            
            # 创建示例配置文件
            if [ -f "$PROJECT_ROOT/env.example" ]; then
                cp "$PROJECT_ROOT/env.example" "$config_dir/env.example"
                print_info "已复制配置示例到: $config_dir/env.example"
            fi
        else
            print_warning "配置目录创建失败"
        fi
    else
        print_info "配置目录已存在: $config_dir"
    fi
}

# 验证安装
verify_installation() {
    local install_dir="$1"
    local install_path="$install_dir/$BINARY_NAME"
    
    print_info "验证安装..."
    
    # 检查文件存在
    if [ ! -f "$install_path" ]; then
        print_error "安装验证失败：文件不存在"
        return 1
    fi
    
    # 检查执行权限
    if [ ! -x "$install_path" ]; then
        print_error "安装验证失败：没有执行权限"
        return 1
    fi
    
    # 检查是否在PATH中
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        print_success "命令在PATH中可用"
        local which_path=$(which "$BINARY_NAME")
        if [ "$which_path" = "$install_path" ]; then
            print_success "PATH指向正确的安装位置"
        else
            print_warning "PATH指向不同位置: $which_path"
            print_warning "可能存在多个安装版本"
        fi
    else
        print_warning "命令不在PATH中"
        print_info "请确保 $install_dir 在您的PATH环境变量中"
        print_info "您可以运行: export PATH=\"$install_dir:\$PATH\""
        print_info "或将其添加到您的 ~/.bashrc 或 ~/.zshrc 文件中"
    fi
    
    # 尝试运行版本命令
    print_info "测试运行..."
    if "$install_path" --version &>/dev/null; then
        local version=$("$install_path" --version 2>/dev/null || echo "unknown")
        print_success "安装验证成功 - 版本: $version"
    else
        print_warning "无法获取版本信息，但文件已安装"
    fi
    
    return 0
}

# 显示安装后信息
show_post_install_info() {
    local install_dir="$1"
    
    cat << EOF

🎉 AI Shell 安装完成！

安装位置: $install_dir/$BINARY_NAME
配置目录: $HOME/.config/aishell/

快速开始:
  1. 设置OpenAI API Key:
     export OPENAI_API_KEY="your_api_key_here"
  
  2. 运行AI Shell:
     $BINARY_NAME
  
  3. 启用调试模式:
     AISHELL_DEBUG=true $BINARY_NAME

更多信息:
  - 查看帮助: $BINARY_NAME --help
  - 项目文档: $PROJECT_ROOT/README.md
  - 配置示例: $HOME/.config/aishell/env.example

EOF
}

# 主函数
main() {
    local install_dir="$DEFAULT_INSTALL_DIR"
    local build_first=false
    local clean_install=false
    local dry_run=false
    local force=false
    local user_install=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -d|--dir)
                install_dir="$2"
                shift 2
                ;;
            --user)
                user_install=true
                install_dir="$HOME/.local/bin"
                shift
                ;;
            --build)
                build_first=true
                shift
                ;;
            --clean-install)
                clean_install=true
                shift
                ;;
            --dry-run)
                dry_run=true
                shift
                ;;
            -f|--force)
                force=true
                shift
                ;;
            *)
                print_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    print_info "🚀 开始安装 AI Shell"
    print_info "安装目录: $install_dir"
    
    if [ "$dry_run" = true ]; then
        print_warning "模拟模式 - 不会实际安装文件"
    fi
    
    # 检查权限
    if [ "$dry_run" != true ]; then
        if ! check_permissions "$install_dir"; then
            exit 1
        fi
    fi
    
    # 重新构建（如果需要）
    if [ "$build_first" = true ]; then
        if [ "$dry_run" != true ]; then
            if ! build_app; then
                exit 1
            fi
        else
            print_info "[模拟] 会构建应用"
        fi
    fi
    
    # 检查二进制文件
    if [ "$dry_run" != true ]; then
        if ! check_binary; then
            exit 1
        fi
    else
        print_info "[模拟] 会检查二进制文件"
    fi
    
    # 清理安装（如果需要）
    if [ "$clean_install" = true ]; then
        if [ "$dry_run" != true ]; then
            if ! uninstall_existing "$install_dir"; then
                exit 1
            fi
        else
            print_info "[模拟] 会卸载现有版本"
        fi
    fi
    
    # 安装二进制文件
    if [ "$dry_run" != true ]; then
        if ! install_binary "$install_dir" "$force"; then
            exit 1
        fi
    else
        print_info "[模拟] 会安装二进制文件到 $install_dir"
    fi
    
    # 创建配置目录
    if [ "$dry_run" != true ]; then
        create_config_dir
    else
        print_info "[模拟] 会创建配置目录"
    fi
    
    # 验证安装
    if [ "$dry_run" != true ]; then
        if verify_installation "$install_dir"; then
            print_success "🎉 安装成功！"
            show_post_install_info "$install_dir"
        else
            print_error "💥 安装验证失败！"
            exit 1
        fi
    else
        print_info "[模拟] 会验证安装"
        print_success "🎉 模拟安装完成！"
    fi
}

# 执行主函数
main "$@"
