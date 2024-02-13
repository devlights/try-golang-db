package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(0)
}

// 09.Columns
//
// *sql.DB.Query() などで *sql.Rows を取得した際に
// カラム名を知りたい場合は *sql.Rows.Columns() を利用する。
//
//	https://pkg.go.dev/database/sql@go1.22.0#Rows.Columns
//
// 結果は []string で返ってくる。
// SELECTで指定した並びで設定されている。
//
// 上記ドキュメントには以下の記載がある。
//
// >Columns returns the column names. Columns returns an error if the rows are closed.
//
// >Columnsはカラム名を返す。Columnsは、行が閉じられている場合はエラーを返します。
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
		return err
	}
	defer db.Close()

	var (
		rows *sql.Rows
	)

	rows, err = db.Query("SELECT * FROM tracks LIMIT 1")
	if err != nil {
		return err
	}
	defer rows.Close()

	var (
		columns []string
	)

	columns, err = rows.Columns()
	if err != nil {
		return err
	}

	log.Println(columns)

	return nil
}
