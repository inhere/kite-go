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
	// 为空时,将当前目录路径hash值作为sessionID
	if sessionID == "" {
		// 用时间会导致产生很多文件
		sessionID = time.Now().Format("20060102_150405")
		// TIP: 用目录也有问题，会按首次打开时生成。。。后续又会切换目录
		// workdir := sysutil.Workdir()
		// // 取出目录名的前4个字符作为前缀
		// prefix := strutil.Substr(fsutil.Name(workdir), 0, 4)
		// sessionID = prefix + "_" + strutil.ShortMd5(workdir)
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

