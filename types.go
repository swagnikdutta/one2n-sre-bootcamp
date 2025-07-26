package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type contextKey string

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
	rows, err := s.db.Query(listSyntax)
	if err != nil {
		log.Printf("Error querying students. Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var id, age int
		var name string

		if err := rows.Scan(&id, &name, &age); err != nil {
			log.Printf("Error scanning student row. Error: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		students = append(students, Student{
			Name: name,
			Age:  age,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(students)
	if err != nil {
		log.Printf("Error encoding response. Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) addStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusNotFound)
	}

	var student Student

	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, "request body missing or invalid", http.StatusBadRequest)
		return
	}

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
	id := strings.Split(strings.Trim(r.URL.Path, "/"), "/")[1]
	ctx := context.WithValue(r.Context(), studentIdKey, id)

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
	studentId, ok := r.Context().Value(studentIdKey).(string)
	if !ok {
		log.Println("Error doing type assertion on student id")
		http.Error(w, "invalid or missing student id", http.StatusBadRequest)
		return
	}

	row := s.db.QueryRow(selectStudentSyntax, studentId)

	var id, age int
	var name string

	if err := row.Scan(&id, &name, &age); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No student found with id %q", studentId)
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}

		log.Printf("Error scanning student row: %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	student := Student{name, age}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(student)
	if err != nil {
		log.Println("Error encoding response. Error: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) updateStudent(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) deleteStudent(w http.ResponseWriter, r *http.Request) {
	studentId, ok := r.Context().Value(studentIdKey).(string)
	if !ok {
		log.Println("Error doing type assertion on student id")
		http.Error(w, "invalid or missing student id", http.StatusBadRequest)
		return
	}

	res, err := s.db.Exec(deleteSyntax, studentId)
	if err != nil {
		log.Printf("Error deleting student with id %q", studentId)
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		log.Println("Student not found, no rows were affected")
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
