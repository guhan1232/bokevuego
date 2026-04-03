@echo off
chcp 65001 >nul
echo ================================
echo   BokeUI Blog - 安装依赖
echo ================================
echo.

set "ROOT_DIR=%~dp0"

:: 检查 Node.js
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未找到 Node.js，请先安装 Node.js
    echo 下载地址: https://nodejs.org/
    pause
    exit /b 1
)

:: 检查 Go
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未找到 Go，请先安装 Go
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
)

echo Node.js 版本:
node -v
echo.
echo Go 版本:
go version
echo.

:: 安装 admin 依赖
echo [1/2] 安装管理后台依赖...
cd /d "%ROOT_DIR%admin"
call npm install --registry=https://registry.npmmirror.com
if %errorlevel% neq 0 (
    echo [错误] 管理后台依赖安装失败
    pause
    exit /b 1
)

:: 安装 web 依赖
echo [2/2] 安装前台页面依赖...
cd /d "%ROOT_DIR%web"
call npm install --registry=https://registry.npmmirror.com
if %errorlevel% neq 0 (
    echo [错误] 前台页面依赖安装失败
    pause
    exit /b 1
)

echo.
echo ================================
echo   依赖安装完成！
echo ================================
echo.
echo 运行 start.bat 构建并启动服务
echo.
pause
