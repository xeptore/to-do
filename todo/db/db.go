package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"

	"github.com/xeptore/to-do/todo/db/migration"
)

func Connect(ctx context.Context) (*sql.DB, error) {
	dbConn, err := sql.Open("postgres", os.Getenv("DSN"))
	if nil != err {
		return nil, fmt.Errorf("db: failed to open database connection: %v", err)
	}
	if err := dbConn.PingContext(ctx); nil != err {
		return nil, fmt.Errorf("db: failed to ping database connection: %v", err)
	}
	log.Info().Msg("successfully connected to database")

	goose.SetLogger(goose.NopLogger())
	goose.SetTableName("migrations")
	goose.SetBaseFS(migration.FS)
	if err := goose.SetDialect("postgres"); nil != err {
		return nil, fmt.Errorf("db: failed to set goose dialect to postgres: %v", err)
	}
	if err := goose.Up(dbConn, "scripts"); nil != err {
		return nil, fmt.Errorf("db: failed to execute goose migrations: %v", err)
	}
	log.Info().Msg("executed database migrations")

	return dbConn, nil
}
