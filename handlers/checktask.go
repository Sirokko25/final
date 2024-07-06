package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"final/database"

)

func GetTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.FormValue("id")
		if id == "" {
			response := map[string]interface{}{
				"error": "Не указан идентификатор",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		task, err := database.Findtask(db, id)
		if err != "" {
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		json.NewEncoder(w).Encode(task)
		return
	}
}
