@echo off
chcp 65001 >nul
echo ================================
echo   BokeUI Blog - 快速启动
echo ================================
echo.

set "ROOT_DIR=%~dp0"
cd /d "%ROOT_DIR%server"

if not exist "bokeui.exe" (
    echo [提示] 后端程序未编译，正在编译...
    go build -o bokeui.exe ./cmd
)

echo 服务启动中...
echo.
echo ================================
echo   前台: http://localhost:9088
echo   后台: http://localhost:9088/admin
echo   账号: admin / admin123
echo ================================
echo.

bokeui.exe
