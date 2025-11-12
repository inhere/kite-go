package netcmd

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/x/ccolor"
)

type nmapCmdOpts struct {
	ports    string
	timeout  int
	threads  int
	host     string
	protocol string
	scanType string
	verbose  bool
	// 常用的服务版本探测
	serviceDetect bool
	// internal fields
	portList []int
}

// NewNMapCmd 实现nmap扫描 tcp 工具命令
func NewNMapCmd() *gcli.Command {
	var nmapOpts = nmapCmdOpts{}

	return &gcli.Command{
		Name:    "nmap",
		Desc: "nmap port scan tool (support tcp/udp)",
		Aliases: []string{"scan"},
		Config: func(c *gcli.Command) {
			c.StrOpt(&nmapOpts.ports, "ports", "p", "80,443,22,21,23,25,53,110,143,993,995,3306,5432,6379,27017",
				`ports to scan, comma separated. port must be between 1 and 65535.
Range: 1-100
Both: 1000-2000,8080,4000-5000
`)
			c.IntOpt(&nmapOpts.timeout, "timeout", "t", 3, "connection timeout in milliseconds")
			c.IntOpt(&nmapOpts.threads, "threads", "c,T", 100, "number of concurrent threads")
			c.StrOpt(&nmapOpts.protocol, "protocol", "P", "tcp", "scan protocol: tcp or udp")
			c.StrOpt(&nmapOpts.scanType, "scan-type", "s", "connect", "scan type: connect, null")
			c.BoolOpt(&nmapOpts.serviceDetect, "service-detect", "sv", false, "detect service versions")
			c.BoolOpt(&nmapOpts.verbose, "verbose", "v", false, "verbose output")
			c.AddArg("host", "target host/IP to scan", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			nmapOpts.host = c.Arg("host").String()

			// 解析端口列表
			if err := parsePorts(&nmapOpts); err != nil {
				return err
			}

			// 验证协议
			if nmapOpts.protocol != "tcp" && nmapOpts.protocol != "udp" {
				return fmt.Errorf("unsupported protocol: %s (only tcp/udp supported)", nmapOpts.protocol)
			}

			// 执行端口扫描
			return scanPorts(&nmapOpts)
		},
		Examples: `
  {$binWithCmd} 192.168.1.1
  {$binWithCmd} -p 80,443,8080 example.com
  {$binWithCmd} -p 1-1000 -t 5 -T 50 192.168.1.1
  {$binWithCmd} -P udp -p 53,123 8.8.8.8
  {$binWithCmd} --sv -p 22,80,443 example.com
`,
	}
}

// parsePorts 解析端口参数
func parsePorts(opts *nmapCmdOpts) error {
	ports := strings.Split(opts.ports, ",")

	for _, portStr := range ports {
		portStr = strings.TrimSpace(portStr)

		// 处理端口范围，例如 "1-100"
		if strings.Contains(portStr, "-") {
			rangeParts := strings.Split(portStr, "-")
			if len(rangeParts) != 2 {
				return fmt.Errorf("invalid port range format: %s", portStr)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return fmt.Errorf("invalid start port: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return fmt.Errorf("invalid end port: %s", rangeParts[1])
			}

			if start > end {
				return fmt.Errorf("start port %d cannot be greater than end port %d", start, end)
			}

			for i := start; i <= end; i++ {
				opts.portList = append(opts.portList, i)
			}
		} else {
			// 处理单个端口
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return fmt.Errorf("invalid port: %s", portStr)
			}

			if port < 1 || port > 65535 {
				return fmt.Errorf("port must be between 1 and 65535: %d", port)
			}

			opts.portList = append(opts.portList, port)
		}
	}

	// 去重端口
	opts.portList = removeDuplicatePorts(opts.portList)
	return nil
}

// removeDuplicatePorts 去除重复端口
func removeDuplicatePorts(ports []int) []int {
	seen := make(map[int]bool)
	result := make([]int, 0)

	for _, port := range ports {
		if !seen[port] {
			seen[port] = true
			result = append(result, port)
		}
	}

	return result
}

// PortScanResult 端口扫描结果
type PortScanResult struct {
	Port    int
	State   string
	Service string
	Reason  string
}

// scanPorts 执行端口扫描
func scanPorts(opts *nmapCmdOpts) (err error) {
	start := time.Now()
	ccolor.Infof("Starting Nmap scan on host: %s, start time: %s\n", opts.host, start.Format("2006-01-02 15:04:05"))

	// 根据协议确定显示名称
	protoDisplay := strings.ToUpper(opts.protocol)
	ccolor.Infof("Scanning [%d] ports with %d threads and %dms timeout (%s)\n\n",
		len(opts.portList), opts.threads, opts.timeout, protoDisplay)

	var ipAddr *net.IPAddr
	// 解析主机地址 - 已经是IP地址则跳过解析
	if ipVal := net.ParseIP(opts.host); ipVal != nil {
		ipAddr = &net.IPAddr{IP: ipVal}
	} else {
		ipAddr, err = net.ResolveIPAddr("ip4", opts.host)
		if err != nil {
			return fmt.Errorf("failed to resolve host %s: %w", opts.host, err)
		}
		fmt.Printf("Resolved %s to %s\n\n", opts.host, ipAddr.IP.String())
	}

	// 创建等待组
	var wg sync.WaitGroup

	// 存储扫描结果的通道
	// 创建通道用于控制并发数量
	semaphore := make(chan struct{}, opts.threads)
	results := make(chan PortScanResult, len(opts.portList))

	// 启动接收扫描结果的goroutine
	go func() {
		for result := range results {
			if result.State == "open" {
				if opts.serviceDetect {
					fmt.Printf("%d/%s\t%s\t%s\n", result.Port, strings.ToLower(protoDisplay), result.State, result.Service)
				} else {
					fmt.Printf("%d/%s\t%s\n", result.Port, strings.ToLower(protoDisplay), result.State)
				}
			} else if opts.verbose {
				fmt.Printf("%d/%s\t%s\t%s\n", result.Port, strings.ToLower(protoDisplay), result.State, result.Reason)
			}
		}
	}()

	// 扫描每个端口
	timeout := time.Duration(opts.timeout) * time.Millisecond
	for _, port := range opts.portList {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(p int) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			state, reason, service := probePort(ipAddr.IP.String(), p, opts.protocol, opts.serviceDetect, timeout)
			results <- PortScanResult{
				Port:    p,
				State:   state,
				Service: service,
				Reason:  reason,
			}
		}(port)
	}

	// 等待所有扫描完成
	wg.Wait()
	close(results)

	ccolor.Successln("\nNmap scan completed. Cost time:", time.Since(start))
	return nil
}

