package shell

import (
	"fmt"
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// XenvScriptGenerator xenv Shell脚本生成器实现
type XenvScriptGenerator struct {
	cfg *models.Configuration
}

// NewScriptGenerator creates a new ShellGenerator
func NewScriptGenerator(cfg *models.Configuration) *XenvScriptGenerator {
	return &XenvScriptGenerator{cfg: cfg}
}

// GenerateScripts 生成Shell hook脚本代码
func (sg *XenvScriptGenerator) GenerateScripts(shellType string) (string, error) {
	shellType = strings.ToLower(shellType)

	switch shellType {
	case "bash":
		return sg.generateBashScripts(), nil
	case "zsh":
		return sg.generateZshScripts(), nil
	case "pwsh", "powershell":
		return sg.generatePwshScripts(), nil
	case "cmd", "clink":
		return sg.generateCmdScripts(), nil
	default:
		return "", fmt.Errorf("unsupported shell type: %s (use bash, zsh, or pwsh)", shellType)
	}
}

func (sg *XenvScriptGenerator) addCommonForLinuxShell(sb *strings.Builder) {
	// 添加全局环境变量
	maputil.EachTypedMap(sg.cfg.GlobalEnv, func(key, value string) {
		sb.WriteString(fmt.Sprintf("export %s=%s\n", strings.ToUpper(key), value))
	})

	// 添加全局PATH条目
	if len(sg.cfg.GlobalPaths) > 0 {
		var fmtPaths []string
		for _, path := range sg.cfg.GlobalPaths {
			// TODO Windows git-bash 将盘符 D:/ 转换成 /d/
			fmtPaths = append(fmtPaths, path)
		}
		sb.WriteString(fmt.Sprintf("export PATH=%s:$PATH\n", strings.Join(fmtPaths, ":")))
	}

	// 添加全局别名
	maputil.EachTypedMap(sg.cfg.ShellAliases, func(key, value string) {
		sb.WriteString(fmt.Sprintf("alias %s='%s'\n", key, value))
	})

}
