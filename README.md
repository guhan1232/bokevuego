# BokeUI - 全栈博客系统

一个基于 **Vue3 + Go (Gin)** 的全栈博客系统，包含管理后台和前台展示，统一在 **9088 端口** 运行。

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go + Gin + SQLite |
| 管理后台 | Vue3 + Element Plus + Pinia + Vite |
| 前台展示 | Vue3 + Marked.js + Vite |

## 快速启动

### 前置要求

- [Go](https://go.dev/) 1.21+
- [Node.js](https://nodejs.org/) 18+
- GCC (SQLite 编译需要)

### 一键启动（推荐）

```bash
# Windows
start.bat

# Linux/Mac
chmod +x start.sh && ./start.sh
```

### 手动启动

```bash
# 1. 构建前端
cd admin && npm install && npx vite build && cd ..
cd web && npm install && npx vite build && cd ..

# 2. 启动服务
cd server && go run cmd/main.go
```

启动后访问：
- **前台首页**: http://localhost:9088
- **管理后台**: http://localhost:9088/admin
- **API 接口**: http://localhost:9088/api

默认管理员: `admin` / `admin123`

## 开发模式

开发时可分别启动前后端热更新：

```bash
# 终端1 - 后端
cd server && go run cmd/main.go

# 终端2 - 管理后台 (http://localhost:5173)
cd admin && npm run dev

# 终端3 - 前台 (http://localhost:5174)
cd web && npm run dev
```

## 功能

- 文章管理 (增删改查、草稿/发布)
- 分类与标签
- 留言板 (前台留言 + 后台回复)
- 站点设置 (标题、副标题、配色)
- 仪表盘统计
- JWT 认证
- 响应式设计 (手机/平板/桌面)

## API

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/login | 登录 | 否 |
| GET | /api/public/articles | 文章列表 | 否 |
| GET | /api/public/articles/:id | 文章详情 | 否 |
| POST | /api/public/messages | 提交留言 | 否 |
| GET | /api/articles | 管理文章列表 | 是 |
| POST | /api/articles | 创建文章 | 是 |
| PUT | /api/articles/:id | 更新文章 | 是 |
| DELETE | /api/articles/:id | 删除文章 | 是 |
| GET | /api/stats | 统计数据 | 是 |
| GET/PUT | /api/configs | 站点配置 | 是 |
| GET | /api/messages | 留言列表 | 是 |
| POST | /api/messages/reply | 回复留言 | 是 |
