package main

import (
	"fmt"
	"net/http"
)

const (
	studentId = "studentId"
)

func NewRequestMultiplexer(server *Server) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/students", server.listStudents)
	mux.HandleFunc("/students/add", server.addStudent)
	mux.HandleFunc("/students/{id}", server.studentHandler)
	return mux
}

func main() {
	server := new(Server)
	httpServer := &http.Server{
		Addr:    ":8000",
		Handler: NewRequestMultiplexer(server),
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
