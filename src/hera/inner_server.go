package hera

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
)

// A GracefulServer maintains a WaitGroup that counts how many in-flight
// requests the server is handling. When it receives a shutdown signal,
// it stops accepting new requests but does not actually shut down until
// all in-flight requests terminate.
//
// GracefulServer embeds the underlying net/http.Server making its non-override
// methods and properties avaiable.
//
// It must be initialized by calling NewWithServer.
type GracefulServer struct {
	*http.Server

	shutdown chan bool
	restart  chan bool
	socket   *net.TCPListener

	wg            *sync.WaitGroup
	lcsmu         sync.RWMutex
	lastConnState map[net.Conn]http.ConnState
}

func NewServer(port int, handler http.Handler) (*GracefulServer, error) {
	srv := &GracefulServer{
		Server:        &http.Server{Handler: handler},
		shutdown:      make(chan bool),
		restart:       make(chan bool),
		socket:        nil,
		wg:            new(sync.WaitGroup),
		lastConnState: make(map[net.Conn]http.ConnState),
	}

	addr := fmt.Sprintf(":%d", port)
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	srv.socket = ln.(*net.TCPListener)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func NewServerFromFD(fd uintptr, handler http.Handler) (*GracefulServer, error) {
	srv := &GracefulServer{
		Server:        &http.Server{Handler: handler},
		shutdown:      make(chan bool),
		restart:       make(chan bool),
		socket:        nil,
		wg:            new(sync.WaitGroup),
		lastConnState: make(map[net.Conn]http.ConnState),
	}

	file := os.NewFile(fd, "/tmp/sock-go-graceful-restart")
	listener, err := net.FileListener(file)
	if err != nil {
		return nil, errors.New("File to recover socket from file descriptor: " + err.Error())
	}
	listenerTCP, ok := listener.(*net.TCPListener)
	if !ok {
		return nil, fmt.Errorf("File descriptor %d is not a valid TCP socket", fd)
	}
	srv.socket = listenerTCP
	return srv, nil
}

// Close stops the server from accepting new requets and begins shutting down.
// It returns true if it's the first time Close is called.
func (s *GracefulServer) Close() bool {
	return <-s.shutdown
}

func (s *GracefulServer) Restart() bool {
	return <-s.restart
}

func (s *GracefulServer) ListenerFD() (uintptr, error) {
	file, err := s.socket.File()
	if err != nil {
		return 0, err
	}
	return file.Fd(), nil
}

func (s *GracefulServer) ShutdownRoutine(closing *int32) {
	s.shutdown <- true
	close(s.shutdown)
	atomic.StoreInt32(closing, 1)
	s.Server.SetKeepAlivesEnabled(false)
	s.socket.Close()
	fmt.Println("socket is closed")
	log.Println(os.Getpid(), "Server gracefully shutdown before wait")
	s.wg.Wait()
	log.Println(os.Getpid(), "Server gracefully shutdown after wait")
	os.Exit(0)
}

func (s *GracefulServer) RestartRoutine() {
	s.restart <- true
	close(s.restart)
	listenerFD, err := s.ListenerFD()
	if err != nil {
		log.Fatalln("Fail to get socket file descriptor:", err)
	}

	// Set a flag for the new process start process
	os.Setenv("_RESTART_", "true")
	execSpec := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), listenerFD},
	}

	// Fork exec the new version of your server
	fork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		log.Fatalln("Fail to fork", err)
	}
	log.Println("SIGHUP received: fork-exec to", fork)
	// Wait for all conections to be finished
	//s.Wait()
	log.Println(os.Getpid(), "Server gracefully restart before wait")
	s.wg.Wait()
	log.Println(os.Getpid(), "Server gracefully restart after wait")
	os.Exit(0)
}

// Serve provides a graceful equivalent net/http.Server.Serve.
func (s *GracefulServer) Serve() error {
	var closing int32
	go s.ShutdownRoutine(&closing)
	go s.RestartRoutine()

	originalConnState := s.Server.ConnState

	// s.ConnState is invoked by the net/http.Server every time a connectiion
	// changes state. It keeps track of each connection's state over time,
	// enabling manners to handle persisted connections correctly.
	s.ConnState = func(conn net.Conn, newState http.ConnState) {
		s.lcsmu.RLock()
		lastConnState := s.lastConnState[conn]
		s.lcsmu.RUnlock()

		switch newState {

		// New connection -> StateNew
		case http.StateNew:
			s.StartRoutine()

			// (StateNew, StateIdle) -> StateActive
		case http.StateActive:
			// The connection transitioned from idle back to active
			if lastConnState == http.StateIdle {
				s.StartRoutine()
			}

			// StateActive -> StateIdle
			// Immediately close newly idle connections; if not they may make
			// one more request before SetKeepAliveEnabled(false) takes effect.
		case http.StateIdle:
			if atomic.LoadInt32(&closing) == 1 {
				conn.Close()
			}
			s.FinishRoutine()

			// (StateNew, StateActive, StateIdle) -> (StateClosed, StateHiJacked)
			// If the connection was idle we do not need to decrement the counter.
		case http.StateClosed, http.StateHijacked:
			if lastConnState != http.StateIdle {
				s.FinishRoutine()
			}

		}

		s.lcsmu.Lock()
		if newState == http.StateClosed || newState == http.StateHijacked {
			delete(s.lastConnState, conn)
		} else {
			s.lastConnState[conn] = newState
		}
		s.lcsmu.Unlock()

		if originalConnState != nil {
			originalConnState(conn, newState)
		}
	}

	// A hook to allow the server to notify others when it is ready to receive
	// requests; only used by tests.
	//if s.up != nil {
	//s.up <- listener
	//}

	//err := s.Server.Serve(listener)
	err := s.Server.Serve(s.socket)

	// This block is reached when the server has received a shut down command
	// or a real error happened.
	if err == nil || atomic.LoadInt32(&closing) == 1 {
		s.wg.Wait()
		return nil
	}

	return err
}

// StartRoutine increments the server's WaitGroup. Use this if a web request
// starts more goroutines and these goroutines are not guaranteed to finish
// before the request.
func (s *GracefulServer) StartRoutine() {
	s.wg.Add(1)
}

// FinishRoutine decrements the server's WaitGroup. Use this to complement
// StartRoutine().
func (s *GracefulServer) FinishRoutine() {
	s.wg.Done()
}
