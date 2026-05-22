package main

import (
	"gotrading/bitflyer"
	"gotrading/config"
	"gotrading/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
	apiClient.GetTicker("BTC_USD")
	
}
