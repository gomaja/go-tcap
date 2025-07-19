package tcap

import (
	"errors"
	"fmt"
)

// Error types for better error handling and debugging
var (
	// ErrInvalidInput indicates invalid input data
	ErrInvalidInput = errors.New("invalid input data")

	// ErrIndefiniteLength indicates non-DER encoding with indefinite length
	ErrIndefiniteLength = errors.New("indefinite length found (not DER)")

	// ErrEmptyData indicates empty or nil data
	ErrEmptyData = errors.New("empty data provided")

	// ErrUnmarshalFailed indicates ASN.1 unmarshal failure
	ErrUnmarshalFailed = errors.New("failed to unmarshal ASN.1 data")

	// ErrMarshalFailed indicates ASN.1 marshal failure
	ErrMarshalFailed = errors.New("failed to marshal ASN.1 data")

	// ErrInvalidTransactionID indicates invalid transaction ID
	ErrInvalidTransactionID = errors.New("invalid transaction ID")

	// ErrInvalidInvokeID indicates invalid invoke ID
	ErrInvalidInvokeID = errors.New("invalid invoke ID")
)

// ParseError represents a parsing error with additional context
type ParseError struct {
	Op    string // operation that failed
	Field string // field that caused the error
	Err   error  // underlying error
}

func (e *ParseError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("parse error in %s.%s: %v", e.Op, e.Field, e.Err)
	}
	return fmt.Sprintf("parse error in %s: %v", e.Op, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// ValidationError represents a validation error with context
type ValidationError struct {
	Field string
	Value interface{}
	Err   error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s (value: %v): %v", e.Field, e.Value, e.Err)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// Convenience functions for creating structured errors
func newParseError(op, field string, err error) error {
	return &ParseError{Op: op, Field: field, Err: err}
}

func newValidationError(field string, value interface{}, err error) error {
	return &ValidationError{Field: field, Value: value, Err: err}
}
