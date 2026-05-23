package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gotrading/config"

	_ "github.com/mattn/go-sqlite3"
)

const (
	tableNameSignalEvents = "signal_events"
)

var DbConnection *sql.DB

func GetCandleTableName(productCode string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", productCode, duration.String())
}

func init() {
	var err error
	DbConnection, err = sql.Open(config.Config.SQLDriver, config.Config.Dbname)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		os.Exit(1)
	}

	cmd := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			time DATETIME PRIMARY KEY NOT NULL,
			product_code STRING,
			side STRING,
			price FLOAT,
			size FLOAT
		)
	`, tableNameSignalEvents)
	DbConnection.Exec(cmd)

	for _, duration := range config.Config.Durations {
		tableName := GetCandleTableName(config.Config.ProductCode, duration)
		c := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				time DATETIME PRIMARY KEY NOT NULL,
				open FLOAT,
				high FLOAT,
				low FLOAT,
				close FLOAT,
				volume FLOAT)`, tableName)
		DbConnection.Exec(c)
	}

}
