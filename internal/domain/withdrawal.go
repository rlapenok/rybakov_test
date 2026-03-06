package domain

import (
	"github.com/google/uuid"
)

type Withdrawal struct {
	id             WithdrawalID
	from           UserID
	amount         Amount
	to             UserID
	idempotencyKey IdempotencyKey
	status         WithdrawalStatus
}

func NewWithdrawal(
	rawUserID uuid.UUID,
	rawAmount string,
	rawDestination uuid.UUID,
	rawIdempotencyKey string,
) (Withdrawal, error) {
	userID := NewUserIDFromUUID(rawUserID)

	amount, err := NewAmount(rawAmount)
	if err != nil {
		return Withdrawal{}, err
	}

	destination := NewUserIDFromUUID(rawDestination)

	idempotencyKey, err := NewIdempotencyKey(rawIdempotencyKey)
	if err != nil {
		return Withdrawal{}, err
	}

	return Withdrawal{
		id:             NewWithdrawalID(),
		from:           userID,
		amount:         amount,
		to:             destination,
		idempotencyKey: idempotencyKey,
		status:         WithdrawalStatusPending,
	}, nil
}

// RehydrateWithdrawal rehydrates a Withdrawal from already-typed domain values (used by repository).
func RehydrateWithdrawal(
	id WithdrawalID,
	userID UserID,
	amount Amount,
	destination UserID,
	idempotencyKey IdempotencyKey,
	status WithdrawalStatus,
) Withdrawal {
	return Withdrawal{
		id:             id,
		from:           userID,
		amount:         amount,
		to:             destination,
		idempotencyKey: idempotencyKey,
		status:         status,
	}
}

// ID returns the ID of the withdrawal
func (w Withdrawal) ID() WithdrawalID {
	return w.id
}

// UserID returns the user ID.
func (w Withdrawal) UserID() UserID {
	return w.from
}

// Amount returns the amount of the withdrawal
func (w Withdrawal) Amount() Amount {
	return w.amount
}

// Destination returns withdrawal destination.
func (w Withdrawal) Destination() UserID {
	return w.to
}

// IdempotencyKey returns the idempotency key of the withdrawal
func (w Withdrawal) IdempotencyKey() IdempotencyKey {
	return w.idempotencyKey
}

// Status returns the status of the withdrawal
func (w Withdrawal) Status() WithdrawalStatus {
	return w.status
}
