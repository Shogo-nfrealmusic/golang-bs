package main

import (
	"fmt"
	"os"

	"github.com/markcheno/go-quote"
	"github.com/markcheno/go-talib"
)

func main() {
	spy, _ := quote.NewQuoteFromTiingo(
		"SPY",
		"2016-01-01",
		"2016-04-01",
		quote.Daily,
		os.Getenv("TIINGO_TOKEN"),
	)
	fmt.Println(spy.CSV())
	rsi2 := talib.Rsi(spy.Close, 2)
	fmt.Println(rsi2)
}