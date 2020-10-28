package storage

import (
	"context"
	"github.com/viile/poker/tools/errors"
	"sync"
)

type Memory struct {
	objects map[interface{}]interface{}
	locks   map[interface{}]sync.Locker
}

func NewMemoryStorage() *Memory {
	return &Memory{
		objects: make(map[interface{}]interface{}, 0),
		locks:   make(map[interface{}]sync.Locker, 0),
	}
}

func (m *Memory) lock(key interface{}) func() {
	val, ok := m.locks[key]
	if !ok {
		val = &sync.Mutex{}
		m.locks[key] = val
	}
	val.Lock()
	return func() {
		val.Unlock()
	}
}

func (m *Memory) Read(ctx context.Context, key interface{}) (val interface{}, err error) {
	var ok bool
	if val, ok = m.objects[key]; !ok {
		err = errors.ErrObjectNotFound
		return
	}

	return
}

func (m *Memory) Write(ctx context.Context, key, val interface{}) (err error) {
	defer m.lock(key)()
	m.objects[key] = val
	return
}

func (m *Memory) Delete(ctx context.Context, key interface{}) (err error) {
	defer m.lock(key)()
	delete(m.objects, key)
	delete(m.locks, key)
	return
}

func (m *Memory) Count(ctx context.Context) (i int, err error) {
	i = len(m.objects)
	return
}

func (m *Memory) List(ctx context.Context, offset, limit int) (objects []interface{}, err error) {
	objects = make([]interface{}, 0)
	o := 0
	for _, v := range m.objects {
		if o >= offset {
			objects = append(objects, v)
		}
		if len(objects) >= limit {
			return
		}
		o++
	}
	return
}

var _ Storage = NewMemoryStorage()
