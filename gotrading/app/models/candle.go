package models

import (
	"fmt"
	"gotrading/bitflyer"
	"gotrading/config"
	"log"
	"time"
)

type Candle struct {
	ProductCode string `json:"product_code"`
	Duration    time.Duration `json:"duration"`
	Time        time.Time `json:"time"`
	Open        float64 `json:"open"`
	High        float64 `json:"high"`
	Low         float64 `json:"low"`
	Close       float64 `json:"close"`
	Volume      float64 `json:"volume"`
}

func formatCandleTime(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05")
}

func truncateCandleTime(t time.Time, duration time.Duration) time.Time {
	t = t.UTC()
	switch duration {
	case config.Duration1Month:
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	case config.Duration1Year:
		return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	default:
		return t.Truncate(duration)
	}
}

func parseCandleTime(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z07:00",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported time format: %s", value)
}

func NewCandle(productCode string, duration time.Duration, timeDate time.Time, open, close, high, low, volume float64) *Candle {
	return &Candle{
		ProductCode: productCode,
		Duration: duration,
		Time: timeDate,
		Open: open,
		Close: close,
		High: high,
		Low: low,
		Volume: volume,
	}
}

func (c *Candle) TableName() string {
	return GetCandleTableName(c.ProductCode, c.Duration)
}

func (c *Candle) Create() error {
	cmd := fmt.Sprintf(`
		INSERT INTO %s (time, open, high, low, close, volume) VALUES (?, ?, ?, ?, ?, ?)
	`, c.TableName())
	_, err := DbConnection.Exec(cmd, formatCandleTime(c.Time), c.Open, c.High, c.Low, c.Close, c.Volume)
	return err
}

func (c *Candle) save() error {
	cmd := fmt.Sprintf(`
		UPDATE %s SET open = ?, high = ?, low = ?, close = ?, volume = ? WHERE time = ?
	`, c.TableName())
	_, err := DbConnection.Exec(cmd, c.Open, c.High, c.Low, c.Close, c.Volume, formatCandleTime(c.Time))
	return err
}

func GetCandle(productCode string, duration time.Duration, datetime time.Time) *Candle {
	tableName := GetCandleTableName(productCode, duration)
	cmd := fmt.Sprintf("SELECT time, open, high, low, close, volume FROM %s WHERE time = ?", tableName)
	row := DbConnection.QueryRow(cmd, formatCandleTime(datetime))
	var (
		timeStr string
		candle  Candle
	)
	err := row.Scan(&timeStr, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Volume)
	if err != nil {
		return nil
	}
	candleTime, err := parseCandleTime(timeStr)
	if err != nil {
		log.Printf("action=GetCandle parse time err=%s value=%s", err, timeStr)
		return nil
	}
	return NewCandle(productCode, duration, candleTime, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
}

func CreateCandleWithDuration(ticker bitflyer.Ticker, productCode string, duration time.Duration) bool {
	candleTime := truncateCandleTime(ticker.DateTime(), duration)
	currentCandle := GetCandle(productCode, duration, candleTime)
	price := ticker.GetMidPrice()
	if currentCandle == nil {
		candle := NewCandle(productCode, duration, candleTime, price, price, price, price, ticker.Volume)
		if err := candle.Create(); err != nil {
			log.Printf("action=CreateCandleWithDuration create err=%s", err)
			return false
		}
		return true
	}
	if price > currentCandle.High {
		currentCandle.High = price
	}
	if price < currentCandle.Low {
		currentCandle.Low = price
	}
	currentCandle.Volume += ticker.Volume
	currentCandle.Close = price
	if err := currentCandle.save(); err != nil {
		log.Printf("action=CreateCandleWithDuration save err=%s", err)
	}
	return false
}

func GetAllCandles(productCode string, duration time.Duration, limit int) (dfCandle * DataFrameCandle, err error) {
	tableName := GetCandleTableName(productCode, duration)
	cmd := fmt.Sprintf(`SELECT * FROM (
		SELECT time, open, high, low, close, volume FROM %s ORDER BY time DESC LIMIT ?
	) ORDER BY time ASC;`, tableName)
	rows, err := DbConnection.Query(cmd, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	dfCandle = &DataFrameCandle{}
	dfCandle.ProductCode = productCode
	dfCandle.Duration = duration
	for rows.Next() {
		var candle Candle
		var timeStr string
		candle.ProductCode = productCode
		candle.Duration = duration
		if err = rows.Scan(&timeStr, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Volume); err != nil {
			return
		}
		candle.Time, err = parseCandleTime(timeStr)
		if err != nil {
			log.Printf("action=GetAllCandles parse time err=%s value=%s", err, timeStr)
			continue
		}
		dfCandle.Candles = append(dfCandle.Candles, candle)
	}
	err = rows.Err()
	if err != nil {
		return
	}
	return dfCandle, nil
}