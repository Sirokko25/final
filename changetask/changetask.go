package changetask
import (
	"final/task"
	"encoding/json"
	"net/http"
	"database/sql"
	"final/nextdate"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func checkId(task task.Task) string {
	if task.ID == ""{
		return "Не указан индентификатор задачи"
	} else {
		return ""
	}
}

func checkTitle(task task.Task) string {
	if task.Title == ""{
		return "Не указан заголовок задачи"
	} else {
		return ""
	}
}

func CheckDate(task task.Task) (task.Task, string) {
	now := time.Now()
			if task.Date == "" {
				task.Date = now.Format("20060102")
				return task, ""
			} else {
				date, err := time.Parse("20060102", task.Date)
				if err != nil {
					return task, "Неправильный формат даты"
				}
				if date.Before(now) {
					if task.Repeat == "" {
						task.Date = now.Format("20060102")
						return task, ""
					} else {
						nowtime := now.Format("20060102")
						nextDate, err := nextdate.CalcNextDate(nowtime, task.Date, task.Repeat)
						if err != nil {
							return task, "Ошибка вычисления даты"
						}
						task.Date = nextDate
						return task, ""
						}
						
				}
				return task, ""
			}
}

func checkRepeate(task task.Task) string {
	if task.Repeat != ""{
		_, err := nextdate.Parsingrepeatrules(task.Repeat)
		if err != nil{
			return "Правило повторения указано в неправильном формате"
		}
	}
	return ""
}

func execQuery(task task.Task) string {
	db, err := sql.Open("sqlite3", "scheduler.db")
		if err != nil {
			return "Ошибка выполнения запроса"
		}
	defer db.Close()
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			return "Ошибка выполнения запроса"
		}
	rowsAffected, err := res.RowsAffected()
		if err != nil {
			return "Ошибка получения результата запроса"
			}

	if rowsAffected == 0 {
		return "Задача не найдена"
	}
	return ""
}

func Changetask(w http.ResponseWriter, r *http.Request){
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
	errstr := checkId(task)
	if errstr != "" {
		response := map[string]interface{}{
			"error": errstr,
			}
		json.NewEncoder(w).Encode(response)
		return
	}
	errstr = checkTitle(task)
	if errstr != ""{
		response := map[string]interface{}{
			"error": errstr,
			}
		json.NewEncoder(w).Encode(response)
		return
	}
	task, errstr = CheckDate(task)
	if errstr != ""{
		response := map[string]interface{}{
			"error": errstr,
			}
		json.NewEncoder(w).Encode(response)
		return
	}
	errstr = checkRepeate(task)
	if errstr != ""{
		response := map[string]interface{}{
			"error": errstr,
			}
		json.NewEncoder(w).Encode(response)
		return
	}
	errstr = execQuery(task)
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