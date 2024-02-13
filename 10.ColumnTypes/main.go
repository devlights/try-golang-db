package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(0)
}

// 10.ColumnTypes
//
// *sql.DB.Query() などで *sql.Rows を取得した際に
// カラム名を知りたい場合は *sql.Rows.Columns() を利用するが
// *sql.Rows.Columns() は、カラム名だけしか取得できない。
//
// カラムに対する補足情報を取得するには *sql.Rows.ColumnTypes() を利用する。
//
//	https://pkg.go.dev/database/sql@go1.22.0#Rows.ColumnTypes
//
// 結果は []*sql.ColumnType で返ってくる。
// SELECTで指定した並びで設定されている。
//
// 上記ドキュメントには以下の記載がある。
//
// >ColumnTypes returns column information such as column type, length, and nullable.
// Some information may not be available from some drivers.
//
// >ColumnTypesは、カラム・タイプ、長さ、Nullableなどのカラム情報を返す。
// ドライバによっては利用できない情報もあります。
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
		cols []*sql.ColumnType
		toI  = func(v int64, ok bool) int64 {
			if !ok {
				return 0
			}

			return v
		}
		toB = func(v bool, ok bool) bool {
			if !ok {
				return false
			}

			return v
		}
	)

	cols, err = rows.ColumnTypes()
	if err != nil {
		return err
	}

	for _, c := range cols {
		log.Printf(
			"NAME=%-15s\tTYPE=%-20s\tLENGTH=%-10d\tNOTNULL=%-10v\tSCAN TYPE=%v",
			c.Name(),
			c.DatabaseTypeName(),
			toI(c.Length()),
			toB(c.Nullable()),
			c.ScanType())
	}

	return nil
}
