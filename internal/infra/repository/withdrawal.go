package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rlapenok/rybakov_test/internal/domain"
	"github.com/rlapenok/rybakov_test/pkg/db/pg"
	"github.com/shopspring/decimal"
)

// WithdrawalRepository is a repository for withdrawals
type WithdrawalRepository struct {
	pool *pg.Pool
}

// NewWithdrawalRepository creates a new withdrawal repository
func NewWithdrawalRepository(pool *pg.Pool) *WithdrawalRepository {
	return &WithdrawalRepository{pool: pool}
}

// CreateWithdrawal creates a new withdrawal inside a single transaction:
//  1. SELECT FOR UPDATE — locks the user's balance row to prevent concurrent double-spend.
//  2. INSERT ... ON CONFLICT (user_id, idempotency_key) DO NOTHING — idempotency gate.
//     - 0 rows returned → key already exists → compare payload hash:
//     - same hash  → return existing ID (idempotent success, no balance change)
//     - diff hash  → 422 ErrIdempotencyPayloadMismatch
//  3. UPDATE users SET balance = balance - $amount WHERE balance >= $amount
//     - 0 rows → 409 ErrInsufficientBalance
//  4. COMMIT — return new WithdrawalID
func (r *WithdrawalRepository) CreateWithdrawal(
	ctx context.Context,
	withdrawal domain.Withdrawal,
	payloadHash string,
) (domain.WithdrawalID, error) {
	// 1. Start transaction
	tx, err := r.pool.Pgx().BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.WithdrawalID{}, err
	}

	// Rollback if error; noop after commit
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// 2. Lock the user's balance row (SELECT FOR UPDATE)
	if err = r.lockUserBalance(ctx, tx, withdrawal.UserID()); err != nil {
		return domain.WithdrawalID{}, err
	}

	// 3. Try to insert withdrawal (idempotency via ON CONFLICT DO NOTHING)
	inserted, err := r.insertWithdrawal(ctx, tx, withdrawal, payloadHash)
	if err != nil {
		return domain.WithdrawalID{}, err
	}

	// 4. ON CONFLICT fired — key already exists
	if !inserted {
		existing, storedHash, getErr := r.getByUserAndIdempotencyKey(ctx, tx, withdrawal.UserID(), withdrawal.IdempotencyKey())
		if getErr != nil {
			err = getErr
			return domain.WithdrawalID{}, err
		}

		if storedHash != payloadHash {
			err = domain.ErrIdempotencyPayloadMismatch
			return domain.WithdrawalID{}, err
		}

		// Same payload — commit and return existing ID (no balance change)
		if err = tx.Commit(ctx); err != nil {
			return domain.WithdrawalID{}, err
		}
		return existing.ID(), nil
	}

	// 5. New withdrawal inserted — debit balance
	if err = r.debitBalance(ctx, tx, withdrawal.UserID(), withdrawal.Amount()); err != nil {
		return domain.WithdrawalID{}, err
	}

	// 6. Commit
	if err = tx.Commit(ctx); err != nil {
		return domain.WithdrawalID{}, err
	}

	return withdrawal.ID(), nil
}

// GetWithdrawalByID returns withdrawal by ID.
func (r *WithdrawalRepository) GetWithdrawalByID(
	ctx context.Context,
	id domain.WithdrawalID,
) (domain.Withdrawal, error) {
	const query = `
		SELECT id, user_id, amount, destination, idempotency_key, status
		FROM withdrawals
		WHERE id = $1
	`

	w, _, err := r.scan(ctx, r.pool.Pgx(), query, false, id.Value())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Withdrawal{}, domain.ErrWithdrawalNotFound
		}
		return domain.Withdrawal{}, err
	}

	return w, nil
}

// ─── private helpers ─────────────────────────────────────────────────────────

// lockUserBalance locks the user balance row for the duration of the transaction.
func (r *WithdrawalRepository) lockUserBalance(
	ctx context.Context,
	tx pgx.Tx,
	userID domain.UserID,
) error {
	const query = `SELECT balance FROM users WHERE id = $1 FOR UPDATE`

	var balance decimal.Decimal
	err := tx.QueryRow(ctx, query, userID.Value()).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}

	return nil
}

