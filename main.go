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
	r.Post("/api/task", api.CreateTask)

	if err := http.ListenAndServe(":7540", r); err != nil {
		log.Printf("Start server error: %s", err.Error())
	}
}
