package db

import (
	"strconv"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Tasks возвращает ближайшие задачи, отсортированные по дате по возрастанию.
func Tasks(limit int) ([]*Task, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := DB.Query(
		`SELECT id, date, title, comment, repeat
		   FROM scheduler
		   ORDER BY date
		   LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*Task
	for rows.Next() {
		var (
			id      int64
			date    string
			title   string
			comment string
			repeat  string
		)
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			return nil, err
		}
		out = append(out, &Task{
			ID:      strconv.FormatInt(id, 10),
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// чтобы в JSON было [] а не null
	if out == nil {
		out = []*Task{}
	}
	return out, nil
}
