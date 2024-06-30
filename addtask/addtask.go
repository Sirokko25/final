package addtask
import (
	"encoding/json"
	"final/nextdate"
	"final/task"
	"net/http"
	"database/sql"
	"errors"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

func checktitle (task task.Task) error{
	if task.Title == "" {
		return errors.New("Пустой заголовок")
	}
	return nil
}
func checkdate (task task.Task) (task.Task, error) {

		now := time.Now()
			if task.Date == "" {
				task.Date = now.Format("20060102")
				return task, nil
			} else {
				date, err := time.Parse("20060102", task.Date)
				if err != nil {
					return task, errors.New("Неправильный формат даты")
				}
				if date.Before(now) {
					if task.Repeat == "" {
						task.Date = now.Format("20060102")
						return task, nil
					} else {
						nowtime := now.Format("20060102")
						nextDate, err := nextdate.CalcNextDate(nowtime, task.Date, task.Repeat)
						if err != nil {
							return task, errors.New("Ошибка вычисления даты")
						}
						task.Date = nextDate
						return task, nil
						}
						
				}
				return task, nil
			}
}


func checkrulerepeat (task task.Task) error{
	if task.Repeat != "" {
		now := time.Now()
		nowtime := now.Format("20060102")
		nextDate, err := nextdate.CalcNextDate(nowtime, task.Date, task.Repeat)
		if err != nil {
			return errors.New("Ошибка вычисления даты")
		}
		task.Date = nextDate
	}
	return nil
} 


func addtasktodb(task task.Task) (int64, error){
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		return 0, errors.New("Ошибка добавления задачи")
	}
	defer db.Close()
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, errors.New("Ошибка добавления задачи")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New("Ошибка добавления задачи")
	}
	return id, nil
}


func Addtask(w http.ResponseWriter, r *http.Request) {
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
	err = checktitle(task)
	if err != nil {
		response := map[string]interface{}{
			"error": err,
		}
		json.NewEncoder(w).Encode(response)
		return	
	}
	taskmod, err := checkdate(task)
	if err != nil {
		response := map[string]interface{}{
			"error": err,
		
		}
		json.NewEncoder(w).Encode(response)
		return	
	}
	err = checkrulerepeat(taskmod)
	if err != nil {
		response := map[string]interface{}{
			"error": err,
		
		}
		json.NewEncoder(w).Encode(response)
		return	
	}
	id, err := addtasktodb(taskmod)	
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
