package session

import (
	"context"
	"fmt"
	"time"
)

// Session struct
type Session struct {
	id string
	name string

	status int

	currRoomID string

	*Conn
}

//
func NewSession(ctx context.Context, c *Conn,i uint32) *Session {
	s := &Session{
		id : fmt.Sprintf("%d%d",time.Now().UnixNano(),i),
		Conn: c,
	}

	return s
}

//
func (s *Session) GetID() string {
	return s.id
}

//
func (s *Session) GetName(ctx context.Context) string {
	return s.name
}

// Logout .
func (s *Session) Logout(ctx context.Context) {
	s.status = 0
}

// Login .
func (s *Session) Login(ctx context.Context, name string) {
	s.name = name
	s.status = 1
}

//  .
func (s *Session) BindRoom(ctx context.Context, id string) {
	s.currRoomID = id
}

func (s *Session) GetRoomID(ctx context.Context) string{
	return s.currRoomID
}

//
func (s *Session) NeedLogin() bool {
	return s.status != 1
}
