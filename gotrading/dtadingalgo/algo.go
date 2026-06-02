package dtadingalgo

import "math"

// IchimokuInput is OHLC series used to calculate Ichimoku lines.
type IchimokuInput struct {
	High  []float64
	Low   []float64
	Close []float64
}

// IchimokuParams uses standard default periods (9, 26, 52, displacement 26).
type IchimokuParams struct {
	TenkanPeriod  int
	KikunPeriod   int
	SenkouBPeriod int
	Displacement  int
}

// IchimokuLines holds each line aligned with the input candle length.
type IchimokuLines struct {
	Tenkan  []float64
	Kikun   []float64
	SenkouA []float64
	SenkouB []float64
	Chikou  []float64
}

// DefaultIchimokuParams returns the commonly used Ichimoku settings.
func DefaultIchimokuParams() IchimokuParams {
	return IchimokuParams{
		TenkanPeriod:  9,
		KikunPeriod:   26,
		SenkouBPeriod: 52,
		Displacement:  26,
	}
}

// Ichimoku calculates the five Ichimoku Kinko Hyo lines.
func Ichimoku(in IchimokuInput, p IchimokuParams) IchimokuLines {
	n := len(in.Close)
	empty := IchimokuLines{
		Tenkan:  make([]float64, n),
		Kikun:   make([]float64, n),
		SenkouA: make([]float64, n),
		SenkouB: make([]float64, n),
		Chikou:  make([]float64, n),
	}
	if n == 0 {
		return empty
	}
	for i := range empty.Tenkan {
		empty.Tenkan[i] = math.NaN()
		empty.Kikun[i] = math.NaN()
		empty.SenkouA[i] = math.NaN()
		empty.SenkouB[i] = math.NaN()
		empty.Chikou[i] = math.NaN()
	}

	tenkan := donchianMid(in.High, in.Low, p.TenkanPeriod)
	kikun := donchianMid(in.High, in.Low, p.KikunPeriod)
	senkouBBase := donchianMid(in.High, in.Low, p.SenkouBPeriod)

	copyLine(empty.Tenkan, tenkan)
	copyLine(empty.Kikun, kikun)

	for i := 0; i < n; i++ {
		if !valid(tenkan[i]) || !valid(kikun[i]) {
			continue
		}
		j := i + p.Displacement
		if j < n {
			empty.SenkouA[j] = (tenkan[i] + kikun[i]) / 2
		}
	}

	for i := 0; i < n; i++ {
		if !valid(senkouBBase[i]) {
			continue
		}
		j := i + p.Displacement
		if j < n {
			empty.SenkouB[j] = senkouBBase[i]
		}
	}

	for i := 0; i < n; i++ {
		j := i + p.Displacement
		if j < n {
			empty.Chikou[i] = in.Close[j]
		}
	}

	return empty
}

func donchianMid(high, low []float64, period int) []float64 {
	n := len(high)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.NaN()
		if period <= 0 || i < period-1 {
			continue
		}
		maxH := high[i-period+1]
		minL := low[i-period+1]
		for j := i - period + 2; j <= i; j++ {
			if high[j] > maxH {
				maxH = high[j]
			}
			if low[j] < minL {
				minL = low[j]
			}
		}
		out[i] = (maxH + minL) / 2
	}
	return out
}

func copyLine(dst, src []float64) {
	for i := range dst {
		if i < len(src) {
			dst[i] = src[i]
		}
	}
}

func valid(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

// HistoricalVolatility calculates annualized historical volatility in percent.
// It uses the standard deviation of log returns over period, scaled by sqrt(periodsPerYear).
func HistoricalVolatility(closes []float64, period int, periodsPerYear float64) []float64 {
	n := len(closes)
	out := make([]float64, n)
	for i := range out {
		out[i] = math.NaN()
	}
	if period <= 1 || n <= period || periodsPerYear <= 0 {
		return out
	}

	logReturns := make([]float64, n)
	for i := 1; i < n; i++ {
		if closes[i-1] > 0 && closes[i] > 0 {
			logReturns[i] = math.Log(closes[i] / closes[i-1])
		} else {
			logReturns[i] = math.NaN()
		}
	}

	annualScale := math.Sqrt(periodsPerYear) * 100
	for i := period; i < n; i++ {
		std := stdDev(logReturns[i-period+1 : i+1])
		if valid(std) {
			out[i] = std * annualScale
		}
	}
	return out
}

func stdDev(values []float64) float64 {
	var sum float64
	count := 0
	for _, v := range values {
		if !valid(v) {
			continue
		}
		sum += v
		count++
	}
	if count <= 1 {
		return math.NaN()
	}

	mean := sum / float64(count)
	var sqDiff float64
	for _, v := range values {
		if !valid(v) {
			continue
		}
		diff := v - mean
		sqDiff += diff * diff
	}
	return math.Sqrt(sqDiff / float64(count-1))
}
