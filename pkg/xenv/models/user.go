package models

import "time"

// User 代表使用xenv工具的单一用户类型，具有对自身配置和环境的完全控制权
type User struct {
	ID         string    `json:"id"`          // 用户唯一标识符
	ConfigPath string    `json:"config_path"` // 用户配置文件路径
	HomeDir    string    `json:"home_dir"`    // 用户主目录
	ShellType  string    `json:"shell_type"`  // 用户使用的shell类型: "bash", "zsh", "pwsh"
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}