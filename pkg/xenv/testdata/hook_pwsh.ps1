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
        [string]$Result,
        [int]$ExitCode
    )

    if ($ExitCode -eq 0) {
        # debug
        Write-Output "----------------in Invoke-XenvResult--------------"
        Write-Output $result

        if ($Result) {
            # 检查结果是否包含 '--Expression--' 分隔符
            if ($Result.Contains('--Expression--')) {
                # 使用 '--Expression--' 分割内容
                $parts = $Result.Split('--Expression--', 2)
                $msgPart = $parts[0].Trim()
                $exprPart = $parts[1].Trim()

                # 后面部分当做代码执行
                if ($exprPart) {
                    Invoke-Expression $exprPart
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

# Override cd command to automatically run: kite xenv init-direnv
function global:cd {
    param(
        [Parameter(Position=0)]
        [string]$Path = $HOME,
        [switch]$PassThru
    )

    # Call original Set-Location
    if ($PassThru) {
        Set-Location $Path -PassThru
    } else {
        Set-Location $Path
    }

    # Check if xenv is available and run init-direnv
    if (Get-Command kite -ErrorAction SilentlyContinue) {
        # Run kite xenv init-direnv, eval result scripts
        $result = & kite xenv init-direnv
        Invoke-XenvResult -Result $result -ExitCode $LASTEXITCODE
    }
}

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

    {{EnvAliases}}

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
                $result = & kite xenv $Command @Arguments
                # Write-Output $result # DEBUG
                Invoke-XenvResult -Result $result -ExitCode $LASTEXITCODE
            }
            { $_ -in @('set', 'unset') } {
                $result = & kite xenv env $Command @Arguments
                Invoke-XenvResult -Result $result -ExitCode $LASTEXITCODE
            }
            default {
                # For other commands, just pass through to xenv
                & kite xenv $Command @Arguments
            }
        }
    }

    # fire xenv hooks to kite, use for generate code to exec TODO
    $result1 = & kite xenv hook-init --type pwsh
    # TODO exec output result
    Invoke-XenvResult -Result $result1 -ExitCode $LASTEXITCODE

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
