package models

import "time"

// EnvironmentVariable 代表系统中的环境变量，具有名称、值、作用域（全局/会话）属性
type EnvironmentVariable struct {
	Name      string    `json:"name"`      // 环境变量名称
	Value     string    `json:"value"`     // 环境变量值
	Scope     string    `json:"scope"`     // 作用域: "global" 或 "session"
	IsActive  bool      `json:"is_active"` // 是否当前激活
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}