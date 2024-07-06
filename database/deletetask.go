package database

import "database/sql"

func Deletetask(db *sql.DB, id string) string {
	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	_, err := db.Exec(deleteQuery, id)
	if err != nil {
		return "Ошибка удаления задачи"
	}
	return ""
}