// probePort 探测端口状态
func probePort(host string, port int, protocol string, detectService bool, timeout time.Duration) (state, reason, service string) {
	switch strings.ToLower(protocol) {
	case "t", "tcp":
		return probeTCPPort(host, port, detectService, timeout)
	case "u", "udp":
		return probeUDPPort(host, port, detectService, timeout)
	default:
		return "unknown", "unsupported_protocol", "unknown"
	}
}

// probeTCPPort 探测TCP端口
func probeTCPPort(host string, port int, detectService bool, timeout time.Duration) (state, reason, service string) {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return "closed", err.Error(), ""
	}
	defer conn.Close()

	service = detectServiceName(port, "tcp")

	// 如果需要服务检测 - 尝试获取更多服务信息
	if detectService {
		serviceDetail := probeService(conn, port, "tcp", timeout)
		if serviceDetail != "" {
			service = serviceDetail
		}
	}

	return "open", "syn-ack", service
}

// probeUDPPort 探测UDP端口
func probeUDPPort(host string, port int, detectService bool, timeout time.Duration) (state, reason, service string) {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("udp", target, timeout)
	if err != nil {
		return "closed", err.Error(), ""
	}
	defer conn.Close()

	// 对于UDP，我们发送一个简单的数据包来测试
	_, err = conn.Write([]byte(""))
	if err != nil {
		return "closed", "write_failed", ""
	}

	// 如果需要服务检测 TODO
	service = detectServiceName(port, "udp")
	return "open|filtered", "no_response", service
}

var (
	serviceDescMap = map[string]string{
		"ftp":    "File Transfer Protocol",
		"ssh":    "Secure Shell",
		"telnet": "Telnet",
		"smtp":   "Simple Mail Transfer Protocol",
		"domain": "Domain Name System",
		"http":   "Hypertext Transfer Protocol",
		"pop3":   "Post Office Protocol",
		"imap":   "Internet Message Access Protocol",
		"https":  "Hypertext Transfer Protocol Secure",
		"imaps":  "Internet Message Access Protocol Secure",
	}

	// 常见端口对应的服务名称
	tcpServiceMap = map[int]string{
		21: "ftp",
		22: "ssh",
		23: "telnet",
		25: "smtp",
		53: "domain",
		80: "http",

		110: "pop3",
		143: "imap",
		443: "https",
		502: "modbus",
		993: "imaps",
		995: "pop3s",

		1433: "mssql",
		3306: "mysql",
		5432: "postgres",
		5555: "adb",
		6379: "redis",

		11740: "codesys",
		27017: "mongodb",
	}

	// 常见端口对应的服务名称
	udpServiceMap = map[int]string{
		53: "domain",
		67: "dhcp",
		68: "dhcp",

		123: "ntp",
		161: "snmp",
		514: "syslog",
	}
)

// detectServiceName 检测服务名称
func detectServiceName(port int, protocol string) string {
	if protocol == "tcp" {
		if service, ok := tcpServiceMap[port]; ok {
			return service
		}
	}

	if protocol == "udp" {
		if service, ok := udpServiceMap[port]; ok {
			return service
		}
	}
	return "unknown"
}

// probeService 探测服务详情
func probeService(conn net.Conn, port int, protocol string, timeout time.Duration) string {
	// 设置读取超时
	err1 := conn.SetReadDeadline(time.Now().Add(timeout))
	if err1 != nil {
		ccolor.Warnln("Failed to set read deadline:", err1)
	}

	// 根据端口发送特定的探测请求
	switch port {
	case 80, 8080:
		// HTTP 请求
		_, _ = fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			response := string(buffer[:n])
			// 简单解析HTTP响应头
			lines := strings.Split(response, "\n")
			if len(lines) > 0 {
				return strings.TrimSpace(lines[0]) // 返回HTTP状态行
			}
		}
	case 443:
		// HTTPS 连接不会返回明文，但可以检测连接是否成功
		return "https"
	case 22:
		// SSH 服务检测
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			response := string(buffer[:n])
			if strings.HasPrefix(response, "SSH-") {
				return strings.TrimSpace(response)
			}
		}
		return "ssh"
	case 21:
		// FTP 服务检测
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			response := string(buffer[:n])
			if strings.HasPrefix(response, "220") {
				return strings.TrimSpace(response)
			}
		}
		return "ftp"
	}

	return ""
}
