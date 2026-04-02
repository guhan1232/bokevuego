#!/bin/bash

# BokeUI 停止脚本

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 停止 systemd 服务
if systemctl is-active --quiet bokeui; then
    echo -e "${GREEN}停止 systemd 服务...${NC}"
    systemctl stop bokeui
    echo -e "${GREEN}已停止 systemd 服务${NC}"
fi

# 停止 PM2 服务
if pm2 list | grep -q "bokeui"; then
    echo -e "${GREEN}停止 PM2 服务...${NC}"
    pm2 stop bokeui
    echo -e "${GREEN}已停止 PM2 服务${NC}"
fi

# 停止直接运行的进程
PID=$(ps aux | grep './server/bokeui' | grep -v grep | awk '{print $2}')
if [ ! -z "$PID" ]; then
    echo -e "${GREEN}停止进程 (PID: $PID)...${NC}"
    kill $PID
    echo -e "${GREEN}已停止进程${NC}"
fi

echo -e "${GREEN}所有服务已停止${NC}"
