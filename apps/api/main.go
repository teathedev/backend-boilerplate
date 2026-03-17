package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	_ "github.com/teathedev/backend-boilerplate/apps/api/controller"
	"github.com/teathedev/backend-boilerplate/apps/api/rest"
	"github.com/teathedev/backend-boilerplate/internal/actions"
	"github.com/teathedev/backend-boilerplate/internal/db"
	_ "github.com/teathedev/backend-boilerplate/internal/validation"
	"github.com/teathedev/pkg/env"
	"github.com/teathedev/pkg/logger"
)

func main() {
	log := logger.New("System")
	ctx := context.Background()
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := db.Client.ExecContext(pingCtx, "SELECT 1;"); err != nil {
		log.Fatal(
			"Failed to ping database!",
			logger.LogParams{
				"Error": err,
			},
		)
		os.Exit(1)
	}

	if err := actions.RefreshJWKs(ctx); err != nil {
		log.Fatal(
			"Failed to refresh jwks!",
			logger.LogParams{
				"Error": err,
			},
		)
		os.Exit(1)
	}

	port := env.GetNumber("PORT", 8888)
	addr := fmt.Sprintf(":%d", port)

	log.Info(
		"API documentation is available.",
		logger.LogParams{
			"URL":  fmt.Sprintf("http://localhost%s/docs", addr),
			"JSON": fmt.Sprintf("http://localhost%s/openapi.json", addr),
			"YAML": fmt.Sprintf("http://localhost%s/openapi.yaml", addr),
		},
	)
	if err := http.ListenAndServe(addr, rest.Router); err != nil {
		log.Fatal(
			"Failed to listen and serve!",
			logger.LogParams{
				"Error": err,
			},
		)
	}
}
