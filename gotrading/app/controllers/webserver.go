package controllers

import (
	"encoding/json"
	"fmt"
	"gotrading/app/models"
	"gotrading/config"
	"html/template"
	"net/http"
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

func StartWebServer() error {
	http.HandleFunc("/chart", viewChartHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}
