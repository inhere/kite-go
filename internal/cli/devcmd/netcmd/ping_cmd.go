package netcmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/gookit/gcli/v3"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type pingCmdOption struct {
	count   int
	timeout int
}

// NewPingCmd 实现ping工具命令，增强支持多个IP
func NewPingCmd() *gcli.Command {
	var pingOpts = pingCmdOption{}

	return &gcli.Command{
		Name: "ping",
		Desc: "Ping one or multi host",
		// Aliases: []string{"p"},
		Config: func(c *gcli.Command) {
			c.IntOpt(&pingOpts.count, "count", "c", 4, "Number of echo requests to send")
			c.IntOpt(&pingOpts.timeout, "timeout", "t", 3, "Timeout in seconds")

			c.AddArg("hosts", "The host(s) to ping", true, true)
		},
		Func: func(c *gcli.Command, args []string) error {
			hosts := c.Arg("hosts").Strings()

			for _, host := range hosts {
				c.Infof("Pinging %s with %d packets:\n", host, pingOpts.count)

				err := pingHost(host, pingOpts.count, pingOpts.timeout)
				if err != nil {
					c.Errorf("Failed to ping %s: %v\n", host, err)
				}
				c.Println()
			}

			return nil
		},
	}
}

// pingHost 执行对单个主机的 ping 操作
func pingHost(host string, count, timeout int) error {
	// 检查系统是否有 ping 命令
	if hasPingCommand() {
		return pingWithSysCmd(host, count, timeout)
	}

	// 使用 ICMP 库实现 ping
	return pingWithICMP(host, count, timeout)
}

// hasPingCommand 检查系统是否有 ping 命令
func hasPingCommand() bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("where", "ping")
	} else {
		cmd = exec.Command("which", "ping")
	}

	err := cmd.Run()
	return err == nil
}

// pingWithSysCmd 使用系统 ping 命令
func pingWithSysCmd(host string, count, timeout int) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Windows: ping -n count -w timeout(ms) host
		timeoutMs := timeout * 1000 // Windows 使用毫秒
		cmd = exec.Command("ping", "-n", strconv.Itoa(count), "-w", strconv.Itoa(timeoutMs), host)
	} else {
		// Unix-like: ping -c count -W timeout(s) host
		cmd = exec.Command("ping", "-c", strconv.Itoa(count), "-W", strconv.Itoa(timeout), host)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// pingWithICMP 使用 ICMP 库实现 ping
func pingWithICMP(host string, count, timeout int) error {
	// 解析主机地址
	ipAddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return fmt.Errorf("failed to resolve host %s: %v", host, err)
	}

	// 创建 ICMP 连接
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return fmt.Errorf("failed to listen ICMP: %v", err)
	}
	defer conn.Close()

	// 设置超时
	timeoutDuration := time.Duration(timeout) * time.Second
	err = conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	if err != nil {
		return err
	}

	fmt.Printf("PING %s (%s): %d data bytes\n", host, ipAddr.IP, 56)

	successCount := 0
	totalTime := time.Duration(0)

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
			return fmt.Errorf("failed to marshal ICMP message: %v", err)
		}

		// 记录发送时间
		start := time.Now()

		// 发送消息
		_, err = conn.WriteTo(msgBytes, &net.IPAddr{IP: ipAddr.IP})
		if err != nil {
			fmt.Printf("Request timeout for seq %d\n", i)
			continue
		}

		// 接收回复
		reply := make([]byte, 1500)
		n, _, err := conn.ReadFrom(reply)
		if err != nil {
			fmt.Printf("Request timeout for seq %d\n", i)
			continue
		}

		// 计算耗时
		duration := time.Since(start)
		totalTime += duration
		successCount++

		// 解析回复
		rm, err := icmp.ParseMessage(1, reply[:n])
		if err != nil {
			fmt.Printf("Failed to parse reply for seq %d\n", i)
			continue
		}

		// 检查回复类型
		if rm.Type == ipv4.ICMPTypeEchoReply {
			fmt.Printf("%d bytes from %s: icmp_seq=%d ttl=%d time=%.2f ms\n",
				n, ipAddr.IP, i, 64, float64(duration.Nanoseconds())/1e6)
		}

		// 等待间隔
		time.Sleep(time.Second)
	}

	// 统计信息
	packetLoss := float64(count-successCount) / float64(count) * 100
	fmt.Printf("\n--- %s ping statistics ---\n", host)
	fmt.Printf("%d packets transmitted, %d packets received, %.1f%% packet loss\n",
		count, successCount, packetLoss)

	if successCount > 0 {
		avgTime := float64(totalTime.Nanoseconds()) / float64(successCount) / 1e6
		fmt.Printf("round-trip min/avg/max = %.3f/%.3f/%.3f ms\n", avgTime, avgTime, avgTime)
	}

	return nil
}
