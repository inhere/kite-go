package sshcmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/interact"
	"github.com/gookit/goutil/envutil"
	"golang.org/x/crypto/ssh"
)

type sshExecOptions struct {
	scriptFile  string
	script      string
	envVars     gflag.KVString
	ips         string
	ipFile      string
	password    string
	username    string
	autoConfirm bool
	concurrency int
	timeout     int
	sshPort     int
}

// getUserName 获取SSH登录用户名
func (o *sshExecOptions) getUserName() string {
	if o.username != "" {
		return o.username
	}
	return envutil.Getenv("SSH_USER", "root")
}

type sshExecResult struct {
	ip         string
	success    bool
	output     string
	errMessage string
	duration   time.Duration
}

func NewSshExecCmd() *gcli.Command {
	opts := &sshExecOptions{
		username:    "root",
		concurrency: 0,
		timeout:     30,
		sshPort:     22,
	}

	return &gcli.Command{
		Name:    "sshexec",
		Desc:    "Execute script on remote machines via SSH",
		Aliases: []string{"ssh-exec", "se"},
		Config: func(c *gcli.Command) {
			c.StrOpt(&opts.scriptFile, "file", "f", "", "Local script file path to execute")
			c.StrOpt(&opts.script, "script", "s", "", "Script content to execute directly")
			c.VarOpt2(&opts.envVars, "env,e", "Environment variables for script execution (KEY=VALUE)")
			c.StrOpt(&opts.ips, "ip", "", "", "Remote machine IP addresses, comma-separated for multiple")
			c.StrOpt(&opts.ipFile, "ip-file", "", "", "File containing IP addresses (one per line)")
			c.StrOpt(&opts.password, "pwd", "p", "", "SSH login password (or set SSH_PWD env var)")
			c.StrOpt(&opts.username, "user", "u", "", "SSH login username, default is 'root' (or set SSH_USER env var)")
			c.BoolOpt(&opts.autoConfirm, "yes", "y", false, "Auto confirm, skip confirmation prompt for each machine")
			c.IntOpt(&opts.concurrency, "concurrency", "c", 5, "Concurrent execution count for multiple IPs")
			c.IntOpt(&opts.timeout, "timeout", "", 30, "SSH connection timeout in seconds")
			c.IntOpt(&opts.sshPort, "port", "", 22, "SSH server port")

			c.AddArg("ips-arg", "IP addresses (alternative to --ip option)", false, true)
		},
		Func: func(c *gcli.Command, args []string) error {
			return runSshExec(c, opts, args)
		},
	}
}

func runSshExec(c *gcli.Command, opts *sshExecOptions, args []string) error {
	script, err := prepareScript(opts)
	if err != nil {
		return err
	}

	if script == "" {
		return fmt.Errorf("no script provided, use --file or --script option")
	}

	ips, err := collectIPs(opts, args)
	if err != nil {
		return err
	}

	if len(ips) == 0 {
		return fmt.Errorf("no IP addresses provided, use --ip, --ip-file or positional arguments")
	}

	password := resolvePassword(opts)
	if password == "" {
		return fmt.Errorf("SSH password required, use --pwd option or set SSH_PWD environment variable")
	}
	opts.password = password
	opts.username = opts.getUserName()

	c.Infof("Will execute script on %d machine(s):\n", len(ips))
	for _, ip := range ips {
		c.Printf("  - %s\n", ip)
	}

	if !opts.autoConfirm && len(ips) > 0 {
		if !interact.Confirm("Are you sure to continue?") {
			c.Println("Execution cancelled")
			return nil
		}
	}

	envMap := parseEnvVars(opts.envVars)

	var results []sshExecResult
	if len(ips) > 1 && opts.concurrency > 1 {
		results = executeConcurrent(c, ips, script, envMap, opts, password)
	} else {
		results = executeSequential(c, ips, script, envMap, opts, password)
	}

	printSummary(c, results)

	return nil
}

func prepareScript(opts *sshExecOptions) (string, error) {
	if opts.script != "" {
		return opts.script, nil
	}

	if opts.scriptFile != "" {
		content, err := os.ReadFile(opts.scriptFile)
		if err != nil {
			return "", fmt.Errorf("failed to read script file: %w", err)
		}
		return string(content), nil
	}

	return "", nil
}

func collectIPs(opts *sshExecOptions, args []string) ([]string, error) {
	var ips []string
	seen := make(map[string]bool)

	addIP := func(ip string) {
		ip = strings.TrimSpace(ip)
		if ip != "" && !seen[ip] {
			seen[ip] = true
			ips = append(ips, ip)
		}
	}

	if opts.ips != "" {
		for _, ip := range strings.Split(opts.ips, ",") {
			addIP(ip)
		}
	}

	if opts.ipFile != "" {
		file, err := os.Open(opts.ipFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open IP file: %w", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			addIP(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to read IP file: %w", err)
		}
	}

	for _, arg := range args {
		addIP(arg)
	}

	return ips, nil
}

func resolvePassword(opts *sshExecOptions) string {
	if opts.password != "" {
		return opts.password
	}
	return os.Getenv("SSH_PWD")
}

func parseEnvVars(envVars gflag.KVString) map[string]string {
	result := make(map[string]string)
	for k, v := range envVars.SMap {
		result[k] = v
	}
	return result
}

