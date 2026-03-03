package aiservice

import (
	"fmt"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/service/aiservice/aiclaude"
)

type AIService struct {
	Config
	inited  bool
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
	KeyName  string // api key name
	Model    string // 指定模型名称
	Scope    string // 作用域 user, project
	Shell    string
	Write    bool
	Show     bool
	List     bool
}

// SetClaudeCode 设置 cc 的 API url 和令牌信息
func (s *AIService) SetClaudeCode(opts SetClaudeCodeParam) error {
	if opts.List {
		ccolor.Magentaln("📄  CC providers configuration:")
		show.MList(s.CcProviders)
		return nil
	}

	runCfg, err1 := aiclaude.LoadConfig(opts.Scope)
	if err1 != nil {
		return fmt.Errorf("failed to read config: %w", err1)
	}
	if opts.Show {
		ccolor.Magentaln("📄  User Claude configuration:")
		show.MList(runCfg)
		return nil
	}

	if opts.Shell == "" {
		shellName := sysutil.CurrentShell(true)
		if shellName == "" && sysutil.IsWindows() {
			shellName = "pwsh"
		}
		opts.Shell = shellName
		ccolor.Infof("")
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
	envs := provider.GetEnvMaps(opts.KeyName, opts.Model)
	if ati := s.AgentTools["claude_code"]; ati != nil {
		envs = maputil.MergeStrMap(envs, ati.Envs)
	}

	// Handle --write flag
	if opts.Write {
		runCfg.Env = envs
		if err = runCfg.Save(); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}

		ccolor.Infof("Configuration written to %s, please reopen cc for use.\n", runCfg.ConfigFile())
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
	_ = sb.WriteByte('\n')

	sb.Writef(`
echo "Claude code ENV variables updated! (provider:%s. can use 'env | grep ANT' check)"
`, name)
	sb.Writeln(`echo "📢 请确认 claude settings.json 中的 env 没有任何模型设置!"`)
	sb.Writeln(`echo "    不然当前 ENV 的设置可能不会生效"`)
	if isPwsh {
		sb.Writef(`
# 📌 Active in pwsh shell
# kite ai cc set --shell pwsh --use %s | Out-String | Invoke-Expression
`, name)
	} else {
		sb.Writef(`
# 📌 Active in bash, zsh shell
# eval "$(kite ai cc set --shell $SHELL --use %s)"
`, name)
	}

	fmt.Print(sb.String())
}

func (s *AIService) ShowConfig(keyName string) {
	show.ATitle("Current AI Config: " + s.cfgFile)
	var data any
	switch strings.ToLower(keyName) {
	case "ps", "providers":
		data = s.Config.Providers
	case "pa", "p-alias", "p-aliaes":
		data = s.Config.ProviderAliases
	case "ma", "m-alias", "m-aliaes":
		data = s.Config.ModelAliases
	case "cc-ps", "cc-providers":
		data = s.Config.CcProviders
	default:
		if keyName != "" {
			ccolor.Warnf("Invalid key name: %s", keyName)
			return
		}
		data = s.Config
	}

	dump.NoLoc(data)
}
