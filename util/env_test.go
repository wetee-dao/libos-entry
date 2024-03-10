package util

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Test when environment variable is set
	os.Setenv("TEST_KEY", "test_value")
	expected := "test_value"
	result := GetEnv("TEST_KEY", "default_value")
	if result != expected {
		t.Errorf("GetEnv() = %s; expected %s", result, expected)
	}

	// Test when environment variable is not set
	os.Unsetenv("TEST_KEY")
	expected = "default_value"
	result = GetEnv("TEST_KEY", "default_value")
	if result != expected {
		t.Errorf("GetEnv() = %s; expected %s", result, expected)
	}
}
