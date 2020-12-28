package tangle

import (
	// ビルド時のみ使用する
	"database/sql"
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
	DbConnection, _ := sql.Open("sqlite3", dbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection.Close()

	// blog テーブルの作成
	cmd := `CREATE TABLE IF NOT EXISTS transaction(
             address STRING,    
             value INT,
	     bundle STRING,
	     tag STRING)`

	// cmd を実行
	// _ -> 受け取った結果に対して何もしないので、_ にする
	_, err := DbConnection.Exec(cmd)

	// エラーハンドリング(Go だと大体このパターン)
	if err != nil {
		// Fatalln は便利
		// エラーが発生した場合、以降の処理を実行しない
		log.Fatalln(err)
	}

	// データを挿入(? には、値が入る)
	cmd2 := "INSERT INTO tsc (address, value, bundle, tag) VALUES (?, ?, ?, ?)"
	DbConnection.Exec(cmd2, "hoge1", 1, "hoge2", "hoge3")
}
