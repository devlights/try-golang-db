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

// 01.Open
//
// データベースを開く (データベースハンドルを取得する) には
//
//   - sql.Open(driverName, datasourceName) *sql.DB
//
// を利用する。
//
// *sql.DB は他の言語と同様にリソースを
// 解放する必要があるため、最終的に Close() を呼び出すようにしておく方がお作法としては良い感じがするが
// 以下の記載にある通り、基本はアプリケーションで一度 *sql.DB を作成したら
// それを使い回すことになる。寿命がアプリケーションのライフサイクルと同じである場合が
// ほとんどなので、明示的に Close() を呼び出すことはあまり無いと思われる。
// (https://pkg.go.dev/database/sql@go1.21.6#DB.Close)
//
// https://pkg.go.dev/database/sql@go1.21.6#DB には以下の記載がある。
//
// > DB is a database handle representing a pool of zero or more underlying connections. It's safe for concurrent use by multiple goroutines.
//
// > (DBは、0個以上の基礎となる接続のプールを表すデータベース・ハンドルです。複数のゴルーチンが同時に使用しても安全です。)
//
// > The sql package creates and frees connections automatically; it also maintains a free pool of idle connections.
//
// > (sqlパッケージは自動的に接続を作成し、解放します。また、アイドル状態の接続の空きプールを維持します。)
//
// つまり、*sql.DB は一つだけあれば良く、ゴルーチンセーフとなっており、コネクションの作成・解放はお任せすることが出来る。
//
// https://pkg.go.dev/database/sql@go1.21.6#Open には以下の記載がある。
//
// > Open may just validate its arguments without creating a connection to the database. To verify that the data source name is valid, call Ping.
//
// > (Open は、データベースへの接続を作成せずに、引数を検証するだけかもしれない。データ・ソース名が有効であることを確認するには、Ping を呼び出します。)
//
// > The returned DB is safe for concurrent use by multiple goroutines and maintains its own pool of idle connections.
// > Thus, the Open function should be called just once. It is rarely necessary to close a DB.
//
// > (返されたDBは、複数のゴルーチンによる同時使用に対して安全であり、アイドル接続のプールを独自に維持する。
// > したがって、Open 関数は一度だけ呼び出す必要があります。DBを閉じる必要はほとんどない。)
//
// db.Ping() は、ネットワークコマンドの ping と同じ意味。
// 接続が有効であるかを確認する際などに利用出来る。
// 接続していない場合は、内部で接続が行われるため、ちゃんと接続できるかどうかのチェックに利用できる。
// (SQLiteを対象に処理している場合、新規データベースを指定している場合は Ping() の呼び出しで生成される)
//
// https://go.dev/doc/tutorial/database-access に以下の記載がある。
//
// > Call DB.Ping to confirm that connecting to the database works. At run time, sql.Open might not immediately connect,
// depending on the driver. You’re using Ping here to confirm that the database/sql package can connect when it needs to.
//
// > (DB.Pingを呼び出して、データベースへの接続がうまくいくことを確認する。
// 実行時に、ドライバによってはsql.Openがすぐに接続しないかもしれません。
// データベース/SQLパッケージが必要なときに接続できることを確認するために、ここでPingを使用しています。)
//
// # REFERENCES
//   - https://go.dev/doc/tutorial/database-access
//   - https://pkg.go.dev/database/sql@go1.21.6#DB
//   - https://github.com/mattn/go-sqlite3
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

const (
	driver     = "sqlite3"
	datasource = "../chinook.db"
)

var (
	db *sql.DB
)

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

	log.Printf("Database Open: driver=%s\tdatasource=%s\n", driver, datasource)

	return nil
}
