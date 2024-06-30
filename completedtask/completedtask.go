package completedtask
import (
	"final/task"
	"encoding/json"
	"net/http"
	"database/sql"
	"time"
	"final/nextdate"
	_ "github.com/mattn/go-sqlite3"
)

func updatetask(date string, id string) string{
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		return "Ошибка запроса к бд"
	}
	defer db.Close()
	updateQuery := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err = db.Exec(updateQuery, date, id)
	if err != nil {
		return "Ошибка обновления задачи"
	}
	return ""
}

func deletetask(id string) string{
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		return "Ошибка запроса к бд"
	}
	defer db.Close()
	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
			_, err = db.Exec(deleteQuery, id)
			if err != nil {
				return "Ошибка удаления задачи"
			}
			return ""
}

func findtask(task task.Task, id string) (task.Task, string) {
	db, err := sql.Open("sqlite3", "scheduler.db")
		if err != nil {
			return task, "Ошибка запроса к бд"
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

func Completedtask (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, `{"error": "Не указан индентификатор"}`, http.StatusInternalServerError)
		return	
	}
	var task task.Task
	task, err := findtask(task, id)
	if err != "" {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusInternalServerError)
		return
	}
		if task.Repeat == "" {
			// Удаляем одноразовую задачу
			err = deletetask(id)
			if err != ""{
			http.Error(w, `{"error": "Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
			} 
		} else {
			// Рассчитываем следующую дату для периодической задачи
			now := time.Now()
			timeNow := now.Format("20060102")
			date, errnotstr := nextdate.CalcNextDate(timeNow, task.Date, task.Repeat)
			if errnotstr != nil {
				http.Error(w, `{"error": "Ошибка вычисления следующей даты"}`, http.StatusInternalServerError)
				return
			}

			// Обновляем дату задачи
			err = updatetask(date, id)
			if err != "" {
				http.Error(w, `{"error": "Ошибка обновления задачи"}`, http.StatusInternalServerError)
				return
			}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
}