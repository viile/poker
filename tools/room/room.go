package room

import (
	"context"
	"github.com/viile/poker/tools/event"
	"github.com/viile/poker/tools/session"
	"time"
)

type Logic interface {
	Join(ctx context.Context, sess *session.Session)
	Exit(ctx context.Context, sess *session.Session)
	Event(ctx context.Context, sess *session.Session, event *event.Event)
}

type Room struct {
	Index int
	Name  string

	createdAt time.Time
	updatedAt time.Time
	count     uint
	typeName  string

	LogicServer Logic
}

func NewRoom(id int, name string) *Room {
	return &Room{
		Index:     id,
		Name:      name,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}

func (r *Room) Join(ctx context.Context, s *session.Session) {

}

func (r *Room) Exit(ctx context.Context, s *session.Session) {

}

func (r *Room) Event(ctx context.Context, sess *session.Session, event *event.Event) {

}

func (r *Room) OnlineCounts(ctx context.Context, s *session.Session) (c uint) {
	return r.count
}
