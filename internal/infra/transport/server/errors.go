package server

import pkgErr "github.com/rlapenok/rybakov_test/pkg/errors"

const (
	// reasonStartServer is the reason for the error when the server fails to start
	reasonStartServer = "FAILED_START_SERVER"
	reasonAuthError   = "AUTH_ERROR"
)

var (
	// ErrStartServer is the error when the server fails to start
	ErrStartServer = pkgErr.NewEmptyError().WithCode(pkgErr.CodeInternalError).WithReason(reasonStartServer)
	ErrAuthError   = pkgErr.NewEmptyError().WithCode(pkgErr.CodeUnauthorized).WithReason(reasonAuthError)
)
