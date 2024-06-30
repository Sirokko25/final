package main

import (
	"final/addtask"
	"final/changetask"
	"final/checktask"
	"final/completedtask"
	"final/database"
	"final/deletetask"
	"final/nextdate"
	"final/receivetasks"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	//для 7 задачи
	r.Delete("/api/task", deletetask.Deletetask)
	r.Post("/api/task/done", completedtask.Completedtask)
	//для 6 задачи
	r.Get("/api/task", checktask.Checktask)
	r.Put("/api/task", changetask.Changetask)
	//для 5 задания
	r.Get("/api/tasks", receivetasks.Receivetasks)
	// для 4 задания
	r.Post("/api/task", addtask.Addtask)
	// для 3 задания
	r.Get("/api/nextdate", nextdate.NextDate)
	// для 2 задания
	database.Createdatabase()
	//для 1 задания
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./web/css/"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./web/js/"))))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))

	if err := http.ListenAndServe(":7540", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
