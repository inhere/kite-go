package aiservice

import (
	"fmt"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/service/aiservice/aiclaude"
)

type AIService struct {
	Config
	inited bool
	cfgFile string
	// 提供者别名映射
	providerAliases map[string]string
}

func New() *AIService {
	return &AIService{}
}

// Init 加载配置并初始化
func (s *AIService) Init() (*AIService, error) {
	if s.inited {
		return s, nil
	}

	s.cfgFile = app.App().ConfigPath("ai/ai-config.yaml")
	cfg := config.NewGeneric("ai-config", config.WithTagName("yaml"))
	cfg.AddDriver(yaml.Driver)

	if err := cfg.LoadFiles(s.cfgFile); err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", s.cfgFile, err)
	}

	if err := cfg.Decode(&s.Config); err != nil {
		return nil, fmt.Errorf("failed to decode AI config: %w", err)
	}

	cfg.ClearAll()
	err := s.Config.Init()
	return s, err
}

// SetClaudeCodeParam 参数
type SetClaudeCodeParam struct {
	Provider string
	KeyName string // api key name
	Shell string
	Write bool
	Show bool
}

// SetClaudeCode 设置 cc 的 API url 和令牌信息
func (s *AIService) SetClaudeCode(opts SetClaudeCodeParam) error {
	runCfg, err1 := aiclaude.ReadUserConfig()
	if err1 != nil {
		return fmt.Errorf("failed to read config: %w", err1)
	}
	if opts.Show {
		ccolor.Magentaln("📄  User Claude configuration:")
		show.MList(runCfg)
		return nil
	}

	useName := opts.Provider
	if useName == "" {
		return fmt.Errorf("used provider name is required, by --use <name>")
	}
	if !s.IsCCProvider(useName) {
		return fmt.Errorf("invalid provider %q, allowed values: %s", useName, strings.Join(s.ProviderNames(), ","))
	}

	// Get provider configuration
	provider, err := s.CCProviderConfig(useName)
	if err != nil {
		return fmt.Errorf("failed to get provider config: %w", err)
	}
	envs := provider.GetEnvMaps(opts.KeyName)

	// Handle --write flag
	if opts.Write {
		runCfg.Env = envs
		if err = runCfg.Save(); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}

		ccolor.Infof("Configuration written to %s, please open new terminal for user.\n", aiclaude.UserConfigFile())
		return nil
	}

	// Generate output based on shell type
	s.printCCShellEnv(provider.Name, opts.Shell, envs)
	return nil
}

func (s *AIService) printCCShellEnv(name, shell string, envs map[string]string) {
	var sb strutil.Builder
	sb.Writef("# Claude Code Configuration (use %s)", name)

	isPwsh := shell == "pwsh" || shell == "powershell"
	for k, v := range envs {
		if isPwsh {
			sb.Writef("\n$env:%s=%q", k, v)
		} else {
			sb.Writef("\n%s=%q", k, v)
		}
	}

	if isPwsh {
		sb.Writef(`
echo "Claude code ENV settings updated! (use %s)"

# 📌 Active in pwsh shell
# kite ai cc set --shell pwsh --use %s | Out-String | Invoke-Expression
`, name, name)
	} else {
		sb.Writef(`
echo "Claude code ENV settings updated! (use %s)"

# 📌 Active in bash, zsh shell
# eval "$(kite ai cc set --shell $SHELL --use %s)"
`, name, name)
	}

	fmt.Print(sb.String())
}

func (s *AIService) ShowConfig() {
	show.ATitle("Current AI Config:")
	dump.NoLoc(s.Config)
}
