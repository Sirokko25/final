package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"final/changetask"
	"final/database"
	"final/task"
)

func ChangeTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var task task.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			response := map[string]interface{}{
				"error": "Ошибка десериализации JSON",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		errstr := changetask.CheckId(task)
		if errstr != "" {
			response := map[string]interface{}{
				"error": errstr,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		errstr = changetask.CheckTitle(task)
		if errstr != "" {
			response := map[string]interface{}{
				"error": errstr,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		task, errstr = changetask.CheckDate(task)
		if errstr != "" {
			response := map[string]interface{}{
				"error": errstr,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		errstr = changetask.CheckRepeate(task)
		if errstr != "" {
			response := map[string]interface{}{
				"error": errstr,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		errstr = database.ExecQuery(db, task)
		if errstr != "" {
			response := map[string]interface{}{
				"error": errstr,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		response := map[string]interface{}{
			"": "",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}
