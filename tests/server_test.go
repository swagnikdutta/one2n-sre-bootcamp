package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		t.Errorf("expected status %d, got %d", http.StatusCreated, statusGot)
	}
}
