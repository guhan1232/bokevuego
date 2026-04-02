#!/bin/bash

# BokeUI 博客系统一键部署脚本
# 适用于 Ubuntu/Debian/CentOS 等 Linux 系统
# 作者：BokeUI Team
# 版本：1.0.0

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否为 root 用户
check_root() {
    if [ "$EUID" -ne 0 ]; then
        log_warning "建议使用 sudo 或 root 用户运行此脚本"
        read -p "是否继续？(y/n): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# 检测系统类型
detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        VER=$VERSION_ID
    elif [ -f /etc/redhat-release ]; then
        OS=rhel
    else
        OS=$(uname -s)
    fi
    log_info "检测到系统: $OS $VER"
}

# 安装依赖
install_dependencies() {
    log_info "开始安装系统依赖..."
    
    if [ "$OS" = "ubuntu" ] || [ "$OS" = "debian" ]; then
        apt-get update
        apt-get install -y wget curl git build-essential
    elif [ "$OS" = "centos" ] || [ "$OS" = "rhel" ]; then
        yum install -y wget curl git gcc gcc-c++ make
    else
        log_warning "未知系统类型，跳过系统依赖安装"
    fi
    
    log_success "系统依赖安装完成"
}

# 安装 Go
install_go() {
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}')
        log_success "Go 已安装: $GO_VERSION"
        return
    fi
    
    log_info "开始安装 Go..."
    
    GO_VERSION="1.21.5"
    GO_TAR="go${GO_VERSION}.linux-amd64.tar.gz"
    GO_URL="https://golang.google.cn/dl/${GO_TAR}"
    
    wget -q $GO_URL -O /tmp/$GO_TAR || {
        log_error "Go 下载失败，尝试使用官方源..."
        GO_URL="https://dl.google.com/go/${GO_TAR}"
        wget -q $GO_URL -O /tmp/$GO_TAR
    }
    
    tar -C /usr/local -xzf /tmp/$GO_TAR
    rm /tmp/$GO_TAR
    
    # 配置环境变量
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    
    log_success "Go 安装完成: $(go version)"
}

# 安装 Node.js
install_nodejs() {
    if command -v node &> /dev/null; then
        NODE_VERSION=$(node -v)
        log_success "Node.js 已安装: $NODE_VERSION"
        return
    fi
    
    log_info "开始安装 Node.js..."
    
    # 安装 Node.js 18.x LTS
    curl -fsSL https://rpm.nodesource.com/setup_18.x | bash - || {
        curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
    }
    
    if [ "$OS" = "ubuntu" ] || [ "$OS" = "debian" ]; then
        apt-get install -y nodejs
    elif [ "$OS" = "centos" ] || [ "$OS" = "rhel" ]; then
        yum install -y nodejs
    fi
    
    log_success "Node.js 安装完成: $(node -v)"
    log_success "npm 版本: $(npm -v)"
}

# 安装 PM2（进程管理器）
install_pm2() {
    if command -v pm2 &> /dev/null; then
        log_success "PM2 已安装"
        return
    fi
    
    log_info "安装 PM2 进程管理器..."
    npm install -g pm2
    log_success "PM2 安装完成"
}

# 构建项目
build_project() {
    log_info "开始构建项目..."
    
    PROJECT_DIR=$(pwd)
    
    # 构建后端
    log_info "编译后端程序..."
    log_info "设置 Go 代理加速..."
    export GOPROXY=https://goproxy.cn,direct
    export GO111MODULE=on
    
    cd server
    log_info "下载 Go 依赖（首次可能需要几分钟）..."
    go mod download 2>&1 | while read line; do
        log_info "  $line"
    done || true
    log_info "开始编译..."
    go build -v -o bokeui ./cmd
    chmod +x bokeui
    cd ..
    log_success "后端编译完成"
    
    # 构建前端
    log_info "构建管理后台..."
    cd admin
    npm install --registry=https://registry.npmmirror.com
    npm run build
    cd ..
    log_success "管理后台构建完成"
    
    log_info "构建前台页面..."
    cd web
    npm install --registry=https://registry.npmmirror.com
    npm run build
    cd ..
    log_success "前台页面构建完成"
}

