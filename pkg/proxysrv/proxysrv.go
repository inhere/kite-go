// Package proxysrv provides a proxy server for Kite.
package proxysrv

// Config struct
type Config struct {
	Port       int
	Rules      []string
	RuleFiles  []string
	RuleDirs   []string
	GlobalVars map[string]string
}

type ConfigFn func(c *Config)

// ProxySrv struct
//
// refer:
//
//   - https://github.com/avwo/whistle
type ProxySrv struct {
	*Config
	Rules []*ProxyRule
}

// NewProxySrv create a new ProxySrv instance
func NewProxySrv(fns ...ConfigFn) *ProxySrv {
	ps := &ProxySrv{
		Config: &Config{
			Port: 8090,
		},
	}

	for _, fn := range fns {
		fn(ps.Config)
	}
	return ps
}

func (s *ProxySrv) Start() error {
	return nil
}

// ProxyRule definition
type ProxyRule struct {
	Name  string // id name
	Group string
	Index int // index in rules
	Sort  int // sort number

	// Protocol name of rule. http, https, ws, wss, file, ...
	Protocol string

	// Pattern from host or path
	Pattern string
	// Operate dist host or path
	Operate string
	// Config for the rule
	Config map[string]string
}
