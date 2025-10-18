package models


// SDKInfo SDK信息
type SDKInfo struct {
	Name      string `json:"name"`      // SDK名称
	Version   string `json:"version"`   // 版本号
	IsActive  bool   `json:"is_active"` // 是否激活
	Path      string `json:"path"`      // 安装路径
	Installed bool   `json:"installed"` // 是否已安装
}
