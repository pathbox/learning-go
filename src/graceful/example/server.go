package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Scalingo/go-graceful-restart-example/logger"
)

type Server struct {
	cm     *ConnectionManager
	socket *net.TCPListener
	logger *logger.Logger
}

func New(logger *logger.Logger, port int) (*Server, error) {
	s := &Server{cm: NewConnectionManager(), logger: logger}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("fail to resolve addr: %v", err)
	}
	sock, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("fail to listen tcp: %v", err)
	}

	s.socket = sock
	return s, nil
}

func NewFromFD(logger *logger.Logger, fd uintptr) (*Server, error) {
	s := &Server{cm: NewConnectionManager(), logger: logger}

	file := os.NewFile(fd, "/tmp/sock-go-graceful-restart")
	listener, err := net.FileListener(file)
	if err != nil {
		return nil, errors.New("File to recover socket from file descriptor: " + err.Error())
	}
	listenerTCP, ok := listener.(*net.TCPListener)
	if !ok {
		return nil, fmt.Errorf("File descriptor %d is not a valid TCP socket", fd)
	}
	s.socket = listenerTCP

	return s, nil
}

func (s *Server) Stop() {
	// Accept will instantly return a timeout error
	s.socket.SetDeadline(time.Now())
}

func (s *Server) ListenerFD() (uintptr, error) {
	file, err := s.socket.File()
	if err != nil {
		return 0, err
	}
	return file.Fd(), nil
}

func (s *Server) Wait() {
	s.cm.Wait()
}

var WaitTimeoutError = errors.New("timeout")

func (s *Server) WaitWithTimeout(duration time.Duration) error {
	timeout := time.NewTimer(duration)
	wait := make(chan struct{})
	go func() {
		s.Wait()
		wait <- struct{}{}
	}()

	select {
	case <-timeout.C:
		return WaitTimeoutError
	case <-wait:
		return nil
	}
}

func (s *Server) StartAcceptLoop() {
	for {
		conn, err := s.socket.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				s.logger.Println("Stop accepting connections")
				return
			}
			s.logger.Println("[Error] fail to accept:", err)
		}
		go func() {
			s.cm.Add(1)
			s.handleConn(conn)
			s.cm.Done()
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) {
	tick := time.NewTicker(time.Second)
	buffer := make([]byte, 64)
	for {
		select {
		case <-tick.C:
			_, err := conn.Write([]byte("ping"))
			if err != nil {
				s.logger.Println("[Error] fail to write 'ping':", err)
				conn.Close()
				return
			}
			s.logger.Printf("[Server] Sent 'ping'\n")

			n, err := conn.Read(buffer)
			if err != nil {
				s.logger.Println("[Error] fail to read from socket:", err)
				conn.Close()
				return
			}

			s.logger.Printf("[Server] OK: read %d bytes: '%s'\n", n, string(buffer[:n]))
		}
	}
}

func (s *Server) Addr() net.Addr {
	return s.socket.Addr()
}

func (s *Server) ConnectionsCounter() int {
	return s.cm.Counter
}
