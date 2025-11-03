package config

import (
	"testing"
)

func TestConfig_Init_LoadsEnvFile(t *testing.T) {
	// Mock environment variables instead of loading .env
	SetMockEnv(map[string]string{
		"DB_HOST":              "mockhost",
		"DB_USER":              "mockuser",
		"DB_PASSWORD":          "mockpass",
		"DB_NAME":              "mockdb",
		"DB_PORT":              "1234",
		"DB_SSLMODE":           "require",
		"DB_TIMEZONE":          "America/Sao_Paulo",
		"HTTP_ALLOWED_ORIGINS": "http://test.com",
		"HTTP_PORT":            "9999",
		"HTTP_URL":             "http://test.com",
		"ENV":                  "test",
		"ADMIN_EMAIL":          "admin@admin.com.br",
		"ADMIN_PASSWORD":       "Pass123@",
		"MAILTRAP_TOKEN":       "mocktoken",
		"MAILTRAP_URL":         "http://mockurl.com",
	})
	defer UnsetMockEnv([]string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT", "DB_SSLMODE", "DB_TIMEZONE", "HTTP_ALLOWED_ORIGINS", "HTTP_PORT", "HTTP_URL", "ENV", "ADMIN_EMAIL", "ADMIN_PASSWORD"})

	cfg, err := Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	if cfg.Database == nil || cfg.Http == nil || cfg.AdminUser == nil {
		t.Error("Config struct not initialized correctly")
	}

	if cfg.Database.Host != "mockhost" {
		t.Errorf("Expected DB_HOST to be 'mockhost', got '%s'", cfg.Database.Host)
	}
	if cfg.Http.Port != "9999" {
		t.Errorf("Expected HTTP_PORT to be '9999', got '%s'", cfg.Http.Port)
	}
	if cfg.Env != "test" {
		t.Errorf("Expected ENV to be 'test', got '%s'", cfg.Env)
	}
}

func TestConfig_Initi_GetEnvReturnsValueOrDefault(t *testing.T) {
	SetMockEnv(map[string]string{"TEST_KEY": "test_value"})
	val := getEnv("TEST_KEY", "default")
	if val != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", val)
	}
	UnsetMockEnv([]string{"TEST_KEY"})

	val = getEnv("TEST_KEY", "default")
	if val != "default" {
		t.Errorf("Expected 'default', got '%s'", val)
	}
}
