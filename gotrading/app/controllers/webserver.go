package controllers

import (
	"encoding/json"
	"fmt"
	"gotrading/app/models"
	"gotrading/config"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var templates = template.Must(template.ParseFiles("app/views/chart.html"))

type chartPageData struct {
	CandlesJSON template.JS
}

func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	limit := 100
	duration := "1s"
	durationTime := config.Config.Durations[duration]
	df, err := models.GetAllCandles(config.Config.ProductCode, durationTime, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type candleRow struct {
		Time  string  `json:"time"`
		Low   float64 `json:"low"`
		Open  float64 `json:"open"`
		Close float64 `json:"close"`
		High  float64 `json:"high"`
	}

	rows := make([]candleRow, 0, len(df.Candles))
	for _, c := range df.Candles {
		rows = append(rows, candleRow{
			Time:  c.Time.Format("2006-01-02 15:04:05"),
			Low:   c.Low,
			Open:  c.Open,
			Close: c.Close,
			High:  c.High,
		})
	}

	raw, err := json.Marshal(rows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "chart.html", chartPageData{
		CandlesJSON: template.JS(raw),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type JSONEncoder struct {
	Error string `json:"error"`
	Code int `json:"code"`
}

func APIError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonError, marshalErr := json.Marshal(JSONEncoder{Error: err.Error(), Code: code})
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}
	w.Write(jsonError)
}

var apiVaidPaths = regexp.MustCompile(`^/api/candle/$`)

func apiMakeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := apiVaidPaths.FindStringSubmatch(r.URL.Path)
		if m == nil {
			APIError(w, fmt.Errorf("invalid path: %s", r.URL.Path), http.StatusNotFound)
			return
		}
		fn(w, r)
	}
}

