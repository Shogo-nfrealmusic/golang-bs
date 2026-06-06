package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	ApiKey string
	ApiSecret string
	LogFile string
	ProductCode string

	TradeDuration time.Duration
	Durations map[string]time.Duration
	Dbname string
	SQLDriver string
	Port int

	BackTest bool
	UsePercet float64
	DataLimit int
	StopLimitPercet float64
	NumRanking int
}

var Config ConfigList

const (
	Duration1Month = 30 * 24 * time.Hour
	Duration1Year  = 365 * 24 * time.Hour
)

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		os.Exit(1)
	}

	durations := map[string]time.Duration{
		"1s":  time.Second,
		"1m":  time.Minute,
		"1h":  time.Hour,
		"1mo": Duration1Month,
		"1y":  Duration1Year,
	}

	Config = ConfigList{
		ApiKey: cfg.Section("bitflyer").Key("api_key").String(),
		ApiSecret: cfg.Section("bitflyer").Key("api_secret").String(),
		LogFile: cfg.Section("gotrading").Key("log_file").String(),
		ProductCode: cfg.Section("gotrading").Key("product_code").String(),
		Durations: durations,
		TradeDuration: durations[cfg.Section("gotrading").Key("trade_duration").String()],
		Dbname: cfg.Section("db").Key("name").String(),
		SQLDriver: cfg.Section("db").Key("driver").String(),
		Port: cfg.Section("web").Key("port").MustInt(8080),
		BackTest: cfg.Section("gotrading").Key("back_test").MustBool(false),
		UsePercet: cfg.Section("gotrading").Key("use_percent").MustFloat64(0.9),
		DataLimit: cfg.Section("gotrading").Key("data_limit").MustInt(365),
		StopLimitPercet: cfg.Section("gotrading").Key("stop_limit_percent").MustFloat64(0.9),
		NumRanking: cfg.Section("gotrading").Key("num_ranking").MustInt(3),
	}
}