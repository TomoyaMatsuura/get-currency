package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Rates struct {
	USD float64 //1.077447
	JPY float64 //149.15042
	BRL float64 //5.358034
	MXN float64 //19.368097
	ARS float64 //253.137514
	CLP float64 //863.035119
	PEN float64 //3.983791
	BOB float64 //7.448222
}

type JsonType struct {
	Success   bool   //true
	Timestamp int    //1684897743
	Base      string //"EUR"
	Date      string //"2023-05-24"
	Rates     Rates  //{"USD":1.078411,"JPY":149.318947,"BRL":5.362724,"MXN":19.361185,"ARS":253.343999,"CLP":863.806891,"PEN":3.987355,"BOB":7.454886}
}

func main() {
	baseURL := "http://data.fixer.io/api/latest?access_key=27ccb6a5175f81a9499b130075a951ad&symbols=USD,JPY,BRL,MXN,ARS,CLP,PEN,BOB"

	res, err := http.Get(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))

	fmt.Println(baseURL)

	var currencyRate JsonType

	if err := json.Unmarshal(body, &currencyRate); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", currencyRate)

}
