#!/bin/bash

# 环境管理工具集成测试脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 测试配置
TEST_BASE_DIR="/tmp/kite-env-test-$$"
KITE_CMD="go run ../cmd/kite/main.go"

# 清理函数
cleanup() {
    log_info "清理测试环境..."
    rm -rf "$TEST_BASE_DIR"
}

# 设置陷阱以确保清理
trap cleanup EXIT

# 初始化测试环境
init_test_env() {
    log_info "初始化测试环境: $TEST_BASE_DIR"
    
    mkdir -p "$TEST_BASE_DIR"/{config/module,data/shell_env}
    
    # 创建测试配置文件
    cat > "$TEST_BASE_DIR/config/module/shell_env.yml" << EOF
add_paths: []
add_envs: {}
remove_envs: []
sdk_dir: $TEST_BASE_DIR/sdk

sdks:
  - name: go
    install_url: https://golang.org/dl/go{version}.{os}-{arch}.tar.gz
    install_dir: $TEST_BASE_DIR/sdk/go{version}
    
  - name: node
    install_url: https://nodejs.org/dist/v{version}/node-v{version}-{os}-{arch}.tar.xz
    install_dir: $TEST_BASE_DIR/sdk/node{version}
    
  - name: test-sdk
    install_dir: $TEST_BASE_DIR/sdk/test-sdk{version}
EOF
    
    export KITE_BASE_DIR="$TEST_BASE_DIR"
}

# 测试版本解析
test_version_parsing() {
    log_info "测试版本解析功能..."
    
    # 这里应该运行单元测试
    go test ./pkg/envmgr -run TestParseVersionSpec -v
    
    if [ $? -eq 0 ]; then
        log_info "版本解析测试通过"
    else
        log_error "版本解析测试失败"
        exit 1
    fi
}

# 测试配置管理
test_config_management() {
    log_info "测试配置管理功能..."
    
    # 测试配置查看
    $KITE_CMD dev env config
    
    if [ $? -eq 0 ]; then
        log_info "配置管理测试通过"
    else
        log_error "配置管理测试失败"
        exit 1
    fi
}

# 测试shell脚本生成
test_shell_script_generation() {
    log_info "测试shell脚本生成功能..."
    
    # 测试bash脚本生成
    local bash_script=$($KITE_CMD dev env shell bash)
    
    if [[ $bash_script == *"ktenv"* ]]; then
        log_info "Bash脚本生成测试通过"
    else
        log_error "Bash脚本生成测试失败"
        exit 1
    fi
    
    # 测试PowerShell脚本生成
    local pwsh_script=$($KITE_CMD dev env shell pwsh)
    
    if [[ $pwsh_script == *"function ktenv"* ]]; then
        log_info "PowerShell脚本生成测试通过"
    else
        log_error "PowerShell脚本生成测试失败"
        exit 1
    fi
}

# 测试SDK列表功能
test_sdk_listing() {
    log_info "测试SDK列表功能..."
    
    $KITE_CMD dev env list
    
    if [ $? -eq 0 ]; then
        log_info "SDK列表测试通过"
    else
        log_error "SDK列表测试失败"
        exit 1
    fi
}

# 模拟SDK安装测试
test_mock_sdk_installation() {
    log_info "测试模拟SDK安装..."
    
    # 创建模拟的SDK目录结构
    mkdir -p "$TEST_BASE_DIR/sdk/test-sdk1.0.0/bin"
    echo '#!/bin/bash\necho "test-sdk version 1.0.0"' > "$TEST_BASE_DIR/sdk/test-sdk1.0.0/bin/test-sdk"
    chmod +x "$TEST_BASE_DIR/sdk/test-sdk1.0.0/bin/test-sdk"
    
    log_info "模拟SDK安装完成"
}

# 测试状态管理
test_state_management() {
    log_info "测试状态管理功能..."
    
    # 运行状态管理单元测试
    go test ./pkg/envmgr -run TestDefaultStateManager -v
    
    if [ $? -eq 0 ]; then
        log_info "状态管理测试通过"
    else
        log_error "状态管理测试失败"
        exit 1
    fi
}

# 测试ktenv独立程序
test_ktenv_standalone() {
    log_info "测试ktenv独立程序..."
    
    # 构建ktenv程序
    go build -o "$TEST_BASE_DIR/ktenv" ./cmd/ktenv
    
    if [ ! -f "$TEST_BASE_DIR/ktenv" ]; then
        log_error "ktenv程序构建失败"
        exit 1
    fi
    
    # 测试帮助命令
    "$TEST_BASE_DIR/ktenv" help > /dev/null
    
    if [ $? -eq 0 ]; then
        log_info "ktenv独立程序测试通过"
    else
        log_error "ktenv独立程序测试失败"
        exit 1
    fi
}

# 主测试流程
main() {
    log_info "开始环境管理工具集成测试..."
    
    init_test_env
    test_version_parsing
    test_config_management
    test_shell_script_generation
    test_sdk_listing
    test_mock_sdk_installation
    test_state_management
    test_ktenv_standalone
    
    log_info "所有测试通过！✅"
}

# 运行主函数
main "$@"