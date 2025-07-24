package app

import (
	"os"
	"strings"

	"github.com/gookit/goutil"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite-go/internal/appconst"
)

// some special chars
var (
	// PathMarkPrefixes allowed path name prefix mark.
	PathMarkPrefixes = []byte{'$', '#', '@'}

	OSPathSepChar = uint8(os.PathSeparator)
	OSPathSepStr  = string(os.PathSeparator)
)

// Config struct of app
//
// Gen by:
//
//	kite go gen st -s @c -t json --name Conf
type Config struct {
	// the dotenv file path
	dotenvFile string
	// the main config file path
	confFile string

	// BaseDir base data dir. default: ~/.kite-go
	BaseDir string `json:"base_dir" default:"${KITE_BASE_DIR}"`
	// TmpDir tmp dir. default in BaseDir/tmp
	TmpDir string `json:"tmp_dir"`
	// DataDir tmp dir. default in BaseDir/data
	DataDir string `json:"data_dir"`
	// CacheDir cache dir. default in BaseDir/tmp/cache
	CacheDir string `json:"cache_dir"`
	// ConfigDir config dir. default in BaseDir/config
	ConfigDir string `json:"config_dir"`
	// ResourceDir resource dir. default in BaseDir/resource
	ResourceDir string `json:"resource_dir"`
	// IncludeConfig include config files.
	// default file path relative the config_dir
	IncludeConfig []string `json:"include_config"`
}

func (c *Config) ensurePaths() {
	// c.prefixes = PathMarkPrefixes
	if c.BaseDir == "" {
		c.BaseDir = appconst.KiteDefaultBaseDir
	}

	// expand base dir. eg "~"
	c.BaseDir = sysutil.ExpandPath(c.BaseDir)
	c.TmpDir = goutil.OrValue(c.TmpDir == "", c.BaseDir+"/tmp", c.PathResolve(c.TmpDir))
	c.DataDir = goutil.OrValue(c.DataDir == "", c.BaseDir+"/data", c.PathResolve(c.DataDir))

	c.CacheDir = goutil.OrValue(c.CacheDir == "", c.BaseDir+"/tmp/cache", c.PathResolve(c.CacheDir))
	c.ConfigDir = goutil.OrValue(c.ConfigDir == "", c.BaseDir+"/config", c.PathResolve(c.ConfigDir))

	c.ResourceDir = goutil.OrValue(c.ResourceDir == "", c.BaseDir+"/resource", c.PathResolve(c.ResourceDir))
}

// Path build and get full path relative the base dir.
func (c *Config) Path(subPaths ...string) string {
	return c.PathBuild(subPaths...)
}

// PathBuild build and get full path relative the base dir.
func (c *Config) PathBuild(subPaths ...string) string {
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

var sysAliases = []string{
	// kite dir
	"data",
	"base", "root",
	"tmp", "temp",
	"cache", "caches",
	"cfg", "conf", "config",
	"res", "resource",
	// os dir
	"user", "home",
	"workdir", "pwd",
}

// IsPathAlias check is path alias name, without prefix char.
func (c *Config) IsPathAlias(name string) bool {
	return arrutil.InStrings(name, sysAliases)
}

// HasAliasMark check has alias mark
func (c *Config) HasAliasMark(path string) bool {
	return HasAliasMark(path)
}

// ResolvePath alias name.
func (c *Config) ResolvePath(path string) string {
	return c.PathResolve(path)
}

// PathResolve resolve path alias. "$base/tmp" => "path/to/base_dir/tmp"
func (c *Config) PathResolve(path string) string {
	if !HasAliasMark(path) {
		return path
	}

	var other string
	name := path[1:]
	sepIdx := strings.IndexByte(path, '/')
	if sepIdx > 0 {
		name = path[1:sepIdx]
		other = path[sepIdx:]
	}

	if prefix := c.PathByName(name); len(prefix) > 0 {
		return prefix + other
	}
	return path[1:] // remove prefix char
}

// PathByName get path by alias name
func (c *Config) PathByName(name string) string {
	switch name {
	case "data":
		return c.DataDir
	case "base", "root":
		return c.BaseDir
	case "tmp", "temp":
		return c.TmpDir
	case "cache", "caches":
		return c.CacheDir
	case "cfg", "conf", "config":
		return c.ConfigDir
	case "res", "resource":
		return c.ResourceDir
	case "user", "home":
		return sysutil.HomeDir()
	case "workdir", "pwd":
		return sysutil.Workdir()
	}
	return ""
}

// ConfFile get main config file
func (c *Config) ConfFile() string {
	return c.confFile
}

// DotenvFile path
func (c *Config) DotenvFile() string {
	return c.dotenvFile
}

// SetDotenvFile path
func (c *Config) SetDotenvFile(dotenvFile string) {
	c.dotenvFile = dotenvFile
}

// PathsMap collects
func (c *Config) PathsMap() map[string]string {
	return map[string]string{
		"base":   c.BaseDir,
		"config": c.ConfigDir,
		"cache":  c.CacheDir,
		"data":   c.DataDir,
		"tmp":    c.TmpDir,
		"res":    c.ResourceDir,
	}
}

func joinPath(basePath string, subPaths []string) string {
	if len(subPaths) == 0 {
		return basePath
	}
	return basePath + OSPathSepStr + strings.Join(subPaths, OSPathSepStr)
}

// HasAliasMark check path string, start one of PathMarkPrefixes
func HasAliasMark(path string) bool {
	return len(path) > 0 && arrutil.In(path[0], PathMarkPrefixes)
}
