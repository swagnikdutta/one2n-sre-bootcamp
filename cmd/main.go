package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/swagnikdutta/one2n-sre-bootcamp/student"
)

func NewRequestMultiplexer(server *student.Server) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/students", server.ListStudents)
	mux.HandleFunc("/api/v1/students/add", server.CreateStudent)
	mux.HandleFunc("/api/v1/students/{id}", server.StudentHandler)
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return mux
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	sqliteStore := student.NewSQLiteDataStore()
	server := student.NewServer(sqliteStore)
	httpServer := &http.Server{
		Addr:    ":8000",
		Handler: NewRequestMultiplexer(server),
	}

	fmt.Println("Running http server on port 8000")
	err = httpServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
