package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

const (
	studentIdKey     contextKey = "studentId"
	sqliteDriverName            = "sqlite3"
	dbPath                      = "DB_PATH"
)

func NewRequestMultiplexer(server *Server) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/students", server.listStudents)
	mux.HandleFunc("/api/v1/students/add", server.addStudent)
	mux.HandleFunc("/api/v1/students/{id}", server.studentHandler)
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

	sqliteStore := NewSQLiteDataStore()
	server := NewServer(sqliteStore)
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
