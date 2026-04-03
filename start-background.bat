@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ================================
echo   BokeUI Blog - 后台启动
echo ================================
echo.

set "ROOT_DIR=%~dp0"
cd /d "%ROOT_DIR%"

:: 检查并编译后端
echo 检查后端程序...
cd /d "%ROOT_DIR%server"
if not exist "bokeui.exe" (
    echo 后端程序未编译，正在编译...
    go build -o bokeui.exe ./cmd
    if %errorlevel% neq 0 (
        echo [错误] 后端编译失败
        pause
        exit /b 1
    )
)

:: 后台启动服务
echo 启动服务中...
start "BokeUI Server" bokeui.exe

echo.
echo ================================
echo   服务已在后台启动！
echo ================================
echo   前台: http://localhost:9088
echo   后台: http://localhost:9088/admin
echo   账号: admin / admin123
echo ================================
echo.
echo 服务运行在新窗口中
echo 关闭窗口或按 Ctrl+C 可停止服务
echo.
pause
