package quickjump

import (
	"runtime"
	"strings"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/jsonutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/slog"
	"github.com/inhere/kite-go/pkg/common"
)

const (
	StatusErr = iota
	StatusAdd // new add
	StatusRpt // repeat add
)

const (
	ShellBash = "bash"
	ShellZsh  = "zsh"
	ShellFish = "fish" // TODO
)

// ShellTplMap shell templates
var ShellTplMap = map[string]string{
	ShellBash: JumpBashTpl,
	ShellZsh:  JumpZshTpl,
}

// IsSupported check shell name is supported
func IsSupported(name string) bool {
	_, ok := ShellTplMap[name]
	return ok
}

// IsStatusAdd check
func IsStatusAdd(st int) bool {
	return st == StatusAdd
}

// QuickJump struct
type QuickJump struct {
	*Metadata `json:"-"`
	common.PathResolver
	init bool

	// DataDir is data dir for save metadata files.
	DataDir string `json:"data_dir"`
	// CheckExist check path is exists
	CheckExist bool `json:"check_exist"`
	// SlashPath if true, will replace the path separator to slash
	//
	// - Useful on Windows
	SlashPath bool `json:"slash_path"`
	// NamedPaths pre-define named paths
	NamedPaths map[string]string `json:"named_paths"`
}

// NewQuickJump new quick jump instance
func NewQuickJump() *QuickJump {
	return &QuickJump{
		Metadata:   NewMetadata(),
		SlashPath:  true,
		CheckExist: true,
		PathResolver: common.PathResolver{
			PathResolve: fsutil.ResolvePath,
		},
	}
}

// Init load data from file
func (j *QuickJump) Init() error {
	if j.init {
		return nil
	}

	j.init = true
	j.checkExist = j.CheckExist
	j.slashPath = j.SlashPath
	j.changedHook = func() {
		slog.ErrorT(j.saveToFile())
	}
	return j.load()
}

func (j *QuickJump) load() (err error) {
	// check data dir
	if j.DataDir == "" {
		return nil
	}

	// load data from file
	dFile := j.Datafile()

	if fsutil.IsFile(dFile) {
		err = jsonutil.ReadFile(dFile, j.Metadata)
		if err != nil {
			return err
		}

		j.AddNamedPaths(j.NamedPaths)
	} else {
		// init file
		j.AddNamedPaths(j.NamedPaths)
		err = jsonutil.WritePretty(dFile, j.Metadata)
	}
	return
}

// Save data to disk
func (j *QuickJump) Save() error {
	return j.saveToFile()
}

func (j *QuickJump) saveToFile() error {
	j.Metadata.Datetime = timex.Now().Datetime()

	// save to file
	return jsonutil.WritePretty(j.Datafile(), j.Metadata)
}

// Datafile get data file path
func (j *QuickJump) Datafile() string {
	return j.PathResolve(j.DataDir) + "/" + j.datafileName()
}

func (j *QuickJump) datafileName() string {
	return "quick-jump." + runtime.GOOS + ".json"
}

// Metadata struct for quick jump
type Metadata struct {
	// datafile string

	Datetime string `json:"datetime"`
	LastPath string `json:"last_path"`
	PrevPath string `json:"prev_path"`

	NamedPaths map[string]string `json:"named_paths"`
	Histories  map[string]string `json:"histories"`

	slashPath   bool
	checkExist  bool
	changedHook func()
}

// NewMetadata new metadata instance
func NewMetadata() *Metadata {
	return &Metadata{
		NamedPaths: make(map[string]string),
		Histories:  make(map[string]string),
	}
}

// FormatKeywords handle
func (m *Metadata) FormatKeywords(keywords []string) []string {
	ln := len(keywords)
	if ln == 0 {
		return keywords
	}

	// from bash/zsh: "php order" => [php, order]
	if ln == 1 {
		return strutil.Split(strings.Trim(keywords[0], ` '"`), " ")
	}
	return keywords
}

// CheckOrMatch path by input keywords(name,path,...)
func (m *Metadata) CheckOrMatch(keywords []string) string {
	keywords = m.FormatKeywords(keywords)

	ln := len(keywords)
	if ln == 0 {
		return ""
	}

	if ln == 1 {
		first := keywords[0]
		if first == "." {
			return sysutil.Workdir()
		}

		// return prev path
		if len(first) == 0 || first == "-" {
			return m.PrevPath
		}

		if fsutil.IsDir(first) {
			return fsutil.Realpath(first)
		}

		if dirPath, ok := m.NamedPaths[first]; ok {
			return dirPath
		}

		ss := m.SearchByString(first, 1, false)
		if len(ss) > 0 {
			return ss[0]
		}
		return ""
	}

	// search by multiple keywords
	ss := m.Search(keywords, 1, false)
	if len(ss) > 0 {
		return ss[0]
	}
	return ""
}

