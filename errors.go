package turso

import (
	"errors"
	"fmt"
)

var (
	// ErrAPITokenNotSet is returned when the API token is not set
	ErrAPITokenNotSet = errors.New("api token not set, but required")

	// ErrInvalidDatabaseName is returned when a database name is invalid
	ErrInvalidDatabaseName = errors.New("invalid database name, can only contain lowercase letters, numbers, dashes with a maximum of 32 characters")

	// ErrExpirationNotSet is returned when the expiration is not set
	ErrExpirationInvalid = errors.New("expiration invalid, must be a valid duration (e.g. 12w) or never")

	// ErrAuthorizationInvalid is returned when the authorization is invalid
	ErrAuthorizationInvalid = errors.New("authorization invalid, valid options are full-access or read-only")
)

// TursoError is returned when a request to the Turso API fails
type TursoError struct {
	// Object is the object that the error occurred on
	Object string
	// Method is the method that the error occurred in
	Method string
	// Status is the status code of the error
	Status int
}

// Error returns the RequiredFieldMissingError in string format
func (e *TursoError) Error() string {
	return fmt.Sprintf("error %s %s: %d", e.Method, e.Object, e.Status)
}

// newBadRequestError returns an error a bad request
func newBadRequestError(object, method string, status int) *TursoError {
	return &TursoError{
		Object: object,
		Method: method,
		Status: status,
	}
}

// MissingRequiredFieldError is returned when a required field was not provided in a request
type MissingRequiredFieldError struct {
	// RequiredField that is missing
	RequiredField string
}

// Error returns the MissingRequiredFieldError in string format
func (e *MissingRequiredFieldError) Error() string {
	return fmt.Sprintf("%s is required", e.RequiredField)
}

// newMissingRequiredField returns an error for a missing required field
func newMissingRequiredFieldError(field string) *MissingRequiredFieldError {
	return &MissingRequiredFieldError{
		RequiredField: field,
	}
}

// InvalidFieldError is returned when a required field does not meet the required criteria
type InvalidFieldError struct {
	Field   string
	Message string
}

// Error returns the InvalidFieldError in string format
func (e *InvalidFieldError) Error() string {
	return fmt.Sprintf("%s is invalid, %s", e.Field, e.Message)
}

// newMissingRequiredField returns an error for a missing required field
func newInvalidFieldError(field, message string) *InvalidFieldError {
	return &InvalidFieldError{
		Field:   field,
		Message: message,
	}
}
