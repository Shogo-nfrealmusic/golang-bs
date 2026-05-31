package models

import (
	"time"

	"github.com/markcheno/go-talib"
)

type DataFrameCandle struct {
	ProductCode string `json:"product_code"`
	Duration time.Duration `json:"duration"`
	Candles []Candle `json:"candles"`
	Smas    []Sma     `json:"smas,omitempty"`
	Emas    []Ema     `json:"emas,omitempty"`
	Bbands  []Bbands  `json:"bbands,omitempty"`
}

type Sma struct {
	Period int       `json:"period,omitempty"`
	Values []float64 `json:"values,omitempty"`
}

type Ema struct {
	Period int       `json:"period,omitempty"`
	Values []float64 `json:"values,omitempty"`
}

type Bbands struct {
	N    int       `json:"n,omitempty"`
	K    float64   `json:"k,omitempty"`
	Up   []float64 `json:"up,omitempty"`
	Mid  []float64 `json:"mid,omitempty"`
	Down []float64 `json:"down,omitempty"`
}

func (df * DataFrameCandle) Times() []time.Time {
	s := make([]time.Time, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Time
	}
	return s
}

func (df * DataFrameCandle) Open() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Open
	}
	return s
}

func (df * DataFrameCandle) High() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.High
	}
	return s
}

func (df * DataFrameCandle) Low() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Low
	}
	return s
}

func (df * DataFrameCandle) Close() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Close
	}
	return s
}

func (df * DataFrameCandle) Volume() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Volume
	}
	return s
}

func (df *DataFrameCandle) AddSma(period int) bool {
	if len(df.Candles) <= period {
		return false
	}
	df.Smas = append(df.Smas, Sma{
		Period: period,
		Values: talib.Sma(df.Close(), period),
	})
	return true
}

func (df *DataFrameCandle) AddEma(period int) bool {
	if len(df.Candles) <= period {
		return false
	}
	df.Emas = append(df.Emas, Ema{
		Period: period,
		Values: talib.Ema(df.Close(), period),
	})
	return true
}

func (df *DataFrameCandle) AddBbands(n int, k float64) bool {
	if len(df.Candles) <= n {
		return false
	}
	up, mid, down := talib.BBands(df.Close(), n, k, k, talib.SMA)
	df.Bbands = append(df.Bbands, Bbands{
		N:    n,
		K:    k,
		Up:   up,
		Mid:  mid,
		Down: down,
	})
	return true
}