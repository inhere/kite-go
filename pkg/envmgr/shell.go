package envmgr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DefaultShellGenerator 默认Shell脚本生成器实现
type DefaultShellGenerator struct {
	kiteCommand string // kite命令路径
}

// NewShellGenerator 创建Shell脚本生成器
func NewShellGenerator(kiteCommand string) *DefaultShellGenerator {
	return &DefaultShellGenerator{
		kiteCommand: kiteCommand,
	}
}

// GenerateScript 生成shell脚本
func (sg *DefaultShellGenerator) GenerateScript(shellType ShellType, state *ActiveState, config *ShellEnvConfig) (string, error) {
	if !shellType.IsValid() {
		return "", fmt.Errorf("unsupported shell type: %s", shellType)
	}
	
	var script strings.Builder
	
	// 生成头部注释
	sg.writeHeader(&script, shellType)
	
	// 生成ktenv函数
	ktenvFunc, err := sg.GenerateKtenvFunction(shellType)
	if err != nil {
		return "", fmt.Errorf("failed to generate ktenv function: %w", err)
	}
	script.WriteString(ktenvFunc)
	script.WriteString("\n")
	
	// 生成环境变量设置
	if len(state.AddEnvs) > 0 {
		envScript, err := sg.GenerateEnvVars(shellType, state.AddEnvs)
		if err != nil {
			return "", fmt.Errorf("failed to generate env vars: %w", err)
		}
		script.WriteString(envScript)
		script.WriteString("\n")
	}
	
	// 生成PATH更新
	if len(state.AddPaths) > 0 {
		pathScript, err := sg.GeneratePathUpdate(shellType, state.AddPaths, PathAdd)
		if err != nil {
			return "", fmt.Errorf("failed to generate path update: %w", err)
		}
		script.WriteString(pathScript)
		script.WriteString("\n")
	}
	
	// 生成自定义脚本加载
	customScript := sg.generateCustomScriptLoading(shellType)
	if customScript != "" {
		script.WriteString(customScript)
		script.WriteString("\n")
	}
	
	return script.String(), nil
}

// GenerateKtenvFunction 生成ktenv函数
func (sg *DefaultShellGenerator) GenerateKtenvFunction(shellType ShellType) (string, error) {
	switch shellType {
	case ShellBash, ShellZsh:
		return sg.generateBashKtenvFunction(), nil
	case ShellPowerShell:
		return sg.generatePowerShellKtenvFunction(), nil
	case ShellCmd:
		return sg.generateCmdKtenvFunction(), nil
	default:
		return "", fmt.Errorf("unsupported shell type: %s", shellType)
	}
}

// GenerateEnvVars 生成环境变量设置
func (sg *DefaultShellGenerator) GenerateEnvVars(shellType ShellType, envs map[string]string) (string, error) {
	if len(envs) == 0 {
		return "", nil
	}
	
	var script strings.Builder
	
	switch shellType {
	case ShellBash, ShellZsh:
		script.WriteString("# Environment variables\n")
		for name, value := range envs {
			script.WriteString(fmt.Sprintf("export %s=\"%s\"\n", name, sg.escapeShellValue(value)))
		}
	case ShellPowerShell:
		script.WriteString("# Environment variables\n")
		for name, value := range envs {
			script.WriteString(fmt.Sprintf("$env:%s = \"%s\"\n", name, sg.escapePowerShellValue(value)))
		}
	case ShellCmd:
		script.WriteString("@REM Environment variables\n")
		for name, value := range envs {
			script.WriteString(fmt.Sprintf("set %s=%s\n", name, sg.escapeCmdValue(value)))
		}
	default:
		return "", fmt.Errorf("unsupported shell type: %s", shellType)
	}
	
	return script.String(), nil
}

// GeneratePathUpdate 生成PATH更新脚本
func (sg *DefaultShellGenerator) GeneratePathUpdate(shellType ShellType, paths []string, operation PathOperation) (string, error) {
	if len(paths) == 0 {
		return "", nil
	}
	
	var script strings.Builder
	
	switch shellType {
	case ShellBash, ShellZsh:
		script.WriteString("# PATH updates\n")
		for _, path := range paths {
			switch operation {
			case PathAdd:
				script.WriteString(fmt.Sprintf("export PATH=\"%s:$PATH\"\n", sg.escapeShellValue(path)))
			case PathRemove:
				script.WriteString(fmt.Sprintf("export PATH=$(echo $PATH | sed -e 's|%s:||g' -e 's|:%s||g' -e 's|^%s$||g')\n", 
					sg.escapeShellValue(path), sg.escapeShellValue(path), sg.escapeShellValue(path)))
			}
		}
	case ShellPowerShell:
		script.WriteString("# PATH updates\n")
		for _, path := range paths {
			switch operation {
			case PathAdd:
				script.WriteString(fmt.Sprintf("$env:PATH = \"%s;\" + $env:PATH\n", sg.escapePowerShellValue(path)))
			case PathRemove:
				script.WriteString(fmt.Sprintf("$env:PATH = $env:PATH -replace [regex]::Escape(\"%s;\"), \"\"\n", sg.escapePowerShellValue(path)))
				script.WriteString(fmt.Sprintf("$env:PATH = $env:PATH -replace [regex]::Escape(\";%s\"), \"\"\n", sg.escapePowerShellValue(path)))
			}
		}
	case ShellCmd:
		script.WriteString("@REM PATH updates\n")
		for _, path := range paths {
			switch operation {
			case PathAdd:
				script.WriteString(fmt.Sprintf("set PATH=%s;%%PATH%%\n", sg.escapeCmdValue(path)))
			}
		}
	default:
		return "", fmt.Errorf("unsupported shell type: %s", shellType)
	}
	
	return script.String(), nil
}

