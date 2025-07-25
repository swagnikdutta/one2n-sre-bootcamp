package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Server struct {
	db *sql.DB
}

func NewServer(db *sql.DB) *Server {
	s := &Server{
		db: db,
	}
	return s
}

func (s *Server) initDB() {
	_, err := s.db.Exec(createTableSyntax)
	if err != nil {
		panic(err)
	}
}

func (s *Server) listStudents(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) addStudent(w http.ResponseWriter, r *http.Request) {
	var student Student

	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Student: %+v\n", student)

	stmt, err := s.db.Prepare(insertSyntax)
	if err != nil {
		http.Error(w, "Failed to prepare insert statement", http.StatusInternalServerError)
		log.Printf("Failed to prepare insert statement. Error: %v", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(nil, student.Name, student.Age)
	if err != nil {
		http.Error(w, "Failed to insert student", http.StatusInternalServerError)
		log.Printf("Failed to insert student. Error: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("Student added successfully"))
}

func (s *Server) studentHandler(w http.ResponseWriter, r *http.Request) {
	print(r.URL.Path)
	ctx := context.WithValue(r.Context(), studentId, r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		s.getStudent(w, r.WithContext(ctx))
	case http.MethodPatch:
		s.updateStudent(w, r.WithContext(ctx))
	case http.MethodDelete:
		s.deleteStudent(w, r.WithContext(ctx))
	}
}

func (s *Server) getStudent(w http.ResponseWriter, r *http.Request) {
	studentId, ok := r.Context().Value(studentId).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Invalid path param"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("Student id got:%q", studentId)))
}

func (s *Server) updateStudent(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) deleteStudent(w http.ResponseWriter, r *http.Request) {

}
