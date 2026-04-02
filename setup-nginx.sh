#!/bin/bash

# BokeUI Nginx 反向代理配置脚本

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

check_nginx() {
    if ! command -v nginx &> /dev/null; then
        log_info "安装 Nginx..."
        if [ -f /etc/os-release ]; then
            . /etc/os-release
            if [ "$ID" = "ubuntu" ] || [ "$ID" = "debian" ]; then
                apt-get update && apt-get install -y nginx
            elif [ "$ID" = "centos" ] || [ "$ID" = "rhel" ]; then
                yum install -y epel-release && yum install -y nginx
            fi
        fi
        log_success "Nginx 安装完成"
    else
        log_success "Nginx 已安装"
    fi
}

configure_domain() {
    echo -e "${GREEN}BokeUI Nginx 配置向导${NC}"
    echo
    
    read -p "域名（如 example.com）: " DOMAIN
    [ -z "$DOMAIN" ] && { log_error "域名不能为空"; exit 1; }
    
    read -p "后端端口 [9088]: " PORT
    PORT=${PORT:-9088}
    
    read -p "启用 HTTPS？(y/n): " -n 1 -r && echo
    ENABLE_HTTPS=$REPLY
    
    SSL_CERT="" SSL_KEY=""
    if [[ $ENABLE_HTTPS =~ ^[Yy]$ ]]; then
        read -p "SSL 证书路径: " SSL_CERT
        read -p "SSL 私钥路径: " SSL_KEY
        [ ! -f "$SSL_CERT" ] || [ ! -f "$SSL_KEY" ] && { log_error "证书文件不存在"; exit 1; }
    fi
    
    create_nginx_config
    enable_config
    test_and_reload
    show_info
}

create_nginx_config() {
    CONFIG_FILE="/etc/nginx/sites-available/bokeui"
    log_info "创建配置文件..."
    
    if [[ $ENABLE_HTTPS =~ ^[Yy]$ ]]; then
        cat > $CONFIG_FILE <<EOF
server {
    listen 80;
    server_name ${DOMAIN};
    return 301 https://\$server_name\$request_uri;
}

server {
    listen 443 ssl http2;
    server_name ${DOMAIN};
    
    ssl_certificate ${SSL_CERT};
    ssl_certificate_key ${SSL_KEY};
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    
    location / {
        proxy_pass http://127.0.0.1:${PORT};
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        client_max_body_size 50M;
    }
    
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff2?)$ {
        proxy_pass http://127.0.0.1:${PORT};
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
EOF
    else
        cat > $CONFIG_FILE <<EOF
server {
    listen 80;
    server_name ${DOMAIN};
    
    location / {
        proxy_pass http://127.0.0.1:${PORT};
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        client_max_body_size 50M;
    }
    
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff2?)$ {
        proxy_pass http://127.0.0.1:${PORT};
        expires 30d;
    }
}
EOF
    fi
    log_success "配置文件已创建"
}

enable_config() {
    log_info "启用配置..."
    ln -sf /etc/nginx/sites-available/bokeui /etc/nginx/sites-enabled/
    rm -f /etc/nginx/sites-enabled/default
    log_success "配置已启用"
}

test_and_reload() {
    log_info "测试配置..."
    nginx -t && systemctl reload nginx && log_success "Nginx 已重载"
}

show_info() {
    echo
    echo -e "${GREEN}配置完成！${NC}"
    echo -e "域名: ${BLUE}http://${DOMAIN}${NC}"
    [[ $ENABLE_HTTPS =~ ^[Yy]$ ]] && echo -e "HTTPS: ${BLUE}https://${DOMAIN}${NC}"
    echo -e "后台: ${BLUE}http://${DOMAIN}/admin${NC}"
}

check_root() {
    [ "$EUID" -ne 0 ] && { log_error "请使用 root 或 sudo 运行"; exit 1; }
}

check_root
check_nginx
configure_domain
