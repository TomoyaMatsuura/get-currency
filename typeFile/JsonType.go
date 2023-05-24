package typeFile

type JsonType struct {
	Success   bool   //true
	Timestamp int    //1684897743
	Base      string //"EUR"
	Date      string //"2023-05-24"
	Rates     Rates  //{"USD":1.078411,"JPY":149.318947,"BRL":5.362724,"MXN":19.361185,"ARS":253.343999,"CLP":863.806891,"PEN":3.987355,"BOB":7.454886}
}
