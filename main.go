package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"

	"final/database"
	"final/handlers"
	"final/tests"
)

func main() {

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	db, err := database.Createdatabase()
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Delete("/api/task", handlers.DeleteTask(db))
	r.Post("/api/task/done", handlers.TaskDone(db))
	r.Get("/api/task", handlers.GetTask(db))
	r.Put("/api/task", handlers.ChangeTask(db))
	r.Get("/api/tasks", handlers.ReceiveTasks(db))
	r.Post("/api/task", handlers.AddTask(db))
	r.Get("/api/nextdate", handlers.NextDate(db))

	r.Handle("/*", http.FileServer(http.Dir("./web")))

	log.Printf("Сервер слушает порт %s", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
