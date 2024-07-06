package receivetasks

import (
	"database/sql"
	"errors"

	"final/constant"
	"final/task"
)

func GetRow(db *sql.DB) ([]task.Task, error) {
	rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`, sql.Named("limit", constant.Limit))
	if err != nil {
		return nil, errors.New("Ошибка выполнения запроса: ")
	}
	defer rows.Close()

	tasks := make([]task.Task, 0, 0)

	for rows.Next() {
		var task task.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, errors.New("Ошибка чтения строки: ")
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("Ошибка обработки результата: ")
	}
	return tasks, nil
}
