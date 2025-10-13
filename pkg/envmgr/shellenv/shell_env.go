package shellenv

// ShellType shell类型枚举
type ShellType string

const (
	ShellZsh        ShellType = "zsh"
	ShellBash       ShellType = "bash"
	ShellFish       ShellType = "fish"
	ShellPowerShell ShellType = "pwsh"
	ShellCmd        ShellType = "cmd"
)

// ShellGenerator Shell脚本生成器接口
type ShellGenerator interface {
	// GenerateScript 生成shell脚本
	GenerateScript(shellType ShellType, state *ActiveState, config *ShellEnvConfig) (string, error)

	// GenerateKtenvFunction 生成ktenv函数
	GenerateKtenvFunction(shellType ShellType) (string, error)

	// GenerateEnvVars 生成环境变量设置
	GenerateEnvVars(shellType ShellType, envs map[string]string) (string, error)

	// GeneratePathUpdate 生成PATH更新脚本
	GeneratePathUpdate(shellType ShellType, paths []string, operation PathOperation) (string, error)
}
