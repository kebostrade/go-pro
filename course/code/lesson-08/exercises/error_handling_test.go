package exercises

import (
	"errors"
	"testing"
)

func TestDivide(t *testing.T) {
	tests := []struct {
		a        float64
		b        float64
		expected float64
		wantErr  bool
	}{
		{10, 2, 5, false},
		{10, 0, 0, true},
		{0, 5, 0, false},
		{-10, 2, -5, false},
		{10, -2, -5, false},
	}

	for _, tt := range tests {
		result, err := Divide(tt.a, tt.b)
		if (err != nil) != tt.wantErr {
			t.Errorf("Divide(%v, %v) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && result != tt.expected {
			t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestValidationError(t *testing.T) {
	err := ValidationError{Field: "email", Message: "invalid format"}
	expected := "email: invalid format"
	if err.Error() != expected {
		t.Errorf("ValidationError.Error() = %s, want %s", err.Error(), expected)
	}
}

func TestValidateUser(t *testing.T) {
	tests := []struct {
		email    string
		username string
		password string
		wantErr  error
	}{
		{"test@example.com", "john", "password123", nil},
		{"", "john", "password123", ErrInvalidEmail},
		{"test@example.com", "", "password123", ErrEmptyUsername},
		{"test@example.com", "john", "short", ErrPasswordTooShort},
		{"", "", "", ErrInvalidEmail}, // First error wins
	}

	for _, tt := range tests {
		err := ValidateUser(tt.email, tt.username, tt.password)
		if err != tt.wantErr {
			t.Errorf("ValidateUser(%q, %q, %q) = %v, want %v", tt.email, tt.username, tt.password, err, tt.wantErr)
		}
	}
}

func TestSafeDivide(t *testing.T) {
	result, err := SafeDivide(10, 2)
	if err != nil {
		t.Errorf("SafeDivide(10, 2) error = %v", err)
	}
	if result != 5 {
		t.Errorf("SafeDivide(10, 2) = %v, want 5", result)
	}

	_, err = SafeDivide(10, 0)
	if err == nil {
		t.Error("SafeDivide(10, 0) expected error")
	}
}

func TestIsValidationError(t *testing.T) {
	customErr := ValidationError{Field: "test", Message: "error"}
	regularErr := errors.New("regular error")
	
	if !IsValidationError(customErr) {
		t.Error("IsValidationError should return true for ValidationError")
	}
	
	if IsValidationError(regularErr) {
		t.Error("IsValidationError should return false for regular error")
	}
}

func TestValidateForm(t *testing.T) {
	err := ValidateForm("", "", "")
	if err == nil {
		t.Error("Expected error for empty form")
	}
	
	err = ValidateForm("John", "john@example.com", "1234567890")
	if err != nil {
		t.Errorf("Expected no error for valid form, got %v", err)
	}
}
