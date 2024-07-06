package database

import (
	"database/sql"
	"errors"

	"final/task"

)

func Addtasktodb(db *sql.DB, task task.Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, errors.New("Ошибка добавления задачи")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New("Ошибка добавления задачи")
	}
	return id, nil
}
