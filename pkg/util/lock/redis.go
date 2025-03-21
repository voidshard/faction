package lock

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/extra/redisotel/v9"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host string
	Port int
}

type RedisLocker struct {
	cfg  *RedisConfig
	pool *goredis.Pool
	sync *redsync.Redsync
}

func NewRedisLocker(cfg *RedisConfig) (*RedisLocker, error) {
	if cfg == nil {
		cfg = &RedisConfig{}
	}
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == 0 {
		cfg.Port = 6379
	}

	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	})
	if err := redisotel.InstrumentTracing(client); err != nil {
		return nil, err
	}
	if err := redisotel.InstrumentMetrics(client); err != nil {
		return nil, err
	}
	pool := goredis.NewPool(client)

	return &RedisLocker{
		cfg:  cfg,
		pool: pool,
		sync: redsync.New(pool),
	}, nil
}

func (r *RedisLocker) Lock(ctx context.Context, key string, ttl time.Duration) error {
	mutex := r.sync.NewMutex(key, redsync.WithExpiry(ttl))
	return mutex.TryLockContext(ctx)
}

func (r *RedisLocker) Unlock(ctx context.Context, key string) error {
	mutex := r.sync.NewMutex(key)
	return mutex.UnlockContext(ctx)
}
