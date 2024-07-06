package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"final/constant"
	"final/database"
	"final/nextdate"

)

func TaskDone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан индентификатор"}`, http.StatusInternalServerError)
			return
		}
		task, err := database.Findtask(db, id)
		if err != "" {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusInternalServerError)
			return
		}
		if task.Repeat == "" {
			// Удаляем одноразовую задачу
			err = database.Deletetask(db, id)
			if err != "" {
				http.Error(w, `{"error": "Ошибка удаления задачи"}`, http.StatusInternalServerError)
				return
			}
		} else {
			// Рассчитываем следующую дату для периодической задачи
			now := time.Now()
			timeNow := now.Format(constant.ParseDate)
			date, errnotstr := nextdate.CalcNextDate(timeNow, task.Date, task.Repeat)
			if errnotstr != nil {
				http.Error(w, `{"error": "Ошибка вычисления следующей даты"}`, http.StatusInternalServerError)
				return
			}

			// Обновляем дату задачи
			err = database.Updatetask(db, date, id)
			if err != "" {
				http.Error(w, `{"error": "Ошибка обновления задачи"}`, http.StatusInternalServerError)
				return
			}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}
