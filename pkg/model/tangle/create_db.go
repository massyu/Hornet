package tangle

import (
	// ビルド時のみ使用する
	"database/sql"
	"fmt"
	"log"
	"strconv"

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
	thash   string
	address string
	tag     string
	value   int
}

// データ格納用
type Coomile struct {
	mindex int
	tag    string
}

func create_db(txBundle string, txHash string, txAddress string, txTag string, txValue string) {
	log.Println("create_db開始")
	// Open(driver,  sql 名(任意の名前))
	DbConnection, _ := sql.Open("sqlite3", dbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection.Close()

	// blog テーブルの作成
	cmd := `CREATE TABLE IF NOT EXISTS tsc(
			 bundle STRING,    
			 thash STRING,    
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

	cmd = "INSERT INTO tsc (bundle, thash, address, tag, value) VALUES (?, ?, ?, ?, ?)"
	_, err = DbConnection.Exec(cmd, txBundle, txHash, txAddress, txTag, txValue)

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
	err = row.Scan(&b.bundle, &b.thash, &b.address, &b.tag, &b.value)
	// エラーハンドリング(共通関数にした方がいいのかな)
	if err != nil {
		// シングルセレクトの場合は、エラーハンドリングが異なる
		if err == sql.ErrNoRows {
			log.Println("There is no row!!!")
		} else {
			log.Println(err)
		}
	}

	//fmt.Println(b.bundle, b.thash, b.address, b.tag, b.value)

	log.Println("create_db")
}

func checkDB(txBundle string, txHash string) int {
	log.Println("check_db開始")
	// Open(driver,  sql 名(任意の名前))
	DbConnection, _ := sql.Open("sqlite3", coodbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection.Close()

	// 出力確認テスト
	sumplebundle := "YKBOJAUHGA9FQWYXBHRDQPAGNYGO9LIFAINJECFUEBUZTO9JHUEJWFMJYFFEGKUHYWDOOIYAGUKZLBHSA"

	// マルチプルセレクト(今度は、_ ではなく、rows)
	// checkBundle := txBundle[0:27]
	checkBundle := sumplebundle[0:27]

	log.Println("checkBundle is " + checkBundle)
	cmd := "SELECT COUNT(*) FROM coomile where thash = ?"
	row := DbConnection.QueryRow(cmd, checkBundle)

	// データ保存領域を確保
	// var b Coomile
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

	var cngValue int
	cngValue = 0

	if count == 0 {
		log.Println("normal transaction")
	} else {
		log.Println("iscanselled transaction")

		// Open(driver,  sql 名(任意の名前))
		DbConnection2, _ := sql.Open("sqlite3", dbPath)

		// Connection をクローズする。(defer で閉じるのが Golang の作法)
		defer DbConnection2.Close()

		// 今のhashを引数に持ってきて問い合わせを行い、valueを読みだす
		cmd2 := "SELECT * FROM tsc where thash = ?"
		// row = DbConnection.QueryRow(cmd, txHash)
		rows, _ := DbConnection2.Query(cmd2, sumplebundle)

		// データ保存領域を確保
		// var b Coomile
		// Scan にて、struct のアドレスにデータを入れる
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

		// valueの加算処理
		for _, b := range bg {
			if b.thash == txHash {
				fmt.Println("cngValueに" + strconv.Itoa(b.value) + "を加算")
				cngValue += b.value
				fmt.Println(cngValue)
			}
			fmt.Println(b.thash, b.value)
		}
	}
	fmt.Println(count)
	log.Println("end_check_db")
	return cngValue
}

func checkDB2(txHash string) int {
	log.Println("check_db開始")
	// Open(driver,  sql 名(任意の名前))
	DbConnection, _ := sql.Open("sqlite3", coodbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection.Close()

	// 出力確認テスト
	sumplebundle := "YKBOJAUHGA9FQWYXBHRDQPAGNYGO9LIFAINJECFUEBUZTO9JHUEJWFMJYFFEGKUHYWDOOIYAGUKZLBHSA"

	// マルチプルセレクト(今度は、_ ではなく、rows)
	// checkBundle := txBundle[0:27]
	checkBundle := sumplebundle[0:27]

	log.Println("checkBundle is " + checkBundle)
	cmd := "SELECT COUNT(*) FROM coomile where tag = ?"
	row := DbConnection.QueryRow(cmd, checkBundle)

	// データ保存領域を確保
	// var b Coomile
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

	var cngValue int
	cngValue = 0

	if count == 0 {
		log.Println("normal transaction")
	} else {
		log.Println("iscanselled transaction")

		// Open(driver,  sql 名(任意の名前))
		DbConnection2, _ := sql.Open("sqlite3", dbPath)

		// Connection をクローズする。(defer で閉じるのが Golang の作法)
		defer DbConnection2.Close()

		// 今のaddressを引数に持ってきて問い合わせを行い、valueを読みだす
		cmd2 := "SELECT * FROM tsc where bundle = ?"
		// row = DbConnection.QueryRow(cmd, txAddress)
		rows, _ := DbConnection2.Query(cmd2, sumplebundle)

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

		// valueの加算処理
		for _, b := range bg {
			if b.thash == txHash {
				fmt.Println("cngValueに" + strconv.Itoa(b.value) + "を加算")
				cngValue += b.value
				fmt.Println(cngValue)
			}
			fmt.Println(b.thash, b.value)
		}
	}
	fmt.Println(count)
	log.Println("end_check_db")
	return cngValue
}

/*
	// データ保存領域を確保
	var b Tsc
	// Scan にて、struct のアドレスにデータを入れる
	log.Println("取り消された取引がDB内に存在するか確認中……")
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
	fmt.Println(b.address, b.value)
*/

/*
	log.Println("coomileの全データ表示")
	// マルチプルセレクト(今度は、_ ではなく、rows)
	cmd = "SELECT * FROM coomile"
	rows, _ := DbConnection.Query(cmd)

	defer rows.Close()

	// データ保存領域を確保
	var bg []Coomile
	for rows.Next() {
		var b Coomile
		// Scan にて、struct のアドレスにデータを入れる
		err := rows.Scan(&b.mindex, &b.tag)
		// エラーハンドリング(共通関数にした方がいいのかな)
		if err != nil {
			log.Println(err)
		}
		// データ取得
		bg = append(bg, b)
	}

	// 操作結果を確認
	for _, b := range bg {
		fmt.Println(b.mindex, b.tag)
	}
*/

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
