package webapp

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gookit/color"
	"github.com/gookit/event"
	"github.com/gookit/goutil/sysutil"
)

// HTTPServer an HTTP web server
type HTTPServer struct {
	srv *http.Server

	pidFile  string
	address  []string
	realAddr string

	processID int
}

// NewHTTPServer create new HTTPServer.
//
// Usage:
// 	srv := NewHTTPServer("127.0.0.1")
// 	srv := NewHTTPServer("127.0.0.1:8090")
// 	srv := NewHTTPServer("127.0.0.1", "8090")
func NewHTTPServer(address ...string) *HTTPServer {
	return &HTTPServer{
		processID: os.Getpid(),

		address:  address,
		realAddr: resolveAddress(address),
	}
}

/*************************************************************
 * Start HTTP server
 *************************************************************/

// Start server, begin handle HTTP request
func (s *HTTPServer) Start() error {
	app := nako.App()

	s.srv = &http.Server{
		Addr: s.realAddr,
	}

	s.srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(500)
			}
		}()

		evtData := event.M{"w": w, "r": r}

		// Fire before route
		app.MustFire(OnBeforeRoute, evtData)

		// Route and dispatch request
		app.Router.ServeHTTP(w, r)

		// Fire after route
		app.MustFire(OnAfterRoute, evtData)
	})

	app.MustFire(OnServerStart, event.M{"addr": s.srv.Addr})

	// Listen signal
	s.handleSignal(s.srv)

	// Save pid to file
	savePidToFile(s.processID, s.pidFile)

	// Start server
	err := s.srv.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	return removePidFile(s.pidFile)
}

// Stop by send signal to exists process
func (s *HTTPServer) Stop(timeout int) error {
	if s.srv != nil {
		return s.Shutdown(timeout)
	}

	pid := s.processID
	if pid <= 0 {
		return fmt.Errorf("invalid process ID value")
	}

	// Server is not started
	if !sysutil.ProcessExists(pid) {
		return fmt.Errorf("cannot stop, the process is not exists(PID: %d)", pid)
	}

	err := sysutil.Kill(pid, syscall.SIGTERM)
	if err == nil {
		return removePidFile(s.pidFile)
	}
	return err
}

// Shutdown from internal
func (s *HTTPServer) Shutdown(timeout int) error {
	if s.srv == nil {
		return fmt.Errorf("server is not running")
	}

	ctx, cancel := context.WithTimeout(
		context.TODO(),
		time.Duration(timeout)*time.Second,
	)

	defer cancel()
	return s.srv.Shutdown(ctx)
}

/*************************************************************
 * Getter/Setter methods
 *************************************************************/

// IsRunning get server is running
func (s *HTTPServer) IsRunning() bool {
	if s.srv != nil {
		return true
	}

	// Get stat by pid
	if s.pidFile == "" {
		return false
	}

	// Get pid from file
	bts, err := ioutil.ReadFile(s.pidFile)
	if err != nil {
		// panic(err)
		return false
	}

	pid, _ := strconv.Atoi(string(bts))
	if pid <= 0 {
		// panic("invalid process ID value, read from pidFile")
		return false
	}

	// Storage pid value
	s.processID = pid

	// Server is not started
	return sysutil.ProcessExists(pid)
}

// RealAddr get resolved read addr
func (s *HTTPServer) RealAddr() string {
	return s.realAddr
}

// ProcessID return
func (s *HTTPServer) ProcessID() int {
	return s.processID
}

// PidFile get pid file path
func (s *HTTPServer) PidFile() string {
	return s.pidFile
}

// SetPidFile set pid file path
func (s *HTTPServer) SetPidFile(pidFile string) {
	s.pidFile = pidFile
}

// handleSignal handles system signal for graceful shutdown.
func (s *HTTPServer) handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		fmt.Printf("Got signal [%s], exiting server now\n", s)
		if err := server.Close(); err != nil {
			fmt.Printf("Server close failed: %s", err.Error())
		}

		nako.App().MustFire(OnServerClose, event.M{"sig": s})
		// service.DisconnectDB()

		color.Infoln("Server exited")
		os.Exit(0)
	}()
}

func savePidToFile(pid int, pidFile string) {
	if pidFile == "" {
		return
	}

	txt := fmt.Sprintf("%d", pid)
	err := ioutil.WriteFile(pidFile, []byte(txt), 0664)
	if err != nil {
		panic(err)
	}
}

func removePidFile(pidFile string) error {
	if pidFile == "" {
		return nil
	}

	return os.Remove(pidFile)
}

func resolveAddress(addr []string) (fullAddr string) {
	ip := "0.0.0.0"
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); len(port) > 0 {
			fmt.Printf("Environment variable PORT=\"%s\"", port)
			return ip + ":" + port
		}
		fmt.Printf("Environment variable PORT is undefined. Using port :8080 by default")
		return ip + ":8080"
	case 1:
		var port string
		if strings.IndexByte(addr[0], ':') != -1 {
			ss := strings.SplitN(addr[0], ":", 2)
			if ss[0] != "" {
				return addr[0]
			}
			port = ss[1]
		} else {
			port = addr[0]
		}

		return ip + ":" + port
	default:
		panic("too much parameters")
	}
}
