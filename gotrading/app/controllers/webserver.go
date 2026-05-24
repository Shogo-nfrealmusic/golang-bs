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

var templates = template.Must(template.ParseFiles("app/views/google.html"))

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

	err = templates.ExecuteTemplate(w, "google.html", chartPageData{
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
