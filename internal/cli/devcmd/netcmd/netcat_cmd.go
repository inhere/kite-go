package netcmd

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/strutil"
)

// NetcatOptions netcat命令选项
type NetcatOptions struct {
	listen   bool
	udp      bool
	port     int
	host     string
	timeout  int
	verbose  bool
}

// NewNetcatCmd 实现linux nc(netcat) 工具命令
//
// Linux nc:
//   # 监听TCP端口
//   nc -l 8080
//   # 监听UDP端口
//   nc -ul 53
//
//   # 连接到远程TCP端口
//   nc example.com 80
//   # 连接到远程UDP端口
//   nc -u example.com 53
//
func NewNetcatCmd() *gcli.Command {
	var ncOpts = NetcatOptions{}

	return &gcli.Command{
		Name:    "nc",
		Desc: "Netcat utility for network connections",
		Aliases: []string{"netcat"},
		Config: func(c *gcli.Command) {
			c.BoolOpt(&ncOpts.listen, "listen", "l", false, "Listen mode, start server")
			c.BoolOpt(&ncOpts.udp, "udp", "u", false, "UDP mode, default is TCP")
			c.IntOpt(&ncOpts.timeout, "timeout", "w", 0, "Timeout in seconds")
			c.BoolOpt(&ncOpts.verbose, "verbose", "v", false, "Verbose output")
		},
		Func: func(c *gcli.Command, args []string) error {
			if ncOpts.listen {
				// 监听模式
				if len(args) != 1 {
					return fmt.Errorf("listen mode requires exactly one argument (port)")
				}

				var err error
				ncOpts.port, err = strutil.ToInt(args[0])
				if err != nil || ncOpts.port < 1 || ncOpts.port > 65535 {
					return fmt.Errorf("invalid port: %s", args[0])
				}

				return listenMode(&ncOpts)
			} else {
				// 连接模式
				if len(args) != 2 {
					return fmt.Errorf("connect mode requires exactly two arguments (host port)")
				}

				ncOpts.host = args[0]
				var err error
				ncOpts.port, err = strutil.ToInt(args[1])
				if err != nil || ncOpts.port < 1 || ncOpts.port > 65535 {
					return fmt.Errorf("invalid port: %s", args[1])
				}

				return connectMode(&ncOpts)
			}
		},
		Examples: `
  {$binWithCmd} -l 8080                 listen on TCP port 8080
  {$binWithCmd} -ul 53                  listen on UDP port 53
  {$binWithCmd} example.com 80          connect to example.com on TCP port 80
  {$binWithCmd} -u example.com 53       connect to example.com on UDP port 53
`,
	}
}

// listenMode 监听模式 - 启动 tcp/udp server
func listenMode(opts *NetcatOptions) error {
	protocol := "tcp"
	if opts.udp {
		protocol = "udp"
	}

	address := fmt.Sprintf(":%d", opts.port)

	if opts.verbose {
		fmt.Printf("Listening on %s port %d (%s)\n", strings.ToUpper(protocol), opts.port, protocol)
	}

	listener, err := net.Listen(protocol, address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", address, err)
	}
	defer listener.Close()

	if opts.verbose {
		fmt.Printf("Connection received\n")
	}

	conn, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept connection: %v", err)
	}
	defer conn.Close()

	// 处理连接
	return handleConnection(conn, opts)
}

// connectMode 连接模式
func connectMode(opts *NetcatOptions) error {
	protocol := "tcp"
	if opts.udp {
		protocol = "udp"
	}

	address := fmt.Sprintf("%s:%d", opts.host, opts.port)
	if opts.verbose {
		fmt.Printf("Connecting to %s port %d (%s)\n", opts.host, opts.port, strings.ToUpper(protocol))
	}

	var err error
	var conn net.Conn

	if opts.timeout > 0 {
		conn, err = net.DialTimeout(protocol, address, time.Duration(opts.timeout)*time.Second)
	} else {
		conn, err = net.Dial(protocol, address)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", address, err)
	}
	defer conn.Close()

	if opts.verbose {
		fmt.Printf("Connected to %s\n", address)
	}

	// 处理连接
	return handleConnection(conn, opts)
}

// handleConnection 处理连接的数据传输 创建两个goroutine分别处理输入和输出
func handleConnection(conn net.Conn, opts *NetcatOptions) error {
	// 从连接读取数据并输出到标准输出
	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				if err != io.EOF {
					if opts.verbose {
						fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
					}
				}
				break
			}

			// 输出到标准输出
			// os.Stdout.Write(buffer[:n])
			fmt.Print("S> ", string(buffer[:n]))
		}
	}()

	// 从标准输入读取数据并发送到连接
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data := scanner.Bytes()
		// 添加换行符（如果不存在）
		if len(data) > 0 && data[len(data)-1] != '\n' {
			data = append(data, '\n')
		}

		_, err := conn.Write(data)
		if err != nil {
			if opts.verbose {
				_, _ = fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
			}
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from stdin: %v", err)
	}
	return nil
}
