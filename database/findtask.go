package database

import (
	"database/sql"

	"final/task"
)

func Findtask(db *sql.DB, id string) (task.Task, string) {
	var task task.Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, "Задача не найдена"
		} else {
			return task, "Ошибка выполнения запроса"
		}

	}
	return task, ""
}
