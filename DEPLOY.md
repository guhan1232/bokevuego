# BokeUI 博客系统 - Linux 部署指南

## 快速部署（推荐）

### 一键部署

将项目上传到服务器后，执行以下命令：

```bash
chmod +x deploy.sh
./deploy.sh
```

部署脚本会自动：
- ✅ 检测系统环境
- ✅ 安装 Go、Node.js 等依赖
- ✅ 编译后端程序
- ✅ 构建前端页面
- ✅ 配置 systemd 服务
- ✅ 开放防火墙端口
- ✅ 启动服务

### 部署时间

- 首次部署：约 10-15 分钟（取决于网络速度）
- 更新部署：约 2-3 分钟

---

## 手动部署

### 1. 环境要求

- **操作系统**：Ubuntu 18.04+ / Debian 10+ / CentOS 7+
- **Go**：1.19 或更高版本
- **Node.js**：16.x 或更高版本
- **内存**：建议 1GB 以上
- **磁盘**：建议 2GB 以上

### 2. 安装依赖

#### Ubuntu/Debian

```bash
# 安装 Go
wget https://golang.google.cn/dl/go1.21.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 安装 Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
apt-get install -y nodejs

# 安装构建工具
apt-get install -y build-essential
```

#### CentOS/RHEL

```bash
# 安装 Go
wget https://golang.google.cn/dl/go1.21.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 安装 Node.js
curl -fsSL https://rpm.nodesource.com/setup_18.x | bash -
yum install -y nodejs

# 安装构建工具
yum install -y gcc gcc-c++ make
```

### 3. 构建项目

```bash
# 克隆项目（或上传项目文件）
git clone <your-repo-url>
cd bokeui

# 构建后端
cd server
go mod download
go build -o bokeui ./cmd
chmod +x bokeui
cd ..

# 构建前端
cd admin
npm install
npm run build
cd ..

cd web
npm install
npm run build
cd ..
```

### 4. 启动服务

#### 方式一：Systemd（推荐）

```bash
# 创建服务文件
cat > /etc/systemd/system/bokeui.service <<EOF
[Unit]
Description=BokeUI Blog System
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/path/to/bokeui/server
ExecStart=/path/to/bokeui/server/bokeui
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
systemctl daemon-reload
systemctl enable bokeui
systemctl start bokeui
```

#### 方式二：PM2

```bash
# 安装 PM2
npm install -g pm2

# 启动服务
pm2 start server/bokeui --name bokeui
pm2 save
pm2 startup
```

#### 方式三：直接运行

```bash
cd server
nohup ./bokeui > ../logs/bokeui.log 2>&1 &
```

---

## 配置说明

### 端口配置

默认端口：`9088`

修改端口：在 `server/cmd/main.go` 中修改端口号

### 防火墙配置

#### Ubuntu/Debian (UFW)

```bash
ufw allow 9088/tcp
ufw reload
```

#### CentOS/RHEL (Firewalld)

```bash
firewall-cmd --permanent --add-port=9088/tcp
firewall-cmd --reload
```

### Nginx 反向代理（可选）

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://127.0.0.1:9088;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## 管理命令

### Systemd 方式

```bash
# 启动服务
systemctl start bokeui

# 停止服务
systemctl stop bokeui

# 重启服务
systemctl restart bokeui

# 查看状态
systemctl status bokeui

# 查看日志
journalctl -u bokeui -f
```

### PM2 方式

```bash
# 启动服务
pm2 start bokeui

# 停止服务
pm2 stop bokeui

# 重启服务
pm2 restart bokeui

# 查看日志
pm2 logs bokeui
```

---

## 更新部署

### 使用更新脚本

```bash
chmod +x update.sh
./update.sh
```

### 手动更新

```bash
# 拉取最新代码
git pull

# 重新构建
cd server && go build -o bokeui ./cmd && cd ..
cd admin && npm run build && cd ..
cd web && npm run build && cd ..

# 重启服务
systemctl restart bokeui
# 或
pm2 restart bokeui
```

---

## 停止服务

```bash
chmod +x stop.sh
./stop.sh
```

---

## 访问地址

部署完成后，通过以下地址访问：

- **前台首页**：`http://your-server-ip:9088`
- **后台管理**：`http://your-server-ip:9088/admin`
- **默认账号**：`admin` / `admin123`

⚠️ **重要提示**：首次登录后请立即修改默认密码！

---

## 常见问题

### 1. 端口被占用

```bash
# 查看端口占用
netstat -tunlp | grep 9088

# 杀死占用进程
kill -9 <PID>
```

### 2. 权限不足

```bash
# 给脚本执行权限
chmod +x *.sh

# 给程序执行权限
chmod +x server/bokeui
```

### 3. Go 编译失败

```bash
# 设置 Go 环境变量
export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on
```

### 4. Node.js 构建失败

```bash
# 清除缓存重新安装
cd admin
rm -rf node_modules package-lock.json
npm install
npm run build
```

### 5. 数据库初始化失败

数据库文件位于 `server/data.db`，首次运行会自动创建。如需重置：

```bash
rm server/data.db*
systemctl restart bokeui
```

---

## 目录结构

```
bokeui/
├── server/              # 后端程序
│   ├── bokeui          # 编译后的可执行文件
│   ├── data.db         # SQLite 数据库
│   ├── data/           # 数据目录
│   │   └── ip2region.xdb  # IP 地理位置数据库
│   └── internal/       # 源代码
├── admin/dist/         # 管理后台构建产物
├── web/dist/           # 前台页面构建产物
├── logs/               # 日志目录
├── deploy.sh           # 一键部署脚本
├── update.sh           # 更新脚本
└── stop.sh             # 停止脚本
```

---

## 性能优化建议

1. **使用 Nginx 反向代理**
   - 处理静态文件
   - SSL/HTTPS 支持
   - 负载均衡

2. **开启 Gzip 压缩**
   - 减少传输体积
   - 提升访问速度

3. **使用 CDN**
   - 加速静态资源访问
   - 减轻服务器负担

4. **定期备份数据库**
   ```bash
   # 备份数据库
   cp server/data.db backups/data.db.$(date +%Y%m%d)
   ```

---

## 技术支持

- GitHub Issues：[提交问题](https://github.com/your-repo/bokeui/issues)
- 文档：[在线文档](https://docs.bokeui.com)

---

## 开源协议

MIT License

Copyright (c) 2024 BokeUI
