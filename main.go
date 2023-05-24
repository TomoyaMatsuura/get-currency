package main

import (
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
	success   bool   //true
	timestamp int    //1684897743
	base      string //"EUR"
	date      string //"2023-05-24"
	Rates     Rates
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
}
