package server

import (
	"github.com/viile/poker/tools/session"
	"context"
)

type Room interface {
	Join(ctx context.Context, sess *session.Session) error
	Exit(ctx context.Context, sess *session.Session)
	TypeName() string
	OnlineCounts() int
	Handle(ctx context.Context, sess *session.Session, msg string) (err error)
}
