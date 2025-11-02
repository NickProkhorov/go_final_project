package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/NickProkhorov/go_final_project/pkg/db"
)

// TaskDoneHandler — POST /api/task/done?id=...
func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// 1) Получаем задачу
	t, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// 2) Если не периодическая — удалить
	if strings.TrimSpace(t.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]any{}) // {}
		return
	}

	// 3) Периодическая — вычислить следующую дату и обновить
	now := time.Now()
	next, err := NextDate(now, t.Date, t.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": "некорректное правило повторения"})
		return
	}
	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]any{}) // {}
}
