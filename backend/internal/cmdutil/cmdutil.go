package cmdutil

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/bjarke-xyz/uber-clone-backend/internal/infra/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewLogger(service string, env string) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	child := logger.With(slog.Group("service_info", slog.String("env", env), slog.String("service", service)))
	return child
}

func NewDatabasePool(ctx context.Context, databaseUrl string, maxConns int) (*pgxpool.Pool, error) {
	if maxConns == 0 {
		maxConns = 1
	}

	queryChar := "?"
	if strings.Contains(databaseUrl, "?") {
		queryChar = "&"
	}
	url := fmt.Sprintf(
		"%s%vpool_max_conns=%d&pool_min_conns=%d",
		databaseUrl,
		queryChar,
		maxConns,
		2,
	)
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	// Setting the build statement cache to nil helps this work with pgbouncer
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Second
	return pgxpool.NewWithConfig(ctx, config)
}

func MigrateDb(databaseUrl string) error {
	err := postgres.Migrate("up", databaseUrl)
	if err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}
	return nil
}