// insertWithdrawal tries to insert a withdrawal row.
// Returns (true, nil) if inserted, (false, nil) if idempotency key already exists.
func (r *WithdrawalRepository) insertWithdrawal(
	ctx context.Context,
	tx pgx.Tx,
	withdrawal domain.Withdrawal,
	payloadHash string,
) (bool, error) {
	const query = `
		INSERT INTO withdrawals (id, user_id, amount, destination, idempotency_key, payload_hash, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, idempotency_key) DO NOTHING
		RETURNING id
	`

	var insertedID uuid.UUID
	err := tx.QueryRow(
		ctx, query,
		withdrawal.ID().Value(),
		withdrawal.UserID().Value(),
		withdrawal.Amount().Value(),
		withdrawal.Destination().Value(),
		withdrawal.IdempotencyKey().Value(),
		payloadHash,
		withdrawal.Status().Value(),
	).Scan(&insertedID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// ON CONFLICT DO NOTHING — no row returned
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// debitBalance subtracts amount from user balance atomically.
// Returns ErrInsufficientBalance if balance < amount.
func (r *WithdrawalRepository) debitBalance(
	ctx context.Context,
	tx pgx.Tx,
	userID domain.UserID,
	amount domain.Amount,
) error {
	const query = `
		UPDATE users
		SET balance = balance - $1
		WHERE id = $2 AND balance >= $1
	`

	tag, err := tx.Exec(ctx, query, amount.Value(), userID.Value())
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrInsufficientBalance
	}

	return nil
}

// getByUserAndIdempotencyKey fetches an existing withdrawal and its payload hash.
func (r *WithdrawalRepository) getByUserAndIdempotencyKey(
	ctx context.Context,
	tx pgx.Tx,
	userID domain.UserID,
	idempotencyKey domain.IdempotencyKey,
) (domain.Withdrawal, string, error) {
	const query = `
		SELECT id, user_id, amount, destination, idempotency_key, status, payload_hash
		FROM withdrawals
		WHERE user_id = $1 AND idempotency_key = $2
	`

	w, hash, err := r.scan(ctx, tx, query, true, userID.Value(), idempotencyKey.Value())
	if err != nil {
		return domain.Withdrawal{}, "", err
	}

	return w, hash, nil
}

// ─── scan helpers ─────────────────────────────────────────────────────────────

type queryRower interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// scan reads one withdrawal row. When withHash is true, expects payload_hash as the last column.
func (r *WithdrawalRepository) scan(
	ctx context.Context,
	q queryRower,
	query string,
	withHash bool,
	args ...any,
) (domain.Withdrawal, string, error) {
	var (
		id             uuid.UUID
		userID         uuid.UUID
		amount         decimal.Decimal
		destination    uuid.UUID
		idempotencyKey string
		status         string
		payloadHash    string
	)

	row := q.QueryRow(ctx, query, args...)

	var err error
	if withHash {
		err = row.Scan(&id, &userID, &amount, &destination, &idempotencyKey, &status, &payloadHash)
	} else {
		err = row.Scan(&id, &userID, &amount, &destination, &idempotencyKey, &status)
	}

	if err != nil {
		return domain.Withdrawal{}, "", err
	}

	domainAmount, err := domain.NewAmount(amount.StringFixed(2))
	if err != nil {
		return domain.Withdrawal{}, "", err
	}

	domainKey, err := domain.NewIdempotencyKey(idempotencyKey)
	if err != nil {
		return domain.Withdrawal{}, "", err
	}

	w := domain.RehydrateWithdrawal(
		domain.RehydrateWithdrawalID(id),
		domain.RehydrateUserID(userID),
		domainAmount,
		domain.RehydrateUserID(destination),
		domainKey,
		domain.WithdrawalStatus(status),
	)

	return w, payloadHash, nil
}
