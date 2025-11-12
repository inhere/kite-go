#!/usr/bin/env bash
#
# kite xenv bash hook
# This script enables xenv to work in bash shells
#
# Usage, .bashrc or .bash_profile add:
#   eval "$(kite xenv shell --type bash)"
#
# Start to set up xenv in the current shell

# 重写 cd 内建命令来实现监听
cd() {
    builtin cd "$@" && {
        if command -v kite >/dev/null 2>&1; then
            kite xenv init-direnv >/dev/null 2>&1
        fi
    }
}

# Helper function to evaluate xenv command results
invoke_xenv_result() {
    local result="$1"
    local exit_code="$2"

    if [ "$exit_code" -eq 0 ]; then
        if [ -n "$result" ]; then
            # 检查结果是否包含 '--Expression--' 分隔符
            if [[ "$result" == *"--Expression--"* ]]; then
                # 使用 '--Expression--' 分割内容
                local msg_part="$(echo "$result" | cut -d'-' -f1)"
                local expr_part="$(echo "$result" | cut -d'-' -f4-)"

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

setup_xenv() {
    # Mark hook enabled
    export XENV_HOOK_SHELL=bash
    export XENV_SESSION_ID="{{SessionId}}"
    # Set up the xenv shims directory in PATH
    local xenv_shims_dir="${XENV_ROOT:-$HOME/.xenv}/shims"

    # Add shims directory to PATH if it's not already there
    case ":$PATH:" in
        *":$xenv_shims_dir:"*) ;;
        *) export PATH="$xenv_shims_dir:$PATH" ;;
    esac

#{{EnvAliases}}

    # Define the xenv function to activate tools
    xenv() {
        local command="$1"
        shift

        case "$command" in
            use|unuse|env|path)
                # 对于这些命令，获取结果并评估
                local result="$(kite xenv "$command" "$@")"
                local exit_code=$?
                invoke_xenv_result "$result" $exit_code
                ;;
            set|unset)
                # 对于环境变量设置/取消设置命令
                local result="$(kite xenv env "$command" "$@")"
                local exit_code=$?
                invoke_xenv_result "$result" $exit_code
                ;;
            *)
                # For other commands, just pass through to xenv
                command kite xenv "$command" "$@"
                ;;
        esac
    }

    # fire xenv hooks to kite, use for generate code to exec TODO
    local result_init = "$(kite xenv shell-init-hook --type bash)"
    local exit_code=$?
    invoke_xenv_result "$result_init" $exit_code

    # Auto-initialize xenv if needed
    if [ -f "$HOME/.xenvrc" ] && [ -z "$XENV_AUTO_INITIALIZED" ]; then
        source "$HOME/.xenvrc"
        export XENV_AUTO_INITIALIZED=1
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

# Enable command completion for xenv
if command -v complete >/dev/null 2>&1; then
    complete -W "use unuse env set unset path list help" xenv
fi
