package student

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type contextKey string

type Student struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Server struct {
	Store  Store
	Logger *slog.Logger
}

func NewServer(s Store) *Server {
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)
	srv := &Server{
		Store:  s,
		Logger: logger,
	}
	return srv
}

func (s *Server) ListStudents(w http.ResponseWriter, r *http.Request) {
	students, err := s.Store.ListStudents()
	if err != nil {
		RespondWithError(w, "Failed to list students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(students)
	if err != nil {
		s.Logger.Error("error encoding response", "error", err)
		RespondWithError(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) CreateStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		RespondWithError(w, "Not Found", http.StatusNotFound)
		return
	}

	var student Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		s.Logger.Error("error unmarshalling request body", "error", err)
		RespondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := s.Store.CreateStudent(student)
	if err != nil {
		s.Logger.Error("error creating student", "error", err)
		RespondWithError(w, "Error creating student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	// TODO: why not set application type content/json
	_, _ = w.Write([]byte("student created successfully"))
}

func (s *Server) StudentHandler(w http.ResponseWriter, r *http.Request) {
	splits := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	studentId, err := strconv.Atoi(splits[len(splits)-1])
	if err != nil {
		s.Logger.Error("error type casting studentId to integer", "studentId", studentId, "error", err)
		http.Error(w, "Invalid studentId %q", studentId)
		return
	}

	ctx := context.WithValue(r.Context(), StudentIdKey, studentId)

	switch r.Method {
	case http.MethodGet:
		s.GetStudent(w, r.WithContext(ctx))
	case http.MethodPatch:
		s.UpdateStudent(w, r.WithContext(ctx))
	case http.MethodDelete:
		s.DeleteStudent(w, r.WithContext(ctx))
	}
}

func (s *Server) GetStudent(w http.ResponseWriter, r *http.Request) {
	studentId, ok := r.Context().Value(StudentIdKey).(int)
	if !ok {
		s.Logger.Error("error asserting type of studentId", "studentId", studentId)
		RespondWithError(w, "Invalid studentId", http.StatusBadRequest)
		return
	}

	student, err := s.Store.GetStudent(studentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msg := fmt.Sprintf("No student found with id %q", studentId)
			s.Logger.Error("student not found", "studentId", studentId, "error", err)
			RespondWithError(w, msg, http.StatusNotFound)
			return
		}

		s.Logger.Error("error getting student", "studentId", studentId, "error", err)
		msg := fmt.Sprintf("Error getting student with id %q. Error: %v", studentId, err)
		RespondWithError(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(student); err != nil {
		s.Logger.Error("error encoding response", "error", err)
		RespondWithError(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	studentId, ok := r.Context().Value(StudentIdKey).(int)
	if !ok {
		s.Logger.Error("error asserting type of studentId", "studentId", studentId)
		RespondWithError(w, "Invalid studentId", http.StatusBadRequest)
		return
	}

	var payload Student
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.Logger.Error("error unmarshalling request body", "error", err)
		RespondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := s.Store.UpdateStudent(studentId, payload)
	if err != nil {
		s.Logger.Error("error updating student", "studentId", studentId, "error", err)
		msg := fmt.Sprintf("Error updating student with id %q. Error: %v", studentId, err)
		RespondWithError(w, msg, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	studentId, ok := r.Context().Value(StudentIdKey).(int)
	if !ok {
		s.Logger.Error("error asserting type of studentId", "studentId", studentId)
		RespondWithError(w, "Invalid studentId", http.StatusBadRequest)
		return
	}

	if err := s.Store.DeleteStudent(studentId); err != nil {
		s.Logger.Error("error deleting student", "studentId", studentId, "error", err)

		errMsg, statusCode := "error deleting student", http.StatusInternalServerError
		if err.Error() == errStudentNotFound {
			errMsg, statusCode = "student not found", http.StatusNotFound
		}

		RespondWithError(w, errMsg, statusCode)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte("student deleted"))
}
