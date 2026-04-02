#!/bin/bash

# BokeUI 恢复脚本

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

PROJECT_DIR=$(cd "$(dirname "$0")" && pwd)
BACKUP_DIR="$PROJECT_DIR/backups"

# 检查备份目录
if [ ! -d "$BACKUP_DIR" ]; then
    echo -e "${RED}备份目录不存在${NC}"
    exit 1
fi

# 列出可用备份
echo -e "${GREEN}可用的备份文件：${NC}"
echo
ls -lh $BACKUP_DIR/backup_*.tar.gz | nl
echo

# 选择备份文件
read -p "请选择要恢复的备份编号: " BACKUP_NUM
BACKUP_FILE=$(ls $BACKUP_DIR/backup_*.tar.gz | sed -n "${BACKUP_NUM}p")

if [ ! -f "$BACKUP_FILE" ]; then
    echo -e "${RED}无效的选择${NC}"
    exit 1
fi

echo -e "${YELLOW}即将恢复备份: $(basename $BACKUP_FILE)${NC}"
read -p "确认恢复？当前数据将被覆盖！(y/n): " -n 1 -r
echo

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "取消恢复"
    exit 0
fi

# 停止服务
echo -e "${GREEN}停止服务...${NC}"
if systemctl is-active --quiet bokeui; then
    systemctl stop bokeui
elif pm2 list | grep -q "bokeui"; then
    pm2 stop bokeui
fi

# 创建临时目录
TEMP_DIR=$(mktemp -d)

# 解压备份
echo -e "${GREEN}解压备份文件...${NC}"
tar -xzf $BACKUP_FILE -C $TEMP_DIR

# 恢复数据库
if [ -f "$TEMP_DIR/data.db.$(basename $BACKUP_FILE | sed 's/backup_\(.*\)\.tar\.gz/data.db.\1/')" ]; then
    DB_FILE=$(ls $TEMP_DIR/data.db.*)
    echo -e "${GREEN}恢复数据库...${NC}"
    cp $DB_FILE $PROJECT_DIR/server/data.db
fi

# 恢复 IP 数据库
if [ -f "$TEMP_DIR/ip2region.xdb.$(basename $BACKUP_FILE | sed 's/backup_\(.*\)\.tar\.gz/ip2region.xdb.\1/')" ]; then
    IP_FILE=$(ls $TEMP_DIR/ip2region.xdb.*)
    echo -e "${GREEN}恢复 IP 数据库...${NC}"
    mkdir -p $PROJECT_DIR/server/data
    cp $IP_FILE $PROJECT_DIR/server/data/ip2region.xdb
fi

# 恢复上传文件
if [ -f "$TEMP_DIR/uploads.$(basename $BACKUP_FILE | sed 's/backup_\(.*\)\.tar\.gz/uploads.\1.tar.gz/')" ]; then
    UPLOAD_FILE=$(ls $TEMP_DIR/uploads.*.tar.gz)
    echo -e "${GREEN}恢复上传文件...${NC}"
    tar -xzf $UPLOAD_FILE -C $PROJECT_DIR/server
fi

# 清理临时文件
rm -rf $TEMP_DIR

echo -e "${GREEN}恢复完成${NC}"

# 重启服务
read -p "是否立即重启服务？(y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if systemctl is-enabled --quiet bokeui; then
        systemctl start bokeui
        echo -e "${GREEN}服务已重启${NC}"
    elif pm2 list | grep -q "bokeui"; then
        pm2 restart bokeui
        echo -e "${GREEN}服务已重启${NC}"
    fi
fi
