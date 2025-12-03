package netcmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/sysutil"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type pingCmdOption struct {
	count   int
	timeout int
	workers int
}

type pingResult struct {
	host         string
	success      bool
	packetLoss   float64
	avgTime      float64
	minTime      float64
	maxTime      float64
	errorMessage string
}

// NewPingCmd 实现ping工具命令，增强支持多个IP和网段检测
func NewPingCmd() *gcli.Command {
	var pingOpts = pingCmdOption{}

	return &gcli.Command{
		Name: "ping",
		Desc: "Ping one or multi host, support CIDR notation (e.g. 192.168.1.0/24)",
		// Aliases: []string{"p"},
		Config: func(c *gcli.Command) {
			c.IntOpt(&pingOpts.count, "count", "c", 4, "Number of echo requests to send")
			c.IntOpt(&pingOpts.timeout, "timeout", "t", 3, "Timeout in seconds")
			c.IntOpt(&pingOpts.workers, "workers", "w", 5, "Number of concurrent workers for multiple hosts")

			c.AddArg("hosts", "The host(s) to ping, support CIDR notation", true, true)
		},
		Func: func(c *gcli.Command, args []string) error {
			inputHosts := c.Arg("hosts").Strings()

			// 解析主机列表，支持CIDR网段
			allHosts := expandHosts(inputHosts)
			c.Infof("Pinging %d hosts with %d packets:\n", len(allHosts), pingOpts.count)

			// 根据主机数量决定是否使用并发
			if len(allHosts) > 3 {
				pingHostsConcurrent(c, allHosts, pingOpts)
			} else {
				pingHostsSequential(c, allHosts, pingOpts)
			}

			return nil
		},
	}
}

// expandHosts 解析主机列表，支持CIDR网段
func expandHosts(hosts []string) []string {
	var result []string

	for _, host := range hosts {
		if strings.Contains(host, "/") {
			// CIDR网段处理
			ips, err := expandCIDR(host)
			if err != nil {
				fmt.Printf("Warning: Failed to parse CIDR %s: %v\n", host, err)
				continue
			}
			result = append(result, ips...)
		} else {
			// 单个主机
			result = append(result, host)
		}
	}

	return result
}

// expandCIDR 解析CIDR网段，返回IP列表
func expandCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
		if len(ips) >= 254 { // 限制最大数量，避免过大的网段
			break
		}
	}

	// 移除网络地址和广播地址
	if len(ips) > 2 {
		return ips[1 : len(ips)-1], nil
	}
	return ips, nil
}

// inc 递增IP地址
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// pingHostsConcurrent 并发ping多个主机
func pingHostsConcurrent(c *gcli.Command, hosts []string, opts pingCmdOption) {
	results := make(chan pingResult, len(hosts))
	var wg sync.WaitGroup

	// 创建工作池
	semaphore := make(chan struct{}, opts.workers)

	for _, host := range hosts {
		wg.Add(1)
		go func(h string) {
			defer wg.Done()
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			result := pingHostForResult(h, opts.count, opts.timeout)
			results <- result
		}(host)
	}

	// 等待所有goroutine完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集并显示结果
	var successCount int
	for result := range results {
		if result.success {
			successCount++
			c.Successf("✓ %s - %.2fms avg (%.1f%% loss)\n",
				result.host, result.avgTime, result.packetLoss)
		} else {
			c.Warnf("✗ %s - %s\n", result.host, result.errorMessage)
		}
	}

	// 显示汇总信息
	c.Printf("\nSummary: %d/%d hosts reachable (%.1f%%)\n",
		successCount, len(hosts), float64(successCount)/float64(len(hosts))*100)
}

// pingHostsSequential 顺序ping多个主机
func pingHostsSequential(c *gcli.Command, hosts []string, opts pingCmdOption) {
	var successCount int

	for _, host := range hosts {
		result := pingHostForResult(host, opts.count, opts.timeout)

		if result.success {
			successCount++
			c.Successf("✓ %s - %.2fms avg (%.1f%% loss)\n",
				host, result.avgTime, result.packetLoss)
		} else {
			c.Warnf("✗ %s - %s\n", host, result.errorMessage)
		}
	}

	// 显示汇总信息
	c.Printf("\nSummary: %d/%d hosts reachable (%.1f%%)\n",
		successCount, len(hosts), float64(successCount)/float64(len(hosts))*100)
}

// pingHostForResult 执行ping并返回结构化结果
func pingHostForResult(host string, count, timeout int) pingResult {
	// 检查系统是否有 ping 命令
	if sysutil.HasExecutable("ping") {
		return pingWithSysCmdResult(host, count, timeout)
	}

	// 使用 ICMP 库实现 ping
	return pingWithICMPResult(host, count, timeout)
}

// pingHost 执行对单个主机的 ping 操作（保留原函数以兼容）
func pingHost(host string, count, timeout int) error {
	result := pingHostForResult(host, count, timeout)
	if !result.success {
		return fmt.Errorf(result.errorMessage)
	}
	return nil
}

