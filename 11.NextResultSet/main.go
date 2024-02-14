package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/devlights/sqlmap"
	"github.com/k0kubun/pp/v3"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(0)
}

// 11.NextResultSet
//
// データベースドライバ側の実装が行われている場合は
// 一度に複数の結果セットを取得することができる。
// その場合、*sql.Rows.NextResultSet() を呼び出し
// 次のカーソルに進める。
//
//   - https://pkg.go.dev/database/sql@go1.22.0#example-DB.Query-MultipleResultSets
//
// 本リポジトリで利用している github.com/mattn/go-sqlite3 は、複数の結果セットに
// 対応していないことが作者さんのブログ記事にて記載されている。
//
//   - https://mattn.kaoriya.net/software/lang/go/20161106232834.htm
//
// 以下、作者さんのブログ記事より引用。
//
// > go-sqlite3 もサポートする予定でしたが、sqlite3 の複数クエリ実行は golang の複数結果セットが期待する物と異なる為、現状は実装を見送りました。
//
// そのため、以下は実装のみを記載しておく。
//
// 恐らく、PostgreSQLとMySQLのドライバは対応しているので可能なはず。
// (MySQLは、multiStatements=trueの指定が必要)
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
		query string
		sb    strings.Builder
	)

	sb.WriteString("SELECT ArtistId,Name FROM artists LIMIT 2;")
	sb.WriteString("SELECT TrackId,Name FROM tracks LIMIT 2;")
	query = sb.String()

	var (
		rows *sql.Rows
	)

	rows, err = db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var (
		m []map[string]any
	)

	//
	// First Result
	//
	m, err = sqlmap.MapRows(rows)
	if err != nil {
		return err
	}
	pp.Println(m)

	//
	// 次の結果セットへ
	//
	if !rows.NextResultSet() {
		return fmt.Errorf("rows.NextResultSet() returns false")
	}

	//
	// Second Result
	//
	m, err = sqlmap.MapRows(rows)
	if err != nil {
		return err
	}
	pp.Println(m)

	return nil
}
