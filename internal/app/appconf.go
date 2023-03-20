package app

import (
	"os"
	"strings"

	"github.com/gookit/goutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite-go/internal/appconst"
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
	// the dotenv file path
	dotenvFile string
	// the main config file path
	confFile string
	// BaseDir base data dir
	BaseDir string `json:"base_dir" default:"${KITE_BASE_DIR}"`
	// TmpDir tmp dir
	TmpDir string `json:"tmp_dir"`
	// DataDir tmp dir
	DataDir string `json:"data_dir"`
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

func (c *Config) ensurePaths() {
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

// IsAliasPath alias
func (c *Config) IsAliasPath(path string) bool {
	return IsAliasPath(path)
}

// ResolvePath alias
func (c *Config) ResolvePath(path string) string {
	return c.PathResolve(path)
}

// PathResolve resolve path alias. "$base/tmp" => "path/to/base_dir/tmp"
func (c *Config) PathResolve(path string) string {
	if path == "" || path[0] != PathAliasPrefix {
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
	return path
}

// PathByName get
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
