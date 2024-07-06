package deletetask

import (
	"database/sql"
)

func DeleteQuery(db *sql.DB, id string) string {
	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(deleteQuery, id)
	if err != nil {
		return "Ошибка выполнения запроса"
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return "Ошибка получения результата запроса"
	}
	if rowsAffected == 0 {
		return "Запись не найдена"
	}
	return ""
}
