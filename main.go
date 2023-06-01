package main

import (
	"encoding/json"
	"fmt"
	"get-currency/typeFile"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-ieproxy"
	"github.com/xuri/excelize/v2"
)

// ログ出力を行う関数
func loggingSettings(filename string) {
	logFile, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	multiLogFile := io.MultiWriter(os.Stdout, logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(multiLogFile)
}

// ドル換算を行う関数(APIフリープランでは基軸通貨がEURなのでドルに直す必要あり…)
func convertToUSD(usd float64, cur float64) float64 {
	return cur / usd
}

// m秒待機する関数
func sleep(m int) {
	time.Sleep(time.Duration(m) * time.Second)
}

func main() {
	//ログファイルを作成
	loggingSettings("ログ.log")
	log.Println("----- Start... -----")

	//為替APIのURL fixerのAPI(free plan)を使用 ※月に100回まで呼び出し可能
	baseURL := "http://data.fixer.io/api/latest?access_key=5eaa513290dc99010a8d2dff7ced9c18&symbols=USD,JPY,BRL,MXN,ARS,CLP,COP,PEN,BOB"

	//プロキシ設定(Internet Optionから確認可能)
	os.Setenv("HTTP_PROXY", "http://w055185.mj.makita.local:8080")
	os.Setenv("HTTPS_PROXY", "http://w055185.mj.makita.local:8080")

	//引数のURLにGETメソッドでのHTTPリクエスト
	res, err := http.Get(baseURL)
	if err != nil {
		log.Fatal(err)
		log.Println(err)
	}
	//Response Bodyはクローズする必要があるので、クローズ処理
	defer res.Body.Close()
	//取得したURLの内容を読み込む
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	//JSONのデータ(byte配列)を構造体に変換
	var currencyRate typeFile.JsonType
	if err := json.Unmarshal(body, &currencyRate); err != nil {
		log.Println(err)
		fmt.Println("API Error, please try again.")
		sleep(2)
		os.Exit(3)
	}

	//Excel出力用の構造体を作成
	var excelData typeFile.ExcelData
	//APIから取得した値をドル換算
	excelData.Date = currencyRate.Date
	excelData.JPY = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.JPY)
	excelData.BRL = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.BRL)
	excelData.MXN = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.MXN)
	excelData.ARS = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.ARS)
	excelData.CLP = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.CLP)
	excelData.PEN = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.PEN)
	excelData.COP = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.COP)
	excelData.BOB = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.BOB)

	//Excelファイル作成
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
			fmt.Println(err)
			fmt.Println("もう一度実行してください")
		}
	}()

	//シート名を月日にする
	page := excelData.Date
	file.SetSheetName("Sheet1", page)

	//セルに記述
	file.SetCellValue(page, "A1", "取得日")
	file.SetCellValue(page, "B1", excelData.Date)
	file.SetCellValue(page, "A3", "日本円")
	file.SetCellValue(page, "B3", excelData.JPY)
	file.SetCellValue(page, "A4", "ブラジルレアル")
	file.SetCellValue(page, "B4", excelData.BRL)
	file.SetCellValue(page, "A5", "メキシコペソ")
	file.SetCellValue(page, "B5", excelData.MXN)
	file.SetCellValue(page, "A6", "アルゼンチンペソ")
	file.SetCellValue(page, "B6", excelData.ARS)
	file.SetCellValue(page, "A7", "チリペソ")
	file.SetCellValue(page, "B7", excelData.CLP)
	file.SetCellValue(page, "A8", "ペルーソル")
	file.SetCellValue(page, "B8", excelData.PEN)
	file.SetCellValue(page, "A9", "コロンビアペソ")
	file.SetCellValue(page, "B9", excelData.COP)
	file.SetCellValue(page, "A10", "ボリビアーノ")
	file.SetCellValue(page, "B10", excelData.BOB)

	//横幅を調整
	file.SetColWidth(page, "A", "B", 20)

	//書式設定
	styleID, err := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		NumFmt:    2,
		Font:      &excelize.Font{Size: 11, Family: "Arial"},
	})
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		fmt.Println("もう一度実行してください")
	}
	file.SetCellStyle(page, "A1", "B11", styleID)

	//名前をつけて保存
	time := time.Now()
	formatTime := time.Format("20060102")
	if err := file.SaveAs("為替レート" + formatTime + ".xlsx"); err != nil {
		log.Println(err)
		fmt.Println(err)
		fmt.Println("もう一度実行してください")
	}

	log.Println("----- Completed!!! -----")
	sleep(3)

}
