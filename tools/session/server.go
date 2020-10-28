package session

import (
	"context"
	"github.com/viile/poker/tools/codec"
	"github.com/viile/poker/tools/event"
	"github.com/viile/poker/tools/log"
	"github.com/viile/poker/tools/template"
	"go.uber.org/zap"
	"net"
)

type EventSession struct {
	*event.Event
	*Session
	*Conn
}

var parser = codec.NewCodec()

// Server struct
type Server struct {
	//
	listener net.Listener
	//
	stopCh chan interface{}

	f func(ctx context.Context, e *EventSession)
}

// NewServer create a new socket service
func NewServer(addr string, f func(ctx context.Context, e *EventSession)) (*Server, error) {
	l, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	s := &Server{
		stopCh:   make(chan interface{}),
		listener: l,
		f:        f,
	}

	return s, nil
}

// Run Start socket service
func (s *Server) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		s.listener.Close()
		cancel()
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
			log.GetLogger().Error("acceptHandler", zap.Error(err))
			continue
		}

		go s.connectHandler(ctx, c)
	}
}

func (s *Server) connectHandler(ctx context.Context, c net.Conn) {
	conn := NewConn(c)
	cctx, cancel := context.WithCancel(ctx)

	defer func() {
		cancel()
		conn.Close()
	}()

	go conn.readCoroutine(cctx)
	go conn.writeCoroutine(cctx)

	sess := NewSession(ctx, conn)

	for {
		select {
		case err := <-conn.done:
			log.GetLogger().Error("connectHandler", zap.Error(err))
			return
		case msg := <-conn.messageCh:
			e, err := parser.Decode(msg)
			if err != nil {
				log.GetLogger().Error("connectHandler", zap.Error(err))
				return
			}
			log.GetLogger().Debug("connectHandler", zap.Any("msg", e))

			if sess.NeedLogin() {
				if e.Match(event.CommandLogin) {
					sess.Login(ctx, e.Argv[1])
					template.LoginSuccess.Execute(conn, nil)
					continue
				} else {
					template.NeedLogin.Execute(conn, nil)
					continue
				}
			}

			ee := &EventSession{e, sess, conn}
			s.f(cctx, ee)
		}
	}
}
