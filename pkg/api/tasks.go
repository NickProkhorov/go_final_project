package api

import (
	"net/http"

	"github.com/NickProkhorov/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	// базовый вариант без поиска, лимит фиксируем 50
	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка получения задач"})
		return
	}
	writeJSON(w, TasksResp{Tasks: tasks})
}