// SearchNamed named paths
func (m *Metadata) SearchNamed(keywords []string, limit int, withName bool) []string {
	var paths []string
	noKw := len(keywords) == 0

	for name, dirPath := range m.NamedPaths {
		if noKw || arrutil.StringsHas(keywords, name) || strutil.ContainsAll(dirPath, keywords) {
			if withName {
				paths = append(paths, name+":"+dirPath)
			} else {
				paths = append(paths, dirPath)
			}

			if limit > 0 && len(paths) >= limit {
				break
			}
		}
	}

	return paths
}

// SearchHistory named paths
func (m *Metadata) SearchHistory(keywords []string, limit int) []string {
	var paths []string
	noKw := len(keywords) == 0

	for _, dirPath := range m.Histories {
		if noKw || strutil.ContainsAll(dirPath, keywords) {
			paths = append(paths, dirPath)
			if limit > 0 && len(paths) >= limit {
				return paths
			}
		}
	}

	return paths
}

// SearchByString search named paths and history paths
func (m *Metadata) SearchByString(keywords string, limit int, withName bool) []string {
	return m.Search(strutil.Split(keywords, " "), limit, withName)
}

// Search named paths and history paths
func (m *Metadata) Search(keywords []string, limit int, withName bool) []string {
	var paths []string
	noKw := len(keywords) == 0

	for name, dirPath := range m.NamedPaths {
		if noKw || arrutil.StringsHas(keywords, name) || strutil.SimpleMatch(dirPath, keywords) {
			if withName {
				paths = append(paths, name+":"+dirPath)
			} else {
				paths = append(paths, dirPath)
			}

			if limit > 0 && len(paths) >= limit {
				return paths
			}
		}
	}

	for _, dirPath := range m.Histories {
		if noKw || strutil.SimpleMatch(dirPath, keywords) {
			paths = append(paths, dirPath)
			if limit > 0 && len(paths) >= limit {
				return paths
			}
		}
	}

	return paths
}

// AddNamed add named path
func (m *Metadata) AddNamed(name string, pathStr string) (ok bool) {
	st := m.addNamed(name, pathStr)
	if st == StatusAdd {
		m.fireHook()
	}
	return st != StatusErr
}

// addNamed add named path
func (m *Metadata) addNamed(name string, pathStr string) (st int) {
	if len(name) > 0 && len(pathStr) > 0 {
		if _, ok := m.NamedPaths[name]; ok {
			return StatusRpt
		}

		pathStr = fsutil.Realpath(pathStr)
		if m.checkExist && !fsutil.IsDir(pathStr) {
			return StatusErr
		}

		if m.slashPath {
			pathStr = fsutil.SlashPath(pathStr)
		}

		// add
		m.NamedPaths[name] = pathStr
		return StatusAdd
	}
	return StatusErr
}

// AddNamedPaths add named path
func (m *Metadata) AddNamedPaths(pathMap map[string]string) (ok bool) {
	var added bool
	for name, path := range pathMap {
		st := m.addNamed(name, path)
		if st == StatusErr {
			return false
		}
		if IsStatusAdd(st) {
			added = true
		}
	}

	if added {
		m.fireHook()
	}
	return true
}

// AddHistory add history path
func (m *Metadata) AddHistory(dirPath string) (string, bool) {
	if len(dirPath) == 0 {
		return "", false
	}

	dirPath = fsutil.Realpath(dirPath)
	if dirPath == m.LastPath {
		return dirPath, true
	}

	if m.slashPath {
		dirPath = fsutil.SlashPath(dirPath)
	}
	m.LastPath, m.PrevPath = dirPath, m.LastPath

	// check is new add
	hisKey := strutil.Md5(dirPath)
	if _, ok := m.Histories[hisKey]; !ok {
		m.Histories[hisKey] = dirPath
		m.fireHook()
	}

	return dirPath, true
}

// CleanHistories refresh histories, remove invalid paths
func (m *Metadata) CleanHistories() (ss []string) {
	for k, v := range m.Histories {
		if !fsutil.IsDir(v) {
			delete(m.Histories, k)
			ss = append(ss, v)
		}
	}

	if len(ss) > 0 {
		m.fireHook()
	}
	return ss
}

// CleanHistories refresh histories, remove invalid paths
func (m *Metadata) fireHook() {
	if m.changedHook != nil {
		m.changedHook()
	}
}
