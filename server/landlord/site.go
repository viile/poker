package landlord

import (
	"github.com/viile/poker/tools/session"
	"context"
)

type Site struct {
	Index  int
	Name   string
	Sess   *session.Session
	Cards  Cards
	Rob    bool
	Status int
}

func (s *Site) bind(ctx context.Context, sess *session.Session) {
	s.Sess = sess
}

func (s *Site) unbind(ctx context.Context, sess *session.Session) {
	s.Name = ""
	s.Sess = nil
}

