@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ================================
echo   BokeUI Blog - 构建并启动
echo ================================
echo.

set "ROOT_DIR=%~dp0"
cd /d "%ROOT_DIR%"

:: 检查 Node.js
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未找到 Node.js，请先安装 Node.js
    pause
    exit /b 1
)

:: 检查 Go
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未找到 Go，请先安装 Go
    pause
    exit /b 1
)

:: 检查并安装 admin 依赖
echo [1/4] 检查管理后台依赖...
cd /d "%ROOT_DIR%admin"
if not exist "node_modules" (
    echo 正在安装管理后台依赖...
    call npm install --registry=https://registry.npmmirror.com
)
echo 构建管理后台...
call npm run build
if %errorlevel% neq 0 (
    echo [错误] 管理后台构建失败
    pause
    exit /b 1
)

:: 检查并安装 web 依赖
echo [2/4] 检查前台页面依赖...
cd /d "%ROOT_DIR%web"
if not exist "node_modules" (
    echo 正在安装前台页面依赖...
    call npm install --registry=https://registry.npmmirror.com
)
echo 构建前台页面...
call npm run build
if %errorlevel% neq 0 (
    echo [错误] 前台页面构建失败
    pause
    exit /b 1
)

:: 编译后端
echo [3/4] 编译后端程序...
cd /d "%ROOT_DIR%server"
if not exist "bokeui.exe" (
    echo 正在编译后端...
    go build -o bokeui.exe ./cmd
    if %errorlevel% neq 0 (
        echo [错误] 后端编译失败
        pause
        exit /b 1
    )
)

:: 启动服务
echo [4/4] 启动服务 (端口 9088)...
echo.
echo ================================
echo   服务启动成功！
echo ================================
echo   前台: http://localhost:9088
echo   后台: http://localhost:9088/admin
echo   账号: admin / admin123
echo ================================
echo.
echo 按 Ctrl+C 可停止服务
echo.

cd /d "%ROOT_DIR%server"
bokeui.exe
