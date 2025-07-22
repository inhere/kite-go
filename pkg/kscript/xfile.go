package kscript

type XFile struct {
	// 离工作目录最近的一个 kitefile 文件
	kf *KiteFile

	// 允许的文件名列表 default: 'kitefile', '.kitefile', 'kitefile.yml', 'kitefile.yaml'
	Filenames []string `json:"filenames"`
	FileExts []string `json:"file_exts"`
	// 除了在当前目录中搜索，还搜索以下目录。找到一个 kitefile 文件后，停止搜索
	ExtraFind string `json:"extra_find"`
	// 向上搜索目录最大深度，默认为 5. 找到的都作为 kf 父级
	MaxFindDepth int `json:"max_find_depth"`
}

func NewXFile() *XFile {
	return &XFile{
		Filenames: []string{"kitefile", ".kitefile"},
		FileExts: DefaultDefineExts,
		MaxFindDepth: 5,
	}
}

type KiteFile struct {
	parent *KiteFile
	// 当前 .kitefile 文件路径
	filePath string

	Name string // 脚本名称，help 时显示
	Desc string // 描述 help 时显示

	// 包含/引用的公共定义文件
	Include []string
}
