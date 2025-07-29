package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type contextKey string

type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Server struct {
	store StudentStore
}

func NewServer(s StudentStore) *Server {
	// TODO: pass in a logger
	srv := &Server{
		store: s,
	}
	return srv
}

func (s *Server) listStudents(w http.ResponseWriter, r *http.Request) {
	students, err := s.store.ListStudents()
	if err != nil {
		RespondWithError(w, "Failed to list students", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(students)
	if err != nil {
		log.Printf("Error encoding response. Error: %v", err)
		RespondWithError(w, "Internal error", http.StatusInternalServerError)
	}
}

func (s *Server) addStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		RespondWithError(w, "", http.StatusNotFound)
	}

	var student Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		log.Printf("Error unmarshalling request body")
		RespondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := s.store.CreateStudent(student)
	if err != nil {
		log.Printf("Error creating student. Error: %v", err)
		RespondWithError(w, "Error creating student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	// TODO: why not set application type content/json
	_, _ = w.Write([]byte(fmt.Sprintf("Student created with id: %s", id)))
}

func (s *Server) studentHandler(w http.ResponseWriter, r *http.Request) {
	splits := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	studentId, err := strconv.Atoi(splits[len(splits)-1])
	if err != nil {
		log.Printf("Error type casting studentId %q to integer. Error: %v", studentId, err)
		http.Error(w, "Invalid studentId %q", studentId)
		return
	}

	ctx := context.WithValue(r.Context(), studentIdKey, studentId)

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
	studentId, ok := r.Context().Value(studentIdKey).(int)
	if !ok {
		log.Printf("Error asserting type of student id %q", studentId)
		RespondWithError(w, "Invalid studentId", http.StatusBadRequest)
		return
	}

	student, err := s.store.GetStudent(studentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msg := fmt.Sprintf("No student found with id %q", studentId)
			log.Println(msg)
			RespondWithError(w, msg, http.StatusNotFound)
			return
		}

		msg := fmt.Sprintf("Error getting student with id %q. Error: %v", studentId, err)
		log.Println(msg)
		RespondWithError(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(student); err != nil {
		log.Printf("Error encoding response. Error: %v", err)
		RespondWithError(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) updateStudent(w http.ResponseWriter, r *http.Request) {
	studentId, ok := r.Context().Value(studentIdKey).(int)
	if !ok {
		log.Printf("Error asserting type of student id %q", studentId)
		RespondWithError(w, "Invalid studentId", http.StatusBadRequest)
		return
	}

	var payload Student
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Error unmarshalling request body")
		RespondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := s.store.UpdateStudent(studentId, payload)
	if err != nil {
		msg := fmt.Sprintf("Error updating student with id %q. Error: %v", studentId, err)
		log.Println(msg)
		RespondWithError(w, msg, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) deleteStudent(w http.ResponseWriter, r *http.Request) {
	studentId, ok := r.Context().Value(studentIdKey).(int)
	if !ok {
		log.Printf("Error asserting type of student id %q", studentId)
		RespondWithError(w, "Invalid studentId", http.StatusBadRequest)
		return
	}

	if err := s.store.DeleteStudent(studentId); err != nil {
		log.Printf("Error deleting student with id %d. Error: %v", studentId, err)
		RespondWithError(w, "Error deleting student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte("student deleted"))
}