func apiCandleHandler(w http.ResponseWriter, r *http.Request) {
	productCode := r.URL.Query().Get("product_code")
	if productCode == "" {
		APIError(w, fmt.Errorf("product_code is required"), http.StatusBadRequest)
		return
	}

	limit := 100
	if config.Config.BackTest && config.Config.DataLimit > 0 {
		limit = config.Config.DataLimit
	}
	if strLimit := r.URL.Query().Get("limit"); strLimit != "" {
		n, err := strconv.Atoi(strLimit)
		if err != nil || n <= 0 || n > 1000 {
			APIError(w, fmt.Errorf("limit must be between 1 and 1000"), http.StatusBadRequest)
			return
		}
		limit = n
	}

	duration := r.URL.Query().Get("duration")
	if duration == "" {
		duration = "1m"
	}
	durationTime, ok := config.Config.Durations[duration]
	if !ok {
		APIError(w, fmt.Errorf("invalid duration: %s", duration), http.StatusBadRequest)
		return
	}

	df, err := models.GetAllCandles(productCode, durationTime, limit)
	if err != nil {
		APIError(w, err, http.StatusInternalServerError)
		return
	}

	sma := r.URL.Query().Get("sma")
	if sma != "" {
		strSmaPeriod1 := r.URL.Query().Get("smaPeriod1")
		strSmaPeriod2 := r.URL.Query().Get("smaPeriod2")
		strSmaPeriod3 := r.URL.Query().Get("smaPeriod3")
		period1, err := strconv.Atoi(strSmaPeriod1)
		if strSmaPeriod1 == "" || err != nil || period1 < 0 {
			period1 = 7
		}
		period2, err := strconv.Atoi(strSmaPeriod2)
		if strSmaPeriod2 == "" || err != nil || period2 < 0 {
			period2 = 14
		}
		period3, err := strconv.Atoi(strSmaPeriod3)
		if strSmaPeriod3 == "" || err != nil || period3 < 0 {
			period3 = 21
		}
		df.AddSma(period1)
		df.AddSma(period2)
		df.AddSma(period3)
	}

	emaPeriod1 := 12
	emaPeriod2 := 26
	ema := r.URL.Query().Get("ema")
	if ema != "" {
		strEmaPeriod1 := r.URL.Query().Get("emaPeriod1")
		strEmaPeriod2 := r.URL.Query().Get("emaPeriod2")
		strEmaPeriod3 := r.URL.Query().Get("emaPeriod3")
		period1, err := strconv.Atoi(strEmaPeriod1)
		if strEmaPeriod1 == "" || err != nil || period1 < 0 {
			period1 = 12
		}
		period2, err := strconv.Atoi(strEmaPeriod2)
		if strEmaPeriod2 == "" || err != nil || period2 < 0 {
			period2 = 26
		}
		period3, err := strconv.Atoi(strEmaPeriod3)
		if strEmaPeriod3 == "" || err != nil || period3 < 0 {
			period3 = 50
		}
		emaPeriod1 = period1
		emaPeriod2 = period2
		df.AddEma(period1)
		df.AddEma(period2)
		df.AddEma(period3)
	}

	bbandsN := 20
	bbandsK := 2.0
	bbands := r.URL.Query().Get("bbands")
	if bbands != "" {
		n := 20
		k := 2.0
		if strN := r.URL.Query().Get("bbandsN"); strN != "" {
			if parsed, err := strconv.Atoi(strN); err == nil && parsed > 0 {
				n = parsed
			}
		}
		if strK := r.URL.Query().Get("bbandsK"); strK != "" {
			if parsed, err := strconv.ParseFloat(strK, 64); err == nil && parsed > 0 {
				k = parsed
			}
		}
		bbandsN = n
		bbandsK = k
		df.AddBbands(n, k)
	}

	ichimoku := r.URL.Query().Get("ichimoku")
	if ichimoku != "" {
		tenkan := 9
		kikun := 26
		senkouB := 52
		displacement := 26
		if v, err := strconv.Atoi(r.URL.Query().Get("ichimokuTenkan")); err == nil && v > 0 {
			tenkan = v
		}
		if v, err := strconv.Atoi(r.URL.Query().Get("ichimokuKikun")); err == nil && v > 0 {
			kikun = v
		}
		if v, err := strconv.Atoi(r.URL.Query().Get("ichimokuSenkouB")); err == nil && v > 0 {
			senkouB = v
		}
		if v, err := strconv.Atoi(r.URL.Query().Get("ichimokuDisplacement")); err == nil && v > 0 {
			displacement = v
		}
		df.AddIchimoku(tenkan, kikun, senkouB, displacement)
	}

	rsi := r.URL.Query().Get("rsi")
	if rsi != "" {
		period := 14
		if v, err := strconv.Atoi(r.URL.Query().Get("rsiPeriod")); err == nil && v > 0 {
			period = v
		}
		df.AddRsi(period)
	}

	macd := r.URL.Query().Get("macd")
	if macd != "" {
		fast := 12
		slow := 26
		signal := 9
		if v, err := strconv.Atoi(r.URL.Query().Get("macdFast")); err == nil && v > 0 {
			fast = v
		}
		if v, err := strconv.Atoi(r.URL.Query().Get("macdSlow")); err == nil && v > 0 {
			slow = v
		}
		if v, err := strconv.Atoi(r.URL.Query().Get("macdSignal")); err == nil && v > 0 {
			signal = v
		}
		df.AddMacd(fast, slow, signal)
	}

	hv := r.URL.Query().Get("hv")
	if hv != "" {
		period := 20
		if v, err := strconv.Atoi(r.URL.Query().Get("hvPeriod")); err == nil && v > 0 {
			period = v
		}
		df.AddHistoricalVolatility(period)
	}

	events := r.URL.Query().Get("events")
	if events != "" && len(df.Candles) > 0 {
		if config.Config.BackTest {
			opts := models.BackTestOptions{
				UsePercent:       config.Config.UsePercet,
				StopLimitPercent: config.Config.StopLimitPercet,
			}
			strategy := r.URL.Query().Get("strategy")
			if strategy == "" {
				strategy = "ema"
			}
			optimize := r.URL.Query().Get("optimize") != "false"
			df.BackTestStrategy = strategy

			switch strategy {
			case "bbands":
				if optimize {
					ranking := df.OptimizeBbandsBackTest(10, 40, config.Config.NumRanking, opts)
					df.BbandsBackTestRanking = ranking
					if len(ranking) > 0 {
						best := ranking[0]
						df.BestBbandsN = best.N
						df.BestBbandsK = best.K
						df.Events = df.BackTestBbands(best.N, best.K, opts)
					}
				} else {
					df.Events = df.BackTestBbands(bbandsN, bbandsK, opts)
					df.BestBbandsN = bbandsN
					df.BestBbandsK = bbandsK
				}
			default:
				if optimize {
					ranking := df.OptimizeEmaBackTest(5, 50, config.Config.NumRanking, opts)
					df.BackTestRanking = ranking
					if len(ranking) > 0 {
						best := ranking[0]
						df.BestEmaPeriod1 = best.Period1
						df.BestEmaPeriod2 = best.Period2
						df.Events = df.BackTestEma(best.Period1, best.Period2, opts)
					}
				} else {
					df.Events = df.BackTestEma(emaPeriod1, emaPeriod2, opts)
					df.BestEmaPeriod1 = emaPeriod1
					df.BestEmaPeriod2 = emaPeriod2
				}
			}
		} else {
			df.AddEvents(df.Candles[0].Time)
		}
	}

	js, err := json.Marshal(df)
	if err != nil {
		APIError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}


func StartWebServer() error {
	http.HandleFunc("/api/candle/", apiMakeHandler(apiCandleHandler))
	http.HandleFunc("/chart", viewChartHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}
