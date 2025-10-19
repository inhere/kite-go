package shell

import (
	"strings"

	"github.com/gookit/goutil/strutil"
)

// generate bash script contents
func (sg *XenvScriptGenerator) generateBashScripts() string {
	// 添加全局环境, PATH, 别名
	var sb strings.Builder
	sg.addCommonForLinuxShell(&sb)

	return strutil.Replaces(BashHookTemplate, map[string]string{
		"{{HooksDir}}":   sg.cfg.ShellHooksDir,
		"{{EnvAliases}}": sb.String(),
	})
}

// GenerateBashHook generates the bash shell hook script

// BashHookTemplate 生成 Bash Hook 的模板
//
// Usage, .bashrc or .bash_profile 新增：
//   eval "$(kite xenv shell --type bash)"
var BashHookTemplate = `# kite xenv bash hook
# This script enables xenv to work in bash shells
# Start to set up xenv in the current shell

setup_xenv() {
    # Mark hook enabled
    export XENV_HOOK_SHELL=bash
    # Set up the xenv shims directory in PATH
    local xenv_shims_dir="${XENV_ROOT:-$HOME/.xenv}/shims"

    # Add shims directory to PATH if it's not already there
    case ":$PATH:" in
        *":$xenv_shims_dir:"*) ;;
        *) export PATH="$xenv_shims_dir:$PATH" ;;
    esac

{{EnvAliases}}

    # Define the xenv function to activate tools
    xenv() {
        local command="$1"
        shift

        case "$command" in
            use)
                # Implementation for switching tool versions
                command kenv use "$@"
                ;;
            unuse)
                # Implementation for unusing tool versions
                command kenv unuse "$@"
                ;;
            shell)
                # Output the shell commands needed to set up xenv
                command kenv shell bash
                ;;
            *)
                # For other commands, just pass through to xenv
                command kenv "$command" "$@"
                ;;
        esac
    }

    # Auto-initialize xenv if needed
    if [ -f "$HOME/.xenvrc" ] && [ -z "$XENV_AUTO_INITIALIZED" ]; then
        source "$HOME/.xenvrc"
        export XENV_AUTO_INITIALIZED=1
    fi

	# Enable command completion for xenv
	if command -v complete >/dev/null 2>&1; then
		complete -W "use unuse add list help" xenv
	fi

    # Load custom hooks script files
	# 使用 glob 获取匹配的文件, 加载所有匹配的脚本
	hook_files=({{HooksDir}}/*.sh)
	for file in "${hook_files[@]}"; do
		if [[ -f "$file" ]] && [[ -r "$file" ]]; then
			source "$file"
		fi
	done
}

# Call setup function to initialize xenv
setup_xenv
`
