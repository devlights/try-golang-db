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

// 07.PreparedQueryInTx
//
// トランザクション (*sql.Tx) からも、Prepared Query を作成することが出来る。
// この場合、その Prepared Query は、当該トランザクションに紐づいた状態となり
// トランザクションの完了（Commit or Rollback）で、自動的にクローズされる。
//
// # REFERENCES
//   - https://go.dev/doc/database/execute-transactions
//   - https://go.dev/doc/database/prepared-statements
//   - https://pkg.go.dev/database/sql@go1.21.6#Tx
//   - https://pkg.go.dev/database/sql@go1.21.6#Tx.Prepare
//   - https://stackoverflow.com/a/25327191
func main() {
	if err := run(); err != nil {
		log.Panic(err)
	}

	/*
	   $ task -d 07.PreparedQueryInTx/
	   task: [default] cp -f ../chinook.db .
	   task: [default] go run main.go
	   id=990  affected=1
	   id=991  affected=1
	   id=992  affected=1
	   id=993  affected=1
	   id=994  affected=1
	   id=995  affected=1
	   id=996  affected=1
	   id=997  affected=1
	   id=998  affected=1
	   id=999  affected=1
	   task: [default] echo "SELECT * FROM artists ORDER BY ArtistId DESC LIMIT 10" | sqlite3 -header -table ./chinook.db
	   +----------+---------+
	   | ArtistId |  Name   |
	   +----------+---------+
	   | 999      | test999 |
	   | 998      | test998 |
	   | 997      | test997 |
	   | 996      | test996 |
	   | 995      | test995 |
	   | 994      | test994 |
	   | 993      | test993 |
	   | 992      | test992 |
	   | 991      | test991 |
	   | 990      | test990 |
	   +----------+---------+
	*/
}

func run() error {
	var (
		db  *sql.DB
		err error
	)

	db, err = sql.Open("sqlite3", "chinook.db")
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer db.Close()

	var (
		tx *sql.Tx
	)

	tx, err = db.Begin()
	if err != nil {
		return fmt.Errorf("db.Begin: %w", err)
	}
	defer tx.Rollback()

	var (
		stmt *sql.Stmt
	)

	stmt, err = tx.Prepare("INSERT INTO artists (ArtistId, Name) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("db.Prepare: %w", err)
	}
	defer stmt.Close() // tx経由で *sql.Stmt を作成した場合、トランザクションと共にクローズされるので無くても良い

	var (
		dropErr = func(v any, _ error) any { return v }
	)

	for i := 990; i < 1000; i++ {
		var (
			rslt sql.Result
		)

		rslt, err = stmt.Exec(i, fmt.Sprintf("test%d", i))
		if err != nil {
			return fmt.Errorf("*sql.Stmt.Exec (in tx): %w", err)
		}

		log.Printf("id=%v\taffected=%v", dropErr(rslt.LastInsertId()), dropErr(rslt.RowsAffected()))
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}
