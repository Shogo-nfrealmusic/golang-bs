package models

import (
	"math"
	"time"

	"gotrading/dtadingalgo"

	"github.com/markcheno/go-talib"
)

type DataFrameCandle struct {
	ProductCode string `json:"product_code"`
	Duration time.Duration `json:"duration"`
	Candles []Candle `json:"candles"`
	Smas    []Sma     `json:"smas,omitempty"`
	Emas    []Ema     `json:"emas,omitempty"`
	Bbands   []Bbands  `json:"bbands,omitempty"`
	Ichimoku *Ichimoku `json:"ichimoku,omitempty"`
	Rsis    []Rsi     `json:"rsis,omitempty"`
	Macds   []Macd    `json:"macds,omitempty"`
	Hvs     []Hv      `json:"hvs,omitempty"`
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

type Ichimoku struct {
	Tenkan  []*float64 `json:"tenkan,omitempty"`
	Kikun   []*float64 `json:"kikun,omitempty"`
	SenkouA []*float64 `json:"senkou_a,omitempty"`
	SenkouB []*float64 `json:"senkou_b,omitempty"`
	Chikou  []*float64 `json:"chikou,omitempty"`
}

type Rsi struct {
	Period int        `json:"period,omitempty"`
	Values []*float64 `json:"values,omitempty"`
}

type Macd struct {
	FastPeriod   int        `json:"fast_period,omitempty"`
	SlowPeriod   int        `json:"slow_period,omitempty"`
	SignalPeriod int        `json:"signal_period,omitempty"`
	Macd         []*float64 `json:"macd,omitempty"`
	Signal       []*float64 `json:"signal,omitempty"`
	Hist         []*float64 `json:"hist,omitempty"`
}

type Hv struct {
	Period int        `json:"period,omitempty"`
	Values []*float64 `json:"values,omitempty"`
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

func (df *DataFrameCandle) AddIchimoku(tenkan, kikun, senkouB, displacement int) bool {
	if tenkan <= 0 {
		tenkan = 9
	}
	if kikun <= 0 {
		kikun = 26
	}
	if senkouB <= 0 {
		senkouB = 52
	}
	if displacement <= 0 {
		displacement = 26
	}
	minLen := senkouB + displacement
	if len(df.Candles) < minLen {
		return false
	}

	lines := dtadingalgo.Ichimoku(dtadingalgo.IchimokuInput{
		High:  df.High(),
		Low:   df.Low(),
		Close: df.Close(),
	}, dtadingalgo.IchimokuParams{
		TenkanPeriod:  tenkan,
		KikunPeriod:   kikun,
		SenkouBPeriod: senkouB,
		Displacement:  displacement,
	})

	df.Ichimoku = &Ichimoku{
		Tenkan:  floatsToNullableJSON(lines.Tenkan),
		Kikun:   floatsToNullableJSON(lines.Kikun),
		SenkouA: floatsToNullableJSON(lines.SenkouA),
		SenkouB: floatsToNullableJSON(lines.SenkouB),
		Chikou:  floatsToNullableJSON(lines.Chikou),
	}
	return true
}

func (df *DataFrameCandle) AddRsi(period int) bool {
	if period <= 0 {
		period = 14
	}
	if len(df.Candles) <= period {
		return false
	}
	df.Rsis = append(df.Rsis, Rsi{
		Period: period,
		Values: floatsToNullableJSON(talib.Rsi(df.Close(), period)),
	})
	return true
}

func (df *DataFrameCandle) AddMacd(fast, slow, signal int) bool {
	if fast <= 0 {
		fast = 12
	}
	if slow <= 0 {
		slow = 26
	}
	if signal <= 0 {
		signal = 9
	}
	if fast >= slow {
		return false
	}
	minLen := slow + signal
	if len(df.Candles) <= minLen {
		return false
	}

	macdLine, signalLine, hist := talib.Macd(df.Close(), fast, slow, signal)
	df.Macds = append(df.Macds, Macd{
		FastPeriod:   fast,
		SlowPeriod:   slow,
		SignalPeriod: signal,
		Macd:         floatsToNullableJSON(macdLine),
		Signal:       floatsToNullableJSON(signalLine),
		Hist:         floatsToNullableJSON(hist),
	})
	return true
}

func periodsPerYear(duration time.Duration) float64 {
	if duration <= 0 {
		return 252
	}
	year := float64(365 * 24 * time.Hour)
	return year / float64(duration)
}

func (df *DataFrameCandle) AddHistoricalVolatility(period int) bool {
	if period <= 0 {
		period = 20
	}
	if len(df.Candles) <= period {
		return false
	}

	values := dtadingalgo.HistoricalVolatility(
		df.Close(),
		period,
		periodsPerYear(df.Duration),
	)
	df.Hvs = append(df.Hvs, Hv{
		Period: period,
		Values: floatsToNullableJSON(values),
	})
	return true
}

func floatsToNullableJSON(values []float64) []*float64 {
	out := make([]*float64, len(values))
	for i, v := range values {
		if !math.IsNaN(v) && !math.IsInf(v, 0) {
			val := v
			out[i] = &val
		}
	}
	return out
}