package xenvcom

import (
	"fmt"
	"os"
	"time"

	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/x/ccolor"
)

var sessionID = os.Getenv(SessIdEnvName)

// SessionID 获取当前会话ID
func SessionID() string {
	if sessionID == "" {
		sessionID = time.Now().Format("20060102_150405")
	}
	return sessionID
}

// SessionFile 获取当前会话状态文件
func SessionFile() string {
	return SessionStateDir + "/" + SessionID() + ".json"
}

// SetSessionID 设置当前会话ID (用于测试)
func SetSessionID(id string) {
	sessionID = id
}

// DebugMode debug mode flag
var DebugMode = envutil.GetBool(XenvDebugEnvName, false)

// Debugf prints debug messages
func Debugf(format string, args ...any) {
	if DebugMode {
		ccolor.Printf("<cyan>DEBUG</>: "+format, args...)
	}
}

func Debugln(args ...any) {
	if DebugMode {
		ccolor.Println("<cyan>DEBUG</>: ", fmt.Sprint(args...))
	}
}

