package aiservice

import (
	"fmt"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/dump"
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
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	if err := cfg.Decode(&s.Config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	cfg.ClearAll()
	return s, s.Config.Init()
}

type SetAuthInfoParam struct {
	Provider string
	Shell string
	Write bool
	Show bool
}

// SetAuthInfo 设置 API url 和令牌信息
func (s *AIService) SetAuthInfo(opts SetAuthInfoParam) error {
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
	if !s.IsProvider(useName) {
		return fmt.Errorf("invalid provider %q, allowed values: %s", useName, strings.Join(s.ProviderNames(), ","))
	}

	// Get provider configuration
	provider, err := s.ProviderConfig(useName)
	if err != nil {
		return fmt.Errorf("failed to get provider config: %w", err)
	}

	// Handle --write flag
	if opts.Write {
		runCfg.Env["ANTHROPIC_BASE_URL"] = provider.BaseURL
		runCfg.Env["ANTHROPIC_AUTH_TOKEN"] = provider.APIKey
		if err = runCfg.Save(); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}

		ccolor.Infof("Configuration written to %s, please open new terminal for user.\n", aiclaude.UserConfigFile())
		return nil
	}

	// Generate output based on shell type
	shellEnv := provider.GenShellEnv(opts.Shell)
	fmt.Println(shellEnv)
	return nil
}

func (s *AIService) ShowConfig() {
	show.ATitle("Current AI Config:")
	dump.NoLoc(s.Config)
}
