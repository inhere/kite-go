package xenvcom

const (
	// HookShellEnvName 当前会话的HOOK SHELL环境变量名称
	HookShellEnvName = "XENV_HOOK_SHELL"
	// SessIdEnvName 当前会话ID环境变量名称
	SessIdEnvName = "XENV_SESSION_ID"
	// XenvDebugEnvName debug环境变量名称
	XenvDebugEnvName = "XENV_DEBUG_MODE"
)

const (
	// GlobalStateFile global state file path
	GlobalStateFile = "~/.config/xenv/global.toml"
	// LocalStateFile local state file path
	LocalStateFile = ".xenv.toml"
	// SessionStateDir 当前SHELL会话状态文件目录 eg: ~/.xenv/session/<session_id>.json
	SessionStateDir = "~/.xenv/session"
)

const (
	InstalledMetaFile = "~/.xenv/tools.local.json"
)

// 升级匹配级别 see config#AllowUpMatch
const (
	UpMatchNone uint8 = iota
	UpMatchOne
	UpMatchTwo
	UpMatchAll uint8 = 9
)
