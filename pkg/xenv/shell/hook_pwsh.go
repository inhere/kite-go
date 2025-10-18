package shell

import (
	"fmt"
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
)

func (sg *XenvScriptGenerator) generatePwshScripts() string {
	var sb strings.Builder
	// 添加全局环境变量
	maputil.EachTypedMap(sg.cfg.GlobalEnv, func(key, value string) {
		sb.WriteString(fmt.Sprintf("$env:%s = '%s'\n", strings.ToUpper(key), value))
	})

	// 添加全局PATH
	for _, path := range sg.cfg.GlobalPaths {
		sb.WriteString(fmt.Sprintf("$env:PATH += ';%s'\n", path))
	}

	// 添加全局别名
	maputil.EachTypedMap(sg.cfg.ShellAliases, func(key, value string) {
		// 复杂 value, 封装为简易方法 eg: function ll { ls.exe -alh $args }
		if strutil.ContainsByte(value, ' ') {
			sb.WriteString(fmt.Sprintf("function %s() { %s $args }\n", key, value))
		} else {
			// 简单 value, 直接使用 Set-Alias
			sb.WriteString(fmt.Sprintf("Set-Alias -name %s -value %s\n", key, value))
		}
	})

	return strutil.Replaces(PwshHookTemplate, map[string]string{
		"{{HooksDir}}": sg.cfg.ShellHooksDir,
		"{{EnvAliases}}": sb.String(),
	})
}

// PwshHookTemplate PowerShell hook模板
//
// Config for pwsh:
//
//	# write to profile. (find by: echo $Profile)
//	# Method 1:
//	Invoke-Expression (&kite xenv shell --type pwsh)
//
//	# Method 2:
//	kite xenv shell --type pwsh | Out-String | Invoke-Expression
var PwshHookTemplate = `# xenv PowerShell hook
# This script enables xenv to work in PowerShell shells

# Function to set up xenv in the current shell
function Setup-Xenv {
    # Mark hook enabled
    $env:XENV_HOOK_SHELL = "pwsh"
    # Set up the xenv shims directory in PATH
    $xenvShimsDir = if ($env:XENV_ROOT) { "$env:XENV_ROOT\shims" } else { "$HOME\.xenv\shims" }

    # Add shims directory to PATH if it's not already there
    if ($env:PATH -notlike "*$xenvShimsDir*") {
        $env:PATH = "$xenvShimsDir;$env:PATH"
    }

    # Add global shell ENV and aliases
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
            "use" {
                # Implementation for switching tool versions
                & kenv use @Arguments
            }
            "unuse" {
                # Implementation for unusing tool versions
                & kenv unuse @Arguments
            }
            "shell" {
                # Output the shell commands needed to set up xenv
                & kenv shell pwsh
            }
            default {
                # For other commands, just pass through to xenv
                & kenv $Command @Arguments
            }
        }
    }

    # Auto-initialize xenv if needed
    $xenvrcPath = "$HOME\.xenvrc.ps1"
    if (Test-Path $xenvrcPath -PathType Leaf -and (-not $env:XENV_AUTO_INITIALIZED)) {
        . $xenvrcPath
        $env:XENV_AUTO_INITIALIZED = "1"
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
`
