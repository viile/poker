package server

import (
	"log"
)

// Session struct
type Session struct {
	id       string
	conn     *Conn
	settings map[string]interface{}

}

// NewSession create a new session
func NewSession(conn *Conn,id string) *Session {
	session := &Session{
		id:      id,
		conn:     conn,
		settings: make(map[string]interface{}),

	}

	return session
}

// GetID get session ID
func (s *Session) GetID() string {
	return s.id
}

// OnDisconnect .
func (s *Session) OnDisconnect(err error) {
	log.Println("[event]",s.id , " lost.",err)
}
// OnConnect .
func (s *Session) OnConnect() {
	log.Println("[event]",s.id , " connected.")
}

// OnHandle .
func (s *Session) OnHandle(buf *[]byte) {
	_ = s.conn.SendMessage([]byte("长度test错误️"))
}