func executeConcurrent(c *gcli.Command, ips []string, script string, envMap map[string]string, opts *sshExecOptions, password string) []sshExecResult {
	results := make([]sshExecResult, len(ips))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, opts.concurrency)

	for i, ip := range ips {
		wg.Add(1)
		go func(idx int, host string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			results[idx] = executeOnHost(c, host, script, envMap, opts, password)
		}(i, ip)
	}

	wg.Wait()
	return results
}

func executeSequential(c *gcli.Command, ips []string, script string, envMap map[string]string, opts *sshExecOptions, password string) []sshExecResult {
	results := make([]sshExecResult, len(ips))

	for i, ip := range ips {
		results[i] = executeOnHost(c, ip, script, envMap, opts, password)
	}

	return results
}

func executeOnHost(c *gcli.Command, ip, script string, envMap map[string]string, opts *sshExecOptions, password string) sshExecResult {
	start := time.Now()
	result := sshExecResult{ip: ip}

	c.Printf("\n[%s] Connecting...\n", ip)

	config := &ssh.ClientConfig{
		User: opts.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(opts.timeout) * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", ip, opts.sshPort)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		result.errMessage = fmt.Sprintf("failed to connect: %v", err)
		result.duration = time.Since(start)
		c.Errorf("[%s] Connection failed: %v\n", ip, err)
		return result
	}
	defer client.Close()

	c.Printf("[%s] Connected, executing script...\n", ip)

	session, err := client.NewSession()
	if err != nil {
		result.errMessage = fmt.Sprintf("failed to create session: %v", err)
		result.duration = time.Since(start)
		c.Errorf("[%s] Failed to create session: %v\n", ip, err)
		return result
	}
	defer session.Close()

	for k, v := range envMap {
		if err := session.Setenv(k, v); err != nil {
			c.Warnf("[%s] Failed to set env %s: %v\n", ip, k, err)
		}
	}

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(script)
	result.duration = time.Since(start)

	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\n[STDERR]\n" + stderr.String()
	}
	result.output = strings.TrimSpace(output)

	if err != nil {
		result.errMessage = fmt.Sprintf("script execution failed: %v", err)
		c.Errorf("[%s] Execution failed: %v\n", ip, err)
		if result.output != "" {
			c.Printf("[%s] Output:\n%s\n", ip, result.output)
		}
		return result
	}

	result.success = true
	c.Successf("[%s] Success (took %v)\n", ip, result.duration)
	if result.output != "" {
		c.Printf("[%s] Output:\n%s\n", ip, result.output)
	}

	return result
}

func printSummary(c *gcli.Command, results []sshExecResult) {
	var successCount, failCount int
	var totalDuration time.Duration

	c.Println("\n" + strings.Repeat("=", 50))
	c.Println("Execution Summary:")
	c.Println(strings.Repeat("=", 50))

	for _, r := range results {
		totalDuration += r.duration
		if r.success {
			successCount++
			c.Successf("  ✓ %s - %v\n", r.ip, r.duration)
		} else {
			failCount++
			c.Errorf("  ✗ %s - %s\n", r.ip, r.errMessage)
		}
	}

	c.Println(strings.Repeat("-", 50))
	c.Printf("Total: %d | Success: %d | Failed: %d\n", len(results), successCount, failCount)
	c.Printf("Total time: %v\n", totalDuration)
}

type SshExecutor struct {
	Username    string
	Password    string
	Port        int
	Timeout     time.Duration
	EnvVars     map[string]string
	AutoConfirm bool
}

func NewSshExecutor(username, password string) *SshExecutor {
	return &SshExecutor{
		Username: username,
		Password: password,
		Port:     22,
		Timeout:  30 * time.Second,
		EnvVars:  make(map[string]string),
	}
}

func (e *SshExecutor) Connect(ip string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: e.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(e.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         e.Timeout,
	}

	addr := fmt.Sprintf("%s:%d", ip, e.Port)
	return ssh.Dial("tcp", addr, config)
}

func (e *SshExecutor) Execute(client *ssh.Client, script string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	for k, v := range e.EnvVars {
		if err := session.Setenv(k, v); err != nil {
		}
	}

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(script)
	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\n[STDERR]\n" + stderr.String()
	}

	return strings.TrimSpace(output), err
}

func (e *SshExecutor) ExecuteOnHost(ip, script string) (string, time.Duration, error) {
	start := time.Now()

	client, err := e.Connect(ip)
	if err != nil {
		return "", time.Since(start), fmt.Errorf("connection failed: %w", err)
	}
	defer client.Close()

	output, err := e.Execute(client, script)
	return output, time.Since(start), err
}

type MultiExecutor struct {
	Executor    *SshExecutor
	Concurrency int
}

func NewMultiExecutor(executor *SshExecutor, concurrency int) *MultiExecutor {
	return &MultiExecutor{
		Executor:    executor,
		Concurrency: concurrency,
	}
}

func (m *MultiExecutor) ExecuteAll(ips []string, script string, writer io.Writer) map[string]sshExecResult {
	results := make(map[string]sshExecResult)
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, m.Concurrency)

	for _, ip := range ips {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			output, duration, err := m.Executor.ExecuteOnHost(host, script)

			mu.Lock()
			results[host] = sshExecResult{
				ip:       host,
				success:  err == nil,
				output:   output,
				duration: duration,
			}
			if err != nil {
				results[host] = sshExecResult{
					ip:         host,
					success:    false,
					errMessage: err.Error(),
					duration:   duration,
				}
			}
			mu.Unlock()
		}(ip)
	}

	wg.Wait()
	return results
}
