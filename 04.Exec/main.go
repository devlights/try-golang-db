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

	/*
	   $ task -d 04.Exec/
	   task: [default] go run main.go
	   LastInsertId: 999       RowsAffected: 1
	   task: [default] echo "SELECT * FROM artists WHERE ArtistId=999" | sqlite3 -header -table ../chinook.db
	   +----------+------+
	   | ArtistId | Name |
	   +----------+------+
	   | 999      | test |
	   +----------+------+
	   task: [default] echo "DELETE FROM artists WHERE ArtistId=999" | sqlite3 -header -table ../chinook.db
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

// 04.Exec
//
// データベースに対してINSERT，UPDATE，DELETEを発行するには *DB.Exec() を利用する。
//
// - https://pkg.go.dev/database/sql@go1.21.6#DB.Exec
//
// *DB.Exec() は、行データを返さず、代わりに sql.Result を返す。
//
// - https://pkg.go.dev/database/sql@go1.21.6#Result
//
// sql.Result からは、最後に追加されたID値と影響を受けた行数が取得できる。
// sql.Result.LastInsertId()は、Auto Incrementな列の場合に利用できる。
//
// 今回は artists テーブルに新たなレコードをINSERTする。
// artistsテーブルのレイアウトは以下。
//
//	$ sqlite3 chinook.db
//	SQLite version 3.37.2 2022-01-06 13:25:41
//	Enter ".help" for usage hints.
//	sqlite> .headers on
//	sqlite> .mode table
//	sqlite> pragma table_info(artists);
//	+-----+----------+---------------+---------+------------+----+
//	| cid |   name   |     type      | notnull | dflt_value | pk |
//	+-----+----------+---------------+---------+------------+----+
//	| 0   | ArtistId | INTEGER       | 1       |            | 1  |
//	| 1   | Name     | NVARCHAR(120) | 0       |            | 0  |
//	+-----+----------+---------------+---------+------------+----+
//	sqlite> SELECT MAX(ArtistId) FROM artists;
//	+---------------+
//	| MAX(ArtistId) |
//	+---------------+
//	| 275           |
//	+---------------+
//
// # REFERENCES
//   - https://go.dev/doc/tutorial/database-access#add_data
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
		result   sql.Result
		lastId   int64
		affected int64
	)

	result, err = db.Exec("INSERT INTO artists (ArtistId, Name) VALUES (?, ?)", 999, "test")
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}

	lastId, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Result.LastInsertId: %w", err)
	}

	affected, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Result.RowsAffected: %w", err)
	}

	log.Printf("LastInsertId: %v\tRowsAffected: %v", lastId, affected)

	return nil
}
