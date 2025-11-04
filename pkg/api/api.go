package api

import (
	"fmt"
	"net/http"
	"time"
)

// NextDateHandler — обработчик GET /api/nextdate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateLayout, nowStr)
		if err != nil {
			http.Error(w, "некорректная дата параметра now", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, next)
}
