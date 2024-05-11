package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/adettelle/go_final_project/pkg/db"
	"github.com/adettelle/go_final_project/pkg/db/repo"
	"github.com/adettelle/go_final_project/pkg/handlers"
	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

var webDir = "./web/"

func MethodSwitch(w http.ResponseWriter, r *http.Request, tRepository repo.TasksRepository) { // , router *chi.Mux, api *handlers.Api
	api := handlers.NewApi(&tRepository)

	router := chi.NewRouter()
	router.Handle("/*", http.FileServer(http.Dir(webDir)))
	router.Get("/api/nextdate", handlers.GetNextDay)

	if r.Method != http.MethodPost { // почему != ??? Разве не наоборот???
		router.Post("/api/task", api.CreateTask)
	} else if r.Method != http.MethodGet { // почему != ??? Разве не наоборот???
		router.Get("/api/task", api.CreateTask)
	}
}

func main() {
	db.DbConnection()

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Println(err)
		return
	}
	tRepository := repo.NewTasksRepository(db)

	api := handlers.NewApi(&tRepository)

	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", handlers.GetNextDay)

	// api, которое ожидается в этом задании
	// 1. POST /api/task создает таск
	// 2. GET /api/task возвращает ошибку (судя по тестам это желаемое поведение)
	// 3. GET /api/tasks возвращает набор тасков без фильтрации
	// 4. GET /api/tasks?search=... возвращает набор тасков с фильрацией по параметру search
	// 5. GET /api/tasks/{id} возвращает таск по id
	r.HandleFunc("/api/task", api.TaskHandler) // get и post и put
	r.Get("/api/tasks", api.GetTasksHandler)
	r.Get("/api/tasks/{id}", api.GetTaskByIdHandler) // http://localhost:7540/api/tasks/257
	// r.Put("/api/task", api.PutHandler)
	// r.Post("/api/task", api.CreateTask)
	// r.Get("/api/task", api.CreateTask)

	if err := http.ListenAndServe(":7540", r); err != nil {
		log.Printf("Start server error: %s", err.Error())
	}
}
