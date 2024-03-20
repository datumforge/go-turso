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
