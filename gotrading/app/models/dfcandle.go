package models

import (
	"math"
	"sort"
	"time"

	"gotrading/config"
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
	Events               *SignalEvents         `json:"events,omitempty"`
	BackTestRanking      []EmaBackTestRank     `json:"back_test_ranking,omitempty"`
	BbandsBackTestRanking []BbandsBackTestRank `json:"bbands_back_test_ranking,omitempty"`
	BestEmaPeriod1       int                   `json:"best_ema_period1,omitempty"`
	BestEmaPeriod2       int                   `json:"best_ema_period2,omitempty"`
	BestBbandsN          int                   `json:"best_bbands_n,omitempty"`
	BestBbandsK          float64               `json:"best_bbands_k,omitempty"`
	BackTestStrategy     string                `json:"back_test_strategy,omitempty"`
}

type BackTestOptions struct {
	UsePercent       float64
	StopLimitPercent float64
}

func (o BackTestOptions) withDefaults() BackTestOptions {
	if o.UsePercent <= 0 {
		o.UsePercent = config.Config.UsePercet
	}
	if o.StopLimitPercent <= 0 {
		o.StopLimitPercent = config.Config.StopLimitPercet
	}
	return o
}

type EmaBackTestRank struct {
	Period1 int     `json:"period1"`
	Period2 int     `json:"period2"`
	Profit  float64 `json:"profit"`
}

type BbandsBackTestRank struct {
	N      int     `json:"n"`
	K      float64 `json:"k"`
	Profit float64 `json:"profit"`
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

func (df *DataFrameCandle) AddEvents(timeTime time.Time) bool {
	signalEvents := GetSignalEventsAfterTime(timeTime)
	if signalEvents != nil && len(signalEvents.Signals) > 0 {
		df.Events = signalEvents
		return true
	}
	return false
}

func (df *DataFrameCandle) tradeSize(price float64, usePercent float64) float64 {
	if price <= 0 || usePercent <= 0 {
		return 0
	}
	base := df.Candles[0].Close
	if base <= 0 {
		base = price
	}
	return base * usePercent / price
}

func (s *SignalEvents) holdingBuyPrice() float64 {
	if len(s.Signals) == 0 {
		return 0
	}
	last := s.Signals[len(s.Signals)-1]
	if last.Side == "BUY" {
		return last.Price
	}
	return 0
}

func (df *DataFrameCandle) BackTestEma(period1, period2 int, opts BackTestOptions) *SignalEvents {
	opts = opts.withDefaults()
	if period1 > period2 {
		period1, period2 = period2, period1
	}
	if period1 <= 0 || period2 <= 0 || period1 >= period2 {
		return nil
	}

	lenCandles := len(df.Candles)
	if lenCandles <= period1 || lenCandles <= period2 {
		return nil
	}

	signalEvents := NewSignalEvents()
	emaValue1 := talib.Ema(df.Close(), period1)
	emaValue2 := talib.Ema(df.Close(), period2)
	for i := 0; i < lenCandles; i++ {
		if i < period1 || i < period2 {
			continue
		}
		candle := df.Candles[i]

		buyPrice := signalEvents.holdingBuyPrice()
		if buyPrice > 0 && candle.Low <= buyPrice*opts.StopLimitPercent {
			stopPrice := buyPrice * opts.StopLimitPercent
			size := df.tradeSize(stopPrice, opts.UsePercent)
			signalEvents.Sell(df.ProductCode, candle.Time, stopPrice, size, false)
			continue
		}

		if emaValue1[i-1] < emaValue2[i-1] && emaValue1[i] >= emaValue2[i] {
			size := df.tradeSize(candle.Close, opts.UsePercent)
			signalEvents.Buy(df.ProductCode, candle.Time, candle.Close, size, false)
		}
		if emaValue1[i-1] > emaValue2[i-1] && emaValue1[i] <= emaValue2[i] {
			size := df.tradeSize(candle.Close, opts.UsePercent)
			signalEvents.Sell(df.ProductCode, candle.Time, candle.Close, size, false)
		}
	}
	return signalEvents
}

func (df *DataFrameCandle) OptimizeEmaBackTest(minPeriod, maxPeriod, topN int, opts BackTestOptions) []EmaBackTestRank {
	opts = opts.withDefaults()
	if topN <= 0 {
		topN = config.Config.NumRanking
	}
	if minPeriod <= 0 {
		minPeriod = 5
	}
	lenCandles := len(df.Candles)
	if lenCandles <= minPeriod {
		return nil
	}
	if maxPeriod <= 0 || maxPeriod > lenCandles/2 {
		maxPeriod = lenCandles / 2
	}
	if maxPeriod > 50 {
		maxPeriod = 50
	}
	if maxPeriod <= minPeriod {
		return nil
	}

	var results []EmaBackTestRank
	for period1 := minPeriod; period1 < maxPeriod; period1++ {
		for period2 := period1 + 1; period2 <= maxPeriod; period2++ {
			events := df.BackTestEma(period1, period2, opts)
			if events == nil || len(events.Signals) == 0 {
				continue
			}
			results = append(results, EmaBackTestRank{
				Period1: period1,
				Period2: period2,
				Profit:  events.Profit(),
			})
		}
	}
	if len(results) == 0 {
		return nil
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Profit > results[j].Profit
	})
	if len(results) > topN {
		results = results[:topN]
	}
	return results
}

