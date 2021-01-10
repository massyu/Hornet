package tangle

import (
	// ビルド時のみ使用する
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
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
	mindex  int
	thash   string
	address string
	bundle  string
	value   int
}

func createDB(txBundle string, txHash string, txAddress string, txTag string, txValue string) {
	log.Println("createDB開始")
	// Open(driver,  sql 名(任意の名前))
	DbConnection, _ := sql.Open("sqlite3", dbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection.Close()

	// tsc テーブルの作成
	/////////////////////////////////////////////////////////////////////////////////////
	cmd := `CREATE TABLE IF NOT EXISTS tsc(
			 bundle STRING,    
			 thash STRING,    
             address STRING,    
             tag STRING,    
             value INT)`

	// cmd を実行
	// _ -> 受け取った結果に対して何もしないので、_ にする
	_, err := DbConnection.Exec(cmd)

	// エラーハンドリング
	if err != nil {
		// エラーが発生した場合、以降の処理を実行しない
		log.Fatalln(err)
	}
	/////////////////////////////////////////////////////////////////////////////////////

	//データ挿入処理
	/////////////////////////////////////////////////////////////////////////////////////
	cmd = "INSERT INTO tsc (bundle, thash, address, tag, value) VALUES (?, ?, ?, ?, ?)"
	_, err = DbConnection.Exec(cmd, txBundle, txHash, txAddress, txTag, txValue)

	if err != nil {
		// golang には、try-catch がない。nil か否かで判定
		log.Fatalln(err)
	}
	/////////////////////////////////////////////////////////////////////////////////////

	//挿入したデータの一覧を出力する処理
	/////////////////////////////////////////////////////////////////////////////////////
	// マルチプルセレクト(今度は、_ ではなく、rows)
	cmd = "SELECT * FROM tsc where thash = ?"
	row := DbConnection.QueryRow(cmd, txHash)

	// データ保存領域を確保
	var b Tsc
	// Scan にて、struct のアドレスにデータを入れる
	err = row.Scan(&b.bundle, &b.thash, &b.address, &b.tag, &b.value)
	// エラーハンドリング
	if err != nil {
		// シングルセレクトの場合は、エラーハンドリングが異なる
		if err == sql.ErrNoRows {
			log.Println("There is no row!!!")
		} else {
			log.Println(err)
		}
	}
	/////////////////////////////////////////////////////////////////////////////////////

	fmt.Println(b.bundle, b.thash, b.address, b.tag, b.value)
	log.Println("をtscDBに挿入")

	log.Println("createDB終了")
}

//CoodbにtxHashが存在するかどうかチェック
func checkHashForCooDB(txHash string) (bool, error) {
	log.Println("checkHashForCooDB開始")
	// Open(driver,  sql 名(任意の名前))
	DbConnection, _ := sql.Open("sqlite3", coodbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection.Close()

	// 出力確認テスト
	sumpleTxHash := "ZKDFASYSQ9VZMRWNTYPV9SGLLFSQZEYMPPSUWDBEXYQJQGFWMHXNQTOVPQPWLIULZDECHYVXOMTLODX99"

	// checkBundle := txBundle[0:27]
	// checkBundle := sumpleTxHash[0:27]
	// log.Println("checkBundle is " + checkBundle)

	/////////////////////////////////////////////////////////////////////////////////////
	cmd := "SELECT COUNT(*) FROM coomile where thash = ?"
	row := DbConnection.QueryRow(cmd, sumpleTxHash)

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
	/////////////////////////////////////////////////////////////////////////////////////

	fmt.Println(count)
	if count == 0 {
		log.Println("normal transaction")
		log.Println("checkHashForCooDB終了")
		return false, nil

	} else {
		log.Println("iscanselled transaction")
		log.Println("checkHashForCooDB終了")
		return true, errors.New("無効化された取引です")

	}
	return false, nil
}

func checkDBForBalances(txHash string) int {
	log.Println("checkDBForBalances開始")
	// Open(driver,  sql 名(任意の名前))
	DbConnection, _ := sql.Open("sqlite3", coodbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection.Close()

	// 出力確認テスト
	sumpleTxHash := "ZKDFASYSQ9VZMRWNTYPV9SGLLFSQZEYMPPSUWDBEXYQJQGFWMHXNQTOVPQPWLIULZDECHYVXOMTLODX99"

	// マルチプルセレクト(今度は、_ ではなく、rows)
	// checkBundle := txBundle[0:27]
	// checkBundle := sumplebundle[0:27]

	/////////////////////////////////////////////////////////////////////////////////////
	cmd := "SELECT COUNT(*) FROM coomile where thash = ?"
	row := DbConnection.QueryRow(cmd, sumpleTxHash)

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
	/////////////////////////////////////////////////////////////////////////////////////

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
		/////////////////////////////////////////////////////////////////////////////////////
		cmd2 := "SELECT * FROM tsc where thash = ?"
		// row = DbConnection.QueryRow(cmd, txHash)
		rows, _ := DbConnection2.Query(cmd2, sumpleTxHash) //通常時はtxHashが引数
		defer rows.Close()

		// データ保存領域を確保
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
		/////////////////////////////////////////////////////////////////////////////////////

		// データ保存領域を確保
		/////////////////////////////////////////////////////////////////////////////////////
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
		/////////////////////////////////////////////////////////////////////////////////////

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
	log.Println("checkDBForBalances終了")
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
