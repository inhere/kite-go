package shell

// GenerateBashHook generates the bash shell hook script
func GenerateBashHook() string {
	return `# xenv bash hook
# This script enables xenv to work in bash shells

# Function to set up xenv in the current shell
setup_xenv() {
    # Mark hook enabled
    export XENV_HOOK_ENABLE=true
    # Set up the xenv shims directory in PATH
    local xenv_shims_dir="${XENV_ROOT:-$HOME/.xenv}/shims"

    # Add shims directory to PATH if it's not already there
    case ":$PATH:" in
        *":$xenv_shims_dir:"*) ;;
        *) export PATH="$xenv_shims_dir:$PATH" ;;
    esac

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
}

# Call setup function to initialize xenv
setup_xenv
`
}

// GenerateZshHook generates the zsh shell hook script
func GenerateZshHook() string {
	return `# xenv zsh hook
# This script enables xenv to work in zsh shells

# Function to set up xenv in the current shell
setup_xenv() {
    # Mark hook enabled
    export XENV_HOOK_ENABLE=true
    # Set up the xenv shims directory in PATH
    local xenv_shims_dir="${XENV_ROOT:-$HOME/.xenv}/shims"

    # Add shims directory to PATH if it's not already there
    case ":$PATH:" in
        *":$xenv_shims_dir:"*) ;;
        *) export PATH="$xenv_shims_dir:$PATH" ;;
    esac

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
                command kenv shell zsh
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
}

# Call setup function to initialize xenv
setup_xenv
`
}

// GeneratePwshHook generates the PowerShell shell hook script
func GeneratePwshHook() string {
	return `# xenv PowerShell hook
# This script enables xenv to work in PowerShell shells

# Function to set up xenv in the current shell
function Setup-Xenv {
    # Mark hook enabled
    $env:XENV_HOOK_ENABLE = "true"
    # Set up the xenv shims directory in PATH
    $xenvShimsDir = if ($env:XENV_ROOT) { "$env:XENV_ROOT\shims" } else { "$HOME\.xenv\shims" }

    # Add shims directory to PATH if it's not already there
    if ($env:PATH -notlike "*$xenvShimsDir*") {
        $env:PATH = "$xenvShimsDir;$env:PATH"
    }

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
}

# Call setup function to initialize xenv
Setup-Xenv
`
}
