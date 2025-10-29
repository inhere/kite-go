package shell

import (
	"fmt"
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

func (sg *XenvScriptGenerator) generatePwshScripts() string {
	cfg := sg.cfg
	var sb strings.Builder
	// 添加全局环境变量
	if len(sg.cfg.GlobalEnv) > 0 {
		sb.WriteString("  # Add global ENV variables from kite xenv\n")
		maputil.EachTypedMap(cfg.GlobalEnv, func(key, value string) {
			sb.WriteString(fmt.Sprintf("  $env:%s = '%s'\n", strings.ToUpper(key), value))
		})
	}

	// 添加全局PATH
	if len(cfg.GlobalPaths) > 0 {
		sb.WriteString("  # Add global PATH from kite xenv\n")
		paths := strings.Join(sg.cfg.GlobalPaths, ";")
		sb.WriteString(fmt.Sprintf("  $env:PATH = '%s;' + $env:PATH\n", paths))
	}

	// 添加全局别名
	if len(sg.cfg.ShellAliases) > 0 {
		sb.WriteString("  # Add global aliases from kite xenv\n")
		maputil.EachTypedMap(cfg.ShellAliases, func(key, value string) {
			// 复杂 value, 封装为简易方法 eg: function ll { ls.exe -alh $args }
			if strutil.ContainsByte(value, ' ') {
				sb.WriteString(fmt.Sprintf("  function %s() { %s $args }\n", key, value))
			} else {
				// 简单 value, 直接使用 Set-Alias
				sb.WriteString(fmt.Sprintf("  Set-Alias -name %s -value %s\n", key, value))
			}
		})
	}

	return strutil.Replaces(PwshHookTemplate, map[string]string{
		"{{HooksDir}}": cfg.ShellHooksDir,
		"{{SessionId}}": models.SessionID(),
		"{{EnvAliases}}": sb.String(),
	})
}

// PwshHookTemplate PowerShell hook模板
//
// PS version:
//  echo $PSVersionTable.PSVersion.ToString()
//
// Config on pwsh:
//
//	# Write to profile.
//	 find by: echo $PROFILE.CurrentUserAllHosts
//
//	# Method 1:
//	Invoke-Expression (&kite xenv shell --type pwsh)
//
//	# Method 2:
//	kite xenv shell --type pwsh | Out-String | Invoke-Expression
var PwshHookTemplate = `# xenv PowerShell hook
# This script enables xenv to work in PowerShell shells

# Helper function to evaluate xenv command results
function Eval-XenvResult {
	param(
		[string]$Result,
		[int]$ExitCode
	)

	if ($ExitCode -eq 0) {
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
				Eval-XenvResult -Result $result -ExitCode $LASTEXITCODE
            }
            { $_ -in @('set', 'unset') } {
				$result = & kite xenv env $Command @Arguments
				Eval-XenvResult -Result $result -ExitCode $LASTEXITCODE
            }
            default {
                # For other commands, just pass through to xenv
                & kite xenv $Command @Arguments
            }
        }
    }

    # fire xenv hooks to kite, use for generate code to exec TODO
    # $result = & kite xenv hook-init --type pwsh
    # TODO exec output result

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
`
