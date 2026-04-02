#!/bin/bash

# BokeUI 更新脚本

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# 拉取最新代码
log_info "拉取最新代码..."
git pull

# 构建后端
log_info "重新编译后端..."
cd server
go build -o bokeui ./cmd
chmod +x bokeui
cd ..

# 构建前端
log_info "重新构建管理后台..."
cd admin
npm run build
cd ..

log_info "重新构建前台页面..."
cd web
npm run build
cd ..

# 重启服务
log_info "重启服务..."
if systemctl is-active --quiet bokeui; then
    systemctl restart bokeui
    log_success "systemd 服务已重启"
elif pm2 list | grep -q "bokeui"; then
    pm2 restart bokeui
    log_success "PM2 服务已重启"
else
    log_info "未检测到运行中的服务，请手动启动"
fi

log_success "更新完成！"
