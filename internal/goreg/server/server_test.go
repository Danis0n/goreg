package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func setupTestServer() *Server {
	logger, _ := zap.NewDevelopment()
	store, _ := NewServerStore(logger)
	server := &Server{
		logger:      logger,
		store:       store,
		errch:       make(chan error),
		closeCh:     make(chan struct{}),
		closeDoneCh: make(chan struct{}),
		port:        8080, // Порт можно задать произвольно
	}
	return server
}

func TestSetHandler(t *testing.T) {
	server := setupTestServer()

	svc := Service{Name: "testService", Callback: "http://callback.url"}
	body, _ := json.Marshal(svc)

	req, err := http.NewRequest(http.MethodPost, "/set", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.SetHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

func TestGetHandler(t *testing.T) {
	server := setupTestServer()

	// Установим сервис, чтобы было что получать
	server.store.Set("testService", "http://callback.url")

	req, err := http.NewRequest(http.MethodGet, "/get?name=testService", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var service Service
	err = json.NewDecoder(rr.Body).Decode(&service)
	if err != nil {
		t.Fatal(err)
	}

	if service.Name != "testService" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			service.Name, "testService")
	}
}

func TestGetAllHandler(t *testing.T) {
	server := setupTestServer()

	// Добавим несколько сервисов
	server.store.Set("service1", "http://callback1.url")
	server.store.Set("service2", "http://callback2.url")

	req, err := http.NewRequest(http.MethodGet, "/getall", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetAllHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var services []*Service
	err = json.NewDecoder(rr.Body).Decode(&services)
	if err != nil {
		t.Fatal(err)
	}

	if len(services) != 2 {
		t.Errorf("handler returned unexpected number of services: got %v want %v",
			len(services), 2)
	}
}

func TestDeleteHandler(t *testing.T) {
	server := setupTestServer()

	// Добавим сервис, который будем удалять
	server.store.Set("testService", "http://callback.url")

	req, err := http.NewRequest(http.MethodDelete, "/delete?name=testService", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.DeleteHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Проверяем, что сервис удален
	_, err = server.store.Get("testService")
	if err == nil {
		t.Errorf("expected error when getting deleted service, got nil")
	}
}

func TestInvalidMethod(t *testing.T) {
	server := setupTestServer()

	req, err := http.NewRequest(http.MethodPut, "/get", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code for invalid method: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}
