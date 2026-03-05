package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"modernc.org/sqlite"
)

const (
	driver = "sqlite"
	dsn    = "./chinook.db"
)

const (
	// initPragmaSQL: 接続ごとに毎回適用するPRAGMA群
	//
	// journal_mode=WAL はDBファイルに永続化されるが、
	// RegisterConnectionHook で毎回発行しても副作用はない。
	// busy_timeout / synchronous 等は接続ごとに再設定が必要。
	//
	// SQLiteはデフォルト設定がシングルユーザー・小規模用途向けのため、
	// サーバーサイドやマルチスレッド環境では以下のチューニングが必要となる。
	//
	// 【前提】
	//
	//	本PRAGMA群はDB接続(sql.Open後のPingやExec)のたびに適用すること。
	//	SQLiteのPRAGMAはコネクション単位で有効なものと、DBファイル単位で
	//	永続化されるものが混在するため、接続プールを使う場合は
	//	sql.DB の SetMaxOpenConns(1) またはConnectHook等で確実に適用すること。
	//
	// 【パラメータの説明】
	//
	//	PRAGMA journal_mode=WAL;
	//
	//		ジャーナルモードをWAL(Write-Ahead Logging)に変更する。
	//		デフォルトのDELETEモード(ロールバックジャーナル)と異なり、
	//		書き込みと読み込みが互いをブロックしないため並行性が大幅に向上する。
	//		具体的にはリーダーはWALファイルの古いスナップショットを参照し続けられるため、
	//		WRITER 1本 + READER 複数 の同時アクセスが可能となる。
	//
	//		注意: WALモードはDBファイル単位で永続化される。
	//		一度設定すれば以降は不要だが、他ツールとDBを共有する場合はWALモード対応を確認すること。
	//		また -wal / -shm の2つの補助ファイルが生成される。
	//
	//	PRAGMA synchronous=NORMAL;
	//
	//		fsync(ディスク同期)の頻度を制御する。
	//
	//		FULL  : コミットごとに必ずfsync → 最も安全だがI/Oコストが高い(デフォルト)
	//		NORMAL: チェックポイント時のみfsync → WALモード時は実用上十分な耐障害性を維持しつつ書き込みスループットがFULLの約2〜5倍向上する(公式ドキュメント記載)
	//		OFF   : fsync一切なし → 最速だがOSクラッシュ時にDBが破損するリスクあり
	//
	//		WALモードではNORMALでも電源断以外のクラッシュに対してはACIDを満たすため、
	//		サーバーサイドではNORMALが推奨される組み合わせとなる。
	//
	//	PRAGMA busy_timeout=2000;
	//
	//		他のプロセス/スレッドがロックを保持している場合に待機する最大時間(ミリ秒)。
	//		デフォルト値は0(即座にBUSYエラーを返す)のため、並行書き込みが発生する
	//		環境では必ず設定すること。2000ms(2秒)はWebアプリ等の一般的な推奨値。
	//		sql.DB側のコンテキストタイムアウトより小さい値に設定するのが望ましい。
	//		なお PRAGMA busy_timeout はコネクション単位で有効。
	//
	//	PRAGMA cache_size=-32000;
	//
	//		ページキャッシュのサイズを指定する。
	//		正の値はページ数、負の値はKiB単位での指定となる。
	//		-32000 = 32,000 KiB = 約32MB のメモリをキャッシュに割り当てる。
	//		デフォルトは -2000(約2MB)であり、それに比べ約16倍のキャッシュ容量となる。
	//		頻繁にアクセスするDBが32MB未満であればほぼメモリ上で完結し、
	//		ディスクI/Oを大幅に削減できる。メモリに余裕がある環境での推奨設定。
	//
	//	PRAGMA temp_store=MEMORY;
	//
	//		一時テーブル・インデックス・ソート用ワーク領域の格納先を指定する。
	//
	//		- DEFAULT: デフォルト(コンパイル時設定に依存、多くの場合ファイル)
	//		- FILE   : 常にディスクファイルに書き出す
	//		- MEMORY : 常にメモリ上に展開する
	//
	//		MEMORYを指定することでソートや集計処理の一時領域がディスクI/Oを発生させず、
	//		クエリパフォーマンスが向上する。ただしメモリ使用量が増加するため
	//		大量データを扱うバッチ処理では注意が必要。
	//		なお本設定はコネクション単位で有効。
	//
	//	PRAGMA mmap_size=268435456;
	//
	//		メモリマップI/O(mmap)の上限サイズ(バイト)を指定する。
	//		268435456 = 256MB。
	//		mmapが有効な場合、OSのページキャッシュを直接アドレス空間にマップするため
	//		read()システムコールのオーバーヘッドを削減し、大規模なREAD処理が高速化する。
	//		0を指定するとmmapは無効となる(デフォルト)。
	//		32bitプロセスではアドレス空間の制約からOOMを引き起こす可能性があるため、
	//		64bitプロセス専用の設定として扱うこと。
	//		WALモードとの組み合わせで読み取りパフォーマンスが特に向上する。
	//
	//	PRAGMA wal_autocheckpoint=1000;
	//
	//		WALファイルのページ数がこの値を超えた際に自動チェックポイントを実行する閾値。
	//		チェックポイントとはWALファイルの内容をメインDBファイルに書き戻す処理。
	//		デフォルト値は1000ページ(通常1ページ=4096バイトのため約4MB相当)。
	//		値を大きくすると書き込みスループットが上がるがリカバリ時間が長くなり、
	//		WALファイルが肥大化する。値を小さくするとその逆のトレードオフとなる。
	//		本設定はDBファイル単位で永続化される。
	//		なお高負荷環境では自動チェックポイントを無効化(=0)して
	//		アプリ側で明示的にsqlite3_wal_checkpoint_v2()を呼ぶ設計も検討すること。
	initPragmaSQL = `
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;
PRAGMA busy_timeout=2000;
PRAGMA cache_size=-32000;
PRAGMA temp_store=MEMORY;
PRAGMA mmap_size=268435456;
PRAGMA wal_autocheckpoint=1000;
`
)

