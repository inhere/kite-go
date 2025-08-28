package kiteext

import (
	"fmt"
	"sort"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/jsonutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/slog"
)

// BinPrefix kite ext 命令文件名称前缀 eg: kite-abc -> ext: abc
const BinPrefix = "kite-"

// MetaSchema kite ext 注册信息文件内容结构
type MetaSchema struct {
	Uptime string     `json:"uptime"`
	Exts   []*KiteExt `json:"exts"`
}

// KiteExt kite ext 基本信息
type KiteExt struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	// 别名列表
	Aliases []string `json:"aliases,omitempty"`
	// 命令文件名称, 没有扩展名。eg: kite-abc
	binName string
	// OsPaths 不同系统平台下的文件路径
	//  - key: 系统平台 windows, darwin, linux. value: 扩展文件路径
	OsPaths map[string]string `json:"os_paths"`
	osPath  string

	// Args 默认运行参数
	Args []string `json:"args,omitempty"`
	// 指定默认运行目录
	Workdir  string `json:"workdir,omitempty"`
	Disable  bool   `json:"disable,omitempty"`
	Author   string `json:"author,omitempty"`
	Version  string `json:"version,omitempty"`
	Homepage string `json:"homepage,omitempty"`

	// BeforeSave 保存前执行 hook.
	//  - 返回 true 继续保存，返回 false 停止保存
	BeforeSave func(ext *KiteExt) bool `json:"-"`
}

func NewExt(name, desc string, path ...string) *KiteExt {
	return &KiteExt{
		Name:   name,
		Desc:   desc,
		osPath: arrutil.FirstOr(path),
	}
}

// 检查扩展 name
func (e *KiteExt) checkName() error {
	if e.Name == "" {
		return errorx.E("ext name can't be empty")
	}

	// check name valid
	if !strutil.IsVarName(e.Name) {
		return errorx.Ef("invalid extension name: %s", e.Name)
	}
	return nil
}

// init ext info
func (e *KiteExt) init() {
	if e.Desc == "" {
		e.Desc = "NO DESCRIPTION"
	}

	// save to pathMap
	osName := sysutil.OsName
	if e.OsPaths == nil {
		e.OsPaths = make(map[string]string)
	}
	e.OsPaths[osName] = e.osPath
}

// IsValid check path is valid
func (e *KiteExt) IsValid() bool {
	return fsutil.IsFile(e.OsPath())
}

// OsPath in the current os
func (e *KiteExt) OsPath() string {
	if e.osPath == "" {
		smp := maputil.StrMap(e.OsPaths)
		e.osPath = smp.Default(sysutil.OsName, "NONE")
	}
	return e.osPath
}

// SetPath set file path in the current os
func (e *KiteExt) SetPath(osPath string) {
	if osPath != "" {
		e.osPath = osPath
	}
}

// SetBinName set bin name
func (e *KiteExt) SetBinName(name string) {
	if name != "" {
		e.binName = name
	}
}

// BinName get the bin name, without an ext name.
func (e *KiteExt) BinName() string {
	if e.binName == "" {
		e.binName = BinPrefix + e.Name
	}
	return e.binName
}

// SetAliases set aliases for the extension.
func (e *KiteExt) SetAliases(str string) {
	if str == "" {
		return
	}

	aliases := strutil.SplitTrimmed(str, ",")
	if len(aliases) > 0 {
		e.Aliases = aliases
	}
}

// ExtManager kite cli 扩展管理器实现
type ExtManager struct {
	Disable bool `json:"disable"`
	// Metafile kite ext 注册记录文件
	Metafile string `json:"metafile"`
	// SearchPaths 除了 env PATH 外，额外搜索ext文件的目录
	SearchPaths []string `json:"search_paths"`
	// PathResolver handler. 用于查找 Metafile 文件
	PathResolver func(path string) string

	schema  *MetaSchema
	extMap  map[string]*KiteExt
	aliasMp maputil.Aliases
}

func NewExtManager() *ExtManager {
	return &ExtManager{
		aliasMp: make(maputil.Aliases),
		extMap:  make(map[string]*KiteExt),
	}
}

// Init 初始化ext manager
func (m *ExtManager) Init() error {
	m.Metafile = m.PathResolver(m.Metafile)

	// 加载 metafile 文件
	ms := &MetaSchema{}
	err := jsonutil.DecodeFile(m.Metafile, ms)
	if err != nil {
		return errorx.Rf("extMgr: load metafile error：%w", err)
	}

	m.schema = ms
	for _, ext := range ms.Exts {
		m.extMap[ext.Name] = ext
		m.aliasMp.AddAliases(ext.Name, ext.Aliases)
	}
	return nil
}

// Ext gets ext by name
func (m *ExtManager) Ext(name string) (ext *KiteExt, ok bool) {
	name = m.aliasMp.ResolveAlias(name)
	ext, ok = m.extMap[name]
	if ok {
		ext.OsPath() // init osPath
	}
	return
}

// Exts gets ext list.
func (m *ExtManager) Exts() []*KiteExt { return m.schema.Exts }

// Each every ext information
func (m *ExtManager) Each(fn func(ext *KiteExt)) {
	for _, ext := range m.extMap {
		fn(ext)
	}
}

