package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(0)
}

// 08.Conn
//
// db.Query()などの関数を利用すると、内部のコネクションプールから
// 任意のコネクションが利用される。単一の接続を取得して処理を行いたい場合は
// db.Conn()を利用して、コネクションを取得する。
//
// >Conn returns a single connection by either opening a new connection or returning an existing connection from the connection pool.
// Conn will block until either a connection is returned or ctx is canceled.
// Queries run on the same Conn will be run in the same database session.
//
// >Connは、新しい接続をオープンするか、接続プールから既存の接続を返すことによって、単一の接続を返します。
// Connは、接続が返されるかctxがキャンセルされるまでブロックします。
// 同じConnで実行されるクエリは、同じデータベース・セッションで実行されます。
//
// 取得したコネクションは Close メソッドを呼び出してリソースを解放する必要がある。
//
// クエリの発行方法などは、*sql.DB と同じである。
func main() {
	if err := run(); err != nil {
		log.Panic(err)
	}

	/*
	   $ task -d 08.Conn/
	   task: [default] cp -f ../chinook.db .
	   task: [default] go run main.go
	   AC/DC
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

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("db.Ping: %w", err)
	}

	var (
		ctx  = context.Background()
		conn *sql.Conn
	)

	conn, err = db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("db.Conn: %w", err)
	}
	defer conn.Close()

	err = conn.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("conn.PingContext: %w", err)
	}

	var (
		row  *sql.Row
		name string
	)

	row = conn.QueryRowContext(ctx, "SELECT Name from artists")
	err = row.Scan(&name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("conn.QueryContext: %w", err)
	}

	log.Println(name)

	return nil
}
