package errors

import "fmt"

// Error is a struct that contains the error information
type Error struct {
	// Code is a code of the error
	code Code
	// Message is a message of the error only for developers
	message string
	// Reason is a reason of the error for users
	reason string
	// Meta is a metadata of the error
	meta map[string]any
}

// NewError creates a new error
func NewError(code Code, message string, reason string, meta map[string]any) *Error {
	if meta == nil {
		meta = make(map[string]any)
	}

	return &Error{
		code:    code,
		message: message,
		reason:  reason,
		meta:    meta,
	}
}

// NewEmptyError creates a new empty error
func NewEmptyError() *Error {
	return &Error{
		code:    0,
		message: "",
		reason:  "",
		meta:    make(map[string]any),
	}
}

// WithCode sets the code of the error
func (e *Error) WithCode(code Code) *Error {
	e.code = code
	return e
}

// WithMessage sets the message of the error
func (e *Error) WithMessage(message string) *Error {
	e.message = message
	return e
}

// WithReason sets the reason of the error
func (e *Error) WithReason(reason string) *Error {
	e.reason = reason
	return e
}

// AddMeta adds a metadata to the error
func (e *Error) AddMeta(key string, value any) *Error {
	e.meta[key] = value
	return e
}

// WithMeta adds metadata to the error
func (e *Error) WithMeta(meta map[string]any) *Error {
	e.meta = meta
	return e
}

// Error returns a string representation of the error
func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s, reason: %s, meta: %v",
		e.code,
		e.message,
		e.reason,
		e.meta,
	)
}

// ToHTTPResponse converts the error to HTTP status and payload.
func (e *Error) ToHTTPResponse() (int, map[string]any) {
	payload := map[string]any{
		"message": e.message,
		"reason":  e.reason,
	}

	if e.meta != nil {
		payload["meta"] = e.meta
	}

	return e.code.ToHTTPStatus(), payload
}
