#!/bin/bash

# BokeUI 启动脚本

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

PROJECT_DIR=$(cd "$(dirname "$0")" && pwd)

echo -e "${GREEN}启动 BokeUI 服务...${NC}"

# 检查是否已运行
if systemctl is-active --quiet bokeui; then
    echo -e "${RED}服务已在运行中（systemd）${NC}"
    systemctl status bokeui --no-pager
    exit 0
fi

if pm2 list | grep -q "bokeui.*online"; then
    echo -e "${RED}服务已在运行中（PM2）${NC}"
    pm2 list
    exit 0
fi

PID=$(ps aux | grep './server/bokeui' | grep -v grep | awk '{print $2}')
if [ ! -z "$PID" ]; then
    echo -e "${RED}服务已在运行中（PID: $PID）${NC}"
    exit 0
fi

# 选择启动方式
echo "选择启动方式："
echo "1) systemd（推荐）"
echo "2) PM2"
echo "3) 直接运行（后台）"
read -p "请选择 (1/2/3): " -n 1 -r
echo

case $REPLY in
    1)
        if [ -f /etc/systemd/system/bokeui.service ]; then
            systemctl start bokeui
            echo -e "${GREEN}服务已启动（systemd）${NC}"
            systemctl status bokeui --no-pager
        else
            echo -e "${RED}systemd 服务未配置，请先运行 deploy.sh${NC}"
            exit 1
        fi
        ;;
    2)
        if command -v pm2 &> /dev/null; then
            cd $PROJECT_DIR/server
            pm2 start ../ecosystem.config.js
            pm2 save
            echo -e "${GREEN}服务已启动（PM2）${NC}"
            pm2 list
        else
            echo -e "${RED}PM2 未安装，请先安装: npm install -g pm2${NC}"
            exit 1
        fi
        ;;
    3)
        cd $PROJECT_DIR/server
        mkdir -p ../logs
        nohup ./bokeui > ../logs/bokeui.log 2>&1 &
        sleep 2
        echo -e "${GREEN}服务已后台启动${NC}"
        echo -e "日志文件: ${PROJECT_DIR}/logs/bokeui.log"
        tail -n 20 ../logs/bokeui.log
        ;;
    *)
        echo -e "${RED}无效选择${NC}"
        exit 1
        ;;
esac
