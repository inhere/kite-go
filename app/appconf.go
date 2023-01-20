package app

import (
	"os"
	"strings"

	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite/internal/appconst"
)

// some special chars
var (
	PathAliasPrefix uint8 = '$'

	OSPathSepChar = uint8(os.PathSeparator)
	OSPathSepStr  = string(os.PathSeparator)
)

// Info for kite app
type Info struct {
	Branch    string
	Version   string
	Revision  string
	GoVersion string
	BuildDate string
	PublishAt string
	UpdatedAt string
}

// Config struct of app
//
// Gen by:
//
//	kite go gen st -s @c -t json --name Conf
type Config struct {
	// the main config file path
	confFile string
	// BaseDir base data dir
	BaseDir string `json:"base_dir"`
	// TmpDir tmp dir
	TmpDir string `json:"tmp_dir"`
	// CacheDir cache dir
	CacheDir string `json:"cache_dir"`
	// ConfigDir config dir
	ConfigDir string `json:"config_dir"`
	// ResourceDir resource dir
	ResourceDir string `json:"resource_dir"`
	// IncludeConfig include config files.
	// default file path relative the config_dir
	IncludeConfig []string `json:"include_config"`
}

// ConfFile config
func (c *Config) ConfFile() string {
	return c.confFile
}

func newDefaultConf() *Config {
	defDataDir := sysutil.ExpandPath(appconst.KiteDefaultDataDir)

	return &Config{
		BaseDir:     defDataDir,
		TmpDir:      defDataDir + "/tmp",
		CacheDir:    defDataDir + "/tmp/cache",
		ConfigDir:   defDataDir + "/config",
		ResourceDir: defDataDir + "/resource",
	}
}

func (c *Config) ensurePaths() {
	if c.BaseDir == "" {
		c.BaseDir = appconst.KiteDefaultDataDir
	}

	// expand base dir. eg "~"
	c.BaseDir = sysutil.ExpandPath(c.BaseDir)

	if c.TmpDir == "" {
		c.TmpDir = c.BaseDir + "/tmp"
	} else if c.TmpDir[0] == PathAliasPrefix {
		c.TmpDir = c.PathResolve(c.TmpDir)
	}

	if c.CacheDir == "" {
		c.CacheDir = c.BaseDir + "/tmp/cache"
	} else if c.CacheDir[0] == PathAliasPrefix {
		c.CacheDir = c.PathResolve(c.CacheDir)
	}

	if c.ConfigDir == "" {
		c.ConfigDir = c.BaseDir + "/config"
	} else if c.ConfigDir[0] == PathAliasPrefix {
		c.ConfigDir = c.PathResolve(c.ConfigDir)
	}

	if c.ResourceDir == "" {
		c.ResourceDir = c.BaseDir + "/resource"
	} else if c.ResourceDir[0] == PathAliasPrefix {
		c.ResourceDir = c.PathResolve(c.ResourceDir)
	}
}

// CfgFile get main config file
func (c *Config) CfgFile() string {
	return c.confFile
}

// Path build and get full path relative the base dir.
func (c *Config) Path(subPaths ...string) string {
	return joinPath(c.BaseDir, subPaths)
}

// TmpPath build and get full path relative the tmp dir.
func (c *Config) TmpPath(subPaths ...string) string {
	return joinPath(c.TmpDir, subPaths)
}

// CachePath build and get full path relative the cache dir.
func (c *Config) CachePath(subPaths ...string) string {
	return joinPath(c.CacheDir, subPaths)
}

// ConfigPath build and get full path relative the config dir.
func (c *Config) ConfigPath(subPaths ...string) string {
	return joinPath(c.ConfigDir, subPaths)
}

// PathResolve resolve path alias. "$base/tmp" => "path/to/base_dir/tmp"
func (c *Config) PathResolve(path string) string {
	if path == "" || path[0] != PathAliasPrefix {
		return path
	}

	var other string
	name := path[1:]
	sepIdx := strings.IndexRune(path, '/')
	if sepIdx > 0 {
		name = path[1:sepIdx]
		other = path[sepIdx:]
	}

	switch name {
	case "base":
		return c.BaseDir + other
	case "tmp", "temp":
		return c.TmpDir + other
	case "cache":
		return c.CacheDir + other
	case "cfg", "config":
		return c.ConfigDir + other
	case "res", "resource":
		return c.ResourceDir + other
	case "user", "home":
		return sysutil.HomeDir() + other
	}
	return path
}

func joinPath(basePath string, subPaths []string) string {
	if len(subPaths) == 0 {
		return basePath
	}
	return basePath + OSPathSepStr + strings.Join(subPaths, OSPathSepStr)
}

// IsAliasPath string
func IsAliasPath(path string) bool {
	return path[0] == PathAliasPrefix
}
