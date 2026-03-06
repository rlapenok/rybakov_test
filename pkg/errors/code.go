package errors

import "net/http"

// Code is a type that represents the error code
type Code int

// Codes are the error codes
const (
	// These codes are used as status codes in HTTP and GRPC servers
	CodeBadRequest         Code = 400
	CodeNotFound           Code = 404
	CodeInternal           Code = 500
	CodeUnauthorized       Code = 401
	CodeForbidden          Code = 403
	CodeConflict           Code = 409
	CodeUnprocessable      Code = 422
	CodeTooManyRequests    Code = 429
	CodeNotImplemented     Code = 501
	CodeServiceUnavailable Code = 503
	CodeGatewayTimeout     Code = 504

	// These codes are used for internal package errors
	CodeInternalError Code = 1001
)

func StatusCodeToCode(statusCode int) Code {
	switch statusCode {
	case http.StatusBadRequest:
		return CodeBadRequest
	case http.StatusUnauthorized:
		return CodeUnauthorized
	case http.StatusForbidden:
		return CodeForbidden
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusConflict:
		return CodeConflict
	case http.StatusUnprocessableEntity:
		return CodeUnprocessable
	case http.StatusTooManyRequests:
		return CodeTooManyRequests
	case http.StatusInternalServerError:
		return CodeInternal
	}

	return CodeInternal
}

// toHTTPStatus returns the HTTP status code for the error code
func (c Code) ToHTTPStatus() int {
	switch c {
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeUnprocessable:
		return http.StatusUnprocessableEntity
	case CodeTooManyRequests:
		return http.StatusTooManyRequests
	case CodeInternal:
		return http.StatusInternalServerError
	case CodeNotImplemented:
		return http.StatusNotImplemented
	case CodeServiceUnavailable:
		return http.StatusServiceUnavailable
	case CodeGatewayTimeout:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}
