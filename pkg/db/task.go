package db

import (
	"database/sql"
	"fmt"
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

// GetTask возвращает задачу по идентификатору (строкой).
func GetTask(id string) (*Task, error) {
	row := DB.QueryRow(
		`SELECT id, date, title, comment, repeat
		   FROM scheduler
		  WHERE id = ?`, id)

	var (
		idi     int64
		date    string
		title   string
		comment string
		repeat  string
	)
	if err := row.Scan(&idi, &date, &title, &comment, &repeat); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Задача не найдена")
		}
		return nil, err
	}
	return &Task{
		ID:      strconv.FormatInt(idi, 10),
		Date:    date,
		Title:   title,
		Comment: comment,
		Repeat:  repeat,
	}, nil
}

// UpdateTask обновляет существующую задачу по ID.
func UpdateTask(t *Task) error {
	res, err := DB.Exec(
		`UPDATE scheduler
		    SET date = ?, title = ?, comment = ?, repeat = ?
		  WHERE id = ?`,
		t.Date, t.Title, t.Comment, t.Repeat, t.ID,
	)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return fmt.Errorf("Задача не найдена")
	}
	return nil
}
