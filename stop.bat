@echo off
chcp 65001 >nul
echo ================================
echo   BokeUI Blog - 停止服务
echo ================================
echo.

echo 正在查找运行中的服务...
tasklist /FI "IMAGENAME eq bokeui.exe" 2>nul | find /I "bokeui.exe" >nul
if %errorlevel% equ 0 (
    echo 找到运行中的服务，正在停止...
    taskkill /F /IM bokeui.exe >nul 2>&1
    if %errorlevel% equ 0 (
        echo 服务已停止
    ) else (
        echo 停止服务失败，可能需要管理员权限
    )
) else (
    echo 未找到运行中的服务
)

echo.
pause