// Exists checks ext exists.
func (m *ExtManager) Exists(name string) bool {
	_, ok := m.extMap[name]
	return ok
}

// Dumpfile 保存扩展元数据
func (m *ExtManager) Dumpfile() error {
	m.schema.Uptime = timex.Now().Datetime()

	// m.schema.Exts 按 Name 排序
	sort.Slice(m.schema.Exts, func(i, j int) bool {
		return m.schema.Exts[i].Name < m.schema.Exts[j].Name
	})

	return jsonutil.WritePretty(m.Metafile, m.schema)
}

//
// region T: add/update/delete ext
//

// QuickAdd 添加一个扩展
func (m *ExtManager) QuickAdd(name, desc string, path ...string) error {
	return m.Add(&KiteExt{
		Name:   name,
		Desc:   desc,
		osPath: arrutil.FirstOr(path),
	})
}

// Add 添加一个扩展
func (m *ExtManager) Add(ext *KiteExt) error {
	// check ext name
	if err := ext.checkName(); err != nil {
		return err
	}

	// check exist
	if _, ok := m.extMap[ext.Name]; ok {
		return errorx.Rf("kite ext '%s' already exists", ext.Name)
	}
	return m.save(ext)
}

// Update 更新扩展信息
func (m *ExtManager) Update(ext *KiteExt) error {
	// check ext name
	if err := ext.checkName(); err != nil {
		return err
	}
	return m.save(ext)
}

// Delete 删除扩展信息
func (m *ExtManager) Delete(name string) error {
	index := -1
	for i, ext := range m.schema.Exts {
		if ext.Name == name {
			index = i
			break
		}
	}
	if index < 0 {
		return errorx.Rf("delete: kite ext '%s' not exists", name)
	}

	// delete ext from map and slice
	delete(m.extMap, name)
	exts := append(m.schema.Exts[:index], m.schema.Exts[index+1:]...)
	m.schema.Exts = exts

	return m.Dumpfile()
}

// CleanInvalid 删除无效的(不存在的)扩展
func (m *ExtManager) CleanInvalid() error {
	exts := make([]*KiteExt, 0)
	for _, ext := range m.schema.Exts {
		if !ext.IsValid() {
			delete(m.extMap, ext.Name)
			continue
		}
		exts = append(exts, ext)
	}

	m.schema.Exts = exts
	return m.Dumpfile()
}

func (m *ExtManager) save(ext *KiteExt) error {
	// 如果 ext.osPath 为空，则搜索ext文件路径
	if ext.osPath == "" {
		ext.osPath = m.findExtFile(ext.BinName())
		if ext.osPath == "" {
			return errorx.Ef("can't find executable file for ext %q", ext.Name)
		}
		slog.Infof("auto found ext file: %s", ext.osPath)
	} else {
		// check the path exists
		if !ext.IsValid() {
			return errorx.Ef("ext path is not a file: %s", ext.osPath)
		}
	}

	ext.init()

	// call before save hook
	if ext.BeforeSave != nil && !ext.BeforeSave(ext) {
		return nil
	}
	ext.BeforeSave = nil

	// add to metadata
	m.extMap[ext.Name] = ext
	m.schema.Exts = append(m.schema.Exts, ext)
	m.aliasMp.AddAliases(ext.Name, ext.Aliases)
	return m.Dumpfile()
}

// ext.Path 为空时，自动搜索ext文件路径
func (m *ExtManager) findExtFile(binName string) (extFile string) {
	// binName := BinPrefix + extName

	// 先搜索 m.SearchPaths
	if len(m.SearchPaths) > 0 {
		names := []string{binName, binName + ".sh"}
		if strutil.ContainsByte(binName, '.') {
			names = []string{binName}
		} else if sysutil.IsWindows() {
			// TODO use env PATHEXT
			names = []string{binName + ".exe", binName + ".bat", binName + ".cmd", binName + ".sh"}
		}

		slog.Infof("search ext file in %v, names: %v", m.SearchPaths, names)
		extFile = fsutil.FileInDirs(m.SearchPaths, names...)
		if extFile != "" {
			return extFile
		}
	}

	// 再搜索 env PATH
	slog.Infof("try search ext file %q in env PATH", binName)
	extFile, _ = sysutil.FindExecutable(binName)
	return extFile
}

//
// region T: run ext commands
//

// RunCtx data
type RunCtx struct {
	Dry bool
	// Dir 设置运行目录
	Dir string
	// Env 设置环境变量
	Env map[string]string
	// Args 运行参数
	Args []string
}

// Run 运行ext命令
func (m *ExtManager) Run(name string, ctx *RunCtx) error {
	name = m.aliasMp.ResolveAlias(name)
	ext, ok := m.extMap[name]
	if !ok {
		return fmt.Errorf("kite: ext '%s' not found", name)
	}

	cArgs := ext.Args
	cArgs = append(cArgs, ctx.Args...)
	dir := strutil.OrElse(ctx.Dir, ext.Workdir)

	fmt.Printf("--------------------------- Run Ext %s, Args %v ---------------------------\n", name, cArgs)
	cmd := cmdr.NewCmd(ext.OsPath(), cArgs...)

	err := cmd.WorkDirOnNE(dir).WithDryRun(ctx.Dry).AppendEnv(ctx.Env).PrintCmdline2().FlushRun()
	return err
}
