package shell

import (
	"fmt"
	"os"
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/inhere/kite-go/pkg/util"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// XenvScriptGenerator xenv Shell脚本生成器实现
type XenvScriptGenerator struct {
	// cfg *models.Configuration
	shell ShType
}

// NewScriptGenerator creates a new ShellGenerator
func NewScriptGenerator(shellType ShType) *XenvScriptGenerator {
	return &XenvScriptGenerator{shell: shellType}
}

// endregion
// region Generate Init Scripts
//

// GenHookScripts 生成 Shell Hook 初始化脚本代码
func (sg *XenvScriptGenerator) GenHookScripts(ps *models.GenInitScriptParams) (string, error) {
	switch sg.shell {
	case Bash:
		return sg.generateBashScripts(ps), nil
	case Zsh:
		return sg.generateZshScripts(ps), nil
	case Pwsh:
		return sg.generatePwshScripts(ps), nil
	default:
		return sg.generateCmdScripts(ps), nil
	}
}

// InstallToProfile 安装 Shell Hook 脚本到配置文件(eg: .bashrc, .zshrc)
func (sg *XenvScriptGenerator) InstallToProfile(pwshProfile string) error {
	switch sg.shell {
	case Bash:
	case Zsh:
	case Pwsh:
		// echo $PROFILE.CurrentUserAllHosts
		// v1: path-to-users\Documents\WindowsPowerShell\profile.ps1
		// v7: path-to-users\Documents\PowerShell\profile.ps1
	default:
		// C:\Users\{username}\AppData\Local\clink\ 创建 profile.lua
	}
	return nil
}

// installScriptsToProfile 安装 Shell Hook 脚本到配置文件(eg: .bashrc, .zshrc)
//  - 检查文件是否存在，如果不存在则创建一个
//  - 检查文件内容是否包含 xenv 脚本，如果存在则返回
//  - 如果不存在内容则添加到文件的末尾
func (sg *XenvScriptGenerator) installScriptsToProfile(script, profile string) error {

	return nil
}

// endregion
// region Generate Snippets
//

// GenSetEnvs 批量生成环境变量设置脚本代码
func (sg *XenvScriptGenerator) GenSetEnvs(envs map[string]string) string {
	var ss []string
	for name, value := range envs {
		ss = append(ss, sg.GenSetEnv(name, value))
	}
	return strings.Join(ss, "\n")
}

// GenUnsetEnvs 批量生成环境变量删除脚本代码
func (sg *XenvScriptGenerator) GenUnsetEnvs(names []string) string {
	var ss []string
	for _, name := range names {
		ss = append(ss, sg.GenUnsetEnv(name))
	}
	return strings.Join(ss, "\n")
}

// GenSetEnv 生成环境变量设置脚本代码
func (sg *XenvScriptGenerator) GenSetEnv(name, value string) string {
	name = strings.ToUpper(name)
	switch sg.shell {
	case Bash, Zsh:
		return fmt.Sprintf("export %s='%s'\n", name, value)
	case Pwsh:
		return fmt.Sprintf("$Env:%s='%s';\n", name, value)
	default:
		return fmt.Sprintf("os.setenv('%s', '%s')\n\n", name, value)
	}
}

// GenUnsetEnv 删除环境变量的脚本代码
func (sg *XenvScriptGenerator) GenUnsetEnv(name string) string {
	name = strings.ToUpper(name)
	switch sg.shell {
	case Bash, Zsh:
		return fmt.Sprintf("unset %s\n", name)
	case Pwsh:
		return fmt.Sprintf("Remove-Item Env:%s\n", name)
	default:
		return fmt.Sprintf("os.unsetenv('%s')\n", name)
	}
}

// GenAddPath 添加 PATH 脚本代码（添加到 PATH 的第一个位置）
func (sg *XenvScriptGenerator) GenAddPath(path string) string {
	switch sg.shell {
	case Bash, Zsh:
		return fmt.Sprintf("export PATH=%s:$PATH\n", path)
	case Pwsh:
		return fmt.Sprintf("$Env:PATH=\"%s;$Env:PATH\"\n", path)
	default:
		return fmt.Sprintf("os.setenv('PATH', '%s;%%PATH%%')\n", path)
	}
}

// GenAddPaths 一次添加多个到 PATH 的脚本代码
func (sg *XenvScriptGenerator) GenAddPaths(paths []string) string {
	newPath := util.JoinPaths(paths)
	switch sg.shell {
	case Bash, Zsh:
		return fmt.Sprintf("export PATH=%s:$PATH\n", newPath)
	case Pwsh:
		// pwsh "" 支持变量插值和表达式求值
		return fmt.Sprintf("$Env:PATH=\"%s;$Env:PATH\"\n", newPath)
	default:
		return fmt.Sprintf("os.setenv('PATH', '%s;%%PATH%%')\n", newPath)
	}
}

// GenSetPath 设置 PATH 脚本代码
func (sg *XenvScriptGenerator) GenSetPath(paths []string) string {
	newPath := util.JoinPaths(paths)
	switch sg.shell {
	case Bash, Zsh:
		return fmt.Sprintf("export PATH='%s'\n", newPath)
	case Pwsh:
		return fmt.Sprintf("$Env:PATH='%s';\n", newPath)
	default:
		return fmt.Sprintf("os.setenv('PATH', '%s')\n\n", newPath)
	}
}

// GenRemovePaths 生成批量删除 PATH 的脚本代码
func (sg *XenvScriptGenerator) GenRemovePaths(paths []string) (script string, notFounds []string) {
	var newPaths []string
	osPathList := util.SplitPath(os.Getenv("PATH"))

	_, newPaths, notFounds = DiffRemovePaths(osPathList, paths)
	if len(newPaths) > 0 {
		script = sg.GenSetPath(newPaths)
	}
	return
}

// GenRemThenAddPaths 生成批量删除后再新增的 PATH 的脚本代码
func (sg *XenvScriptGenerator) GenRemThenAddPaths(rmPaths, addPaths []string) (script string) {
	var newPaths []string
	osPathList := util.SplitPath(os.Getenv("PATH"))

	_, newPaths, _ = DiffRemovePaths(osPathList, rmPaths)
	if len(newPaths) > 0 {
		if len(addPaths) > 0 {
			newPaths = append(addPaths, newPaths...)
		}
		script = sg.GenSetPath(newPaths)
	} else if len(addPaths) > 0 {
		script = sg.GenAddPaths(addPaths)
	}
	return
}


// endregion
// region Helper methods
//

func (sg *XenvScriptGenerator) addCommonForLinuxShell(sb *strings.Builder, ps *models.GenInitScriptParams) {
	// 添加全局环境变量
	if len(ps.Envs) > 0 {
		sb.WriteString("  # Add global ENV variables from kite xenv\n")
		maputil.EachTypedMap(ps.Envs, func(key, value string) {
			sb.WriteString(fmt.Sprintf("  export %s=%s\n", strings.ToUpper(key), value))
		})
	}

	// 添加全局PATH条目
	if len(ps.Paths) > 0 {
		sb.WriteString("  # Add global PATH from kite xenv\n")
		var fmtPaths []string
		for _, path := range ps.Paths {
			// TODO Windows git-bash 将盘符 D:/ 转换成 /d/
			fmtPaths = append(fmtPaths, path)
		}
		sb.WriteString(fmt.Sprintf("  export PATH=%s:$PATH\n", strings.Join(fmtPaths, ":")))
	}

	// 添加全局别名
	if len(ps.ShellAliases) > 0 {
		sb.WriteString("  # Add global aliases from kite xenv\n")
		maputil.EachTypedMap(ps.ShellAliases, func(key, value string) {
			sb.WriteString(fmt.Sprintf("  alias %s='%s'\n", key, value))
		})
	}

}