# 创建必要目录
create_directories() {
    log_info "创建必要目录..."
    
    mkdir -p server/data
    mkdir -p logs
    
    log_success "目录创建完成"
}

# 创建 systemd 服务
create_systemd_service() {
    log_info "创建 systemd 服务..."
    
    PROJECT_DIR=$(pwd)
    SERVICE_FILE="/etc/systemd/system/bokeui.service"
    
    cat > $SERVICE_FILE <<EOF
[Unit]
Description=BokeUI Blog System
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=${PROJECT_DIR}/server
ExecStart=${PROJECT_DIR}/server/bokeui
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF
    
    systemctl daemon-reload
    systemctl enable bokeui
    
    log_success "systemd 服务创建完成"
}

# 创建 PM2 配置
create_pm2_config() {
    log_info "创建 PM2 配置..."
    
    PROJECT_DIR=$(pwd)
    
    cat > ecosystem.config.js <<EOF
module.exports = {
  apps: [{
    name: 'bokeui',
    script: './server/bokeui',
    cwd: '${PROJECT_DIR}',
    instances: 1,
    autorestart: true,
    watch: false,
    max_memory_restart: '1G',
    env: {
      NODE_ENV: 'production',
      PORT: 9088
    }
  }]
}
EOF
    
    log_success "PM2 配置创建完成"
}

# 配置防火墙
configure_firewall() {
    log_info "配置防火墙..."
    
    if command -v ufw &> /dev/null; then
        ufw allow 9088/tcp
        ufw reload
        log_success "UFW 防火墙已开放 9088 端口"
    elif command -v firewall-cmd &> /dev/null; then
        firewall-cmd --permanent --add-port=9088/tcp
        firewall-cmd --reload
        log_success "Firewalld 防火墙已开放 9088 端口"
    else
        log_warning "未检测到防火墙，请手动开放 9088 端口"
    fi
}

# 启动服务
start_service() {
    log_info "选择启动方式："
    echo "1) systemd（推荐）"
    echo "2) PM2"
    echo "3) 直接运行"
    read -p "请选择 (1/2/3): " -n 1 -r
    echo
    
    case $REPLY in
        1)
            systemctl start bokeui
            systemctl status bokeui --no-pager
            ;;
        2)
            pm2 start ecosystem.config.js
            pm2 save
            pm2 startup
            ;;
        3)
            cd server
            nohup ./bokeui > ../logs/bokeui.log 2>&1 &
            cd ..
            log_success "服务已后台启动，日志: logs/bokeui.log"
            ;;
        *)
            log_error "无效选择"
            exit 1
            ;;
    esac
}

# 显示部署信息
show_info() {
    SERVER_IP=$(curl -s ifconfig.me || hostname -I | awk '{print $1}')
    
    echo
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}    BokeUI 博客系统部署完成！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo
    echo -e "访问地址："
    echo -e "  前台：${BLUE}http://${SERVER_IP}:9088${NC}"
    echo -e "  后台：${BLUE}http://${SERVER_IP}:9088/admin${NC}"
    echo
    echo -e "默认账号："
    echo -e "  用户名：${YELLOW}admin${NC}"
    echo -e "  密码：${YELLOW}admin123${NC}"
    echo
    echo -e "管理命令："
    echo -e "  启动服务：systemctl start bokeui"
    echo -e "  停止服务：systemctl stop bokeui"
    echo -e "  重启服务：systemctl restart bokeui"
    echo -e "  查看状态：systemctl status bokeui"
    echo -e "  查看日志：journalctl -u bokeui -f"
    echo
    echo -e "${YELLOW}提示：请及时修改默认密码！${NC}"
    echo
}

# 主函数
main() {
    echo -e "${GREEN}"
    echo "======================================"
    echo "    BokeUI 博客系统一键部署脚本"
    echo "======================================"
    echo -e "${NC}"
    
    # 检查运行环境
    check_root
    detect_os
    
    # 安装依赖
    install_dependencies
    install_go
    install_nodejs
    install_pm2
    
    # 构建项目
    build_project
    
    # 创建必要目录
    create_directories
    
    # 配置服务
    create_systemd_service
    create_pm2_config
    
    # 配置防火墙
    configure_firewall
    
    # 启动服务
    start_service
    
    # 显示信息
    show_info
}

# 执行主函数
main
