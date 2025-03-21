package controller

import (
	"context"
	"time"
)

type nullLocker struct{}

func (n *nullLocker) Lock(ctx context.Context, key string, ttl time.Duration) error {
	return nil
}

func (n *nullLocker) Unlock(ctx context.Context, key string) error {
	return nil
}
