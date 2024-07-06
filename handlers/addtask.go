package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"final/addtask"
	"final/database"
	"final/task"
)

func AddTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var task task.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		err = addtask.Checktitle(task)
		if err != nil {
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		taskmod, err := addtask.Checkdate(task)
		if err != nil {
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		err = addtask.Checkrulerepeat(taskmod)
		if err != nil {
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		id, err := database.Addtasktodb(db, taskmod)
		if err != nil {
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		response := map[string]interface{}{
			"id": id,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}
