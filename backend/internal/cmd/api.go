package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/bjarke-xyz/uber-clone-backend/internal/auth"
	"github.com/bjarke-xyz/uber-clone-backend/internal/cfg"
	"github.com/bjarke-xyz/uber-clone-backend/internal/cmdutil"
	"github.com/bjarke-xyz/uber-clone-backend/internal/infra/http"
	"github.com/bjarke-xyz/uber-clone-backend/internal/infra/pubsub"
	"github.com/bjarke-xyz/uber-clone-backend/internal/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(cfg.AuthGrpcUrl, opts...)
	if err != nil {
		return fmt.Errorf("failed to dial auth grpc: %w", err)
	}
	defer conn.Close()
	authClient := auth.NewAuthClient(conn)

	ps := pubsub.NewInMemoryPubsub()

	api := http.NewAPI(ctx, logger, cfg, authClient, db, osrClient, ps)
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
