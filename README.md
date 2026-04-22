# try-golang-db

[try-golang](https://github.com/devlights/try-golang) プロジェクトの姉妹版。データベースに関連しているサンプルが配置されています。

## 参考情報

### ブログ記事など

- [Tutorial: Accessing a relational database](https://go.dev/doc/tutorial/database-access)
- [Accessing relational databases](https://go.dev/doc/database/index)
- [Opening a database handle](https://go.dev/doc/database/open-handle)
- [Executing SQL statements that don't return data](https://go.dev/doc/database/change-data)
- [Querying for data](https://go.dev/doc/database/querying)
- [Using prepared statements](https://go.dev/doc/database/prepared-statements)
- [Executing transactions](https://go.dev/doc/database/execute-transactions)
- [Canceling in-progress database operations](https://go.dev/doc/database/cancel-operations)
- [Managing connections](https://go.dev/doc/database/manage-connections)
- [Avoiding SQL injection risk](https://go.dev/doc/database/sql-injection)
- [Golang with database operations](https://dev.to/burrock/golang-with-database-operations-3jl0)
- [Go言語でSQLite3を使う](https://zenn.dev/teasy/articles/go-sqlite3-sample)
- [Go ORMs Compared](https://dev.to/encore/go-orms-compared-2c8g)
- [Go database/sql の操作ガイドあったんかい](https://sourjp.github.io/posts/go-db/)
- [Golangでsqlxを使う](https://zenn.dev/robon/articles/ff2419b7f5a76c)
- [【Go】ORM、Bun について](https://zenn.dev/wasuwa/articles/f691f589da591c)

### ドライバやORMなど

- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
  - Native Driver (Cgoを必要とする)
- [modernc.org/sqlite](https://gitlab.com/cznic/sqlite)
  - Pure Driver (Cgoが必要ない)
- [glebarez/go-sqlite](https://github.com/glebarez/go-sqlite)
  - Pure Driver (modernc.org/sqlite のラッパー。GORM用のドライバなどが追加されている)
- [sqlx](https://github.com/jmoiron/sqlx)
- [sqlc](https://github.com/sqlc-dev/sqlc)
- [uptrace/bun](https://github.com/uptrace/bun)
- [sqlboiler](https://github.com/aarondl/sqlboiler)
  - This package is currently in maintenance mode
- [bob](https://github.com/stephenafamo/bob)
- [textql](https://github.com/dinedal/textql)
- [xlsxsql](https://github.com/noborus/xlsxsql)
- [trdsql](https://github.com/noborus/trdsql)
- [scan](https://github.com/blockloop/scan)
- [gorm](https://github.com/go-gorm/gorm)
- [ent](https://github.com/ent/ent)
- [upper/db](https://github.com/upper/db)
