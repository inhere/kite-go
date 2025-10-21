package httptpl

import (
	"path/filepath"
	"strings"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/pkg/util/bizutil"
)

// EnvsMap type
type EnvsMap map[string]maputil.Data

// Merge other env map
func (m EnvsMap) Merge(s EnvsMap) {
	for name, data := range s {
		if md, ok := m[name]; ok {
			md.Load(data)
		} else {
			m[name] = data
		}
	}
}

// LoadEnvsByFile handle
func LoadEnvsByFile(envFile string) (EnvsMap, error) {
	ec := bizutil.NewConfig()
	if err := ec.LoadFiles(envFile); err != nil {
		return nil, errorx.Stacked(err)
	}

	se := make(EnvsMap)
	for name, val := range ec.Data() {
		if mp, ok := val.(map[string]any); ok {
			se[name] = mp
		}
	}
	return se, nil
}

// DomainConfig struct.
// name: domain, topic, collection
type DomainConfig struct {
	// Name domain name
	Name string `json:"name"`
	// ConfigFile domain config file path
	ConfigFile string `json:"config_file"`
	// PathResolver handler
	PathResolver func(path string) string `json:"-"`

	// TplDir path for request templates.
	//
	// default is: {ConfigFile Dir}/{domainName}
	TplDir string `json:"tpl_dir"`
	// TplExt allowed template file ext list
	TplExt []string `json:"tpl_ext"`

	// Vars global variables
	Vars maputil.Data `json:"vars"`
	// Header global headers for send request
	Header map[string]string `json:"header"`
	// DefaultEnv name
	DefaultEnv string `json:"default_env"`
	// EnvFile The ide http-client.env.json file path
	EnvFile string `json:"env_file"`
	// Envs definition
	Envs EnvsMap `json:"envs"`

	// templates map
	//
	// TIP: default group name is Name
	//
	// Example:
	//
	// 	{
	//		// default group: load from template dir
	//		"jenkins": Templates{set: [Template, ...]}
	//		// load from jenkins.http file
	//		"jenkins.http": Templates{set: [Template, ...]}
	// 	}
	tsMap map[string]*Templates
}

// NewDomainConfig instance
func NewDomainConfig(name, configFile string) *DomainConfig {
	return &DomainConfig{
		Name:       name,
		ConfigFile: configFile,
		// default sets
		// PathResolver: fsutil.ResolvePath,
	}
}

// Init some info, load domain config
func (d *DomainConfig) Init() error {
	d.TplExt = []string{"json", "json5", "yaml", "ini"}
	d.tsMap = make(map[string]*Templates)
	if d.Envs == nil {
		d.Envs = make(EnvsMap)
	}

	if len(d.ConfigFile) > 0 {
		d.ConfigFile = d.PathResolver(d.ConfigFile)
	}

	if len(d.TplDir) == 0 {
		if len(d.ConfigFile) > 0 {
			d.TplDir = filepath.Dir(d.ConfigFile) + "/" + d.Name
		}
	} else {
		d.TplDir = d.PathResolver(d.TplDir)
	}

	return d.LoadConfig()
}

// LoadConfig file and env file
func (d *DomainConfig) LoadConfig() error {
	cfg := bizutil.NewConfig()

	if len(d.ConfigFile) > 0 {
		if err := cfg.LoadFiles(d.ConfigFile); err != nil {
			return errorx.Withf(err, "load domain %q config fail, file: %s", d.Name, d.ConfigFile)
		}

		if err := cfg.Decode(d); err != nil {
			return errorx.Stacked(err)
		}
	}

	if len(d.EnvFile) > 0 {
		d.EnvFile = d.PathResolver(d.EnvFile)
		if len(d.TplDir) == 0 {
			d.TplDir = filepath.Dir(d.EnvFile)
		}

		se, err := LoadEnvsByFile(d.EnvFile)
		if err != nil {
			return err
		}
		d.Envs.Merge(se)
	}

	return nil
}

// Load templates
func (d *DomainConfig) Load() error {

	// TODO
	return nil
}

// SendTemplate request by name, will find from tsMap[DomainConfig.Name]
// func (d *DomainConfig) SendTemplate(name, env string) (*Template, error) {
// 	t := d.Template(name)
// 	if t == nil {
// 		return nil, errorx.Rawf("not found template by name %q", name)
// 	}
//
// 	vars := d.Vars
// 	if svs, ok := d.Envs[env]; ok {
// 		vars.Load(svs)
// 	}
//
// 	err := t.Send(vars)
// 	return t, err
// }

// Template get by file name, will find from tsMap[DomainConfig.Name]
func (d *DomainConfig) Template(name string) *Template {
	group := d.Name
	if strutil.ContainsByte(name, ':') {
		group, name = strutil.MustCut(name, ":")
	}

	if t, err := d.Lookup(group, name); err == nil {
		return t
	}
	return nil
}

// Lookup get by name and group
func (d *DomainConfig) Lookup(group, name string) (*Template, error) {
	ts, ok := d.Templates(group)
	if !ok {
		return nil, errorx.Rawf("templates group %q not found", group)
	}
	return ts.Lookup(name)
}

// Templates get by name.
//
// Usage:
//
//	d.Templates("") // get default
//	d.Templates("jenkins.http")
func (d *DomainConfig) Templates(group string) (*Templates, bool) {
	if len(group) == 0 {
		group = d.Name
	}
	if ts, ok := d.tsMap[group]; ok {
		return ts, true
	}

	ts := NewTemplates(group)

	if group == d.Name {
		ts.path = d.TplDir
		ts.exts = d.TplExt
		d.tsMap[group] = ts
	} else {
		// try load from hc-file
		ts.typ = TypeHttpClient
		fName := strutil.OrCond(strings.HasSuffix(group, HcFileExt), group, group+HcFileExt)
		err := ts.FromHCFile(d.TplDir + "/" + fName)
		if err != nil {
			return nil, false
		}
	}

	return ts, true
}

// DefaultTemplates get by name
func (d *DomainConfig) DefaultTemplates() *Templates {
	return d.tsMap[d.Name]
}

// BuildVars by env name and file
func (d *DomainConfig) BuildVars(envName, envFile string) (maputil.Data, error) {
	var vs maputil.Data
	// append global vars
	if len(d.Vars) > 0 {
		vs = d.Vars
	}

	// load env vars
	if envName != "" {
		if envFile != "" {
			em, err := LoadEnvsByFile(d.PathResolver(envFile))
			if err != nil {
				return nil, err
			}
			vs.Load(em[envName])
		} else if ev, ok := d.LookupVars(envName); ok {
			vs.Load(ev)
		}
	}

	return vs, nil
}

// LookupVars get by env name
func (d *DomainConfig) LookupVars(env string) (maputil.Data, bool) {
	vs, ok := d.Envs[env]
	return vs, ok
}

// EnvVars get by env name
func (d *DomainConfig) EnvVars(env string) maputil.Data {
	return d.Envs[env]
}
