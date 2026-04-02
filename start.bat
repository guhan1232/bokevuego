@echo off
echo ================================
echo   BokeUI Blog - 构建并启动
echo ================================
echo.

echo [1/3] 构建管理后台...
cd admin && call npx vite build && cd ..
echo [2/3] 构建前台页面...
cd web && call npx vite build && cd ..
echo [3/3] 启动服务 (端口 9088)...
cd server && go run cmd/main.go
