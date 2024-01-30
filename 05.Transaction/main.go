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

// 05.Transaction
//
// トランザクションを開始する場合、 *sql.DB.Begin() を利用する。
// トランザクションは *sql.Tx で表される。
//
// 基本的な使い方は、他の言語と同様で
//
//   - *sql.Tx.Query()
//   - *sql.Tx.QueryRow()
//   - *sql.Tx.Exec()
//   - *sql.Tx.Rollback()
//   - *sql.Tx.Commit()
//
// を用いてトランザクションを操作する。
//
// 定型文として、トランザクションを開始したら
//
//	defer tx.Rollback()
//
// を呼び出しておく。これにより、エラー発生などに
// ロールバックが行われる。（コミットした後のロールバックは何も影響しない）
//
// https://go.dev/doc/database/execute-transactions に `Best Practice` として以下が記載されている。
//
// > Use the APIs described in this section to manage transactions.
// Do not use transaction-related SQL statements such as BEGIN and COMMIT directly—doing so can leave your database in an unpredictable state,
// especially in concurrent programs.
//
// > トランザクションを管理するには、このセクションで説明するAPIを使用してください。
// BEGINやCOMMITのようなトランザクション関連のSQL文を直接使用しないでください。
//
// > When using a transaction, take care not to call the non-transaction sql.DB methods directly, too, as those will execute outside the transaction,
// giving your code an inconsistent view of the state of the database or even causing deadlocks.
//
// > トランザクションを使用する場合、トランザクション以外のsql.DBメソッドも直接呼び出さないように注意してください。
// これらのメソッドはトランザクションの外で実行されるため、
// コードの中でデータベースの状態に一貫性がなくなったり、デッドロックの原因になったりします。
//
// https://pkg.go.dev/database/sql@go1.21.6#Tx には、以下の記載がある。
//
// > A transaction must end with a call to Commit or Rollback.
//
// > トランザクションは必ず Commit もしくは Rollback で完了する必要があります。
//
// > After a call to Commit or Rollback, all operations on the transaction fail with ErrTxDone.
//
// > コミットまたはロールバックを呼び出した後、トランザクションに対するすべての操作は ErrTxDone で失敗する。
//
// > The statements prepared for a transaction by calling the transaction's Prepare or Stmt methods are closed by the call to Commit or Rollback.
//
// > トランザクションのPrepareメソッドまたはStmtメソッドを呼び出してトランザクションに準備されたステートメントは、CommitまたはRollbackの呼び出しによって閉じられます。
//
// # REFERENCES
//   - https://go.dev/doc/database/execute-transactions
//   - https://pkg.go.dev/database/sql@go1.21.6#DB.Begin
//   - https://pkg.go.dev/database/sql@go1.21.6#Tx
//   - https://stackoverflow.com/a/25327191
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
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("db.Ping: %w", err)
	}

	var (
		tx *sql.Tx
	)

	tx, err = db.Begin()
	if err != nil {
		return fmt.Errorf("db.Begin: %w", err)
	}
	defer tx.Rollback()

	for i := 990; i < 1000; i++ {
		_, err = tx.Exec("INSERT INTO artists (ArtistId, Name) VALUES (?, ?)", i, fmt.Sprintf("test%d", i))
		if err != nil {
			return fmt.Errorf("tx.Exec: %w (%d)", err, i)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}
