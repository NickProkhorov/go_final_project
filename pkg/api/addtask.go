package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NickProkhorov/go_final_project/pkg/db"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка чтения JSON"})
		return
	}

	if strings.TrimSpace(t.Title) == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	if err := checkDate(&t); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&t)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка добавления задачи"})
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(v)
}

// валидация/нормализация даты

func checkDate(t *db.Task) error {
	now := time.Now()

	if strings.TrimSpace(t.Date) == "" {
		t.Date = now.Format(DateLayout)
	}

	d, err := time.Parse(DateLayout, t.Date)
	if err != nil {
		return errors.New("некорректная дата (ожидается 20060102)")
	}

	if strings.TrimSpace(t.Repeat) != "" {
		next, err := NextDate(now, t.Date, t.Repeat)
		if err != nil {
			return errors.New("некорректное правило повторения")
		}
		// Сдвигаем ТОЛЬКО если дата в прошлом (строго меньше сегодня)
		if isBeforeDate(d, now) {
			t.Date = next
		}
	} else {
		// Без правила: если дата в прошлом — ставим сегодняшнюю
		if isBeforeDate(d, now) {
			t.Date = now.Format(DateLayout)
		}
	}

	return nil
}

// сравнивает только даты (без времени): true, если d > now
func isBeforeDate(d, now time.Time) bool {
	yd, md, dd := d.Date()
	yn, mn, dn := now.Date()
	d0 := time.Date(yd, md, dd, 0, 0, 0, 0, time.UTC)
	n0 := time.Date(yn, mn, dn, 0, 0, 0, 0, time.UTC)
	return d0.Before(n0)
}
