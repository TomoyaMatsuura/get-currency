package main

import (
	"encoding/json"
	"fmt"
	"get-currency/typeFile"
	"io"
	"log"
	"net/http"

	"github.com/xuri/excelize/v2"
)

func main() {
	baseURL := "http://data.fixer.io/api/latest?access_key=27ccb6a5175f81a9499b130075a951ad&symbols=USD,JPY,BRL,MXN,ARS,CLP,PEN,BOB"

	//引数のURLにGETリクエスト
	res, err := http.Get(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	//Response Bodyはクローズする必要があるので、クローズ処理
	defer res.Body.Close()
	//取得したURLの内容を読み込む
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))

	//TODO あとで削除
	fmt.Println(baseURL)

	var currencyRate typeFile.JsonType

	//JSONのデータ(byte配列)を構造体に変換
	if err := json.Unmarshal(body, &currencyRate); err != nil {
		fmt.Println(err)
	}
	//TODO あとで削除
	fmt.Printf("%+v\n", currencyRate)
	//Excelファイル作成
	file := excelize.NewFile()

	//名前をつけて保存
	if err := file.SaveAs("為替レート.xlsx"); err != nil {
		fmt.Println(err)
	}

}
