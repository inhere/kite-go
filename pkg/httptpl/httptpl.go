package httptpl

import (
	"fmt"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
)

// Manager HTTP tpl manager
type Manager struct {
	// DataDirs []string `json:"data_dirs"`

	// DefaultExt default config file ext, for auto load domain config
	//
	// Allow: json, json5, yaml
	DefaultExt string `json:"default_ext"`
	// DefaultDir default domain config file dir, for auto load domain config
	//
	// path format: {DefaultDir}/{domainName}-domain.{DefaultExt}
	DefaultDir string `json:"default_dir"`
	// Domains config definitions
	Domains map[string]*DomainConfig `json:"domains"`
	// PathResolver handler
	PathResolver func(path string) string
}

// NewManager instance
func NewManager() *Manager {
	return &Manager{
		DefaultExt:   "json",
		PathResolver: fsutil.ResolvePath,
	}
}

// Init some info. should call it before use
func (m *Manager) Init() error {
	if m.DefaultExt == "" {
		m.DefaultExt = "json"
	}

	for name, dc := range m.Domains {
		dc.Name = name
		dc.PathResolver = m.PathResolver

		// init dc
		if err := dc.Init(); err != nil {
			return err
		}
	}
	return nil
}

// Domain get
func (m *Manager) Domain(name string) (*DomainConfig, error) {
	dc, ok := m.Domains[name]
	if ok {
		return dc, nil
	}

	// auto load domain config from m.DefaultDir
	if dir := m.DefaultDir; dir != "" {
		confFile := fmt.Sprintf("%s/%s-domain.%s", dir, name, m.DefaultExt)

		if fsutil.IsFile(confFile) {
			dc = NewDomainConfig(name, confFile)
			dc.PathResolver = m.PathResolver
			if err := dc.Init(); err != nil {
				return nil, err
			}

			m.Domains[name] = dc
			return dc, nil
		}
	}

	return nil, errorx.Rawf("not found domain config of the %q", name)
}
