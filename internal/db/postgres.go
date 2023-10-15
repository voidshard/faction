package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/voidshard/faction/pkg/config"
)

const (
	retryFailConnect = 5
)

type Postgres struct {
	*sqlDB
}

func NewPostgres(cfg *config.Database) (*Postgres, error) {
	var db *sqlx.DB
	var err error
	for i := 1; i < retryFailConnect+1; i++ {
		db, err = sqlx.Connect("postgres", cfg.Location)
		if err != nil {
			time.Sleep(time.Duration(i*i) * time.Second)
			continue
		}
		me := &Postgres{&sqlDB{conn: db}}
		return me, err
	}
	return nil, err
}
