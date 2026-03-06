package domain

import pkgErr "github.com/rlapenok/rybakov_test/pkg/errors"

const (
	// reasonNewAmount is the reason for the error when the amount is invalid
	reasonNewAmount = "INVALID_AMOUNT_FORMAT"
	// reasonAmountNotGreaterThanZero is the reason for the error when amount is less than or equal to zero
	reasonAmountNotGreaterThanZero = "AMOUNT_MUST_BE_GREATER_THAN_ZERO"
	// reasonAmountScaleExceeded is the reason for the error when amount has more than 2 decimal places
	reasonAmountScaleExceeded = "AMOUNT_SCALE_EXCEEDED"
	// reasonAmountOutOfRange is the reason for the error when amount exceeds DECIMAL(18,2) range
	reasonAmountOutOfRange = "AMOUNT_OUT_OF_RANGE"
	// reasonNewCurrencyType is the reason for the error when the currency type is invalid
	reasonNewCurrencyType = "INVALID_CURRENCY_TYPE_FORMAT"
	// reasonNewIdempotencyKey is the reason for the error when the idempotency key is invalid
	reasonNewIdempotencyKey = "INVALID_IDEMPOTENCY_KEY_FORMAT"
	// reasonUserNotFound is the reason for the error when user is not found
	reasonUserNotFound = "USER_NOT_FOUND"
	// reasonWithdrawalNotFound is the reason for the error when withdrawal is not found
	reasonWithdrawalNotFound = "WITHDRAWAL_NOT_FOUND"
	// reasonInsufficientBalance is the reason for the error when user balance is insufficient
	reasonInsufficientBalance = "INSUFFICIENT_BALANCE"
	// reasonIdempotencyPayloadMismatch is the reason for the error when idempotency key is reused with different payload
	reasonIdempotencyPayloadMismatch = "IDEMPOTENCY_PAYLOAD_MISMATCH"
	// reasonInvalidRequestPayload is the reason for malformed request payload
	reasonInvalidRequestPayload = "INVALID_REQUEST_PAYLOAD"
)

var (
	// ErrNewAmount is the error when the amount is invalid
	ErrNewAmount = pkgErr.NewError(pkgErr.CodeBadRequest, "Invalid amount", reasonNewAmount, nil)
	// ErrAmountNotGreaterThanZero is the error when the amount is less than or equal to zero
	ErrAmountNotGreaterThanZero = pkgErr.NewError(pkgErr.CodeBadRequest, "Amount must be greater than zero", reasonAmountNotGreaterThanZero, nil)
	// ErrAmountScaleExceeded is the error when the amount has more than 2 decimal places
	ErrAmountScaleExceeded = pkgErr.NewError(pkgErr.CodeBadRequest, "Amount has more than 2 decimal places", reasonAmountScaleExceeded, nil)
	// ErrAmountOutOfRange is the error when the amount exceeds DECIMAL(18,2) range
	ErrAmountOutOfRange = pkgErr.NewError(pkgErr.CodeBadRequest, "Amount exceeds DECIMAL(18,2) range", reasonAmountOutOfRange, nil)
	// ErrNewCurrencyType is the error when the currency type is invalid
	ErrNewCurrencyType = pkgErr.NewError(pkgErr.CodeBadRequest, "Invalid currency type", reasonNewCurrencyType, nil)
	// ErrNewIdempotencyKey is the error when the idempotency key is invalid
	ErrNewIdempotencyKey = pkgErr.NewError(pkgErr.CodeBadRequest, "Invalid idempotency key", reasonNewIdempotencyKey, nil)
	// ErrUserNotFound is the error when user is not found
	ErrUserNotFound = pkgErr.NewError(pkgErr.CodeNotFound, "User not found", reasonUserNotFound, nil)
	// ErrWithdrawalNotFound is the error when withdrawal is not found
	ErrWithdrawalNotFound = pkgErr.NewError(pkgErr.CodeNotFound, "Withdrawal not found", reasonWithdrawalNotFound, nil)
	// ErrInsufficientBalance is the error when user balance is insufficient
	ErrInsufficientBalance = pkgErr.NewError(pkgErr.CodeConflict, "Insufficient balance", reasonInsufficientBalance, nil)
	// ErrIdempotencyPayloadMismatch is the error when idempotency key is reused with different payload
	ErrIdempotencyPayloadMismatch = pkgErr.NewError(pkgErr.CodeUnprocessable, "Idempotency key was used with different payload", reasonIdempotencyPayloadMismatch, nil)
	// ErrInvalidRequestPayload is the error for invalid request payload
	ErrInvalidRequestPayload = pkgErr.NewError(pkgErr.CodeBadRequest, "Invalid request payload", reasonInvalidRequestPayload, nil)
)
