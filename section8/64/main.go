package main

import (
	"fmt"
	"log"
	"os"

	"github.com/markcheno/go-quote"
	"github.com/markcheno/go-talib"
)

func main() {
	spy, err := quote.NewQuoteFromTiingo("spy", "2016-01-01", "2016-04-01", quote.Daily, os.Getenv("TIINGO_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(spy.CSV())
	rsi2 := talib.Rsi(spy.Close, 2)
	fmt.Println(rsi2)
}
