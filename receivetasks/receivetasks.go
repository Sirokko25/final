package receivetasks
import (
	"encoding/json"
	"net/http"
	"final/task"
	"database/sql"
	"errors"
)

func getrow() ([]task.Task, error){
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		return nil, errors.New("Ошибка выполнения запроса")
	}
	defer db.Close()
	rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50`)
	if err != nil {
		return nil, errors.New("Ошибка выполнения запроса: ")
	}
	defer rows.Close()

	tasks := make([]task.Task, 0 , 0)

	for rows.Next() {
		var task task.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, errors.New("Ошибка чтения строки: ")
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("Ошибка обработки результата: ")
	}
	return tasks, nil
}

func Receivetasks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	
	tasks, err := getrow()
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