// pingWithSysCmdResult 使用系统 ping 命令并返回结构化结果
func pingWithSysCmdResult(host string, count, timeout int) pingResult {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Windows: ping -n count -w timeout(ms) host
		timeoutMs := timeout * 1000 // Windows 使用毫秒
		cmd = exec.Command("ping", "-n", strconv.Itoa(count), "-w", strconv.Itoa(timeoutMs), host)
	} else {
		// Unix-like: ping -c count -W timeout(s) host
		cmd = exec.Command("ping", "-c", strconv.Itoa(count), "-W", strconv.Itoa(timeout), host)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return pingResult{
			host:         host,
			success:      false,
			errorMessage: strings.TrimSpace(string(output)),
		}
	}

	return parsePingOutput(host, string(output))
}

// parsePingOutput 解析ping命令输出，提取关键信息
func parsePingOutput(host, output string) pingResult {
	result := pingResult{host: host}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 查找丢包率信息
		if strings.Contains(line, "packet loss") || strings.Contains(line, "丢失") {
			if strings.Contains(line, "0%") {
				result.packetLoss = 0
				result.success = true
			} else {
				// 简单解析丢包率
				if idx := strings.Index(line, "("); idx != -1 {
					lossStr := line[idx+1:]
					if endIdx := strings.Index(lossStr, "%"); endIdx != -1 {
						lossStr = lossStr[:endIdx]
						if loss, err := strconv.ParseFloat(lossStr, 64); err == nil {
							result.packetLoss = loss
							result.success = loss < 100
						}
					}
				}
			}
		}

		// 查找平均时间信息
		if strings.Contains(line, "Average") || strings.Contains(line, "平均") {
			if idx := strings.Index(line, "="); idx != -1 {
				timeStr := line[idx+1:]
				timeStr = strings.TrimSpace(timeStr)
				if endIdx := strings.Index(timeStr, "ms"); endIdx != -1 {
					timeStr = timeStr[:endIdx]
					if avg, err := strconv.ParseFloat(timeStr, 64); err == nil {
						result.avgTime = avg
					}
				}
			}
		}
	}

	// 如果没有找到丢包信息，假设成功
	if result.packetLoss == 0 && !strings.Contains(output, "timeout") && !strings.Contains(output, "unreachable") {
		result.success = true
	}

	return result
}

// pingWithICMPResult 使用 ICMP 库实现 ping 并返回结构化结果
func pingWithICMPResult(host string, count, timeout int) pingResult {
	result := pingResult{host: host}

	// 解析主机地址
	ipAddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		result.success = false
		result.errorMessage = fmt.Sprintf("failed to resolve host %s: %v", host, err)
		return result
	}

	// 创建 ICMP 连接
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		result.success = false
		result.errorMessage = fmt.Sprintf("failed to listen ICMP: %v", err)
		return result
	}
	defer conn.Close()

	// 设置超时
	timeoutDuration := time.Duration(timeout) * time.Second
	err = conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	if err != nil {
		result.success = false
		result.errorMessage = fmt.Sprintf("failed to set deadline: %v", err)
		return result
	}

	successCount := 0
	var totalTime float64
	var minTime, maxTime float64 = 999999, 0

	// 发送 ICMP Echo 请求
	for i := 1; i <= count; i++ {
		// 构造 ICMP 消息
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  i,
				Data: make([]byte, 56),
			},
		}

		// 编码消息
		msgBytes, err := msg.Marshal(nil)
		if err != nil {
			continue
		}

		// 记录发送时间
		start := time.Now()

		// 发送消息
		_, err = conn.WriteTo(msgBytes, &net.IPAddr{IP: ipAddr.IP})
		if err != nil {
			continue
		}

		// 接收回复
		reply := make([]byte, 1500)
		n, _, err := conn.ReadFrom(reply)
		if err != nil {
			continue
		}

		// 计算耗时
		duration := time.Since(start)
		durationMs := float64(duration.Nanoseconds()) / 1e6
		totalTime += durationMs
		successCount++

		// 更新最小最大时间
		if durationMs < minTime {
			minTime = durationMs
		}
		if durationMs > maxTime {
			maxTime = durationMs
		}

		// 解析回复
		rm, err := icmp.ParseMessage(1, reply[:n])
		if err != nil {
			continue
		}

		// 检查回复类型
		if rm.Type == ipv4.ICMPTypeEchoReply {
			// 成功收到回复
		}

		// 等待间隔
		time.Sleep(time.Second)
	}

	// 填充结果
	result.packetLoss = float64(count-successCount) / float64(count) * 100
	result.success = successCount > 0

	if successCount > 0 {
		result.avgTime = totalTime / float64(successCount)
		result.minTime = minTime
		result.maxTime = maxTime
	} else {
		result.errorMessage = "All packets lost"
	}

	return result
}
