package database

import (
	"database/sql"
	"final/task"
)

func ExecQuery(db *sql.DB, task task.Task) string {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return "Ошибка выполнения запроса"
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return "Ошибка получения результата запроса"
	}

	if rowsAffected == 0 {
		return "Задача не найдена"
	}
	return ""
}
