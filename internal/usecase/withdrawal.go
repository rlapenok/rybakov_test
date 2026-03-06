package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
	"github.com/rlapenok/rybakov_test/internal/domain"
)

// WithdrawalUseCase is a use case for withdrawals
type WithdrawalUseCase struct {
	repo domain.WithdrawalRepo
}

// NewWithdrawalUseCase creates a new withdrawal use case
func NewWithdrawalUseCase(repo domain.WithdrawalRepo) *WithdrawalUseCase {
	return &WithdrawalUseCase{repo: repo}
}

// CreateWithdrawalInput is the input for the CreateWithdrawal use case
type CreateWithdrawalInput struct {
	UserID         uuid.UUID `json:"user_id"`
	Amount         string    `json:"amount"`
	Currency       string    `json:"currency"`
	Destination    uuid.UUID `json:"destination"`
	IdempotencyKey string    `json:"idempotency_key"`
}

// CreateWithdrawalOutput is the output for the CreateWithdrawal use case
type CreateWithdrawalOutput struct {
	ID string `json:"id"`
}

// GetWithdrawalOutput is the output for the GetWithdrawalByID use case
type GetWithdrawalOutput struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	Destination    string `json:"destination"`
	IdempotencyKey string `json:"idempotency_key"`
	Status         string `json:"status"`
}

// CreateWithdrawal creates a new withdrawal and returns its ID.
func (u *WithdrawalUseCase) CreateWithdrawal(
	ctx context.Context,
	input CreateWithdrawalInput,
) (CreateWithdrawalOutput, error) {
	//1. Validate currency
	_, err := domain.NewCurrency(input.Currency)
	if err != nil {
		return CreateWithdrawalOutput{}, err
	}

	//2. Build and validate withdrawal domain object
	withdrawal, err := domain.NewWithdrawal(
		input.UserID,
		input.Amount,
		input.Destination,
		input.IdempotencyKey,
	)
	if err != nil {
		return CreateWithdrawalOutput{}, err
	}

	//3. Build payload hash
	payloadHash := buildPayloadHash(
		withdrawal.UserID().String(),
		withdrawal.Amount().String(),
		withdrawal.Destination().String(),
	)

	//4. Persist
	id, err := u.repo.CreateWithdrawal(ctx, withdrawal, payloadHash)
	if err != nil {
		return CreateWithdrawalOutput{}, err
	}

	return CreateWithdrawalOutput{ID: id.Value().String()}, nil
}

// GetWithdrawalByID returns a withdrawal by its ID
func (u *WithdrawalUseCase) GetWithdrawalByID(
	ctx context.Context,
	id uuid.UUID,
) (GetWithdrawalOutput, error) {
	withdrawalID := domain.NewWithdrawalIDFromUUID(id)

	w, err := u.repo.GetWithdrawalByID(ctx, withdrawalID)
	if err != nil {
		return GetWithdrawalOutput{}, err
	}

	return GetWithdrawalOutput{
		ID:             w.ID().Value().String(),
		UserID:         w.UserID().String(),
		Amount:         w.Amount().String(),
		Currency:       string(domain.CurrencyTypeUSDT),
		Destination:    w.Destination().String(),
		IdempotencyKey: w.IdempotencyKey().Value(),
		Status:         string(w.Status()),
	}, nil
}

func buildPayloadHash(values ...string) string {
	joined := strings.Join(values, "|")
	hash := sha256.Sum256([]byte(joined))
	return hex.EncodeToString(hash[:])
}
