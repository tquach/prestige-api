package trips

import (
	"time"

	pg "gopkg.in/pg.v4"
)

var (
	db *pg.DB
)

func setupTest() {
	db = pg.Connect(pgOptions())
}

func pgOptions() *pg.Options {
	return &pg.Options{
		User:               "prestige",
		Password:           "changeme",
		Database:           "test",
		Addr:               "192.168.99.100:5432",
		DialTimeout:        30 * time.Second,
		ReadTimeout:        10 * time.Second,
		WriteTimeout:       10 * time.Second,
		PoolSize:           10,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        10 * time.Second,
		IdleCheckFrequency: 100 * time.Millisecond,
	}
}
