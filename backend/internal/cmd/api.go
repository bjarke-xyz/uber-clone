package cmd

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/bjarke-xyz/uber-clone-backend/internal/cfg"
	"github.com/bjarke-xyz/uber-clone-backend/internal/cmdutil"
	"github.com/bjarke-xyz/uber-clone-backend/internal/infra/http"
	"github.com/bjarke-xyz/uber-clone-backend/internal/infra/logging"
	"github.com/bjarke-xyz/uber-clone-backend/internal/infra/pubsub"
	"github.com/bjarke-xyz/uber-clone-backend/internal/service"
	"github.com/joho/godotenv"
)

func APICmd(ctx context.Context) error {
	godotenv.Load()
	cfg := cfg.NewConfig()
	port := 7000
	if cfg.Port != "" {
		port, _ = strconv.Atoi(cfg.Port)
	}
	logger := logging.NewLogger("api", cfg.Env)

	db, err := cmdutil.NewDatabasePool(ctx, cfg.DatabaseConnectionPoolUrl, 16)
	if err != nil {
		return err
	}
	cmdutil.MigrateDb(cfg.DatabaseConnectionPoolUrl)
	defer db.Close()

	osrClient := service.NewOpenRouteServiceClient(cfg.OSRApiKey)

	ps := pubsub.NewInMemoryPubsub()

	api := http.NewAPI(ctx, logger, cfg, db, osrClient, ps)
	srv := api.Server(port)

	go http.ServeMetrics(":9091")
	go api.PubsubSubscribe(ctx)
	go func() {
		_ = srv.ListenAndServe()
	}()
	logger.Info("started api", slog.Int("port", port))
	<-ctx.Done()
	_ = srv.Shutdown(ctx)
	return nil
}
