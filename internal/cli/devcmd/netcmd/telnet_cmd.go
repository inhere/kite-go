package netcmd

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/x/ccolor"
)

// NewTelnetClientCmd creates a new Telnet client cmd
func NewTelnetClientCmd() *gcli.Command {
	var tcCmdOpts = struct {
		timeout int // 连接超时时间(s) 默认超时5秒
	}{}

	return &gcli.Command{
		Name: "telnet",
		Desc: "start a telnet client",
		Config: func(c *gcli.Command) {
			c.IntOpt(&tcCmdOpts.timeout, "timeout", "t", 5, "connection timeout seconds")
			c.AddArg("host", "telnet server host").WithDefault("127.0.0.1")
			c.AddArg( "port", "telnet server port").WithDefault(23)
		},
		Func: func(c *gcli.Command, args []string) error {
			// 构建服务器地址
			host := c.Arg("host").String()
			port := c.Arg("port").Int()
			addr := fmt.Sprintf("%s:%d", host, port)
			fmt.Printf("Connecting to %s ...\n", addr)

			// 连接到telnet服务器
			conn, err := net.DialTimeout("tcp", addr, time.Duration(tcCmdOpts.timeout)*time.Second)
			if err != nil {
				return fmt.Errorf("failed to connect to server: %w", err)
			}
			defer conn.Close()

			ccolor.Infoln(`Connect successful! Input "quit" OR "exit" to exit.`)
			ccolor.Println("Type commands and press Enter to send them to the server.\n")

			// 同时设置读写超时
			conn.SetDeadline(time.Time{})
			// conn.SetDeadline(time.Now().Add(10 * time.Second))
			var exited atomic.Bool

			// 创建1个goroutine用于双向通信
			// 1. 从服务器读取数据并显示在终端上
			go func() {
				buf := make([]byte, 1024)
				for {
					n, err1 := conn.Read(buf)
					if err1 != nil {
						if err1 != io.EOF && !exited.Load() {
							ccolor.Fprintf(os.Stderr, "<red1>Error</> reading from server: %v\n", err1)
						}
						return
					}
					fmt.Print("S> ", string(buf[:n]))
				}
			}()

			// 2. 从标准输入读取数据并发送到服务器
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "quit" || line == "exit" {
					exited.Store(true)
					ccolor.Magentaln("Exit, Bye!")
					break
				}

				_, err1 := conn.Write([]byte(line + "\n"))
				if err1 != nil {
					return fmt.Errorf("failed to send data to server: %w", err1)
				}
			}

			// if err := scanner.Err(); err != nil {
			// 	return fmt.Errorf("error reading from stdin: %w", err)
			// }
			return nil
		},
	}
}
