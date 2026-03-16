package main

import (
	"database/sql"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "embed"

	embedpsql "github.com/fergusstrange/embedded-postgres"
	_ "github.com/lib/pq"
)

const (
	dsn = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
)

var (
	//go:embed northwind.sql
	northwindSQL string // ビルド時にSQLファイルを埋め込む
)

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
}

func main() {
	//
	// fergusstrange/embedded-postgres は、PostgreSQLを埋め込みで利用できるようにしてくれるライブラリ。
	// テスト用のデータベースで確認したい場合などにとても便利。
	//
	// データベースは起動する度にまっさらな状態で起動する。
	//
	// ダウンロードされた PostgreSQL は ~/.embedded-postgres-go に配置される。
	//
	log.Println("=> START")
	defer func() { log.Println("=> END") }()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	//
	// embedded postgres の設定と起動
	//
	// ログバッファ
	var (
		logBuf = io.Discard
	)
	if os.Getenv("PG_DEBUG") != "" {
		logBuf = os.Stdout
	}

	// 設定
	var (
		conf embedpsql.Config
	)
	conf = embedpsql.DefaultConfig().
		Username("postgres").
		Password("postgres").
		Database("postgres").
		Port(5432).
		Version(embedpsql.V18).
		StartTimeout(30 * time.Second).
		Logger(logBuf)

	// データベース起動
	var (
		pg = embedpsql.NewDatabase(conf)
	)
	if err := pg.Start(); err != nil {
		return err
	}

	// シグナルハンドリング
	var (
		sigCh = make(chan os.Signal, 1)
	)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("\n==> Shutting down...")

		if err := pg.Stop(); err != nil {
			log.Printf("pg.Stop error: %v", err)
		}

		os.Exit(0)
	}()

	defer func() {
		if err := pg.Stop(); err != nil {
			log.Printf("==> pg.Stop error: %v", err)
		}

		log.Println("==> embedded-postgres stopped")
	}()

	log.Println("==> embedded-postgres started")

	//
	// database/sql で接続
	//
	var (
		db  *sql.DB
		err error
	)
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}
	log.Println("==> ping OK")

	//
	// スキーマとデータ投入
	//
	if _, err = db.Exec(northwindSQL); err != nil {
		return err
	}

	//
	// クエリ発行
	//
	const (
		query = `
				SELECT ship_country, COUNT(*) AS order_count
				FROM orders
				GROUP BY ship_country
				ORDER BY order_count DESC
				LIMIT 10
                `
	)
	var (
		rows *sql.Rows
	)
	if rows, err = db.Query(query); err != nil {
		return err
	}
	defer rows.Close()

	var (
		country string
		count   int
	)
	for rows.Next() {
		rows.Scan(&country, &count)
		log.Printf("%-20s %d", country, count)
	}

	return nil
}
