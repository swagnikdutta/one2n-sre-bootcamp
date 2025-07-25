package main

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
}

func (s *Server) listStudents(w http.ResponseWriter, request *http.Request) {

}

func (s *Server) addStudent(w http.ResponseWriter, request *http.Request) {

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
