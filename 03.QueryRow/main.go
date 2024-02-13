package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(0)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

	/*
	   $ task -d 03.QueryRow/
	   task: [default] go run main.go
	   id=275, name=Philip Glass Ensemble
	*/
}

const (
	driver     = "sqlite3"
	datasource = "./chinook.db"
)

type (
	Artist struct {
		Id   int
		Name string
	}
)

var (
	db *sql.DB
)

// 03.QueryRow
//
// クエリを発行し結果を１件取得するには DB.QueryRow() を利用する。
//
// - https://pkg.go.dev/database/sql@go1.21.6#DB.QueryRow
//
// 結果は *sql.Row で返ってくる。結果セットが複数行となっている場合、１行目が返ってくる。
// DB.QueryRowのドキュメントには以下のように記載がある。
//
// > QueryRow executes a query that is expected to return at most one row.
//
// > (QueryRow は、最大でも1行を返すと予想されるクエリを実行します。)
//
// > QueryRow always returns a non-nil value.
//
// > (QueryRowは、常に non-nil な値を返します。)
//
// > Errors are deferred until Row's Scan method is called. If the query selects no rows, the *Row's Scan will return ErrNoRows.
//
// > (エラー発生は *sql.Row.Scan() が呼ばれるまで遅延されます。行が存在しない場合、*sql.Row.Scan() は、sql.ErrNoRows を返します。)
//
// # REFERENCES
//   - https://go.dev/doc/tutorial/database-access
//   - https://pkg.go.dev/database/sql@go1.21.6#DB
func run() error {
	var (
		err error
	)

	db, err = sql.Open(driver, datasource)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("db.Ping: %w", err)
	}

	var (
		artist Artist
		row    *sql.Row
	)

	row = db.QueryRow("SELECT ArtistId, Name FROM artists ORDER BY ArtistId DESC LIMIT 5")
	err = row.Scan(&artist.Id, &artist.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("NOT FOUND: %w", err)
		}

		return fmt.Errorf("row.Scan: %w", err)
	}

	log.Printf("id=%v, name=%v", artist.Id, artist.Name)

	return nil
}
