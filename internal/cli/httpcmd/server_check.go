package httpcmd

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/x/ccolor"
)

var accessCheckOpts = struct {
	// hostPattern ip pattern. eg: 192.168.1.0/24 OR 192.168.1.22-255 OR 192.168.1.*
	hostPattern string
	// server  port. eg: 80, 443
	port uint
	// timeout  setting
	timeout time.Duration
	workers int
	verbose  bool
}{}

// NewAccessCheckCmd instance
//
//	实现批量的基于IP段的 http server 访问检测.
func NewAccessCheckCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "check",
		Desc:    "batch discover and check http server access status",
		Aliases: []string{"scan"},
		Examples: `
{$fullCmd} -h 192.168.1.*
{$fullCmd} -h 192.168.0-1.*
`,
		Config: func(c *gcli.Command) {
			c.DurationOpt(&accessCheckOpts.timeout, "timeout", "t", 2*time.Second, "set http server check timeout")
			c.IntOpt(&accessCheckOpts.workers, "workers", "w", 10, "set http server check workers")
			c.StrOpt(&accessCheckOpts.hostPattern, "host", "H", "192.168.1.0/24", "host ip pattern")
			c.UintOpt(&accessCheckOpts.port, "port", "p", 80, "server port")
			c.BoolOpt(&accessCheckOpts.verbose, "verbose", "v",  false,"show more details")
		},
		Func: func(c *gcli.Command, args []string) error {
			start := time.Now()
			// 解析主机模式，生成IP地址列表
			ips, err := parseHostPattern(accessCheckOpts.hostPattern)
			if err != nil {
				return err
			}

			serverPort := accessCheckOpts.port
			timeout := accessCheckOpts.timeout
			ccolor.Infof("Checking %d IP addresses...(workers: %d, timeout: %s)\n", len(ips), accessCheckOpts.workers, timeout)

			// 创建工作池进行并发检查
			jobs := make(chan string, len(ips))
			results := make(chan checkResult, len(ips))

			// 启动worker goroutines
			for i := 0; i < accessCheckOpts.workers; i++ {
				go worker(jobs, results, serverPort, timeout)
			}

			// 发送任务到jobs channel
			for _, ip := range ips {
				jobs <- ip
			}
			close(jobs)

			// 收集结果
			for i := 0; i < len(ips); i++ {
				result := <-results
				if result.err != nil {
					if accessCheckOpts.verbose {
						gcli.Printf("❌ %s:%d - %v\n", result.ip, serverPort, result.err)
					}
				} else {
					gcli.Printf("✅ %s:%d - OK (took %v)\n", result.ip, serverPort, result.duration)
				}
			}

			ccolor.Infof("Check completed in %s\n", time.Since(start))
			return nil
		},
	}
}


// checkResult 存储单个IP检查的结果
type checkResult struct {
	ip       string
	duration time.Duration
	err      error
}

// worker 执行HTTP连接检查的工作函数
func worker(jobs <-chan string, results chan<- checkResult, port uint, timeout time.Duration) {
	for ip := range jobs {
		start := time.Now()
		err := checkHTTPServer(ip, port, timeout)
		duration := time.Since(start)

		results <- checkResult{
			ip:       ip,
			duration: duration,
			err:      err,
		}
	}
}

// checkHTTPServer 检查指定IP和端口的HTTP服务器是否可访问
func checkHTTPServer(ip string, port uint, timeout time.Duration) error {
	addr := fmt.Sprintf("http://%s:%d", ip, port)

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(addr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP status code: %d", resp.StatusCode)
	}

	return nil
}

// parseHostPattern 解析主机模式并返回IP地址列表
func parseHostPattern(pattern string) ([]string, error) {
	var ips []string

	// 处理CIDR格式 (e.g., 192.168.1.0/24)
	if strings.Contains(pattern, "/") {
		_, ipnet, err := net.ParseCIDR(pattern)
		if err != nil {
			return nil, err
		}

		// 遍历CIDR范围内的所有IP
		for ip := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ip); incIP(ip) {
			ips = append(ips, ip.String())
		}
		return ips, nil
	}

	// 处理范围格式 (e.g., 192.168.1.22-255)
	if strings.Contains(pattern, "-") {
		parts := strings.Split(pattern, ".")
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid IP range format")
		}

		rangePart := parts[3]
		rangeParts := strings.Split(rangePart, "-")
		if len(rangeParts) != 2 {
			return nil, fmt.Errorf("invalid range format")
		}

		start, err1 := strconv.Atoi(rangeParts[0])
		end, err2 := strconv.Atoi(rangeParts[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid range values")
		}

		for i := start; i <= end; i++ {
			ip := fmt.Sprintf("%s.%s.%s.%d", parts[0], parts[1], parts[2], i)
			ips = append(ips, ip)
		}
		return ips, nil
	}

	// 处理通配符格式 (e.g., 192.168.1.*)
	if strings.HasSuffix(pattern, ".*") {
		base := strings.TrimSuffix(pattern, ".*")
		for i := 0; i <= 255; i++ {
			ip := fmt.Sprintf("%s.%d", base, i)
			ips = append(ips, ip)
		}
		return ips, nil
	}

	// 单个IP地址
	if net.ParseIP(pattern) != nil {
		return []string{pattern}, nil
	}

	return nil, fmt.Errorf("unsupported IP pattern format")
}

// incIP 将IP地址递增1
func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
