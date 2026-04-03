@echo off
chcp 65001 >nul
echo ================================
echo   BokeUI - 调试启动
echo ================================
echo.

cd /d "%~dp0server"

echo 当前工作目录:
cd
echo.

echo 检查前端文件:
if exist "..\admin\dist\index.html" (
    echo [OK] admin\dist\index.html 存在
) else (
    echo [错误] admin\dist\index.html 不存在
)

if exist "..\web\dist\index.html" (
    echo [OK] web\dist\index.html 存在
) else (
    echo [错误] web\dist\index.html 不存在
)

echo.
echo 启动服务...
echo.

bokeui.exe
