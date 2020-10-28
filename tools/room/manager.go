package room

import (
	"context"
	"github.com/viile/poker/tools/errors"
	"github.com/viile/poker/tools/session"
	"github.com/viile/poker/tools/storage"
)

type Manager struct {
	se3s storage.Storage
}

func NewManager() *Manager {
	return &Manager{
		se3s: storage.NewMemoryStorage(),
	}
}

func (m *Manager) Join(ctx context.Context, id int, s *session.Session) (err error) {
	var val interface{}
	if val, err = m.se3s.Read(ctx, id); err != nil {
		err = errors.ErrRoomNotExist.With(err)
		return
	}

	r, ok := val.(*Room)
	if !ok {
		err = errors.ErrType
		return
	}

	r.Join(ctx, s)
	return
}

func (m *Manager) Exit(ctx context.Context, id int, s *session.Session) (err error) {
	return
}

func (m *Manager) Create(ctx context.Context, id int, t string, s *session.Session) (err error) {
	return
}

func (m *Manager) List(ctx context.Context, s *session.Session) (err error) {
	return
}
