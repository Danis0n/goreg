package server

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getTestLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = ""
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	return logger
}

func TestNewServerStore(t *testing.T) {
	logger := getTestLogger()

	store, err := NewServerStore(logger)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if store == nil {
		t.Fatalf("expected non-nil store")
	}

	if store.logger == nil {
		t.Fatalf("expected logger to be initialized")
	}
}

func TestServerStore_SetAndGet(t *testing.T) {
	logger := getTestLogger()
	store, _ := NewServerStore(logger)

	serviceName := "testService"
	callbackURL := "http://callback.url"

	// Test Set
	err := store.Set(serviceName, callbackURL)
	if err != nil {
		t.Fatalf("expected no error on Set, got %v", err)
	}

	// Test Set for existing service
	err = store.Set(serviceName, callbackURL)
	if err == nil {
		t.Fatalf("expected error on Set for existing service, got nil")
	}

	// Test Get
	service, err := store.Get(serviceName)
	if err != nil {
		t.Fatalf("expected no error on Get, got %v", err)
	}

	if service.Name != serviceName {
		t.Fatalf("expected service name %v, got %v", serviceName, service.Name)
	}

	if service.Callback != callbackURL {
		t.Fatalf("expected service callback %v, got %v", callbackURL, service.Callback)
	}

	// Test Get for non-existent service
	_, err = store.Get("nonExistentService")
	if err == nil {
		t.Fatalf("expected error on Get for non-existent service, got nil")
	}
}

func TestServerStore_GetAll(t *testing.T) {
	logger := getTestLogger()
	store, _ := NewServerStore(logger)

	store.Set("service1", "http://callback1.url")
	store.Set("service2", "http://callback2.url")

	services := store.GetAll()
	if len(services) != 2 {
		t.Fatalf("expected 2 services, got %v", len(services))
	}

	names := map[string]bool{
		"service1": false,
		"service2": false,
	}

	for _, svc := range services {
		if _, exists := names[svc.Name]; exists {
			names[svc.Name] = true
		}
	}

	for name, found := range names {
		if !found {
			t.Fatalf("expected service %v to be found", name)
		}
	}
}

func TestServerStore_Delete(t *testing.T) {
	logger := getTestLogger()
	store, _ := NewServerStore(logger)

	serviceName := "serviceToDelete"
	callbackURL := "http://callback.url"
	store.Set(serviceName, callbackURL)

	// Test successful Delete
	err := store.Delete(serviceName)
	if err != nil {
		t.Fatalf("expected no error on Delete, got %v", err)
	}

	_, err = store.Get(serviceName)
	if err == nil {
		t.Fatalf("expected error on Get for deleted service, got nil")
	}

	// Test Delete for non-existent service
	err = store.Delete(serviceName)
	if err == nil {
		t.Fatalf("expected error on Delete for non-existent service, got nil")
	}
}
