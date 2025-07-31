package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/swagnikdutta/one2n-sre-bootcamp/mocks"
	"github.com/swagnikdutta/one2n-sre-bootcamp/student"
	"go.uber.org/mock/gomock"
)

func TestCreateStudent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	payload := student.Student{
		Name: "Swagnik",
		Age:  32,
	}

	// json.NewEncoder is usually used when you have a writer, like in http handlers.
	// Since we don't have one here, we write it to a buffer — which of course implements
	// the io.Writer interface
	buf := new(bytes.Buffer) // note that buf is pointer
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		t.Fatalf("Error encoding: %v", err)
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/v1/students", buf)
	response := httptest.NewRecorder()

	mockStore := mocks.NewMockStore(ctrl)
	mockStore.EXPECT().CreateStudent(payload).Return(10, nil)

	s := &student.Server{
		Store: mockStore,
	}
	s.CreateStudent(response, request)

	statusWant := http.StatusCreated
	statusGot := response.Code

	if statusWant != statusGot {
		t.Errorf("expected status %d, got %d", statusWant, statusGot)
	}
}

func TestCreateStudent_Failure_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	payload := student.Student{
		Name: "Swagnik",
		Age:  32,
	}

	// json.NewEncoder is usually used when you have a writer, like in http handlers.
	// Since we don't have one here, we write it to a buffer — which of course implements
	// the io.Writer interface
	buf := new(bytes.Buffer) // note that buf is pointer
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		t.Fatalf("Error encoding: %v", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/api/v1/students", buf)
	response := httptest.NewRecorder()

	mockStore := mocks.NewMockStore(ctrl)
	// mockStore.EXPECT().CreateStudent(payload).Return(10, nil)

	s := &student.Server{
		Store: mockStore,
	}
	s.CreateStudent(response, request)

	statusWant := http.StatusNotFound
	statusGot := response.Code

	responseBodyWant := "Not Found\n"
	responseBodyGot := response.Body.String()

	if statusWant != statusGot {
		t.Errorf("expected status %d, got %d", statusWant, statusGot)
	}

	if responseBodyWant != responseBodyGot {
		t.Errorf("expected response body %q, got %q", responseBodyWant, responseBodyGot)
	}
}

func TestCreateStudent_Failure_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	payload := ""

	// json.NewEncoder is usually used when you have a writer, like in http handlers.
	// Since we don't have one here, we write it to a buffer — which of course implements
	// the io.Writer interface
	buf := new(bytes.Buffer) // note that buf is pointer
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		t.Fatalf("Error encoding: %v", err)
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/v1/students", buf)
	response := httptest.NewRecorder()

	mockStore := mocks.NewMockStore(ctrl)

	s := &student.Server{
		Store: mockStore,
	}
	s.CreateStudent(response, request)

	statusWant := http.StatusBadRequest
	statusGot := response.Code

	responseBodyWant := "Invalid request body\n"
	responseBodyGot := response.Body.String()

	if statusWant != statusGot {
		t.Errorf("expected status %d, got %d", statusWant, statusGot)
	}

	if responseBodyWant != responseBodyGot {
		t.Errorf("expected response body %q, got %q", responseBodyWant, responseBodyGot)
	}
}

func TestGetStudent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	studentId := 100
	mockResponse := student.Student{100, "Swagnik", 32}

	request, _ := http.NewRequest(http.MethodGet, "/api/v1/students/"+strconv.Itoa(studentId), nil)
	ctx := context.WithValue(request.Context(), student.StudentIdKey, studentId)
	response := httptest.NewRecorder()

	mockStore := mocks.NewMockStore(ctrl)
	mockStore.EXPECT().GetStudent(studentId).Return(&mockResponse, nil)

	s := &student.Server{
		Store: mockStore,
	}
	s.GetStudent(response, request.WithContext(ctx))

	statusWant := http.StatusOK
	statusGot := response.Code
	// {\"id\":100,\"name\":\"Swagnik\",\"age\":32}\\n
	// {\"id\":100,\"name\":\"Swagnik\",\"age\":32}\n

	responseBodyWant := `{"id":100,"name":"Swagnik","age":32}` + "\n"
	responseBodyGot := response.Body.String()

	if statusWant != statusGot {
		t.Errorf("expected status %q, got %q", statusWant, statusGot)
	}

	if responseBodyWant != responseBodyGot {
		t.Errorf("expected response body %q, got %q", responseBodyWant, responseBodyGot)
	}
}

func TestGetStudent_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	studentId := "100"
	request, _ := http.NewRequest(http.MethodGet, "/api/v1/students/"+studentId, nil)
	response := httptest.NewRecorder()

	mockStore := mocks.NewMockStore(ctrl)

	s := &student.Server{
		Store: mockStore,
	}
	s.GetStudent(response, request)

	statusWant := http.StatusBadRequest
	statusGot := response.Code

	responseBodyWant := "Invalid studentId\n"
	responseBodyGot := response.Body.String()

	if statusWant != statusGot {
		t.Errorf("expected status %d, got %d", statusWant, statusGot)
	}

	if responseBodyWant != responseBodyGot {
		t.Errorf("expected response body %q, got %q", responseBodyWant, responseBodyGot)
	}
}
