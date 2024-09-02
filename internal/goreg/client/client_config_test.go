package client

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewClientConfigWithDefaults(t *testing.T) {
	registrator := "http://registrator.url"
	callback := "http://callback.url"
	port := 8080

	cfg, err := NewClientConfigWithDefaults(registrator, callback, port)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Registrator != registrator {
		t.Errorf("Expected registrator %s, got %s", registrator, cfg.Registrator)
	}

	if cfg.Callback != callback {
		t.Errorf("Expected callback %s, got %s", callback, cfg.Callback)
	}

	if _, err := uuid.Parse(cfg.Name); err != nil {
		t.Errorf("Expected valid UUID for name, got %s", cfg.Name)
	}

	if cfg.Port != port {
		t.Errorf("Expected port %d, got %d", port, cfg.Port)
	}
}

func TestNewClientConfigWithDefaults_InvalidSettings(t *testing.T) {
	_, err := NewClientConfigWithDefaults("", "", -1)
	if err == nil {
		t.Fatal("Expected error for invalid settings, got nil")
	}
}

func TestNewClientConfigWithName(t *testing.T) {
	registrator := "http://registrator.url"
	callback := "http://callback.url"
	port := 8080
	name := "test-client"

	cfg, err := NewClientConfigWithName(registrator, callback, port, name)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Registrator != registrator {
		t.Errorf("Expected registrator %s, got %s", registrator, cfg.Registrator)
	}

	if cfg.Callback != callback {
		t.Errorf("Expected callback %s, got %s", callback, cfg.Callback)
	}

	if cfg.Name != name {
		t.Errorf("Expected name %s, got %s", name, cfg.Name)
	}

	if cfg.Port != port {
		t.Errorf("Expected port %d, got %d", port, cfg.Port)
	}
}

func TestNewClientConfigWithName_InvalidSettings(t *testing.T) {
	_, err := NewClientConfigWithName("", "", -1, "test-client")
	if err == nil {
		t.Fatal("Expected error for invalid settings, got nil")
	}
}

func TestValidateClientConfig(t *testing.T) {
	cfg := ClientConfig{
		Registrator: "http://registrator.url",
		Callback:    "http://callback.url",
		Name:        "test-client",
		Port:        8080,
	}

	err := ValidateClientConfig(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestValidateClientConfig_Invalid(t *testing.T) {
	cfg := ClientConfig{
		Registrator: "",
		Callback:    "",
		Name:        "test-client",
		Port:        -1,
	}

	err := ValidateClientConfig(cfg)
	if err == nil {
		t.Fatal("Expected error for invalid config, got nil")
	}
}
