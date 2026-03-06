package repository_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rlapenok/rybakov_test/internal/domain"
	"github.com/rlapenok/rybakov_test/internal/infra/repository"
	"github.com/rlapenok/rybakov_test/pkg/db/pg"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// destID — seed user pre-inserted by migration, used as withdrawal destination.
var destID = uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")

var (
	testRepo *repository.WithdrawalRepository
	testPool *pg.Pool
)

// TestMain spins up a PostgreSQL container, runs migrations once,
// and tears everything down after all tests finish.
func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx, "postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("start postgres container: %v", err)
	}
	defer container.Terminate(ctx) //nolint:errcheck

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("get container host: %v", err)
	}
	natPort, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("get container port: %v", err)
	}
	portNum, _ := strconv.Atoi(natPort.Port())

	testPool, err = pg.NewPool(ctx, &pg.PgPoolConfig{
		PgConnectionConfig: pg.PgConnectionConfig{
			Host:     host,
			Port:     uint16(portNum),
			User:     "test",
			Password: "test",
			Database: "testdb",
			SSLMode:  "disable",
		},
		MinConns: 1,
		MaxConns: 10,
	})
	if err != nil {
		log.Fatalf("create pg pool: %v", err)
	}
	defer testPool.Close()

	if err = testPool.Migrate(ctx, "../../../migrations"); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	testRepo = repository.NewWithdrawalRepository(testPool)

	os.Exit(m.Run())
}

// ─── helpers ──────────────────────────────────────────────────────────────────

// insertTestUser inserts a fresh user with the given balance and registers
// cleanup to remove that user (and its withdrawals) after the test.
func insertTestUser(t *testing.T, balance string) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	id := uuid.New()

	_, err := testPool.Pgx().Exec(ctx,
		`INSERT INTO users (id, balance) VALUES ($1, $2)`, id, balance,
	)
	if err != nil {
		t.Fatalf("insertTestUser: %v", err)
	}

	t.Cleanup(func() {
		testPool.Pgx().Exec(ctx, `DELETE FROM withdrawals WHERE user_id = $1`, id) //nolint:errcheck
		testPool.Pgx().Exec(ctx, `DELETE FROM users WHERE id = $1`, id)            //nolint:errcheck
	})
	return id
}

// makeWithdrawal builds a domain.Withdrawal and computes its payload hash,
// mirroring the logic inside WithdrawalUseCase.
func makeWithdrawal(t *testing.T, userID uuid.UUID, amount, key string) (domain.Withdrawal, string) {
	t.Helper()
	w, err := domain.NewWithdrawal(userID, amount, destID, key)
	if err != nil {
		t.Fatalf("domain.NewWithdrawal: %v", err)
	}
	joined := strings.Join([]string{w.UserID().String(), w.Amount().String(), w.Destination().String()}, "|")
	sum := sha256.Sum256([]byte(joined))
	hash := hex.EncodeToString(sum[:])
	return w, hash
}

// ─── tests ────────────────────────────────────────────────────────────────────

// TestCreateWithdrawal_Success verifies that a valid withdrawal is persisted
// and a non-zero ID is returned.
func TestCreateWithdrawal_Success(t *testing.T) {
	ctx := context.Background()
	userID := insertTestUser(t, "100.00")

	w, hash := makeWithdrawal(t, userID, "50.00", "key-success")

	id, err := testRepo.CreateWithdrawal(ctx, w, hash)
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	if id.Value() == (uuid.UUID{}) {
		t.Fatal("expected non-zero withdrawal ID")
	}
}

// TestCreateWithdrawal_InsufficientBalance verifies that withdrawing more
// than the current balance returns ErrInsufficientBalance (HTTP 409).
func TestCreateWithdrawal_InsufficientBalance(t *testing.T) {
	ctx := context.Background()
	userID := insertTestUser(t, "10.00") // balance < requested amount

	w, hash := makeWithdrawal(t, userID, "50.00", "key-insuf")

	_, err := testRepo.CreateWithdrawal(ctx, w, hash)
	if !errors.Is(err, domain.ErrInsufficientBalance) {
		t.Fatalf("expected ErrInsufficientBalance, got: %v", err)
	}
}

// TestCreateWithdrawal_Idempotency covers two sub-cases:
//  1. Same key + same payload → returns the original ID without re-debiting.
//  2. Same key + different payload → ErrIdempotencyPayloadMismatch (HTTP 422).
func TestCreateWithdrawal_Idempotency(t *testing.T) {
	ctx := context.Background()
	userID := insertTestUser(t, "100.00")

	// ── sub-case 1: repeat with identical payload ──────────────────────────
	w1, hash1 := makeWithdrawal(t, userID, "30.00", "key-idem")

	id1, err := testRepo.CreateWithdrawal(ctx, w1, hash1)
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	// Build a second withdrawal with a new UUID but the same hash (same payload)
	w1dup, _ := makeWithdrawal(t, userID, "30.00", "key-idem")
	id1dup, err := testRepo.CreateWithdrawal(ctx, w1dup, hash1)
	if err != nil {
		t.Fatalf("idempotent repeat: %v", err)
	}
	if id1.Value() != id1dup.Value() {
		t.Fatalf("idempotent repeat must return same ID: got %s and %s",
			id1.Value(), id1dup.Value())
	}

	// ── sub-case 2: same key, different amount → different hash ────────────
	w2, hash2 := makeWithdrawal(t, userID, "99.00", "key-idem") // amount differs
	_, err = testRepo.CreateWithdrawal(ctx, w2, hash2)
	if !errors.Is(err, domain.ErrIdempotencyPayloadMismatch) {
		t.Fatalf("expected ErrIdempotencyPayloadMismatch, got: %v", err)
	}
}

// TestCreateWithdrawal_Concurrent sends two simultaneous withdrawal requests
// for the full balance (100.00) from a single account. Exactly one must
// succeed and the other must return ErrInsufficientBalance, proving that
// SELECT FOR UPDATE prevents double-spend.
func TestCreateWithdrawal_Concurrent(t *testing.T) {
	ctx := context.Background()
	userID := insertTestUser(t, "100.00")

	type result struct {
		id  domain.WithdrawalID
		err error
	}
	results := make([]result, 2)

	var wg sync.WaitGroup
	for i := range 2 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			w, hash := makeWithdrawal(t, userID, "100.00", fmt.Sprintf("key-concurrent-%d", i))
			id, err := testRepo.CreateWithdrawal(ctx, w, hash)
			results[i] = result{id, err}
		}(i)
	}
	wg.Wait()

	successes, insufficientCount := 0, 0
	for _, r := range results {
		switch {
		case r.err == nil:
			successes++
		case errors.Is(r.err, domain.ErrInsufficientBalance):
			insufficientCount++
		default:
			t.Errorf("unexpected error: %v", r.err)
		}
	}

	if successes != 1 {
		t.Errorf("expected exactly 1 success, got %d", successes)
	}
	if insufficientCount != 1 {
		t.Errorf("expected exactly 1 ErrInsufficientBalance, got %d", insufficientCount)
	}
}
