package celeritas

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func (c *Celeritas) OpenDB(cfg databaseConfig) (*sql.DB, error) {

	var dbType string

	if cfg.dbType == "postgres" || cfg.dbType == "postgresql" {
		dbType = "pgx"		
	}

	db, err := sql.Open(dbType, cfg.dsn)
	if err != nil{
		return nil,err
	}

	db.SetMaxOpenConns(cfg.maxOpenConns)
	db.SetMaxIdleConns(cfg.maxIdleConns)

	duration, err := time.ParseDuration(cfg.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil{
		return nil,err
	}

	return db,nil
}