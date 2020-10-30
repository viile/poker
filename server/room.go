package server

import (
	"github.com/viile/poker/tools/event"
	"github.com/viile/poker/tools/session"
	"context"
)

type Room interface {
	Join(ctx context.Context, s *session.Session) error
	Exit(ctx context.Context, s *session.Session)
	TypeName() string
	OnlineCounts() int
	Handle(ctx context.Context, s *session.Session, e *event.Event) (err error)
}
