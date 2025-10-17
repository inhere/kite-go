package models

import "time"

// PathEntry 代表添加到PATH环境变量中的路径条目，具有路径值和优先级属性
type PathEntry struct {
	Path      string    `json:"path"`      // 添加到PATH的路径
	Priority  int       `json:"priority"`  // 优先级，数值越小优先级越高
	Scope     string    `json:"scope"`     // 作用域: "global" 或 "session"
	IsActive  bool      `json:"is_active"` // 是否当前激活
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}