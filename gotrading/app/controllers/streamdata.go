package controllers

import (
	"log"
	"sync"

	"gotrading/app/models"
	"gotrading/bitflyer"
	"gotrading/config"
)

var (
	emaPeriod1 int
	emaPeriod2 int
	emaOnce    sync.Once
)

func loadBestEmaPeriods() (int, int) {
	opts := models.BackTestOptions{
		UsePercent:       config.Config.UsePercet,
		StopLimitPercent: config.Config.StopLimitPercet,
	}
	limit := config.Config.DataLimit
	if limit <= 0 {
		limit = 365
	}
	return models.ResolveBestEmaPeriods(
		config.Config.ProductCode,
		config.Config.TradeDuration,
		limit,
		opts,
	)
}

func triggerTradeSignal() {
	emaOnce.Do(func() {
		emaPeriod1, emaPeriod2 = loadBestEmaPeriods()
		log.Printf("action=triggerTradeSignal best_ema_period1=%d best_ema_period2=%d", emaPeriod1, emaPeriod2)
	})

	opts := models.BackTestOptions{
		UsePercent:       config.Config.UsePercet,
		StopLimitPercent: config.Config.StopLimitPercet,
	}
	limit := config.Config.DataLimit
	if limit <= 0 {
		limit = 365
	}
	models.TriggerLatestEmaSignal(
		config.Config.ProductCode,
		config.Config.TradeDuration,
		limit,
		emaPeriod1,
		emaPeriod2,
		opts,
	)
}

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
					triggerTradeSignal()
				}
			}
		}
	}()
}
