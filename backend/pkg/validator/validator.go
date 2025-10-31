// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package validator provides request validation utilities.
package validator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	"go-pro-backend/internal/errors"
)

// Validator interface defines validation methods.
type Validator interface {
	Validate(data interface{}) error
	ValidateJSON(r *http.Request, data interface{}) error
}

// validator implements the Validator interface.
type validator struct {
	rules map[reflect.Type][]ValidationRule
}

// ValidationRule defines a validation rule.
type ValidationRule struct {
	Field    string
	Required bool
	MinLen   int
	MaxLen   int
	Pattern  *regexp.Regexp
	Custom   func(interface{}) error
}

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidationErrors represents multiple validation errors.
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (ve ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "validation failed"
	}

	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}

	return strings.Join(messages, "; ")
}

// New creates a new validator.
func New() Validator {
	return &validator{
		rules: make(map[reflect.Type][]ValidationRule),
	}
}

// AddRule adds a validation rule for a specific type.
func (v *validator) AddRule(dataType reflect.Type, rule ValidationRule) {
	v.rules[dataType] = append(v.rules[dataType], rule)
}

// Validate validates the given data according to defined rules.
func (v *validator) Validate(data interface{}) error {
	if data == nil {
		return errors.NewValidationError("data cannot be nil", nil)
	}

	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)

	// Handle pointers.
	if dataType.Kind() == reflect.Ptr {
		if dataValue.IsNil() {
			return errors.NewValidationError("data cannot be nil", nil)
		}
		dataType = dataType.Elem()
		dataValue = dataValue.Elem()
	}

	rules, exists := v.rules[dataType]
	if !exists {
		// If no specific rules, perform basic validation.
		return v.validateBasic(dataValue)
	}

	var validationErrors []ValidationError

	for _, rule := range rules {
		if err := v.validateField(dataValue, rule); err != nil {
			if ve, ok := err.(ValidationErrors); ok {
				validationErrors = append(validationErrors, ve.Errors...)
			} else {
				validationErrors = append(validationErrors, ValidationError{
					Field:   rule.Field,
					Message: err.Error(),
				})
			}
		}
	}

	if len(validationErrors) > 0 {
		return errors.NewValidationError("validation failed", ValidationErrors{Errors: validationErrors})
	}

	return nil
}

// ValidateJSON validates JSON request body.
func (v *validator) ValidateJSON(r *http.Request, data interface{}) error {
	if r.Body == nil {
		return errors.NewBadRequestError("request body is required")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Strict JSON parsing

	if err := decoder.Decode(data); err != nil {
		return errors.NewValidationError("invalid JSON format", err)
	}

	return v.Validate(data)
}

// validateField validates a specific field according to a rule.
func (v *validator) validateField(dataValue reflect.Value, rule ValidationRule) error {
	fieldValue := dataValue.FieldByName(rule.Field)
	if !fieldValue.IsValid() {
		return fmt.Errorf("field %s not found", rule.Field)
	}

	// Check if required.
	if rule.Required && v.isEmpty(fieldValue) {
		return ValidationErrors{Errors: []ValidationError{{
			Field:   rule.Field,
			Message: "field is required",
			Value:   fieldValue.Interface(),
		}}}
	}

	// Skip further validation if field is empty and not required.
	if v.isEmpty(fieldValue) {
		return nil
	}

	// String validations.
	if fieldValue.Kind() == reflect.String {
		str := fieldValue.String()

		// Length validation.
		if rule.MinLen > 0 && utf8.RuneCountInString(str) < rule.MinLen {
			return ValidationErrors{Errors: []ValidationError{{
				Field:   rule.Field,
				Message: fmt.Sprintf("minimum length is %d", rule.MinLen),
				Value:   str,
			}}}
		}

		if rule.MaxLen > 0 && utf8.RuneCountInString(str) > rule.MaxLen {
			return ValidationErrors{Errors: []ValidationError{{
				Field:   rule.Field,
				Message: fmt.Sprintf("maximum length is %d", rule.MaxLen),
				Value:   str,
			}}}
		}

		// Pattern validation.
		if rule.Pattern != nil && !rule.Pattern.MatchString(str) {
			return ValidationErrors{Errors: []ValidationError{{
				Field:   rule.Field,
				Message: "format is invalid",
				Value:   str,
			}}}
		}
	}

	// Custom validation.
	if rule.Custom != nil {
		if err := rule.Custom(fieldValue.Interface()); err != nil {
			return ValidationErrors{Errors: []ValidationError{{
				Field:   rule.Field,
				Message: err.Error(),
				Value:   fieldValue.Interface(),
			}}}
		}
	}

	return nil
}

// validateBasic performs basic validation.
func (v *validator) validateBasic(dataValue reflect.Value) error {
	if dataValue.Kind() != reflect.Struct {
		return nil // Only validate structs by default
	}

	var validationErrors []ValidationError
	dataType := dataValue.Type()

	for i := 0; i < dataValue.NumField(); i++ {
		field := dataValue.Field(i)
		fieldType := dataType.Field(i)

		// Skip unexported fields.
		if !field.CanInterface() {
			continue
		}

		// Check for basic validation tags.
		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		if err := v.validateTag(field, fieldType.Name, tag); err != nil {
			if ve, ok := err.(ValidationErrors); ok {
				validationErrors = append(validationErrors, ve.Errors...)
			} else {
				validationErrors = append(validationErrors, ValidationError{
					Field:   fieldType.Name,
					Message: err.Error(),
				})
			}
		}
	}

	if len(validationErrors) > 0 {
		return ValidationErrors{Errors: validationErrors}
	}

	return nil
}

// validateTag validates a field based on validation tags.
func (v *validator) validateTag(fieldValue reflect.Value, fieldName, tag string) error {
	rules := strings.Split(tag, ",")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)

		switch {
		case rule == "required":
			if v.isEmpty(fieldValue) {
				return ValidationErrors{Errors: []ValidationError{{
					Field:   fieldName,
					Message: "field is required",
				}}}
			}
		case strings.HasPrefix(rule, "min="):
			// Handle min length validation.
			// This is a simplified version - you'd expand this for full tag support.
		}
	}

	return nil
}

// isEmpty checks if a value is considered empty.
func (v *validator) isEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.String() == ""
	case reflect.Ptr, reflect.Interface:
		return value.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		return value.Len() == 0
	default:
		return false
	}
}

// Common validation patterns.
var (
	EmailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	UUIDPattern  = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	SlugPattern  = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

// Common validation functions.
func ValidateEmail(email string) error {
	if !EmailPattern.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func ValidateUUID(uuid string) error {
	if !UUIDPattern.MatchString(uuid) {
		return fmt.Errorf("invalid UUID format")
	}

	return nil
}

func ValidateSlug(slug string) error {
	if !SlugPattern.MatchString(slug) {
		return fmt.Errorf("invalid slug format")
	}

	return nil
}
