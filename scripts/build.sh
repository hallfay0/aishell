#!/bin/bash

# AI Shell - 构建脚本
# 用于编译AI Shell项目

set -e  # 出错时退出

# 脚本配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BUILD_DIR="$PROJECT_ROOT/build"
BINARY_NAME="aishell"

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
AI Shell 构建脚本

用法: $0 [选项]

选项:
    -h, --help          显示此帮助信息
    -c, --clean         清理构建输出目录
    -r, --race          启用竞态检测
    -t, --tags TAGS     构建标签
    -o, --output FILE   输出文件名
    -v, --verbose       详细输出
    --dev              开发模式构建（包含调试信息）
    --release          发布模式构建（优化体积）

示例:
    $0                  # 标准构建
    $0 --clean         # 清理后重新构建
    $0 --dev           # 开发模式构建
    $0 --release       # 发布模式构建
    $0 -o myai         # 指定输出文件名

EOF
}

# 清理构建目录
clean_build() {
    print_info "清理构建目录..."
    if [ -d "$BUILD_DIR" ]; then
        rm -rf "$BUILD_DIR"
        print_success "构建目录已清理"
    else
        print_info "构建目录不存在，跳过清理"
    fi
}

# 检查Go环境
check_go_env() {
    print_info "检查Go环境..."
    
    if ! command -v go &> /dev/null; then
        print_error "未找到Go编译器，请先安装Go语言环境"
        exit 1
    fi
    
    local go_version=$(go version | awk '{print $3}' | sed 's/go//')
    print_info "Go版本: $go_version"
    
    # 检查Go版本是否满足要求（1.21+）
    local required_version="1.21"
    if [ "$(printf '%s\n' "$required_version" "$go_version" | sort -V | head -n1)" != "$required_version" ]; then
        print_warning "建议使用Go $required_version或更高版本"
    fi
}

# 检查依赖
check_dependencies() {
    print_info "检查项目依赖..."
    
    cd "$PROJECT_ROOT"
    
    if [ ! -f "go.mod" ]; then
        print_error "未找到go.mod文件，请确保在正确的项目目录中"
        exit 1
    fi
    
    # 下载依赖
    print_info "下载Go模块依赖..."
    go mod download
    
    # 整理依赖
    print_info "整理依赖关系..."
    go mod tidy
    
    print_success "依赖检查完成"
}

# 运行测试
run_tests() {
    print_info "运行单元测试..."
    
    cd "$PROJECT_ROOT"
    
    if ! go test ./... -v; then
        print_error "单元测试失败"
        exit 1
    fi
    
    print_success "所有测试通过"
}

# 构建应用
build_app() {
    local build_flags="$1"
    local output_file="$2"
    local build_tags="$3"
    
    print_info "开始构建应用..."
    
    cd "$PROJECT_ROOT"
    
    # 创建构建目录
    mkdir -p "$BUILD_DIR"
    
    # 获取构建信息
    local version=$(git describe --tags --always --dirty 2>/dev/null || echo "unknown")
    local commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    local build_time=$(date -u '+%Y-%m-%d_%H:%M:%S')
    local go_version=$(go version | awk '{print $3}')
    
    print_info "版本信息: $version"
    print_info "提交哈希: $commit"
    print_info "构建时间: $build_time"
    
    # 设置ldflags
    local ldflags="-X main.Version=$version"
    ldflags="$ldflags -X main.Commit=$commit"
    ldflags="$ldflags -X main.BuildTime=$build_time"
    ldflags="$ldflags -X main.GoVersion=$go_version"
    
    # 添加构建标志
    if [[ "$build_flags" == *"-s -w"* ]]; then
        ldflags="$ldflags -s -w"
    fi
    
    # 构建命令
    local build_cmd="go build"
    
    if [ -n "$build_tags" ]; then
        build_cmd="$build_cmd -tags=$build_tags"
    fi
    
    if [ "$RACE_DETECT" = "true" ]; then
        build_cmd="$build_cmd -race"
    fi
    
    if [ "$VERBOSE" = "true" ]; then
        build_cmd="$build_cmd -v"
    fi
    
    build_cmd="$build_cmd -ldflags=\"$ldflags\""
    build_cmd="$build_cmd -o $BUILD_DIR/$output_file"
    build_cmd="$build_cmd ./cmd/aishell"
    
    print_info "执行构建命令: $build_cmd"
    
    # 执行构建
    if eval $build_cmd; then
        print_success "构建完成: $BUILD_DIR/$output_file"
        
        # 显示文件信息
        local file_info=$(ls -lh "$BUILD_DIR/$output_file")
        print_info "文件信息: $file_info"
        
        # 显示可执行文件大小
        local file_size=$(du -h "$BUILD_DIR/$output_file" | cut -f1)
        print_info "文件大小: $file_size"
        
        return 0
    else
        print_error "构建失败"
        return 1
    fi
}

# 主函数
main() {
    local clean=false
    local run_test=true
    local build_mode="standard"
    local output_file="$BINARY_NAME"
    local build_tags=""
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -c|--clean)
                clean=true
                shift
                ;;
            -r|--race)
                RACE_DETECT=true
                shift
                ;;
            -t|--tags)
                build_tags="$2"
                shift 2
                ;;
            -o|--output)
                output_file="$2"
                shift 2
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            --dev)
                build_mode="dev"
                shift
                ;;
            --release)
                build_mode="release"
                shift
                ;;
            --no-test)
                run_test=false
                shift
                ;;
            *)
                print_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    print_info "🚀 开始构建 AI Shell"
    print_info "构建模式: $build_mode"
    
    # 清理（如果需要）
    if [ "$clean" = true ]; then
        clean_build
    fi
    
    # 检查环境
    check_go_env
    check_dependencies
    
    # 运行测试（如果需要）
    if [ "$run_test" = true ]; then
        run_tests
    fi
    
    # 根据构建模式设置标志
    local build_flags=""
    case $build_mode in
        "dev")
            build_flags=""  # 保留调试信息
            print_info "开发模式：保留调试信息"
            ;;
        "release")
            build_flags="-s -w"  # 去除调试信息，减小体积
            print_info "发布模式：优化文件体积"
            ;;
        "standard")
            build_flags=""
            print_info "标准模式：默认构建选项"
            ;;
    esac
    
    # 构建应用
    if build_app "$build_flags" "$output_file" "$build_tags"; then
        print_success "🎉 构建成功！"
        print_info "可执行文件路径: $BUILD_DIR/$output_file"
        print_info "运行命令: $BUILD_DIR/$output_file"
        
        # 创建符号链接到项目根目录（方便使用）
        ln -sf "build/$output_file" "$PROJECT_ROOT/$output_file"
        print_info "已创建符号链接: $PROJECT_ROOT/$output_file"
        
        exit 0
    else
        print_error "💥 构建失败！"
        exit 1
    fi
}

# 执行主函数
main "$@"
