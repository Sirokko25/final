package database

import "database/sql"

func Updatetask(db *sql.DB, date string, id string) string {
	updateQuery := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err := db.Exec(updateQuery, date, id)
	if err != nil {
		return "Ошибка обновления задачи"
	}
	return ""
}
