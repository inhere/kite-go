package shell

import (
	"fmt"
	"os"
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// XenvScriptGenerator xenv Shell脚本生成器实现
type XenvScriptGenerator struct {
	cfg *models.Configuration
	shell ShellType
}

// NewScriptGenerator creates a new ShellGenerator
func NewScriptGenerator(shellType ShellType, cfg *models.Configuration) *XenvScriptGenerator {
	return &XenvScriptGenerator{cfg: cfg, shell: shellType}
}

// GenHookScripts 生成 Shell Hook 初始化脚本代码
func (sg *XenvScriptGenerator) GenHookScripts() (string, error) {
	switch sg.shell {
	case Bash:
		return sg.generateBashScripts(), nil
	case Zsh:
		return sg.generateZshScripts(), nil
	case Pwsh:
		return sg.generatePwshScripts(), nil
	default:
		return sg.generateCmdScripts(), nil
	}
}

// GenSetEnvs 批量生成环境变量设置脚本代码
func (sg *XenvScriptGenerator) GenSetEnvs(envs map[string]string) string {
	var sb strings.Builder
	for name, value := range envs {
		sb.WriteString(sg.GenSetEnv(name, value))
	}
	return sb.String()
}

// GenSetEnv 生成环境变量设置脚本代码
func (sg *XenvScriptGenerator) GenSetEnv(name, value string) string {
	name = strings.ToUpper(name)
	switch sg.shell {
	case Bash, Zsh:
		return fmt.Sprintf("export %s='%s'\n", name, value)
	case Pwsh:
		return fmt.Sprintf("$Env:%s = '%s'\n", name, value)
	default:
		return fmt.Sprintf("os.setenv('%s', '%s')\n", name, value)
	}
}

// GenUnsetEnvs 批量生成环境变量删除脚本代码
func (sg *XenvScriptGenerator) GenUnsetEnvs(names []string) string {
	var sb strings.Builder
	for _, name := range names {
		sb.WriteString(sg.GenUnsetEnv(name))
	}
	return sb.String()
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
		return fmt.Sprintf("$Env:PATH = \"%s;$Env:PATH\"\n", path)
	default:
		return fmt.Sprintf("os.setenv('PATH', '%s;%%PATH%%')\n", path)
	}
}

// GenAddPaths 一次添加多个到 PATH 的脚本代码
func (sg *XenvScriptGenerator) GenAddPaths(paths []string) string {
	newPath := JoinPaths(paths)
	switch sg.shell {
	case Bash, Zsh:
		return fmt.Sprintf("export PATH=%s:$PATH\n", newPath)
	case Pwsh:
		// pwsh "" 支持变量插值和表达式求值
		return fmt.Sprintf("$Env:PATH = \"%s;$Env:PATH\"\n", newPath)
	default:
		return fmt.Sprintf("os.setenv('PATH', '%s;%%PATH%%')\n", newPath)
	}
}

// GenRemovePaths 生成批量删除 PATH 的脚本代码
func (sg *XenvScriptGenerator) GenRemovePaths(paths []string) (script string, notFounds []string) {
	var newPaths []string
	osPathList := SplitPath(os.Getenv("PATH"))

	_, newPaths, notFounds = DiffRemovePaths(osPathList, paths)
	if len(newPaths) > 0 {
		script = sg.GenSetPath(newPaths)
	}
	return
}

// GenSetPath 设置 PATH 脚本代码
func (sg *XenvScriptGenerator) GenSetPath(paths []string) string {
	newPath := JoinPaths(paths)
	switch sg.shell {
	case Bash, Zsh:
		return fmt.Sprintf("export PATH=%s\n", newPath)
	case Pwsh:
		return fmt.Sprintf("$Env:PATH = '%s'\n", newPath)
	default:
		return fmt.Sprintf("os.setenv('PATH', '%s')\n", newPath)
	}
}

func (sg *XenvScriptGenerator) addCommonForLinuxShell(sb *strings.Builder) {
	// 添加全局环境变量
	if len(sg.cfg.GlobalEnv) > 0 {
		sb.WriteString("  # Add global ENV variables from kite xenv\n")
		maputil.EachTypedMap(sg.cfg.GlobalEnv, func(key, value string) {
			sb.WriteString(fmt.Sprintf("  export %s=%s\n", strings.ToUpper(key), value))
		})
	}

	// 添加全局PATH条目
	if len(sg.cfg.GlobalPaths) > 0 {
		sb.WriteString("  # Add global PATH from kite xenv\n")
		var fmtPaths []string
		for _, path := range sg.cfg.GlobalPaths {
			// TODO Windows git-bash 将盘符 D:/ 转换成 /d/
			fmtPaths = append(fmtPaths, path)
		}
		sb.WriteString(fmt.Sprintf("  export PATH=%s:$PATH\n", strings.Join(fmtPaths, ":")))
	}

	// 添加全局别名
	if len(sg.cfg.ShellAliases) > 0 {
		sb.WriteString("  # Add global aliases from kite xenv\n")
		maputil.EachTypedMap(sg.cfg.ShellAliases, func(key, value string) {
			sb.WriteString(fmt.Sprintf("  alias %s='%s'\n", key, value))
		})
	}

}
