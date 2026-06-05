package main

import (
	"gotrading/app/models"
	"time"
)

// import (
// 	"gotrading/app/controllers"
// 	"gotrading/config"
// 	"gotrading/utils"
// )

// func main() {
// 	utils.LoggingSettings(config.Config.LogFile)
// 	go controllers.StreamINgestionData()
// 	if err := controllers.StartWebServer(); err != nil {
// 		panic(err)
// 	}
// }

func main() {
	s := models.NewSignalEvents()
	df, _ := models.GetAllCandles("BTC_JPY", time.Minute, 10)
	c1 := df.Candles[0]
	c2 := df.Candles[5]
	s.Buy("BTC_JPY", c1.Time, c1.Close, 1.0, true)
	s.Sell("BTC_JPY", c2.Time, c2.Close, 1.0, true)
}
