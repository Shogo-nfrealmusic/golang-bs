package controllers

import (
	"log"
	"sync"

	"gotrading/app/models"
	"gotrading/bitflyer"
	"gotrading/config"
)

func StreamINgestionData() {
	var tickerChannel = make(chan bitflyer.Ticker)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
	go apiClient.GetRealTimeTicker(config.Config.ProductCode, tickerChannel)

	go func() {
	var once sync.Once
	for ticker := range tickerChannel {
		once.Do(func() {
			log.Printf("action=StreamINgestionData first ticker ltp=%f timestamp=%s", ticker.Ltp, ticker.Timestamp)
		})
		for _, duration := range config.Config.Durations {
			isCreated := models.CreateCandleWithDuration(ticker, config.Config.ProductCode, duration)
			if isCreated && duration == config.Config.TradeDuration {
				_ = isCreated // TODO: trigger trade signal
			}
		}
	}
	}()
}
