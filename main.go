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

func concertToUSD(usd float64, cur float64) float64 {
	return cur / usd
}

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

	//Excel出力用の構造体を作成
	var excelData typeFile.ExcelData
	excelData.Date = currencyRate.Date
	excelData.JPY = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.JPY)
	excelData.BRL = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.BRL)
	excelData.MXN = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.MXN)
	excelData.ARS = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.ARS)
	excelData.CLP = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.CLP)
	excelData.PEN = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.PEN)
	excelData.BOB = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.BOB)

	//TODO あとで削除
	fmt.Printf("%+v", excelData)

	//Excelファイル作成
	file := excelize.NewFile()

	//名前をつけて保存
	if err := file.SaveAs("為替レート.xlsx"); err != nil {
		fmt.Println(err)
	}

}
