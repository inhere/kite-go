package httpcmd

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/netutil"
	"github.com/gookit/goutil/x/ccolor"
)

type httpAccessCheck struct {
	// hostPatterns ip pattern list.
	//  - one item eg: 192.168.1.0/24 OR 192.168.1.22-255 OR 192.168.1.*
	HostPatterns gcli.Strings
	// Port server  port. eg: 80, 443
	Port uint
	// BaseURI server base uri. default is empty.
	BaseURI string
	// timeout for check request
	Timeout time.Duration
	Verbose bool
	// number of workers to perform checks
	Workers int

	// -- internal fields
	hostIps []string
}

// NewAccessCheckCmd instance
//
//	实现批量的基于IP段的 http server 访问检测.
func NewAccessCheckCmd() *gcli.Command {
	var accessCheck = httpAccessCheck{}

	return &gcli.Command{
		Name:    "check",
		Desc:    "batch discover and check http server access status",
		Aliases: []string{"scan"},
		Examples: `
{$fullCmd} -H 192.168.1.*
{$fullCmd} -H 192.168.0-1.* -p 8080
`,
		Config: func(c *gcli.Command) {
			c.DurationOpt(&accessCheck.Timeout, "timeout", "t", 2*time.Second, "set http server check timeout")
			c.IntOpt(&accessCheck.Workers, "workers", "w", 50, "set http server check workers")
			c.VarOpt(&accessCheck.HostPatterns, "host", "H", `host ip pattern, allow multi. one pattern eg:
192.168.1.23
192.168.1.*
192.168.1.0/24
192.168.1.23-251
192.168.0-1.23-251
`)
			c.UintOpt(&accessCheck.Port, "port", "p", 80, "server port")
			c.StringOpt(&accessCheck.BaseURI, "base-uri", "u", "", "server base uri, default access /")
			c.BoolOpt(&accessCheck.Verbose, "verbose", "v", false, "show more details")
		},
		Func: func(c *gcli.Command, args []string) error {
			return accessCheck.Run()
		},
	}
}

func (hac *httpAccessCheck) Run() error {
	start := time.Now()
	// 解析主机模式，生成IP地址列表
	err := hac.parseHostPatterns()
	if err != nil {
		return err
	}

	ipLen := len(hac.hostIps)
	ccolor.Infof("Checking %d IP addresses...(workers: %d, timeout: %s)\n", ipLen, hac.Workers, hac.Timeout)

	// 创建工作池进行并发检查
	jobs := make(chan string, ipLen)
	results := make(chan checkResult, ipLen)

	// 启动worker goroutines
	for i := 0; i < hac.Workers; i++ {
		go hac.worker(jobs, results)
	}

	// 发送任务到jobs channel
	for _, ip := range hac.hostIps {
		jobs <- ip
	}
	close(jobs)

	// 收集结果
	var successCount int
	for i := 0; i < ipLen; i++ {
		result := <-results
		if result.err != nil {
			if hac.Verbose {
				gcli.Printf("❌ %s:%d - %v\n", result.ip, hac.Port, result.err)
			}
		} else {
			successCount++
			gcli.Printf("✅ %s:%d - OK (took %v)\n", result.ip, hac.Port, result.duration)
		}
	}

	ccolor.Infof("Check completed in %s, can access %d\n", time.Since(start), successCount)
	return nil
}

func (hac *httpAccessCheck) parseHostPatterns() error {
	for _, pattern := range hac.HostPatterns {

		// 增强：倒数第二段也支持范围和多个值
		// eg: 192.168.1-3.132 -> 192.168.1.132,192.168.2.132,192.168.3.132
		// eg: 192.168.1,3.234 -> 192.168.1.132,192.168.3.132
		parts := strings.Split(pattern, ".")
		if len(parts) >= 4 {
			strVal := parts[len(parts)-2]
			if strings.Contains(strVal, "-") || strings.Contains(strVal, ",") {
				// 倒数第二段: "1-3" OR "1,3"
				ints, err := parseIntRangeOrEnum(strVal)
				if err != nil {
					return err
				}

				for _, i := range ints {
					// 替换倒数第二段
					parts[len(parts)-2] = strconv.Itoa(i)
					ipPattern := strings.Join(parts, ".")

					ips, err := parseHostPattern(ipPattern)
					if err != nil {
						return err
					}
					hac.hostIps = append(hac.hostIps, ips...)
				}
				continue
			}
		}

		// 最后一段有枚举或范围
		ips, err := parseHostPattern(pattern)
		if err != nil {
			return err
		}

		hac.hostIps = append(hac.hostIps, ips...)
	}

	return nil
}

// 将 "1-3" OR "1,3" 转换为 int 列表
func parseIntRangeOrEnum(value string) ([]int, error) {
	var ints []int
	// 处理范围格式 eg: 1-23
	if strings.Contains(value, "-") {
		parts := strings.Split(value, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range format")
		}

		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid range start value")
		}

		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid range end value")
		}
		for i := start; i <= end; i++ {
			ints = append(ints, i)
		}
	} else {
		// 处理枚举格式 eg: 1,2,3
		for _, part := range strings.Split(value, ",") {
			iVal, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid integer value: %s", part)
			}
			ints = append(ints, iVal)
		}
	}
	return ints, nil
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
		for ip := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ip); netutil.IncrIP(ip) {
			ips = append(ips, ip.String())
		}
		return ips, nil
	}

	// 最后一段支持范围格式(e.g., 192.168.1.22-255)
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
		// 默认排除 0,1
		for i := 2; i <= 255; i++ {
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

// checkResult 存储单个IP检查的结果
type checkResult struct {
	ip       string
	duration time.Duration
	err      error
}

// worker 执行HTTP连接检查的工作函数
func (hac *httpAccessCheck) worker(jobs <-chan string, results chan<- checkResult) {
	for ip := range jobs {
		start := time.Now()
		err := checkHTTPServer(ip, hac.Port, hac.Timeout)
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
