package main

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	Port int
	DbName string
	SQLDriver string
}

var Config ConfigList

func init() {
	cfg, _ := ini.Load("config.ini")
	Config = ConfigList{
		Port: cfg.Section("web").Key("port").MustInt(8080),
		DbName: cfg.Section("db").Key("name").String(),
		SQLDriver: cfg.Section("db").Key("driver").String(),
	}
}

func main() {
	fmt.Println(Config.Port)
	fmt.Println(Config.DbName)
	fmt.Println(Config.SQLDriver)
}