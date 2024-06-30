package deletetask
import (
	"encoding/json"
	"net/http"
	"database/sql"
)

func deleteQuery(id string) string{
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		return "Ошибка выполнения запроса"
	}
	defer db.Close()
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

func Deletetask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "Не указан индентификатор"}`, http.StatusInternalServerError)
		return	
	}
	err := deleteQuery(id)
	if err != ""{
		http.Error(w, `{"error": "Ошибка удаления задачи"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{})
	return
}