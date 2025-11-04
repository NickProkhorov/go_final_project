package api

import (
	"net/http"

	"github.com/NickProkhorov/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(db.DefaultTasksLimit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "ошибка получения задач"})
		return
	}
	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}
