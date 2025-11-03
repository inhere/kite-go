package shell

import (
	"strings"

	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// generateBashScripts generates the zsh shell hook script
func (sg *XenvScriptGenerator) generateZshScripts(ps *models.GenInitScriptParams) string {
	// 添加全局环境, PATH, 别名
	var sb strings.Builder
	sg.addCommonForLinuxShell(&sb, ps)

	return strutil.Replaces(ZshHookTemplate, map[string]string{
		"{{HooksDir}}": ps.ShellHooksDir,
		"{{SessionId}}": models.SessionID(),
		"{{EnvAliases}}": sb.String(),
	})
}

// ZshHookTemplate 生成 zsh hook 的模板
//
// Usage, .zshrc or .zsh_profile 新增：
//
//	eval "$(kite xenv shell --type bash)"
var ZshHookTemplate = `# xenv zsh hook
# This script enables xenv to work in zsh shells

# Helper function to evaluate xenv command results
eval_xenv_result() {
    local result="$1"
    local exit_code="$2"

    if [ "$exit_code" -eq 0 ]; then
        if [ -n "$result" ]; then
            # 检查结果是否包含 '--Expression--' 分隔符
            if [[ "$result" == *"--Expression--"* ]]; then
                # 使用 '--Expression--' 分割内容
                local msg_part="${result%%--Expression--*}"
                local expr_part="${result##*--Expression--}"

                # 后面部分当做代码执行
                if [ -n "$expr_part" ]; then
                    eval "$expr_part"
                fi
                # 前面部分直接输出
                if [ -n "$msg_part" ]; then
                    echo "$msg_part"
                    # echo "$result"  # 输出完整结果用于调试
                fi
            else
                # 否则直接输出内容
                echo "$result"
            fi
        fi
    else
        echo "$result" >&2
        return $exit_code
    fi
}

# Function to set up xenv in the current shell
setup_xenv() {
    # Mark hook enabled
    export XENV_HOOK_SHELL=zsh
    export XENV_SESSION_ID="{{SessionId}}"
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
            use|unuse|env|path)
                # 对于这些命令，获取结果并评估
                local result="$(kite xenv "$command" "$@")"
                local exit_code=$?
                eval_xenv_result "$result" $exit_code
                ;;
            set|unset)
                # 对于环境变量设置/取消设置命令
                local result="$(kite xenv env "$command" "$@")"
                local exit_code=$?
                eval_xenv_result "$result" $exit_code
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
    if command -v compctl >/dev/null 2>&1; then
        compctl -k "use unuse env set unset path list help" xenv
    fi

	# Load custom hooks script files
	hook_files={{HooksDir}}/*.sh
	for file in "${hook_files[@]}"; do
		if [[ -f "$file" && -r "$file" ]]; then
			source "$file"
		fi
	done
}

# Call setup function to initialize xenv
setup_xenv
`
