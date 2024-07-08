package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"final/constant"
	"final/nextdate"
	"final/storage"
	"final/task"
)

type Handlers struct {
	Db storage.DB
}

func (h *Handlers) AddTask() http.HandlerFunc {
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
		err = task.Checktitle()
		if err != nil {
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		taskmod, err := task.Checkdate()
		if err != nil {
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		err = taskmod.Countdate()
		if err != nil {
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		id, err := h.Db.Addtasktodb(taskmod)
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
func (h *Handlers) ChangeTask() http.HandlerFunc {
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
		errstr := task.CheckId()
		if errstr != "" {
			response := map[string]interface{}{
				"error": errstr,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		err = task.Checktitle()
		if err != nil {
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		task, err = task.Checkdate()
		if err != nil {
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		errstr = task.CheckRepeate()
		if errstr != "" {
			response := map[string]interface{}{
				"error": errstr,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		errstr = h.Db.Update(task)
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

func (h *Handlers) GetTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.FormValue("id")
		if id == "" {
			response := map[string]interface{}{
				"error": "Не указан идентификатор",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		task, err := h.Db.Findtask(id)
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
}

func (h *Handlers) TaskDone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан индентификатор"}`, http.StatusInternalServerError)
			return
		}
		task, err := h.Db.Findtask(id)
		if err != "" {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusInternalServerError)
			return
		}
		if task.Repeat == "" {
			// Удаляем одноразовую задачу
			err = h.Db.Deletetask(id)
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
			err = h.Db.Updatetask(date, id)
			if err != "" {
				http.Error(w, `{"error": "Ошибка обновления задачи"}`, http.StatusInternalServerError)
				return
			}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}

func (h *Handlers) DeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан индентификатор"}`, http.StatusInternalServerError)
			return
		}
		err := h.Db.DeleteQuery(id)
		if err != "" {
			http.Error(w, `{"error": "Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}

func (db *Handlers) NextDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// получаю данные из запроса
		now := r.FormValue("now")
		date := r.FormValue("date")
		repeat := r.FormValue("repeat")
		//проверка корректности полученных данных
		if repeat == "" || now == "" || date == "" {
			http.Error(w, "Указаны некорректные данные в запросе", http.StatusBadRequest)
			return
		}
		nextdate, err := nextdate.CalcNextDate(now, date, repeat)
		if err != nil {
			http.Error(w, "Указаны некорректные данные в запросе", http.StatusBadRequest)
			return
		}
		w.Write([]byte(nextdate))
	}
}

func (h *Handlers) ReceiveTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		tasks, err := h.Db.GetTasks()
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
