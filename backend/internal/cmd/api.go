package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	firebase "firebase.google.com/go/v4"
	apiPkg "github.com/bjarke-xyz/uber-clone-backend/internal/api"
	"github.com/bjarke-xyz/uber-clone-backend/internal/cmdutil"
	"github.com/bjarke-xyz/uber-clone-backend/internal/service"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
	"google.golang.org/api/option"
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

	credentialsJson := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_CONTENT")
	credentialsJsonBytes := []byte(credentialsJson)
	opt := option.WithCredentialsJSON(credentialsJsonBytes)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %w", err)
	}

	osrApiKey := os.Getenv("OSR_API_KEY")
	if osrApiKey == "" {
		return fmt.Errorf("OSR_API_KEY environment variable was empty")
	}
	osrClient := service.NewOpenRouteServiceClient(osrApiKey)

	api := apiPkg.NewAPI(ctx, logger, db, app, osrClient)
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
