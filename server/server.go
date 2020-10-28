package server

import (
	"context"
	"fmt"
	"github.com/viile/poker/server/landlord"
	"github.com/viile/poker/tools/event"
	"github.com/viile/poker/tools/log"
	"github.com/viile/poker/tools/session"
	"github.com/viile/poker/tools/errors"
	"github.com/viile/poker/tools/storage"
	"github.com/viile/poker/tools/template"
	"go.uber.org/zap"
	"strconv"
	"sync/atomic"
)

// Server struct
type Server struct {
	//
	sessions storage.Storage
	//
	rooms storage.Storage
	//
	listener *session.Server
	//
	counter uint32
}

func NewServer(addr string) (*Server, error) {
	s := &Server{
		sessions: storage.NewMemoryStorage(),
		rooms:    storage.NewMemoryStorage(),
	}

	l, err := session.NewServer(addr, s.Handle)
	if err != nil {
		return nil, err
	}
	s.listener = l

	return s, nil
}

func (s *Server) getRoom(ctx context.Context,id string) (room Room,err error) {
	var index int
	if index,err = strconv.Atoi(id);err != nil {
		return
	}

	var r interface{}
	if r,err = s.rooms.Read(ctx,index);err != nil {
		return
	}

	var ok bool
	if room,ok = r.(Room);!ok {
		err = errors.ErrRoomNotExist
		return
	}

	return
}

func (s *Server) Handle(ctx context.Context, e *session.EventSession) (err error) {
	log.GetLogger().Info("Handle", zap.Any("e", e))
	if e.Match(event.CommandList) {
		var objects []interface{}
		if objects, err = s.rooms.List(ctx, 0, 20);err != nil {
			log.GetLogger().Error("Handle", zap.Error(err))
			return
		}
		return template.RoomList.Execute(e.Conn, objects)
	} else if e.Match(event.CommandCreate) {
		i := atomic.AddUint32(&s.counter, 1)
		logic := landlord.NewRoom(ctx, int(i), fmt.Sprintf("%s的斗地主房间", e.Session.GetName(ctx)), e.Session.GetName(ctx))

		if err = s.rooms.Write(ctx, i, logic); err != nil {
			return
		}

		return e.SendMsg("房间创建成功\n")
	} else if e.Match(event.CommandJoin) {
		var room Room
		if room,err = s.getRoom(ctx,e.Argv[1]);err != nil {
			return
		}

		if err = room.Join(ctx,e.Session);err != nil {
			return
		}

	} else {

	}

	return
}

func (s *Server) Run() {
	s.listener.Run()
}

/*
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
	conn := network.NewConn(c)
	session := session2.NewSession(conn,fmt.Sprintf("%s-%d",c.RemoteAddr().String(),atomic.AddUint32(&s.counter,1)))

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

	session.Login()

	session.conn.SendMessage([]byte(`
欢乐斗地主终端版V1.0
输入[list]查看当前存在的房间
输入[create]创建新房间
输入[join]加入房间,例如 join 123
`))

	for {
		select {
		case err := <-conn.done:
			session.Logout(err)
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
*/
