package coordinator

import (
	// ビルド時のみ使用する
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DB Path(相対パスでも大丈夫かと思うが、筆者の場合、絶対パスでないと実行できなかった)
const dbPath = "/home/mash/hornet/db.sql"
const coodbPath = "/home/mash/hornet/coodb.sql"

// コネクションプールを作成
var DbConnection *sql.DB

// データ格納用
type Coomile struct {
	mindex  int
	thash   string
	address string
	bundle  string
	value   int
}

// データ格納用
type Tsc struct {
	bundle  string
	thash   string
	address string
	tag     string
	value   int
}

func createCoodb(txIndex int, txTag string, txBundle string) {
	log.Println("createCoodb開始")
	// Open(driver,  sql 名(任意の名前))
	DbConnection, _ := sql.Open("sqlite3", coodbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection.Close()

	// coodbテーブルがなかったら作成
	/////////////////////////////////////////////////////////////////////////////////////
	// blog テーブルの作成
	cmd := `CREATE TABLE IF NOT EXISTS coomile(
			mindex INT,    
			thash STRING,    
			address STRING,    
			bundle STRING,    
			value INT)`

	// cmd を実行
	// _ -> 受け取った結果に対して何もしないので、_ にする
	_, err := DbConnection.Exec(cmd)

	// エラーハンドリング(Go だと大体このパターン)
	if err != nil {
		// Fatalln は便利
		// エラーが発生した場合、以降の処理を実行しない
		log.Fatalln(err)
	}
	/////////////////////////////////////////////////////////////////////////////////////

	// tscからbundleの値を持つデータをもらってくる
	/////////////////////////////////////////////////////////////////////////////////////
	// Open(driver,  sql 名(任意の名前))
	DbConnection2, _ := sql.Open("sqlite3", dbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection2.Close()

	// 今のbundleを引数に持ってきて問い合わせを行い、valueを読みだす
	cmd2 := "SELECT * FROM tsc where bundle = ?"
	// row = DbConnection.QueryRow(cmd, txHash)
	rows, _ := DbConnection2.Query(cmd2, txBundle)

	defer rows.Close()

	// データ保存領域を確保
	var bg []Tsc
	var nextCount int
	nextCount = 0
	for rows.Next() {
		nextCount++
		log.Println(nextCount)
		log.Println("回目")

		var b Tsc
		// Scan にて、struct のアドレスにデータを入れる
		err := rows.Scan(&b.bundle, &b.thash, &b.address, &b.tag, &b.value)
		// エラーハンドリング(共通関数にした方がいいのかな)
		if err != nil {
			log.Println(err)
		}
		// データ取得
		bg = append(bg, b)
	}

	log.Println("bundleに一致するTransaction一覧を出力")
	var txHash string
	var txAddress string
	var txValue int
	// valueの加算処理
	for _, b := range bg {
		txHash = b.thash
		txAddress = b.address
		txValue = b.value
		cmd = "INSERT INTO coomile (mindex, thash, address, bundle, value) VALUES (?, ?, ?, ?, ?)"
		_, err = DbConnection.Exec(cmd, txIndex, txHash, txAddress, txBundle, txValue)

		if err != nil {
			// golang には、try-catch がない。nil か否かで判定
			log.Fatalln(err)
		}
		fmt.Println(txIndex, b.thash, b.address, b.bundle, b.value)
	}
	/////////////////////////////////////////////////////////////////////////////////////

	log.Println("createCoodb終了")
}

/*
	/////////////////////////////////////////////////////////////////////////////////////
	// ここから挿入したデータの一覧を出力する処理
	// マルチプルセレクト(今度は、_ ではなく、rows)
	cmd3 := "SELECT * FROM coomile where mindex = ?"
	rows, _ := DbConnection3.Query(cmd3, txIndex)

	defer rows.Close()

	// データ保存領域を確保
	var bg []Coomile
	var nextCount int
	nextCount = 0
	for rows.Next() {
		var b Coomile
		// Scan にて、struct のアドレスにデータを入れる
		err := rows.Scan(&b.bundle, &b.thash, &b.address, &b.tag, &b.value)
		// エラーハンドリング(共通関数にした方がいいのかな)
		if err != nil {
			log.Println(err)
		}
		// データ取得
		bg = append(bg, b)
	}

	// valueの加算処理
	for _, b := range bg {
		fmt.Println(b.mindex, b.thash, b.address, b.bundle, b.value)
	}
	log.Println("bundleに一致するTransaction一覧")
	/////////////////////////////////////////////////////////////////////////////////////
*/
