package handlers

import (
	"database/sql"
	"net/http"

	"final/nextdate"
)

func NextDate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// получаю данные из запроса
		now := r.FormValue("now")
		date := r.FormValue("date")
		repeat := r.FormValue("repeat")
		//проверка корректности полученных данных
		if repeat == "" || now == "" || date == "" {
			http.Error(w, "Указаны некорректные данные в запросе", http.StatusBadRequest)
			return
		}
		nextdate, err := nextdate.CalcNextDate(now, date, repeat)
		if err != nil {
			http.Error(w, "Указаны некорректные данные в запросе", http.StatusBadRequest)
			return
		}
		w.Write([]byte(nextdate))
	}
}