func TriggerLatestEmaSignal(productCode string, duration time.Duration, limit, period1, period2 int, opts BackTestOptions) {
	df, err := GetAllCandles(productCode, duration, limit)
	if err != nil || len(df.Candles) == 0 {
		return
	}

	events := df.BackTestEma(period1, period2, opts)
	if events == nil || len(events.Signals) == 0 {
		return
	}

	lastSignal := events.Signals[len(events.Signals)-1]
	lastCandle := df.Candles[len(df.Candles)-1]
	if !lastSignal.Time.Equal(lastCandle.Time) {
		return
	}

	stored := GetSignalEventsByCount(100)
	if stored == nil {
		stored = NewSignalEvents()
	}
	if lastSignal.Side == "BUY" {
		stored.Buy(productCode, lastSignal.Time, lastSignal.Price, lastSignal.Size, true)
	} else {
		stored.Sell(productCode, lastSignal.Time, lastSignal.Price, lastSignal.Size, true)
	}
}

func ResolveBestEmaPeriods(productCode string, duration time.Duration, limit int, opts BackTestOptions) (int, int) {
	df, err := GetAllCandles(productCode, duration, limit)
	if err != nil || len(df.Candles) == 0 {
		return 12, 26
	}
	ranking := df.OptimizeEmaBackTest(5, 50, 1, opts)
	if len(ranking) == 0 {
		return 12, 26
	}
	return ranking[0].Period1, ranking[0].Period2
}

func (df *DataFrameCandle) applyStopLoss(signalEvents *SignalEvents, candle Candle, opts BackTestOptions) bool {
	buyPrice := signalEvents.holdingBuyPrice()
	if buyPrice > 0 && candle.Low <= buyPrice*opts.StopLimitPercent {
		stopPrice := buyPrice * opts.StopLimitPercent
		size := df.tradeSize(stopPrice, opts.UsePercent)
		signalEvents.Sell(df.ProductCode, candle.Time, stopPrice, size, false)
		return true
	}
	return false
}

func (df *DataFrameCandle) BackTestBbands(n int, k float64, opts BackTestOptions) *SignalEvents {
	opts = opts.withDefaults()
	if n <= 0 || k <= 0 {
		return nil
	}

	lenCandles := len(df.Candles)
	if lenCandles <= n {
		return nil
	}

	signalEvents := NewSignalEvents()
	closes := df.Close()
	up, _, down := talib.BBands(closes, n, k, k, talib.SMA)

	for i := 0; i < lenCandles; i++ {
		if i < n || i == 0 {
			continue
		}
		if math.IsNaN(up[i]) || math.IsNaN(down[i]) || math.IsNaN(up[i-1]) || math.IsNaN(down[i-1]) {
			continue
		}

		candle := df.Candles[i]
		if df.applyStopLoss(signalEvents, candle, opts) {
			continue
		}

		// 下限バンドからの反発で買い、上限バンドからの反落で売り（逆張り）
		if closes[i-1] <= down[i-1] && closes[i] > down[i] {
			size := df.tradeSize(candle.Close, opts.UsePercent)
			signalEvents.Buy(df.ProductCode, candle.Time, candle.Close, size, false)
		}
		if closes[i-1] >= up[i-1] && closes[i] < up[i] {
			size := df.tradeSize(candle.Close, opts.UsePercent)
			signalEvents.Sell(df.ProductCode, candle.Time, candle.Close, size, false)
		}
	}
	return signalEvents
}

func (df *DataFrameCandle) OptimizeBbandsBackTest(minN, maxN int, topN int, opts BackTestOptions) []BbandsBackTestRank {
	opts = opts.withDefaults()
	if topN <= 0 {
		topN = config.Config.NumRanking
	}
	if minN <= 0 {
		minN = 10
	}
	lenCandles := len(df.Candles)
	if lenCandles <= minN {
		return nil
	}
	if maxN <= 0 || maxN > lenCandles/2 {
		maxN = lenCandles / 2
	}
	if maxN > 40 {
		maxN = 40
	}
	if maxN < minN {
		return nil
	}

	var results []BbandsBackTestRank
	for n := minN; n <= maxN; n++ {
		for k10 := 15; k10 <= 30; k10 += 5 {
			k := float64(k10) / 10.0
			events := df.BackTestBbands(n, k, opts)
			if events == nil || len(events.Signals) == 0 {
				continue
			}
			results = append(results, BbandsBackTestRank{
				N:      n,
				K:      k,
				Profit: events.Profit(),
			})
		}
	}
	if len(results) == 0 {
		return nil
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Profit > results[j].Profit
	})
	if len(results) > topN {
		results = results[:topN]
	}
	return results
}

func TriggerLatestBbandsSignal(productCode string, duration time.Duration, limit, n int, k float64, opts BackTestOptions) {
	df, err := GetAllCandles(productCode, duration, limit)
	if err != nil || len(df.Candles) == 0 {
		return
	}

	events := df.BackTestBbands(n, k, opts)
	if events == nil || len(events.Signals) == 0 {
		return
	}

	lastSignal := events.Signals[len(events.Signals)-1]
	lastCandle := df.Candles[len(df.Candles)-1]
	if !lastSignal.Time.Equal(lastCandle.Time) {
		return
	}

	stored := GetSignalEventsByCount(100)
	if stored == nil {
		stored = NewSignalEvents()
	}
	if lastSignal.Side == "BUY" {
		stored.Buy(productCode, lastSignal.Time, lastSignal.Price, lastSignal.Size, true)
	} else {
		stored.Sell(productCode, lastSignal.Time, lastSignal.Price, lastSignal.Size, true)
	}
}

func ResolveBestBbandsParams(productCode string, duration time.Duration, limit int, opts BackTestOptions) (int, float64) {
	df, err := GetAllCandles(productCode, duration, limit)
	if err != nil || len(df.Candles) == 0 {
		return 20, 2.0
	}
	ranking := df.OptimizeBbandsBackTest(10, 40, 1, opts)
	if len(ranking) == 0 {
		return 20, 2.0
	}
	return ranking[0].N, ranking[0].K
}