#!/bin/bash

# BokeUI 备份脚本

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_DIR=$(cd "$(dirname "$0")" && pwd)
BACKUP_DIR="$PROJECT_DIR/backups"
DATE=$(date +%Y%m%d_%H%M%S)

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# 创建备份目录
mkdir -p $BACKUP_DIR

log_info "开始备份..."

# 备份数据库
if [ -f "$PROJECT_DIR/server/data.db" ]; then
    log_info "备份数据库..."
    cp $PROJECT_DIR/server/data.db $BACKUP_DIR/data.db.$DATE
    log_success "数据库备份完成"
fi

# 备份配置文件
if [ -f "$PROJECT_DIR/server/data/ip2region.xdb" ]; then
    log_info "备份 IP 数据库..."
    cp $PROJECT_DIR/server/data/ip2region.xdb $BACKUP_DIR/ip2region.xdb.$DATE
    log_success "IP 数据库备份完成"
fi

# 备份上传文件（如果有）
if [ -d "$PROJECT_DIR/server/uploads" ]; then
    log_info "备份上传文件..."
    tar -czf $BACKUP_DIR/uploads.$DATE.tar.gz -C $PROJECT_DIR/server uploads
    log_success "上传文件备份完成"
fi

# 压缩所有备份
log_info "压缩备份文件..."
tar -czf $BACKUP_DIR/backup_$DATE.tar.gz -C $BACKUP_DIR \
    data.db.$DATE \
    ip2region.xdb.$DATE \
    uploads.$DATE.tar.gz 2>/dev/null

# 清理临时备份文件
rm -f $BACKUP_DIR/data.db.$DATE
rm -f $BACKUP_DIR/ip2region.xdb.$DATE
rm -f $BACKUP_DIR/uploads.$DATE.tar.gz

log_success "备份完成: $BACKUP_DIR/backup_$DATE.tar.gz"

# 显示备份文件大小
BACKUP_SIZE=$(du -h $BACKUP_DIR/backup_$DATE.tar.gz | awk '{print $1}')
log_success "备份文件大小: $BACKUP_SIZE"

# 清理旧备份（保留最近 7 天）
log_info "清理旧备份文件（保留最近 7 天）..."
find $BACKUP_DIR -name "backup_*.tar.gz" -mtime +7 -delete
log_success "清理完成"

# 列出当前备份
echo
log_info "当前备份列表:"
ls -lh $BACKUP_DIR/backup_*.tar.gz | tail -10
