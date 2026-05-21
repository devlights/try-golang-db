package pgsql

import (
	"io"
	"time"

	embedpsql "github.com/fergusstrange/embedded-postgres"
)

const (
	Dsn = "postgres://postgres:postgres@localhost:5555/postgres?sslmode=disable"
)

var (
	pg *embedpsql.EmbeddedPostgres
)

func Start(logger io.Writer) error {
	if pg != nil {
		if err := Stop(); err != nil {
			return err
		}
	}

	conf := embedpsql.DefaultConfig().
		Username("postgres").
		Password("postgres").
		Database("postgres").
		Port(5555).
		Version(embedpsql.V18).
		StartTimeout(30 * time.Second).
		Logger(logger)

	pg = embedpsql.NewDatabase(conf)

	if err := pg.Start(); err != nil {
		return err
	}

	return nil
}

func Stop() error {
	if pg == nil {
		return nil
	}

	if err := pg.Stop(); err != nil {
		return err
	}

	pg = nil

	return nil
}
