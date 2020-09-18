package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

// Server struct
type Server struct {
	// network
	sessions *sync.Map
	//
	rooms *sync.Map
	//
	listener net.Listener
	//
	stopCh   chan interface{}
	//
	counter uint32
}

// NewServer create a new socket service
func NewServer(addr string) (*Server, error) {
	l, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	s := &Server{
		sessions: &sync.Map{},
		stopCh:   make(chan interface{}),
		listener: l,
	}

	return s, nil
}

// Run Start socket service
func (s *Server) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
		_ = s.listener.Close()
	}()

	go s.acceptHandler(ctx)

	for {
		select {
		case <-s.stopCh:
			return
		}
	}
}
func (s *Server) stop(ctx context.Context) {
	s.stopCh <- nil
}
func (s *Server) acceptHandler(ctx context.Context) {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go s.connectHandler(ctx, c)
	}
}

func (s *Server) connectHandler(ctx context.Context, c net.Conn) {
	conn := NewConn(c)
	session := NewSession(conn,fmt.Sprintf("%s-%d",c.RemoteAddr().String(),atomic.AddUint32(&s.counter,1)))

	s.sessions.Store(session.GetSessionID(), session)

	connctx, cancel := context.WithCancel(ctx)

	defer func() {
		cancel()
		conn.Close()
		s.sessions.Delete(session.GetSessionID())
	}()

	go conn.readCoroutine(connctx)
	go conn.writeCoroutine(connctx)

	session.OnConnect()

	for {
		select {
		case err := <-conn.done:
			session.OnDisconnect(err)
			return
		case msg := <-conn.messageCh:
			log.Println("[debug]","id:",session.GetSessionID(),"rev:",string(*msg))
			session.OnHandle(msg)
		}
	}
}

