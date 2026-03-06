package domain

import "context"

type WithdrawalRepo interface {
	CreateWithdrawal(ctx context.Context, withdrawal Withdrawal, payloadHash string) (WithdrawalID, error)
	GetWithdrawalByID(ctx context.Context, id WithdrawalID) (Withdrawal, error)
}
