package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"final/receivetasks"
)

func ReceiveTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		tasks, err := receivetasks.GetRow(db)
		if err != nil {
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		response := map[string]interface{}{
			"tasks": tasks,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}
