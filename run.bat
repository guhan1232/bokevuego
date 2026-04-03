@echo off
chcp 65001 >nul
echo ================================
echo   BokeUI - 启动服务
echo ================================
echo.

set "ROOT_DIR=%~dp0"
cd /d "%ROOT_DIR%server"

echo 正在启动服务...
echo 请查看下方日志确认路径是否正确
echo.

bokeui.exe
