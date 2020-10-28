package storage

import (
	"context"
)

type Storage interface {
	Read(ctx context.Context, key interface{}) (interface{}, error)
	Write(ctx context.Context, key, val interface{}) error
	Delete(ctx context.Context, key interface{}) error
	Count(ctx context.Context) (i int, err error)
	List(ctx context.Context, offset, limit int) (objects []interface{}, err error)
}
