package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"get-currency/typeFile"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

// ログ出力
func loggingSettings(filename string) {
	logFile, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	multiLogFile := io.MultiWriter(os.Stdout, logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(multiLogFile)
}

//ドル換算を行う関数
func concertToUSD(usd float64, cur float64) float64 {
	return cur / usd
}

func sleep(m int) {
	time.Sleep(time.Duration(m) * time.Second)
}

func main() {
	loggingSettings("ログ.log")
	log.Println("----- Start... -----")
	//為替API
	baseURL := "http://data.fixer.io/api/latest?access_key=5eaa513290dc99010a8d2dff7ced9c18&symbols=USD,JPY,BRL,MXN,ARS,CLP,COP,PEN,BOB"

	//引数のURLにGETリクエスト
	res, err := http.Get(baseURL)
	if err != nil {
		log.Fatal(err)
		log.Println(err)
	}
	//Response Bodyはクローズする必要があるので、クローズ処理
	defer res.Body.Close()
	//取得したURLの内容を読み込む
	body, _ := io.ReadAll(res.Body)

	//JSONのデータ(byte配列)を構造体に変換
	var currencyRate typeFile.JsonType
	if err := json.Unmarshal(body, &currencyRate); err != nil {
		log.Println(err)
		fmt.Println(err)
		fmt.Println("JSON Error, please try again.")
		sleep(3)
		os.Exit(3)
	}

	//Excel出力用の構造体を作成
	var excelData typeFile.ExcelData
	excelData.Date = currencyRate.Date
	excelData.JPY = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.JPY)
	excelData.BRL = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.BRL)
	excelData.MXN = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.MXN)
	excelData.ARS = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.ARS)
	excelData.CLP = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.CLP)
	excelData.PEN = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.PEN)
	excelData.COP = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.COP)
	excelData.BOB = concertToUSD(currencyRate.Rates.USD, currencyRate.Rates.BOB)

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
	fmt.Println("Please Press Enter!")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		in := scanner.Text()
		switch in {
		case "":
			log.Println("----- Completed!!! -----")
			log.Println("If the file is empty, please try again. It may be the error of API...")
			sleep(2)
			goto L
		default:
			fmt.Println("Command Error")
			continue
		}
	}
L:
}
