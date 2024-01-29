package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(0)
}

func main() {
	if err := run(); err != nil {
		log.Panic(err)
	}
}

func run() error {
	var (
		db  *sql.DB
		err error
	)

	db, err = sql.Open("sqlite3", "./chinook.db")
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("db.Ping: %w", err)
	}

	var (
		tx *sql.Tx
	)

	tx, err = db.Begin()
	if err != nil {
		return fmt.Errorf("db.Begin: %w", err)
	}
	defer tx.Rollback()

	for i := 990; i < 1000; i++ {
		_, err = tx.Exec("INSERT INTO artists (ArtistId, Name) VALUES (?, ?)", i, fmt.Sprintf("test%d", i))
		if err != nil {
			return fmt.Errorf("tx.Exec: %w (%d)", err, i)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}