func main() {
	log.SetFlags(0)

	var (
		rootCtx     = context.Background()
		ctx, cancel = context.WithCancel(rootCtx)
	)
	defer cancel()

	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(pCtx context.Context) (err error) {
	// 接続がオープンした際のフック関数を登録し、新規接続が払い出される度に確実にPRAGMAが設定されるようにする。
	//
	// フック設定の関数はドライバ毎に多少異なる。
	//
	// - mattn/go-sqlite3  : ConnectHook で sql.Register() でドライバを新規登録する
	// - modernc.org/sqlite: RegisterConnectionHook で グローバル関数として呼ぶ
	sqlite.RegisterConnectionHook(func(conn sqlite.ExecQuerierContext, _ string) error {
		_, err := conn.ExecContext(context.WithoutCancel(pCtx), initPragmaSQL, nil)
		if err != nil {
			log.Printf("PRAGMA setup failed: %v", err)
		} else {
			log.Printf("PRAGMA setup (%s)", initPragmaSQL)
		}

		return err
	})

	var (
		db *sql.DB
	)
	db, err = sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer func() {
		err = db.Close()
		log.Printf("Open Connections=%d", db.Stats().OpenConnections)
	}()

	// 複数のgoroutineが同じconnにアクセスする可能性がある場合は以下も合ったほうが良い。
	// (プロセス内でWriterは同時に1つのみ という制約がSQLite本体の仕様であるため)
	//
	// Writer を1接続に絞る。
	// sql.DB はコネクションプールのため、MaxOpenConns未設定だと
	// 複数のgoroutineが並走してWriter競合 → database locked が多発する。
	// SetMaxOpenConns(1) にすることで sql.DB レベルでシリアライズされる。
	//
	// ※ 複数プロセスからのアクセスはbusy_timeoutの設定でカバーする。
	db.SetMaxOpenConns(1)    // オープン可能接続数は最大1本
	db.SetMaxIdleConns(1)    // アイドル接続が1本 = 使用中接続が最大1本
	db.SetConnMaxLifetime(0) // 接続を使い回す（再接続コスト回避）

	var (
		timeout     = 100 * time.Millisecond
		ctx, cancel = context.WithTimeout(pCtx, timeout)
	)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	log.Printf("Open Connections=%d", db.Stats().OpenConnections)

	return nil
}
