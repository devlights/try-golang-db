package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(0)
}

// 06.PreparedQuery
//
// Prepared Queryを利用する場合は、*sql.DB.Prepare() を使う。
// *sql.DB.Prepare() は、*sql.Stmt を返し、これにパラメータを指定することにより実行できる。
// *sql.Stmt は、利用が終わったら Close() を呼び出す必要がある。
//
// *sql.Stmt は、複数のgoroutineにて同時に利用することが出来る。
//
// > Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned statement.
// The caller must call the statement's Close method when the statement is no longer needed.
//
// > Prepareは、後のクエリーや実行のために準備されたステートメントを作成します。
// 返されたステートメントから複数のクエリや実行を同時に実行することができます。
// ステートメントが不要になったら、呼び出し元はステートメントの Close メソッドを呼び出さなければなりません。
//
// # REFERENCES
//   - https://go.dev/doc/database/prepared-statements
//   - https://pkg.go.dev/database/sql@go1.21.6#DB.Prepare
//   - https://stackoverflow.com/a/25327191
func main() {
	if err := run(); err != nil {
		log.Panic(err)
	}

	/*
	   $ task -d 06.PreparedQuery/
	   task: [default] cp -f ../chinook.db .
	   task: [default] go run main.go
	   id=10   name=Billy Cobham
	   id=8    name=Audioslave
	   id=9    name=BackBeat
	   id=6    name=Antônio Carlos Jobim
	   id=1    name=AC/DC
	   id=3    name=Aerosmith
	   id=7    name=Apocalyptica
	   id=4    name=Alanis Morissette
	   id=2    name=Accept
	   id=5    name=Alice In Chains
	*/
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
		stmt *sql.Stmt
	)

	stmt, err = db.Prepare("SELECT * FROM artists WHERE ArtistId = ?")
	if err != nil {
		return fmt.Errorf("db.Prepare: %w", err)
	}
	defer stmt.Close()

	const LOOP_COUNT = 10
	var (
		wg    sync.WaitGroup
		errCh = make(chan error, LOOP_COUNT)
	)

	wg.Add(LOOP_COUNT)

	for i := 0; i < LOOP_COUNT; i++ {
		go func(i int) {
			defer wg.Done()

			var (
				row  *sql.Row
				id   int
				name string
			)

			row = stmt.QueryRow(i)
			err = row.Scan(&id, &name)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					errCh <- fmt.Errorf("ErrNoRows: %d", i)
					return
				}

				errCh <- fmt.Errorf("sql.Row.Scan: %w", err)
				return
			}

			log.Printf("id=%v\tname=%v", id, name)
		}(i + 1)
	}

	wg.Wait()
	close(errCh)

	for e := range errCh {
		return e
	}

	return nil
}
