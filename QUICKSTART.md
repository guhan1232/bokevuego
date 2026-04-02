# BokeUI 快速部署指南

## 📦 脚本说明

### 一键部署脚本

| 脚本 | 功能 | 使用场景 |
|------|------|----------|
| `deploy.sh` | 一键部署 | 首次部署或重新部署 |
| `start.sh` | 启动服务 | 手动启动服务 |
| `stop.sh` | 停止服务 | 手动停止服务 |
| `update.sh` | 更新部署 | 代码更新后重新部署 |
| `backup.sh` | 备份数据 | 定期备份 |
| `restore.sh` | 恢复数据 | 数据恢复 |

---

## 🚀 快速开始

### 1. 首次部署

```bash
# 赋予执行权限
chmod +x deploy.sh

# 执行一键部署
./deploy.sh
```

部署过程约 10-15 分钟，完成后会显示访问地址和默认账号。

### 2. 启动服务

```bash
./start.sh
```

选择启动方式：
- `1` - systemd（推荐，开机自启）
- `2` - PM2（进程管理）
- `3` - 直接运行（后台运行）

### 3. 停止服务

```bash
./stop.sh
```

自动检测并停止所有运行中的服务。

### 4. 更新部署

```bash
./update.sh
```

自动拉取最新代码、重新构建并重启服务。

### 5. 备份数据

```bash
./backup.sh
```

备份内容包括：
- SQLite 数据库
- IP 地理位置数据库
- 上传的文件（如果有）

备份文件保存在 `backups/` 目录，自动清理 7 天前的备份。

### 6. 恢复数据

```bash
./restore.sh
```

交互式选择要恢复的备份文件，自动停止服务、恢复数据、重启服务。

---

## 🌐 访问地址

部署完成后：

- **前台首页**：`http://your-server-ip:9088`
- **后台管理**：`http://your-server-ip:9088/admin`
- **默认账号**：`admin` / `admin123`

⚠️ **重要**：首次登录后立即修改默认密码！

---

## 📋 系统要求

- **操作系统**：Ubuntu 18.04+ / Debian 10+ / CentOS 7+
- **内存**：建议 1GB 以上
- **磁盘**：建议 2GB 以上
- **网络**：需要访问外网下载依赖

---

## 🔧 常用命令

### systemd 方式

```bash
systemctl start bokeui      # 启动
systemctl stop bokeui       # 停止
systemctl restart bokeui    # 重启
systemctl status bokeui     # 状态
journalctl -u bokeui -f     # 日志
```

### PM2 方式

```bash
pm2 start bokeui            # 启动
pm2 stop bokeui             # 停止
pm2 restart bokeui          # 重启
pm2 logs bokeui             # 日志
pm2 list                    # 列表
```

---

## 🔥 常见问题

### 端口被占用

```bash
# 查看 9088 端口占用
netstat -tunlp | grep 9088

# 杀死占用进程
kill -9 <PID>
```

### 权限不足

```bash
# 给所有脚本执行权限
chmod +x *.sh

# 给程序执行权限
chmod +x server/bokeui
```

### Go 编译失败

```bash
# 设置国内代理
export GOPROXY=https://goproxy.cn,direct
```

### Node.js 构建失败

```bash
# 清除缓存重新安装
cd admin
rm -rf node_modules package-lock.json
npm install
```

---

## 📚 详细文档

完整部署文档请查看：[DEPLOY.md](./DEPLOY.md)

---

## 💡 提示

1. **首次部署**：使用 `deploy.sh` 一键部署
2. **日常管理**：使用 `start.sh` / `stop.sh` 管理服务
3. **代码更新**：使用 `update.sh` 快速更新
4. **数据安全**：定期执行 `backup.sh` 备份

---

## 📞 技术支持

遇到问题？查看 [常见问题](./DEPLOY.md#常见问题) 或提交 [Issue](https://github.com/your-repo/bokeui/issues)
