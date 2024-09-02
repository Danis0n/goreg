package server

import (
	"testing"
)

// TestNewServerConfig проверяет создание ServerConfig с различными значениями портов.
func TestNewServerConfig(t *testing.T) {
	tests := []struct {
		name      string
		port      int
		wantError bool
	}{
		{
			name:      "Valid port",
			port:      8080,
			wantError: false,
		},
		{
			name:      "Invalid port (negative)",
			port:      -1,
			wantError: true,
		},
		{
			name:      "Invalid port (zero)",
			port:      0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewServerConfig(tt.port)
			if (err != nil) != tt.wantError {
				t.Errorf("NewServerConfig() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && cfg.Port != tt.port {
				t.Errorf("NewServerConfig() cfg.Port = %v, want %v", cfg.Port, tt.port)
			}
		})
	}
}

// TestValidateServerConfig проверяет валидность существующего ServerConfig.
func TestValidateServerConfig(t *testing.T) {
	tests := []struct {
		name      string
		cfg       ServerConfig
		wantError bool
	}{
		{
			name:      "Valid config",
			cfg:       ServerConfig{Port: 8080},
			wantError: false,
		},
		{
			name:      "Invalid config (negative port)",
			cfg:       ServerConfig{Port: -1},
			wantError: true,
		},
		{
			name:      "Invalid config (zero port)",
			cfg:       ServerConfig{Port: 0},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateServerConfig(tt.cfg)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateServerConfig() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestValidateServerSettings проверяет валидность различных значений порта.
func TestValidateServerSettings(t *testing.T) {
	tests := []struct {
		name      string
		port      int
		wantError bool
	}{
		{
			name:      "Valid port",
			port:      8080,
			wantError: false,
		},
		{
			name:      "Invalid port (negative)",
			port:      -1,
			wantError: true,
		},
		{
			name:      "Invalid port (zero)",
			port:      0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateServerSettings(tt.port)
			if (err != nil) != tt.wantError {
				t.Errorf("validateServerSettings() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
