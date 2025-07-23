package kscript

type XFileManager struct {
	// 离工作目录最近的一个 kitefile 文件
	top *XFile

	// 允许的文件名列表 default: 'kitefile', '.kitefile', 'kitefile.yml', 'kitefile.yaml'
	Filenames []string `json:"filenames"`
	FileExts []string `json:"file_exts"`
	// 除了在当前目录中搜索，还搜索以下目录。找到一个 kitefile 文件后，停止搜索
	ExtraFind string `json:"extra_find"`
	// 向上搜索目录最大深度，默认为 5. 找到的都作为 kf 父级
	MaxFindDepth int `json:"max_find_depth"`
}

func NewXFile() *XFileManager {
	return &XFileManager{
		Filenames: []string{"kitefile", ".kitefile"},
		FileExts:  DefaultDefineExts,
		MaxFindDepth: 5,
	}
}

// XAction struct
type XAction struct {
	Name    string   `json:"name"`
	Desc    string   `json:"desc"`
	User    string   `json:"user"`
	Workdir string   `json:"workdir"`
	Cmds    []string `json:"cmds"`
}

type XFile struct {
	parent *XFile
	// 当前 .kitefile 文件路径
	filePath string
	// mode: top_cfg, setting
	Mode string `json:"__xfile_mode"`

	// 脚本名称，help 时显示
	Name string `json:"name"`
	// 描述信息 help 时显示
	Desc    string `json:"desc"`
	Version string `json:"version"`
	Author  string `json:"author"`
	// Homepage message
	Homepage string `json:"homepage"`
	// 包含/引用的公共定义文件
	Include []string
	// Env setting for run
	Env map[string]string
	// 运行时上下文
	Context struct {
		Workdir string `json:"workdir"`
	} `json:"context"`
	// DefaultAction on run
	DefaultAction string `json:"default_action"`

	// Actions 定义
	Actions map[string]*XAction `json:"actions"`
}
