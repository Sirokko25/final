package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"final/deletetask"
)

func DeleteTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан индентификатор"}`, http.StatusInternalServerError)
			return
		}
		err := deletetask.DeleteQuery(db, id)
		if err != "" {
			http.Error(w, `{"error": "Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}
