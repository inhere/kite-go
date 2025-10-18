package shell

import (
	"fmt"
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
)

func (sg *XenvScriptGenerator) generateCmdScripts() string {
	// clink 通过 os.execute('doskey ll=dir /a $*') 实现别名
	var sb strings.Builder
	maputil.EachTypedMap(sg.cfg.ShellAliases, func(key, value string) {
		sb.WriteString(fmt.Sprintf("doskey %s=%s\n", key, value))
	})

	return strutil.Replaces(CmdLuaHookTemplate, map[string]string{
		"{{HooksDir}}":     sg.cfg.ShellHooksDir,
		"{{EnvAliases}}": sb.String(),
	})
}

// CmdLuaHookTemplate CMD 需要基于 clink lua 脚本实现自定义 hooks
//
// 使用:
//
// 在 C:\Users\{username}\AppData\Local\clink 创建 cmdrc.lua 文件。
// 添加内容：
//
//	load(io.popen('kite xenv shell --type cmd'):read("*a"))()
var CmdLuaHookTemplate = `-- xenv CMD hook
-- This script enables xenv to work in CMD shells

-- Function to set up xenv in the current shell
function Setup-Xenv()
{
    -- Mark hook enabled
    os.setenv("XENV_HOOK_SHELL", "cmd")
    -- Set up the xenv shims directory in PATH
    local xenv_shims_dir = os.getenv("XENV_ROOT") or os.getenv("USERPROFILE") .. "\\.xenv\\shims"

    -- Add shims directory to PATH if it's not already there
    local path = os.getenv("PATH") or ""
    if not string.match(path, xenv_shims_dir) then
        os.setenv("PATH", xenv_shims_dir .. ";" .. path)
    end

    -- Add global shell ENV and aliases
    {{EnvAliases}}

    -- Define the xenv function to activate tools
    function xenv(command)
        if command == "use" then
            -- Implementation for switching tool versions
            os.execute("kenv use " .. table.concat(arg, " "))
        elseif command == "unuse" then
            -- Implementation for unusing tool versions
            os.execute("kenv unuse " .. table.concat(arg, " "))
        elseif command == "shell" then
            -- Output the shell commands needed to set up xenv
            os.execute("kenv shell cmd")
        else
            -- For other commands, just pass through to xenv
            os.execute("kenv " .. command .. " " .. table.concat(arg, " "))
        end
    end
}

-- Call setup function to initialize xenv
Setup-Xenv()
`
