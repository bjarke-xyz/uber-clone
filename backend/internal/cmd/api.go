package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	apiPkg "github.com/bjarke-xyz/uber-clone-backend/internal/api"
	"github.com/bjarke-xyz/uber-clone-backend/internal/cmdutil"
	"github.com/bjarke-xyz/uber-clone-backend/internal/service"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func APICmd(ctx context.Context) error {
	godotenv.Load()
	port := 7000
	port = 7000
	if os.Getenv("PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("PORT"))
	}
	logger := cmdutil.NewLogger("api")

	db, err := cmdutil.NewDatabasePool(ctx, 16)
	if err != nil {
		return err
	}
	defer db.Close()

	osrApiKey := os.Getenv("OSR_API_KEY")
	if osrApiKey == "" {
		return fmt.Errorf("OSR_API_KEY environment variable was empty")
	}
	osrClient := service.NewOpenRouteServiceClient(osrApiKey)

	api := apiPkg.NewAPI(ctx, logger, db, osrClient)
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
