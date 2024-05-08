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

	if r.Method != http.MethodPost {
		router.Post("/api/task", api.CreateTask)
	} else if r.Method != http.MethodGet {
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

	// if http.Request.Method.Post {
	// 	r.Post("/api/task", api.CreateTask)
	// } else if http.Request.Method != http.MethodGet {
	// 	r.Get("/api/task", api.CreateTask)
	// }
	r.Get("/api/tasks", api.MyHandle)
	r.Post("/api/task", api.MyHandle)
	// r.Post("/api/task", api.CreateTask)
	// r.Get("/api/task", api.CreateTask)

	if err := http.ListenAndServe(":7540", r); err != nil {
		log.Printf("Start server error: %s", err.Error())
	}
}
