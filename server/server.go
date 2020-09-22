package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// Server struct
type Server struct {
	// 维护网络连接
	sessions *sync.Map
	//
	rooms *sync.Map
	//
	relation *sync.Map
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
		rooms: &sync.Map{},
		relation: &sync.Map{},
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
			log.Println("[error]","acceptHandler:",err)
			continue
		}

		go s.connectHandler(ctx, c)
	}
}

func (s *Server) connectHandler(ctx context.Context, c net.Conn) {
	conn := NewConn(c)
	session := NewSession(conn,fmt.Sprintf("%s-%d",c.RemoteAddr().String(),atomic.AddUint32(&s.counter,1)))

	s.sessions.Store(session.GetID(), session)

	connctx, cancel := context.WithCancel(ctx)

	defer func() {
		cancel()
		conn.Close()
		s.sessions.Delete(session.GetID())
		s.relation.Delete(session.GetID())
	}()

	go conn.readCoroutine(connctx)
	go conn.writeCoroutine(connctx)

	session.OnConnect()

	session.conn.SendMessage([]byte(`
欢乐斗地主终端版V1.0
输入[list]查看当前存在的房间
输入[create]创建新房间
输入[join]加入房间,例如 join 123
`))

	for {
		select {
		case err := <-conn.done:
			session.OnDisconnect(err)
			if i,ok := s.relation.Load(session.GetID());ok {
				s.relation.Delete(session.GetID())
				if room,ok := s.rooms.Load(i);ok{
					if _,ok := room.(*Room);ok {
						room.(*Room).Exit(session)
					}
				}
			}
			return
		case msg := <-conn.messageCh:
			m := strings.ToLower(strings.TrimSpace(string(*msg)))
			if len(m) <= 0 {
				continue
			}
			if m == "list" {
				s.rooms.Range(func(key, value interface{}) bool {
					session.conn.SendMessage([]byte(key.(string) + "\n"))
					return true
				})
			} else if m == "create" {
				i := strconv.Itoa(int(atomic.AddUint32(&s.counter, 1)))
				room := NewRoom(i, session.GetID())
				s.rooms.Store(i, room)
				s.relation.Store(session.GetID(), i)
				session.conn.SendMessage([]byte(fmt.Sprintf("房间创建成功,序号:%s\n",i)))
			} else if len(m) >= 6 && m[:5] == "join " {
				var ret bool
				if room, ok := s.rooms.Load(m[5:]); ok {
					if _,ok := room.(*Room);ok{
						ret = true
						if err := room.(*Room).Join(session);err != nil {
							session.conn.SendMessage([]byte(err.Error() + "\n"))
							continue
						}
						s.relation.Store(session.GetID(), m[5:])
						session.conn.SendMessage([]byte("加入成功.\n"))
					}
				}
				if !ret {
					session.conn.SendMessage([]byte("房间不存在!\n"))
				}
			} else {
				if i,ok := s.relation.Load(session.GetID());ok {
					if room,ok := s.rooms.Load(i);ok{
						if _,ok := room.(*Room);ok {
							if err := room.(*Room).Handle(session,m);err != nil {
								session.conn.SendMessage([]byte(err.Error() + "\n"))
								continue
							}
						}
					}
				}
			}
		}
	}
}

