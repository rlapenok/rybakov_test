package domain

import (
	"strings"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CurrencyType string

const (
	CurrencyTypeUSDT CurrencyType = "USDT"
)

// UserID is a user UserID
type UserID uuid.UUID

// FromUUID creates a new ID from a uuid.UUID
func NewUserIDFromUUID(u uuid.UUID) UserID {
	return UserID(u)
}

// String returns the string representation of the ID
func (id UserID) String() string {
	return uuid.UUID(id).String()
}

// Value returns the value of the ID
func (id UserID) Value() uuid.UUID {
	return uuid.UUID(id)
}

// RehydrateUserID rehydrates a user ID from a uuid.UUID
func RehydrateUserID(rawUserID uuid.UUID) UserID {
	return UserID(rawUserID)
}

// Amount is a type that represents the amount of money
type Amount decimal.Decimal

var maxAmountDecimal182 = decimal.RequireFromString("9999999999999999.99")

// NewAmount creates a new amount from a string
func NewAmount(amount string) (Amount, error) {
	amount = strings.TrimSpace(amount)

	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return Amount(amountDecimal), ErrNewAmount.AddMeta("original_error", err.Error())
	}

	if amountDecimal.LessThanOrEqual(decimal.Zero) {
		return Amount(amountDecimal), ErrAmountNotGreaterThanZero
	}

	roundedAmount := amountDecimal.Round(2)
	if !amountDecimal.Equal(roundedAmount) {
		return Amount(amountDecimal), ErrAmountScaleExceeded
	}

	if roundedAmount.GreaterThan(maxAmountDecimal182) {
		return Amount(amountDecimal), ErrAmountOutOfRange
	}

	return Amount(roundedAmount), nil
}

// String returns the string representation of the amount
func (a Amount) String() string {
	return decimal.Decimal(a).StringFixed(2)
}

// Value returns the value of the amount
func (a Amount) Value() decimal.Decimal {
	return decimal.Decimal(a)
}

// RehydrateAmount rehydrates an amount from a string
func RehydrateAmount(rawAmount string) Amount {
	return Amount(decimal.RequireFromString(rawAmount))
}

// Currency is a type that represents the currency type
type Currency CurrencyType

// NewCurrency creates a new currency type from a string
func NewCurrency(rawCurrency string) (Currency, error) {
	currencyType := strings.TrimSpace(rawCurrency)

	if currencyType != string(CurrencyTypeUSDT) {
		return Currency(currencyType), ErrNewCurrencyType
	}

	return Currency(currencyType), nil
}

// String returns the string representation of the currency type
func (c Currency) String() string {
	return string(c)
}

// Value returns the value of the currency type
func (c Currency) Value() CurrencyType {
	return CurrencyType(c)
}

// IdempotencyKey is a type that represents the idempotency key
type IdempotencyKey string

// NewIdempotencyKey creates a new idempotency key from a string
func NewIdempotencyKey(idempotencyKey string) (IdempotencyKey, error) {
	idempotencyKey = strings.TrimSpace(idempotencyKey)

	if idempotencyKey == "" {
		return IdempotencyKey(idempotencyKey), ErrNewIdempotencyKey
	}
	return IdempotencyKey(idempotencyKey), nil
}

// Value returns the value of the idempotency key
func (i IdempotencyKey) Value() string {
	return string(i)
}

// RehydrateIdempotencyKey rehydrates an idempotency key from a string
func RehydrateIdempotencyKey(rawIdempotencyKey string) IdempotencyKey {
	return IdempotencyKey(rawIdempotencyKey)
}

// WithdrawalStatus is a type that represents the withdrawal status
type WithdrawalStatus string

const (
	WithdrawalStatusPending WithdrawalStatus = "pending"
)

func (w WithdrawalStatus) Value() string {
	return string(w)
}

// WithdrawalID is a type that represents the withdrawal ID
type WithdrawalID uuid.UUID

// NewWithdrawalIDFromUUID creates a new withdrawal ID from a uuid.UUID
func NewWithdrawalIDFromUUID(u uuid.UUID) WithdrawalID {
	return WithdrawalID(u)
}

// NewWithdrawalID creates a new withdrawal ID from a uuid.UUID
func NewWithdrawalID() WithdrawalID {
	return WithdrawalID(uuid.New())
}

// Value returns the value of the withdrawal ID
func (id WithdrawalID) Value() uuid.UUID {
	return uuid.UUID(id)
}

// RehydrateWithdrawalID rehydrates a withdrawal ID from a uuid.UUID
func RehydrateWithdrawalID(rawID uuid.UUID) WithdrawalID {
	return WithdrawalID(rawID)
}
