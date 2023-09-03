package cmd

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	apiPkg "github.com/bjarke-xyz/uber-clone-backend/internal/api"
	"github.com/bjarke-xyz/uber-clone-backend/internal/auth"
	"github.com/bjarke-xyz/uber-clone-backend/internal/cfg"
	"github.com/bjarke-xyz/uber-clone-backend/internal/cmdutil"
	"github.com/bjarke-xyz/uber-clone-backend/internal/service"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func APICmd(ctx context.Context) error {
	godotenv.Load()
	cfg := cfg.NewConfig()
	port := 7000
	if cfg.Port != "" {
		port, _ = strconv.Atoi(cfg.Port)
	}
	logger := cmdutil.NewLogger("api", cfg.Env)

	db, err := cmdutil.NewDatabasePool(ctx, cfg.DatabaseConnectionPoolUrl, 16)
	if err != nil {
		return err
	}
	cmdutil.MigrateDb(cfg.DatabaseConnectionPoolUrl)
	defer db.Close()

	osrClient := service.NewOpenRouteServiceClient(cfg.OSRApiKey)

	auth.WarmCache(ctx)

	api := apiPkg.NewAPI(ctx, logger, cfg, db, osrClient)
	srv := api.Server(port)

	// metrics server
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9091", mux)
	}()

	go func() {
		_ = srv.ListenAndServe()
	}()
	logger.Info("started api", slog.Int("port", port))
	<-ctx.Done()
	_ = srv.Shutdown(ctx)
	return nil
}
