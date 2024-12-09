@echo off

@REM 定义 jump 函数
@REM 示例调用
@REM call :jump path/to/dir
:jump
if "%~1"=="" (
  echo Please input terget dir path or name for jump.
  echo.
  echo Usage: j [PATH or DIRNAME]
  echo.
  echo Example:
  echo   j home
  goto :eof
)

set "input=%*"
set "cmdline=kite.exe tool jump get '%input%'"

:: 执行命令并将输出存储在变量中
for /f "delims=" %%i in ('%cmdline%') do (
    set "output=%%i"
)

:: 执行命令并将输出存储在变量中 - 不行
:: kite.exe tool jump get '%input%' | set /p output=

:: 跳转目录
:: echo Dest: %output%
cd "%output%" || (
  echo Error: %output% is not a valid directory.
  goto :eof
)
kite.exe tool jump chdir "%output%"
goto :eof

call :jump %*
