package server

import (
	"sync"
	"time"
)

type Room struct {
	name string
	sessions map[string]*Session
	// 0 未开始, 1 进行中
	status int
	owner string
	wait string
	created time.Time
	sync.Locker
}

func NewRoom(name,owner string) *Room {
	return &Room{
		name :name,
		sessions:make(map[string]*Session,0),
		owner: owner,
		created: time.Now(),
	}
}

// OnConnect .
func (r *Room) Join(name string,sess *Session) {

}

// OnDisconnect .
func (r *Room) Exit(name string) {

}

// OnHandle .
func (r *Room) Handle(name string) {

}