// generateBashKtenvFunction 生成Bash/Zsh的ktenv函数
func (sg *DefaultShellGenerator) generateBashKtenvFunction() string {
	return fmt.Sprintf(`# ktenv function for environment management
ktenv() {
    local cmd="$1"
    shift
    
    case "$cmd" in
        use|unuse|add|list)
            # Call kite command and evaluate the result
            local result=$(%s dev env ktenv "$cmd" "$@")
            if [ $? -eq 0 ]; then
                if [ -n "$result" ]; then
                    eval "$result"
                fi
            else
                echo "$result" >&2
                return 1
            fi
            ;;
        help|--help|-h)
            %s dev env ktenv help
            ;;
        *)
            echo "ktenv: unknown command '$cmd'" >&2
            echo "Available commands: use, unuse, add, list, help" >&2
            return 1
            ;;
    esac
}

# Enable command completion for ktenv
if command -v complete >/dev/null 2>&1; then
    complete -W "use unuse add list help" ktenv
fi
`, sg.kiteCommand, sg.kiteCommand)
}

// generatePowerShellKtenvFunction 生成PowerShell的ktenv函数
func (sg *DefaultShellGenerator) generatePowerShellKtenvFunction() string {
	return fmt.Sprintf(`# ktenv function for environment management
function ktenv {
    param(
        [Parameter(Position=0, Mandatory=$true)]
        [string]$Command,
        
        [Parameter(Position=1, ValueFromRemainingArguments=$true)]
        [string[]]$Arguments
    )
    
    switch ($Command) {
        { $_ -in @('use', 'unuse', 'add', 'list') } {
            # Call kite command and evaluate the result
            $result = & %s dev env ktenv $Command @Arguments
            if ($LASTEXITCODE -eq 0) {
                if ($result) {
                    Invoke-Expression $result
                }
            } else {
                Write-Error $result
                return
            }
        }
        { $_ -in @('help', '--help', '-h') } {
            & %s dev env ktenv help
        }
        default {
            Write-Error "ktenv: unknown command '$Command'"
            Write-Host "Available commands: use, unuse, add, list, help" -ForegroundColor Yellow
            return
        }
    }
}

# Enable command completion for ktenv
Register-ArgumentCompleter -CommandName ktenv -ParameterName Command -ScriptBlock {
    param($commandName, $parameterName, $wordToComplete, $commandAst, $fakeBoundParameters)
    @('use', 'unuse', 'add', 'list', 'help') | Where-Object { $_ -like "$wordToComplete*" }
}
`, sg.kiteCommand, sg.kiteCommand)
}

// generateCmdKtenvFunction 生成CMD的ktenv函数
func (sg *DefaultShellGenerator) generateCmdKtenvFunction() string {
	return fmt.Sprintf(`@echo off
REM ktenv function for environment management (CMD batch implementation)

if "%%1"=="" (
    echo ktenv: missing command
    echo Available commands: use, unuse, add, list, help
    exit /b 1
)

set "KTENV_CMD=%%1"
shift

if "%%KTENV_CMD%%"=="use" goto :ktenv_use
if "%%KTENV_CMD%%"=="unuse" goto :ktenv_unuse  
if "%%KTENV_CMD%%"=="add" goto :ktenv_add
if "%%KTENV_CMD%%"=="list" goto :ktenv_list
if "%%KTENV_CMD%%"=="help" goto :ktenv_help
if "%%KTENV_CMD%%"=="--help" goto :ktenv_help
if "%%KTENV_CMD%%"=="-h" goto :ktenv_help

echo ktenv: unknown command '%%KTENV_CMD%%'
echo Available commands: use, unuse, add, list, help
exit /b 1

:ktenv_use
:ktenv_unuse
:ktenv_add
:ktenv_list
    REM Call kite command and execute the result
    for /f "delims=" %%%%i in ('%s dev env ktenv %%KTENV_CMD%% %%*') do %%%%i
    exit /b %%ERRORLEVEL%%

:ktenv_help
    %s dev env ktenv help
    exit /b 0
`, sg.kiteCommand, sg.kiteCommand)
}

