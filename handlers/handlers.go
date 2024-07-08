package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"final/addtask"
	"final/database"
	"final/task"
)

type Handlers struct { // В структуре Handlers хранится структура DB
	db database.DB
}

func (h *Handlers) AddTask() http.HandlerFunc { // Функция AddTask становится методом структуры Handlers
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
		id, err := h.db.Addtasktodb(taskmod) // Вызов функции меняется на вызов метода от DB
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
