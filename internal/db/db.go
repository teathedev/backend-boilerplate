// Package db contains database initialization and connection
package db

import (
	"fmt"
	"log"
	"os"

	"github.com/teathedev/backend-boilerplate/internal/ent"
	"github.com/teathedev/pkg/env"
)

var Client *ent.Client

func init() {
	dbURL := env.GetString("POSTGRES_DB_URL", "")
	if dbURL == "" {
		fmt.Println("Postgres Database URL is missing!")
		os.Exit(1)
		return
	}

	var err error
	Client, err = ent.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
}
