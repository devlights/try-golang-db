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
	   $ task
	   task: [default] go run main.go
	   id=275, name=Philip Glass Ensemble
	   id=274, name=Nash Ensemble
	   id=273, name=C. Monteverdi, Nigel Rogers - Chiaroscuro; London Baroque; London Cornett & Sackbu
	   id=272, name=Emerson String Quartet
	   id=271, name=Mela Tenenbaum, Pro Musica Prague & Richard Kapp
	*/

}

const (
	driver     = "sqlite3"
	datasource = "../chinook.db"
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

// 02.Query
//
// クエリを発行し結果を取得するには DB.Query() を利用する。
//
// - https://pkg.go.dev/database/sql@go1.21.6#DB.Query
//
// 結果は *sql.Rows で返ってくる。
// *sql.Rows は、Close() メソッドを呼び出してリソースを開放する必要がある。
// (*sql.DBと異なり、こちらはClose()をちゃんと呼び出して解放しないと駄目)
//
// *sql.Rows は、イテレータのような作りとなっており
// 取得時のカーソルの位置は「最初の行データの前」となっている。
// そのため、データを取得するには、まず *sql.Rows.Next() を呼び出す必要がある。
// 基本的に複数行データを取得しているため、ループ処理をすることになる。
// なので、ループの継続条件に *sql.Rows.Next() を使うのが一般的。
//
// 各行のデータ取得には、*sql.Rows.Scan() を利用する。
// 少しクセのある構造になっており、構造体をポイっと渡して値を設定してくれるようにはなっていない。
// 引数にクエリで取得対象としたカラムと同じ並びでポインタを渡していく必要がある。
// (https://pkg.go.dev/database/sql@go1.21.6#example-DB.Query-MultipleResultSets)
//
// ORM系ライブラリを利用すれば、構造体に一気に設定出来るが
// ちょっと利用するためだけであれば、https://github.com/devlights/sqlmap も手軽に使える。
// こちらは []map[string]any の形で一気に読み取って取得するようになっている。
//
// *sql.Rows には、Err() というメソッドが用意されており、イテレーション中に発生したエラーを返す。
// なので、このメソッドを読み取りが終わった後に確認するのも定型文のように行われる。
// (https://pkg.go.dev/database/sql@go1.21.6#Rows.Err)
func run() error {
	var (
		err error
	)

	err = open()
	if err != nil {
		return fmt.Errorf("open(): %w", err)
	}

	var (
		rows *sql.Rows
	)

	rows, err = db.Query("SELECT ArtistId, Name FROM artists ORDER BY ArtistId DESC LIMIT 5")
	if err != nil {
		return fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			artist Artist
		)

		err = rows.Scan(&artist.Id, &artist.Name) // ポインタで渡す必要がある点に注意
		if err != nil {
			return fmt.Errorf("rows.Scan: %w", err)
		}

		log.Printf("id=%v, name=%v", artist.Id, artist.Name)
	}

	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows.Err: %w", err)
	}

	return nil
}

func open() error {
	var (
		err error
	)

	db, err = sql.Open(driver, datasource)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}
