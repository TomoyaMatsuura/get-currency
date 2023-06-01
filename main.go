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
	baseURL := "http://data.fixer.io/api/latest?access_key=5eaa513290dc99010a8d2dff7ced9c18&symbols=USD,JPY,CAD,AUD,NZD,BRL,MXN,ARS,CLP,COP,PEN,BOB,INR,TRY,RUB,GBP,EUR,MAD"

	//プロキシ設定(Internet Optionから確認可能)
	//os.Setenv("HTTP_PROXY", "http://w055185.mj.makita.local:8080")
	//os.Setenv("HTTPS_PROXY", "http://w055185.mj.makita.local:8080")

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
	excelData.CAD = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.CAD)
	excelData.AUD = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.AUD)
	excelData.NZD = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.NZD)
	excelData.BRL = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.BRL)
	excelData.MXN = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.MXN)
	excelData.ARS = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.ARS)
	excelData.CLP = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.CLP)
	excelData.PEN = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.PEN)
	excelData.COP = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.COP)
	excelData.BOB = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.BOB)
	excelData.INR = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.INR)
	excelData.TRY = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.TRY)
	excelData.RUB = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.RUB)
	excelData.GBP = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.GBP)
	excelData.EUR = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.EUR)
	excelData.MAD = convertToUSD(currencyRate.Rates.USD, currencyRate.Rates.MAD)

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

	//aColumn := [19]string{"取得日","","日本円","カナダドル","オーストラリアドル","ニュージーランドドル","ブラジルレアル","メキシコペソ","アルゼンチンペソ","チリペソ","ペルーソル","コロンビアペソ","ボリビアーノ","インドルピー(参考)","トルコリラ(参考)","ロシアルーブル","英国ポンド(参考)","ユーロ","モロッコディルハム"}
	//bColumn := [19]string{}

	//セルに記述
	for i := 1; i <= 17; i++ {
		file.SetCellValue(page, "A", "")
		file.SetCellValue(page, "B", "")
	}
	file.SetCellValue(page, "A1", "取得日")
	file.SetCellValue(page, "B1", excelData.Date)
	file.SetCellValue(page, "A3", "日本円")
	file.SetCellValue(page, "B3", excelData.JPY)
	file.SetCellValue(page, "A4", "カナダドル")
	file.SetCellValue(page, "B4", excelData.CAD)
	file.SetCellValue(page, "A5", "オーストラリアドル")
	file.SetCellValue(page, "B5", excelData.AUD)
	file.SetCellValue(page, "A6", "ニュージーランドドル")
	file.SetCellValue(page, "B6", excelData.NZD)
	file.SetCellValue(page, "A7", "ブラジルレアル")
	file.SetCellValue(page, "B7", excelData.BRL)
	file.SetCellValue(page, "A8", "メキシコペソ")
	file.SetCellValue(page, "B8", excelData.MXN)
	file.SetCellValue(page, "A9", "アルゼンチンペソ")
	file.SetCellValue(page, "B9", excelData.ARS)
	file.SetCellValue(page, "A10", "チリペソ")
	file.SetCellValue(page, "B10", excelData.CLP)
	file.SetCellValue(page, "A11", "ペルーソル")
	file.SetCellValue(page, "B11", excelData.PEN)
	file.SetCellValue(page, "A12", "コロンビアペソ")
	file.SetCellValue(page, "B12", excelData.COP)
	file.SetCellValue(page, "A13", "ボリビアーノ")
	file.SetCellValue(page, "B13", excelData.BOB)
	file.SetCellValue(page, "A14", "インドルピー(参考)")
	file.SetCellValue(page, "B14", excelData.INR)
	file.SetCellValue(page, "A15", "トルコリラ(参考)")
	file.SetCellValue(page, "B15", excelData.TRY)
	file.SetCellValue(page, "A16", "ロシアルーブル")
	file.SetCellValue(page, "B16", excelData.RUB)
	file.SetCellValue(page, "A17", "英国ポンド(参考)")
	file.SetCellValue(page, "B17", excelData.GBP)
	file.SetCellValue(page, "A18", "ユーロ")
	file.SetCellValue(page, "B18", excelData.EUR)
	file.SetCellValue(page, "A19", "モロッコディルハム")
	file.SetCellValue(page, "B19", excelData.MAD)

	//横幅を調整
	file.SetColWidth(page, "A", "A", 25)
	file.SetColWidth(page, "B", "B", 20)

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
	file.SetCellStyle(page, "A1", "B19", styleID)

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
