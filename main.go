package main

import (
	"log"
	"net/http"

	"github.com/adettelle/go_final_project/handlers"
	"github.com/adettelle/go_final_project/pkg/scheduler"
	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

var webDir = "./web/"

func main() {
	scheduler.DbConnection()
	r := chi.NewRouter()

	r.Get("/api/nextdate", handlers.GetNextDay)

	if err := http.ListenAndServe(":7540", r); err != nil {
		log.Printf("Start server error: %s", err.Error())
	}
}
