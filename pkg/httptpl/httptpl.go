package httptpl

import (
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

// Init some info
func (m *Manager) Init() error {
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
	if !ok {
		return nil, errorx.Rawf("not found domain config of the %q", name)
	}

	return dc, nil
}
