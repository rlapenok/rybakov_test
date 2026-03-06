package app

import (
	"context"

	"github.com/rlapenok/rybakov_test/internal/config"
	"github.com/rlapenok/rybakov_test/internal/infra/repository"
	"github.com/rlapenok/rybakov_test/internal/infra/transport/handlers"
	"github.com/rlapenok/rybakov_test/internal/infra/transport/server"
	"github.com/rlapenok/rybakov_test/internal/usecase"
	"github.com/rlapenok/rybakov_test/pkg/db/pg"
	"github.com/rlapenok/rybakov_test/pkg/logger"
)

func Run(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	log := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Format)

	pool, err := pg.NewPool(ctx, cfg.Database.ToPgPoolConfig())
	if err != nil {
		return err
	}
	defer pool.Close()

	if err = pool.Migrate(ctx, cfg.Database.MigrationPath); err != nil {
		return err
	}

	withdrawalRepo := repository.NewWithdrawalRepository(pool)
	withdrawalUseCase := usecase.NewWithdrawalUseCase(withdrawalRepo)
	withdrawalHandler := handlers.NewWithdrawalHandler(withdrawalUseCase)

	httpServer := server.NewServer(cfg.Server, log, cfg.Auth.BearerToken, withdrawalHandler)

	return httpServer.Start()
}
