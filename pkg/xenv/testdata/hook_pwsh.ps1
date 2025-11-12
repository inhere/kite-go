# xenv PowerShell hook
# This script enables xenv to work in PowerShell shells
#
# Config for pwsh:
#
#	# Write to profile.
#	 find by: echo $PROFILE.CurrentUserAllHosts
#
#	# Method 1:
#	Invoke-Expression (&kite xenv shell --type pwsh)
#
#	# Method 2:
#	kite xenv shell --type pwsh | Out-String | Invoke-Expression

# Helper function to evaluate xenv command results
function Invoke-XenvResult {
    param(
        [string]$CallFrom,
        [string]$Result,
        [int]$ExitCode
    )

    if ($ExitCode -eq 0) {
        if ($Result) {
            # debug
            Write-Host "----------------in Invoke-XenvResult($CallFrom)--------------" -ForegroundColor Green
            Write-Output $Result

            # TODO 使用 '--Expression--' 分割结果
            #  $parts = $Result -split '--Expression--', 2
            # if ($parts.Count -eq 2) {# 前面部分直接输出
            #     Write-Host $parts[0].Trim()
            #     # 后面部分动态执行
            #     $script = $parts[1].Trim()
            #     Write-Host ">>> 动态执行脚本：" -Fore Magenta
            #     Write-Host $script -Fore Cyan
            #     # Invoke-Expression $script
            #     [scriptblock]::Create($script).Invoke()
            # } else {
            #     # 没发现分隔符，原样输出
            #     Write-Host $parts[0]
            # }

            # 检查结果是否包含 '--Expression--' 分隔符
            if ($Result.Contains('--Expression--')) {
                # 使用 '--Expression--' 分割内容
                $parts = $Result.Split('--Expression--', 2)
                $msgPart = $parts[0].Trim()
                $exprPart = $parts[1].Trim()

                # 后面部分当做代码执行
                if ($exprPart) {
                    Invoke-Expression $exprPart
                    # [scriptblock]::Create($script).Invoke()
                }
                # 前面部分直接输出
                if ($msgPart) {
                    Write-Output $msgPart
                    # Write-Output $Result # DEBUG
                }
            } else {
                # 否则直接输出内容
                Write-Output $Result
            }
        }
    } else {
        Write-Error $Result
    }
}

# 创建一个全局变量来保存上一次的目录
#$global:lastPath = $null

# 保存原始的 Set-Location
$originalSetLocation = Get-Command Set-Location -CommandType Cmdlet
#$originalSetLocation = $function:Set-Location

# 重写 cd 命令
function Set-Location {
    param(
        [Parameter(Mandatory=$false, Position=0)]
        [string]$Path = $HOME,
        [switch]$PassThru
    )

    # 如果 Path=-, 回到最近的 lastPath 目录
#    if ($Path -eq "-") {  }

    # 保存最近的目录到ENV
    $currentPath = $PWD.Path
#    if ($currentPath -ne $Path) {
#        # TODO 处理离开目录时的逻辑，删除之前配置的ENV,PATH
#    }

#    $global:lastPath = $currentPath
    $env:PREV_PWD = $currentPath
    # 调用原始命令
    Write-Host "🔧 Goto $Path" -ForegroundColor Cyan
    # & $originalSetLocation @args
    if ($PassThru) {
        & $originalSetLocation $Path -PassThru
    } else {
        & $originalSetLocation $Path
    }

    # 获取当前目录
    # $currentPath = (Get-Location).Path
    $currentPath = $PWD.Path
    Write-Host "- PWD: $currentPath" -ForegroundColor Cyan

    # Check if xenv is available and run init-direnv
    if (Get-Command kite -ErrorAction SilentlyContinue) {
        # Run kite xenv init-direnv, eval result scripts
        $result = (& kite xenv init-direnv | Out-String)
        # Write-Output "DEBUG: \n$result"
        Invoke-XenvResult -CallFrom "Set-Location.init-direnv" -Result $result -ExitCode $LASTEXITCODE
    }
}

#Set-Alias -Name cd -Value Set-Location -Force -Option AllScope

# Function to set up xenv in the current shell
function Setup-Xenv {
    # Mark hook enabled
    $env:XENV_HOOK_SHELL = "pwsh"
    $env:XENV_SESSION_ID = "{{SessionId}}"
    # Set up the xenv shims directory in PATH
    $xenvShimsDir = if ($env:XENV_ROOT) { "$env:XENV_ROOT\shims" } else { "$HOME\.xenv\shims" }

    # Add shims directory to PATH if it's not already there
    if ($env:PATH -notlike "*$xenvShimsDir*") {
        $env:PATH = "$xenvShimsDir;$env:PATH"
    }

#{{EnvAliases}}

    # Define the xenv function to activate tools
    function global:xenv {
        param(
            [Parameter(Position=0)]
            [string]$Command,

            [Parameter(ValueFromRemainingArguments)]
            [string[]]$Arguments
        )

        switch ($Command) {
            { $_ -in @('use', 'unuse', 'env', 'path') } {
                # Call kite command and evaluate the result
                $result = (& kite xenv $Command @Arguments | Out-String)
                # Write-Output $result # DEBUG
                Invoke-XenvResult -CallFrom "xenv.$Command" -Result $result -ExitCode $LASTEXITCODE
            }
            { $_ -in @('set', 'unset') } {
                $result = (& kite xenv env $Command @Arguments | Out-String)
                Invoke-XenvResult -CallFrom "xenv.$Command" -Result $result -ExitCode $LASTEXITCODE
            }
            default {
                # For other commands, just pass through to xenv
                & kite xenv $Command @Arguments
            }
        }
    }

    # fire xenv hooks to kite, use for generate code to exec TODO
    $result_init_hook = & kite xenv shell-init-hook --type pwsh
    Invoke-XenvResult -CallFrom "Setup-Xenv.shell-init-hook" -Result $result_init_hook -ExitCode $LASTEXITCODE

    # Auto-initialize xenv if needed
    $xenvrcPath = "$HOME\.xenvrc.ps1"
    if ((Test-Path $xenvrcPath -PathType Leaf) -and (-not $env:XENV_AUTO_INIT)) {
        . $xenvrcPath
        $env:XENV_AUTO_INIT = "1"
    }

    # Load custom hooks script files
    $hookFiles = Get-ChildItem -Path "{{HooksDir}}" -Filter "*.ps1" -ErrorAction SilentlyContinue
    foreach ($file in $hookFiles) {
        if (Test-Path $file.FullName -PathType Leaf) {
            . $file.FullName
        }
    }
}

# Call setup function to initialize xenv
Setup-Xenv

# Enable command completion for xenv
Register-ArgumentCompleter -CommandName xenv -ParameterName Command -ScriptBlock {
    param($commandName, $parameterName, $wordToComplete, $commandAst, $fakeBoundParameters)
    @('use', 'unuse', 'env', 'set', 'unset', 'path', 'list', '--help') | Where-Object { $_ -like "$wordToComplete*" }
}