// writeHeader 写入脚本头部注释
func (sg *DefaultShellGenerator) writeHeader(script *strings.Builder, shellType ShellType) {
	switch shellType {
	case ShellBash, ShellZsh:
		script.WriteString("#!/bin/bash\n")
		script.WriteString("# Kite Environment Management Shell Integration\n")
		script.WriteString("# Generated automatically - do not edit manually\n\n")
	case ShellPowerShell:
		script.WriteString("# Kite Environment Management Shell Integration\n")
		script.WriteString("# Generated automatically - do not edit manually\n\n")
	case ShellCmd:
		script.WriteString("@REM Kite Environment Management Shell Integration\n")
		script.WriteString("@REM Generated automatically - do not edit manually\n\n")
	}
}

// generateCustomScriptLoading 生成自定义脚本加载
func (sg *DefaultShellGenerator) generateCustomScriptLoading(shellType ShellType) string {
	switch shellType {
	case ShellBash, ShellZsh:
		return `# Load custom scripts
if [ -d "$HOME/.kite-go/data/shell_env" ]; then
    for script in "$HOME/.kite-go/data/shell_env"/*.sh; do
        if [ -f "$script" ]; then
            source "$script"
        fi
    done
fi`
	case ShellPowerShell:
		return `# Load custom scripts
$customDir = "$env:USERPROFILE\.kite-go\data\shell_env"
if (Test-Path $customDir) {
    Get-ChildItem "$customDir\*.ps1" | ForEach-Object {
        if (Test-Path $_.FullName) {
            . $_.FullName
        }
    }
}`
	case ShellCmd:
		return `@REM Load custom scripts
if exist "%USERPROFILE%\.kite-go\data\shell_env\*.bat" (
    for %%f in ("%USERPROFILE%\.kite-go\data\shell_env\*.bat") do (
        call "%%f"
    )
)`
	}
	return ""
}

// escapeShellValue 转义Shell值
func (sg *DefaultShellGenerator) escapeShellValue(value string) string {
	// 简单的转义处理
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "$", "\\$")
	value = strings.ReplaceAll(value, "`", "\\`")
	return value
}

// escapePowerShellValue 转义PowerShell值
func (sg *DefaultShellGenerator) escapePowerShellValue(value string) string {
	// PowerShell字符串转义
	value = strings.ReplaceAll(value, "\"", "`\"")
	value = strings.ReplaceAll(value, "$", "`$")
	value = strings.ReplaceAll(value, "`", "``")
	return value
}

// escapeCmdValue 转义CMD值
func (sg *DefaultShellGenerator) escapeCmdValue(value string) string {
	// CMD特殊字符转义
	value = strings.ReplaceAll(value, "%", "%%")
	value = strings.ReplaceAll(value, "&", "^&")
	value = strings.ReplaceAll(value, "|", "^|")
	value = strings.ReplaceAll(value, "<", "^<")
	value = strings.ReplaceAll(value, ">", "^>")
	return value
}

// GenerateShellActivationScript 生成shell激活脚本
func (sg *DefaultShellGenerator) GenerateShellActivationScript(shellType ShellType) (string, error) {
	switch shellType {
	case ShellBash:
		return `# Add to ~/.bashrc:
eval "$(kite dev env shell bash)"`, nil
	case ShellZsh:
		return `# Add to ~/.zshrc:
eval "$(kite dev env shell zsh)"`, nil
	case ShellPowerShell:
		return `# Add to $PROFILE:
Invoke-Expression (kite dev env shell pwsh)`, nil
	case ShellCmd:
		return `REM Add to your batch script or run manually:
kite dev env shell cmd > ktenv.bat && call ktenv.bat`, nil
	default:
		return "", fmt.Errorf("unsupported shell type: %s", shellType)
	}
}

// DetectShellType 检测当前shell类型
func DetectShellType() ShellType {
	// 从环境变量检测
	if shell := filepath.Base(strings.ToLower(GetEnvVar("SHELL", ""))); shell != "" {
		switch {
		case strings.Contains(shell, "bash"):
			return ShellBash
		case strings.Contains(shell, "zsh"):
			return ShellZsh
		case strings.Contains(shell, "pwsh") || strings.Contains(shell, "powershell"):
			return ShellPowerShell
		}
	}
	
	// 从PSModulePath检测PowerShell
	if GetEnvVar("PSModulePath", "") != "" {
		return ShellPowerShell
	}
	
	// 从ComSpec检测CMD
	if comspec := GetEnvVar("ComSpec", ""); strings.Contains(strings.ToLower(comspec), "cmd") {
		return ShellCmd
	}
	
	// 默认返回bash
	return ShellBash
}

// GetEnvVar 获取环境变量值
func GetEnvVar(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}