package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var webDir = "./web/"

func main() {
	// Get the TODO_PORT environment variable
	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "7540"
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}

}
