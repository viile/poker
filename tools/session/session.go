package session

import (
	"context"
)

// Session struct
type Session struct {
	name string

	// 0 offline 1 online
	status int

	*Conn
}

//
func NewSession(ctx context.Context, c *Conn) *Session {
	s := &Session{
		Conn: c,
	}

	return s
}

//
func (s *Session) GetID() string {
	return s.name
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

//
func (s *Session) NeedLogin() bool {
	return s.status != 1
}
