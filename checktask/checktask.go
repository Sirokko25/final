package checktask
import(
	"final/task"
	"encoding/json"
	"net/http"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Findtask(id string) (task.Task, string) {
	var task task.Task
	db, err := sql.Open("sqlite3", "scheduler.db")
		if err != nil {
			return task, "Ошибка выполнения запроса"
		}
		defer db.Close()
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err = db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, "Задача не найдена"
		} else {
			return task, "Ошибка выполнения запроса"
		}
		
	}
	return task, ""
}

func Checktask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	id := r.FormValue("id")
	if id == "" {
		response := map[string]interface{}{
			"error": "Не указан идентификатор",
		}
		json.NewEncoder(w).Encode(response)
		return	
	}
	task, err :=Findtask(id)
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