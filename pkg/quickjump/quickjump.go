package quickjump

import (
	"runtime"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/jsonutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/slog"
	"github.com/inhere/kite-go/pkg/common"
)

// QuickJump struct
type QuickJump struct {
	*Metadata `json:"-"`
	common.PathResolver
	init bool

	// DataDir is data dir for save metadata files.
	DataDir string `json:"data_dir"`
	// CheckExist check path is exists
	CheckExist bool `json:"check_exist"`
	// NamedPaths pre-define named paths
	NamedPaths map[string]string `json:"named_paths"`
}

// NewQuickJump new quick jump instance
func NewQuickJump() *QuickJump {
	return &QuickJump{
		Metadata: NewMetadata(),
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
	j.changedHook = func() {
		slog.ErrorT(j.saveToFile())
	}
	return j.load()
}

func (j *QuickJump) load() (err error) {
	dFile := j.Datafile()

	// load data
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
	datafile string

	Datetime string `json:"datetime"`
	LastPath string `json:"last_path"`
	PrevPath string `json:"prev_path"`

	NamedPaths map[string]string `json:"named_paths"`
	Histories  map[string]string `json:"histories"`

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

// Match path by input name
func (m *Metadata) Match(name string) string {
	if name == "." {
		return sysutil.Workdir()
	}

	// return prev path
	if len(name) == 0 || name == "-" {
		return m.PrevPath
	}

	if fsutil.IsDir(name) {
		return name
	}

	if dirPath, ok := m.NamedPaths[name]; ok {
		return dirPath
	}

	for _, dirPath := range m.Histories {
		if strutil.IContains(dirPath, name) {
			return dirPath
		}
	}
	return ""
}

// SearchByString search named paths and history paths
func (m *Metadata) SearchByString(keywords string, limit int) []string {
	return m.Search(strutil.Split(keywords, " "), limit)
}

// Search named paths and history paths
func (m *Metadata) Search(keywords []string, limit int) []string {
	var paths []string
	for _, dirPath := range m.NamedPaths {
		if strutil.ContainsAll(dirPath, keywords) {
			paths = append(paths, dirPath)
			if limit > 0 && len(paths) >= limit {
				return paths
			}
		}
	}

	for _, dirPath := range m.Histories {
		if strutil.ContainsAll(dirPath, keywords) {
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
	if ok = m.addNamed(name, pathStr); ok {
		m.fireHook()
	}
	return ok
}

// addNamed add named path
func (m *Metadata) addNamed(name string, pathStr string) (ok bool) {
	if len(name) > 0 && len(pathStr) > 0 {
		pathStr = fsutil.Realpath(pathStr)
		if m.checkExist && !fsutil.IsDir(pathStr) {
			return false
		}

		ok = true
		m.NamedPaths[name] = pathStr
	}
	return
}

// AddNamedPaths add named path
func (m *Metadata) AddNamedPaths(pathMap map[string]string) (ok bool) {
	for name, path := range pathMap {
		if ok = m.addNamed(name, path); !ok {
			return false
		}
	}

	m.fireHook()
	return
}

// AddHistory add history path
func (m *Metadata) AddHistory(pathStr string) bool {
	if len(pathStr) == 0 {
		return false
	}

	pathStr = fsutil.Realpath(pathStr)
	if pathStr == m.LastPath {
		return false
	}

	m.LastPath, m.PrevPath = pathStr, m.LastPath
	m.Histories[strutil.Md5(pathStr)] = pathStr
	m.fireHook()

	return true
}

// CleanHistories refresh histories, remove invalid paths
func (m *Metadata) CleanHistories() (n int) {
	for k, v := range m.Histories {
		if !fsutil.IsDir(v) {
			delete(m.Histories, k)
			n++
		}
	}

	if n > 0 {
		m.fireHook()
	}
	return n
}

// CleanHistories refresh histories, remove invalid paths
func (m *Metadata) fireHook() {
	if m.changedHook != nil {
		m.changedHook()
	}
}
