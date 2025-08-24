@rem 关闭命令回显，且当前行也不显示（@ 符号抑制该行自身的回显），使输出更简洁
@echo off

rem 创建一个局部环境，确保变量只在这个批处理文件中有效
setlocal

rem 操作系统
for /f "delims=" %%i in ('go env GOOS') do set OS=%%i
echo OS      : %OS%

rem CPU 架构
for /f "delims=" %%i in ('go env GOARCH') do set ARCH=%%i
echo ARCH    : %ARCH%

rem 当前目录
set CUR_DIR=%cd%
echo CUR_DIR : %CUR_DIR%

rem 输出目录
set OUT_DIR=%CUR_DIR%\build
echo OUT_DIR : %OUT_DIR%

rem 删除输出目录
if exist "%OUT_DIR%" rd /s /q "%OUT_DIR%"
rem 创建输出目录
mkdir "%OUT_DIR%"

rem 拷贝文件
rem 隐藏无用输出：> nul（标准输出），2> nul（错误输出）
copy /Y "config.ini" "%OUT_DIR%\" > nul

rem 构建
echo BUILDING ...
set OUT_NAME=gmon-%OS%-%ARCH%.exe
set OUT_PATH=%OUT_DIR%\%OUT_NAME%
cd "%CUR_DIR%" && go build -ldflags="-s -w" -o "%OUT_PATH%"

rem 压缩可执行文件
::upx -9 --brute --backup "%OUT_PATH%"

rem 启动命令
(
  echo @echo off
  echo setlocal
  echo title GMon
  echo %OUT_NAME%
  echo endlocal
  echo pause
) > "%OUT_DIR%\start.cmd"

endlocal

pause
