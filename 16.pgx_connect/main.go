package main

import (
	"context"
	"io"
	"log"

	"github.com/devlights/try-golang-db/internal/pgsql"
	"github.com/jackc/pgx/v5"
)

func main() {
	if err := pgsql.Start(io.Discard); err != nil {
		panic(err)
	}
	defer pgsql.Stop()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, pgsql.Dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	log.Println("OK")
}
