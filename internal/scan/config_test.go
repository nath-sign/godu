package scan

import (
	"testing"
)

func TestConfigValidate(t *testing.T) {
	config := NewConfig("test", true, true, true)
	if err := config.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestConfigValidateError(t *testing.T) {
	config := NewConfig("", true, true, true)
	if err := config.Validate(); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
