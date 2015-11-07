package election

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

type TcpProxy struct {
	done chan bool
	port int
	to   func() string
}

func NewTcpProxy(port int, to func() string) *TcpProxy {
	return &TcpProxy{
		port: port,
		to:   to,
		done: make(chan bool),
	}
}

func (t *TcpProxy) Close() error {
	t.done <- true
	return nil
}

func (t *TcpProxy) Forward() error {
	strAddr := fmt.Sprintf("0.0.0.0:%d", t.port)

	logrus.Infof("Listening on %s", strAddr)
	addr, err := net.ResolveTCPAddr("tcp", strAddr)
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()

	wg := sync.WaitGroup{}

	for {
		switch {
		case <-t.done:
			return nil
		default:
		}

		l.SetDeadline(time.Now().Add(1 * time.Second))
		conn, err := l.Accept()
		if acceptErr, ok := err.(*net.OpError); ok && acceptErr.Timeout() {
			continue
		}

		if err != nil {
			return err
		}

		wg.Add(1)
		go func(conn net.Conn) {
			defer wg.Done()
			if err := t.forward(conn); err != nil {
				logrus.Errorf("Failed handling TCP forwarding: %v", err)
			}
		}(conn)
	}

	wg.Wait()
	return nil
}

func (t *TcpProxy) forward(conn net.Conn) error {
	defer conn.Close()
	ip := t.to()
	if ip == "" {
		return errors.New("Target unknown")
	}

	clientConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, t.port))
	if err != nil {
		return err
	}
	defer clientConn.Close()

	_, err = io.Copy(TimeoutConn{clientConn, 2}, TimeoutConn{conn, 2})
	return err
}
