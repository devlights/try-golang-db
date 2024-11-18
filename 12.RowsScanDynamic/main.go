package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	driver     = "sqlite3"
	datasource = "./chinook.db"
)

var (
	db *sql.DB
)

func main() {
	log.SetFlags(0)

	//
	// 12.RowsScanDynamic
	//
	// sql.Query()で取得できる *sql.Rows から、カラム情報を取得する場合
	// 02.Query/main.go にある通り、正しい順序でデータ格納先となる変数を
	// ポインタで渡す必要がある。
	//
	// しかし、ちょっとしたツールなどをパパッと作成したい時に
	// この仕様は若干面倒臭い。そのような場合は以下のようにすると
	// 巷でよく見る「行データがマップになってて、それのリスト」という
	// データ構造で処理することもできる。
	//
	// 本サンプルの処理は [sqlmap](https://github.com/devlights/sqlmap)
	// として公開している。
	//
	if err := run(); err != nil {
		log.Fatal(err)
	}

	/*
	   $ task
	   task: [default] cp -f ../chinook.db .
	   task: [default] go run main.go
	   map[ArtistId:275 Name:Philip Glass Ensemble]
	   map[ArtistId:274 Name:Nash Ensemble]
	   map[ArtistId:273 Name:C. Monteverdi, Nigel Rogers - Chiaroscuro; London Baroque; London Cornett & Sackbu]
	   map[ArtistId:272 Name:Emerson String Quartet]
	   map[ArtistId:271 Name:Mela Tenenbaum, Pro Musica Prague & Richard Kapp]
	*/
}

func run() error {
	//
	// 普通にクエリ発行
	//
	var (
		err error
	)
	if err = open(); err != nil {
		return fmt.Errorf("open(): %w", err)
	}

	const (
		QUERY = "SELECT ArtistId, Name FROM artists ORDER BY ArtistId DESC LIMIT 5"
	)
	var (
		rows *sql.Rows
	)
	if rows, err = db.Query(QUERY); err != nil {
		return fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()

	//
	// マッピング
	//
	var (
		results []map[string]any
	)
	if results, err = mapRows(rows); err != nil {
		return fmt.Errorf("mapRow: %w", err)
	}

	//
	// 表示
	//
	for _, r := range results {
		log.Printf("%v", r)
	}

	return nil
}

func open() error {
	var (
		err error
	)
	if db, err = sql.Open(driver, datasource); err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}

func mapRows(rows *sql.Rows) ([]map[string]any, error) {
	//
	// 結果のカラム名リストを取得
	//
	var (
		cols []string
		err  error
	)
	if cols, err = rows.Columns(); err != nil {
		return nil, err
	}

	//
	// 結果をマッピング
	//
	var (
		results []map[string]any
	)
	for rows.Next() {
		//
		// *sql.Rows.Scan() には、ポインタを渡す必要があるため
		// 予め値の器を用意し、更にそのポインタのリストを構築
		//
		var (
			cv = make([]any, len(cols)) // 各カラム値の格納用
			cp = make([]any, len(cols)) // ポインタリスト
		)
		for i := 0; i < len(cols); i++ {
			cp[i] = &cv[i]
		}

		//
		// 可変長引数の展開演算子（Spread Operator) を使って一気に指定
		//
		if err = rows.Scan(cp...); err != nil {
			return nil, err
		}

		var (
			result = make(map[string]any)
			value  *any
		)
		for i, c := range cols {
			value = cp[i].(*any)
			result[c] = *value
		}

		results = append(results, result)
	}

	return results, nil
}
