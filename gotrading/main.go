package main

import (
	"gotrading/app/controllers"
	"gotrading/config"
	"gotrading/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	go controllers.StreamINgestionData()
	if err := controllers.StartWebServer(); err != nil {
		panic(err)
	}
}
