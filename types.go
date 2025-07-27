package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
		logMessage := LogMessage(fmt.Sprintf("Error querying students from database. Error: %v", err))
		clientMessage := ClientMessage("Error fetching students")
		RespondWithError(w, logMessage, clientMessage, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var id, age int
		var name string

		if err := rows.Scan(&id, &name, &age); err != nil {
			logMessage := LogMessage(fmt.Sprintf("Error scanning student row. Error: %v", err.Error()))
			clientMessage := ClientMessage("Internal error")
			RespondWithError(w, logMessage, clientMessage, http.StatusInternalServerError)
			return
		}

		students = append(students, Student{name, age})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(students)
	if err != nil {
		logMessage := LogMessage(fmt.Sprintf("Error encoding response. Error: %v", err))
		clientMessage := ClientMessage("Internal error")
		RespondWithError(w, logMessage, clientMessage, http.StatusInternalServerError)
	}
}

func (s *Server) addStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusNotFound)
	}

	var student Student

	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		logMessage := LogMessage(fmt.Sprintf("Error unmarshalling request body. Error: %v", err))
		clientMessage := ClientMessage("Invalid request body")
		RespondWithError(w, logMessage, clientMessage, http.StatusInternalServerError)
		return
	}

	stmt, err := s.db.Prepare(insertSyntax)
	if err != nil {
		logMessage := LogMessage(fmt.Sprintf("Error preparing insert statement. Error: %v", err))
		clientMessage := ClientMessage("Internal error")
		RespondWithError(w, logMessage, clientMessage, http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(nil, student.Name, student.Age)
	if err != nil {
		logMessage := LogMessage(fmt.Sprintf("Failed to insert student. Error: %v", err))
		clientMessage := ClientMessage("Failed to insert student")
		RespondWithError(w, logMessage, clientMessage, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("Student added successfully"))
}

func (s *Server) studentHandler(w http.ResponseWriter, r *http.Request) {
	splits := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	ctx := context.WithValue(r.Context(), studentIdKey, splits[len(splits)-1])

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
		logMessage := LogMessage("Error asserting type of studentId")
		clientMessage := ClientMessage("Internal error")
		RespondWithError(w, logMessage, clientMessage, http.StatusBadRequest)
		return
	}

	row := s.db.QueryRow(selectStudentSyntax, studentId)

	var id, age int
	var name string

	if err := row.Scan(&id, &name, &age); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logMessage := LogMessage(fmt.Sprintf("No student found with id %q", studentId))
			clientMessage := ClientMessage("Student not found")
			RespondWithError(w, logMessage, clientMessage, http.StatusNotFound)
			return
		}

		logMessage := LogMessage(fmt.Sprintf("Error scanning student row. Error: %v", err))
		clientMessage := ClientMessage("Internal error")
		RespondWithError(w, logMessage, clientMessage, http.StatusInternalServerError)
		return
	}

	student := Student{name, age}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(student)
	if err != nil {
		logMessage := LogMessage(fmt.Sprintf("Error encoding response. Error: %v", err.Error()))
		clientMessage := ClientMessage("Internal error")
		RespondWithError(w, logMessage, clientMessage, http.StatusInternalServerError)
		return
	}
}

func (s *Server) updateStudent(w http.ResponseWriter, r *http.Request) {
	studentId, ok := r.Context().Value(studentIdKey).(string)
	if !ok {
		logMessage := LogMessage("Error asserting type of studentId")
		clientMessage := ClientMessage("Internal error")
		RespondWithError(w, logMessage, clientMessage, http.StatusBadRequest)
		return
	}

	var student Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		logMessage := LogMessage(fmt.Sprintf("Error reading request body. Error: %v", err))
		clientMessage := ClientMessage("Invalid request body")
		RespondWithError(w, logMessage, clientMessage, http.StatusBadRequest)
		return
	}

	res, err := s.db.Exec(updateSyntax, student.Name, student.Age, studentId)
	if err != nil {
		logMessage := LogMessage(fmt.Sprintf("Error updating student. Error: %v", err))
		clientMessage := ClientMessage("Error updating student")
		RespondWithError(w, logMessage, clientMessage, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		logMessage := LogMessage("Student not found, no rows were affected")
		clientMessage := ClientMessage("Student not found")
		RespondWithError(w, logMessage, clientMessage, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) deleteStudent(w http.ResponseWriter, r *http.Request) {
	studentId, ok := r.Context().Value(studentIdKey).(string)
	if !ok {
		logMessage := LogMessage("Error asserting type of studentId")
		clientMessage := ClientMessage("Internal error")
		RespondWithError(w, logMessage, clientMessage, http.StatusBadRequest)
		return
	}

	res, err := s.db.Exec(deleteSyntax, studentId)
	if err != nil {
		logMessage := LogMessage(fmt.Sprintf("Error deleting student with id %q", studentId))
		clientMessage := ClientMessage("Failed to delete student")
		RespondWithError(w, logMessage, clientMessage, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		logMessage := LogMessage("Student not found, no rows were affected")
		clientMessage := ClientMessage("Student not found")
		RespondWithError(w, logMessage, clientMessage, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
