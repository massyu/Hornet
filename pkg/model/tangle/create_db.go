package tangle

import (
	// ビルド時のみ使用する
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// DB Path(相対パスでも大丈夫かと思うが、筆者の場合、絶対パスでないと実行できなかった)
const dbPath = "/home/mash/hornet/db.sql"

// コネクションプールを作成
var DbConnection *sql.DB

// データ格納用
type Transaction struct {
	address string
	value   int
	bundle  string
	tag     string
}

func create_db() {
	// Open(driver,  sql 名(任意の名前))
	log.Println("create_db1")
	//DbConnection, _ := sql.Open("sqlite3", dbPath)
	log.Println("create_db2")

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	//defer DbConnection.Close()
	log.Println("create_db3")
	// データを挿入(? には、値が入る)
	// cmd := "INSERT INTO tsc (address, value, bundle, tag) VALUES (?, ?, ?, ?)"
	// DbConnection.Exec(cmd, "hoge1", 1, "hoge2", "hoge3")
}
