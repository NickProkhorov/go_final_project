// pkg/api/addtask.go
package api

import (
	"database/sql"
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
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "метод не поддерживается"})
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан идентификатор"})
		return
	}

	t, err := db.GetTask(id)
	if err != nil {
		// 404 если записи нет, иначе 500
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, db.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "задача не найдена"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, t)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "некорректный JSON"})
		return
	}

	if strings.TrimSpace(t.ID) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан идентификатор"})
		return
	}
	if strings.TrimSpace(t.Title) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	// та же валидация/нормализация, что и при добавлении
	if err := checkDate(&t); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&t); err != nil {
		// 404 если строка не обновлена из-за неверного id
		if errors.Is(err, db.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "задача не найдена"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{}) // {}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "некорректный JSON"})
		return
	}

	if strings.TrimSpace(t.Title) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	if err := checkDate(&t); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&t)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": strconv.FormatInt(id, 10)})
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
		// Сдвигаем только если дата в прошлом
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

// сравнивает только даты (без времени): true, если d < now (в прошлом)
func isBeforeDate(d, now time.Time) bool {
	yd, md, dd := d.Date()
	yn, mn, dn := now.Date()
	d0 := time.Date(yd, md, dd, 0, 0, 0, 0, time.UTC)
	n0 := time.Date(yn, mn, dn, 0, 0, 0, 0, time.UTC)
	return d0.Before(n0)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан идентификатор"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		if errors.Is(err, db.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "задача не найдена"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}
