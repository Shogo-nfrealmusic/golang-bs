package main

import (
	"log"

	"gotrading/bitflyer"
	"gotrading/config"
	"gotrading/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)

	tickerCh := make(chan bitflyer.Ticker)
	go apiClient.GetRealTimeTicker(config.Config.ProductCode, tickerCh)

	for ticker := range tickerCh {
		log.Printf(
			"product=%s ltp=%f bid=%f ask=%f",
			ticker.ProductCode,
			ticker.Ltp,
			ticker.BestBid,
			ticker.BestAsk,
		)
	}
}
