package tangle

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
type Tsc struct {
	bundle  string
	address string
	tag     string
	value   int
}

// データ格納用
type Coomile struct {
	mindex int
	tag    string
}

func create_db(txBundle string, txAddress string, txTag string, txValue string) {
	log.Println("create_db開始")
	// Open(driver,  sql 名(任意の名前))
	DbConnection, _ := sql.Open("sqlite3", dbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection.Close()

	// blog テーブルの作成
	cmd := `CREATE TABLE IF NOT EXISTS tsc(
             bundle STRING,    
             address STRING,    
             tag STRING,    
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

	cmd = "INSERT INTO tsc (bundle, address, tag, value) VALUES (?, ?, ?, ?)"
	_, err = DbConnection.Exec(cmd, txBundle, txAddress, txTag, txValue)

	if err != nil {
		// golang には、try-catch がない。nil か否かで判定
		log.Fatalln(err)
	}

	//ここから挿入したデータの一覧を出力する処理
	// マルチプルセレクト(今度は、_ ではなく、rows)
	cmd = "SELECT * FROM tsc where address = ?"
	row := DbConnection.QueryRow(cmd, txAddress)

	// データ保存領域を確保
	var b Tsc
	// Scan にて、struct のアドレスにデータを入れる
	err = row.Scan(&b.bundle, &b.address, &b.tag, &b.value)
	// エラーハンドリング(共通関数にした方がいいのかな)
	if err != nil {
		// シングルセレクトの場合は、エラーハンドリングが異なる
		if err == sql.ErrNoRows {
			log.Println("There is no row!!!")
		} else {
			log.Println(err)
		}
	}

	fmt.Println(b.bundle, b.address, b.tag, b.value)

	log.Println("create_db")
}

func check_db(txBundle string) {
	log.Println("check_db開始")
	// Open(driver,  sql 名(任意の名前))
	DbConnection, _ := sql.Open("sqlite3", coodbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection.Close()

	// マルチプルセレクト(今度は、_ ではなく、rows)
	checkBundle := txBundle[0:27]
	log.Println("checkBundle is " + checkBundle)
	cmd := "SELECT COUNT(*) FROM coomile where tag = ?"
	row := DbConnection.QueryRow(cmd, checkBundle)

	// データ保存領域を確保
	var b Coomile
	// Scan にて、struct のアドレスにデータを入れる
	var count int
	log.Println("取り消された取引がDB内に存在するか確認中……")
	err := row.Scan(&count)
	// エラーハンドリング(共通関数にした方がいいのかな)
	if err != nil {
		// シングルセレクトの場合は、エラーハンドリングが異なる
		if err == sql.ErrNoRows {
			log.Println("There is no row!!!")
		} else {
			log.Println(err)
		}
	}

	if count == 0 {
		log.Println("normal transaction")
	} else {
		log.Println("iscanselled transaction")
	}
	fmt.Println(count)
	log.Println("end_check_db")
}

/* DBに全データを表示させる場合
func create_db(txBundle string, txAddress string, txTag string, txValue string) {
	//ここから挿入したデータの一覧を出力する処理
	// マルチプルセレクト(今度は、_ ではなく、rows)
	cmd = "SELECT * FROM tsc"
	rows, _ := DbConnection.Query(cmd)

	defer rows.Close()

	// データ保存領域を確保
	var bg []Tsc
	for rows.Next() {
		var b Tsc
		// Scan にて、struct のアドレスにデータを入れる
		err := rows.Scan(&b.bundle, &b.address, &b.tag, &b.value)
		// エラーハンドリング(共通関数にした方がいいのかな)
		if err != nil {
			log.Println(err)
		}
		// データ取得
		bg = append(bg, b)
	}

	// 操作結果を確認
	for _, b := range bg {
		fmt.Println(b.bundle, b.address, b.tag, b.value)
	}

	log.Println("create_db")
}
*/
