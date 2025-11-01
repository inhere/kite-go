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
	// internal fields
	portList []int
}

// NewNMapCmd 实现nmap扫描工具命令
func NewNMapCmd() *gcli.Command {
	var nmapOpts = nmapCmdOpts{}

	return &gcli.Command{
		Name:    "nmap",
		Desc:    "nmap port scan tool",
		Aliases: []string{"scan"},
		Config: func(c *gcli.Command) {
			c.StrOpt(&nmapOpts.ports, "ports", "p", "80,443,22,21,23,25,53,110,143,993,995,3306,5432,6379,27017",
				`ports to scan, comma separated. port must be between 1 and 65535.
Range: 1-100
Both: 1000-2000,8080,4000-5000
`)
			c.IntOpt(&nmapOpts.timeout, "timeout", "t", 3, "connection timeout in milliseconds")
			c.IntOpt(&nmapOpts.threads, "threads", "c,T", 100, "number of concurrent threads")
			c.AddArg("host", "target host/IP to scan", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			nmapOpts.host = c.Arg("host").String()

			// 解析端口列表
			if err := parsePorts(&nmapOpts); err != nil {
				return err
			}

			// 执行端口扫描
			return scanPorts(&nmapOpts)
		},
		Examples: `
  {$binWithCmd} 192.168.1.1
  {$binWithCmd} -p 80,443,8080 example.com
  {$binWithCmd} -p 1-1000 -t 5 -T 50 192.168.1.1
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

// scanPorts 执行端口扫描
func scanPorts(opts *nmapCmdOpts) (err error) {
	start := time.Now()
	ccolor.Infof("Starting scan on host: %s, start time: %s\n", opts.host, start.Format("2006-01-02 15:04:05"))
	ccolor.Infof("Scanning [%d] ports with %d threads and %dms timeout\n\n",
		len(opts.portList), opts.threads, opts.timeout)

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

	// 创建通道用于控制并发数量
	semaphore := make(chan struct{}, opts.threads)

	// 创建等待组
	var wg sync.WaitGroup

	// 存储开放端口的通道
	openPorts := make(chan int, len(opts.portList))

	// 启动接收开放端口的goroutine
	go func() {
		for port := range openPorts {
			fmt.Printf("Port %d/tcp is open\n", port)
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

			if isPortOpen(ipAddr.IP.String(), p, timeout) {
				openPorts <- p
			}
		}(port)
	}

	// 等待所有扫描完成
	wg.Wait()
	close(openPorts)

	ccolor.Successln("\nScan completed. Cost time:", time.Since(start))
	return nil
}

// isPortOpen 检查端口是否开放
func isPortOpen(host string, port int, timeout time.Duration) bool {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
