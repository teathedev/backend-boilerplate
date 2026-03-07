//go:build ignore

// Programmatic migration generation for ent (versioned migrations with Atlas).
// Run from repo root: go run -mod=mod ./internal/ent/migrate/main.go <migration_name>
// Requires a running Postgres (e.g. task docker-up). Uses POSTGRES_DB_URL or builds from POSTGRES_* env vars.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"

	"github.com/teathedev/backend-boilerplate/internal/ent/migrate"
)

const migrationsDir = "internal/ent/migrate/migrations"

func main() {
	ctx := context.Background()

	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		log.Fatalf("creating migration directory: %v", err)
	}

	dir, err := atlas.NewLocalDir(migrationsDir)
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}

	devURL := getDevURL()
	opts := []schema.MigrateOption{
		schema.WithDir(dir),
		schema.WithMigrationMode(schema.ModeReplay),
		schema.WithDialect(dialect.Postgres),
		schema.WithFormatter(atlas.DefaultFormatter),
	}

	name := "changes"
	if len(os.Args) >= 2 {
		name = os.Args[1]
	}

	if err := migrate.NamedDiff(ctx, devURL, name, opts...); err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}

func getDevURL() string {
	if u := os.Getenv("POSTGRES_DB_URL"); u != "" {
		if !strings.Contains(u, "sslmode=") {
			if strings.Contains(u, "?") {
				u += "&sslmode=disable"
			} else {
				u += "?sslmode=disable"
			}
		}
		return u
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "postgres"
	}
	pass := os.Getenv("POSTGRES_PASSWORD")
	if pass == "" {
		pass = "postgres"
	}
	db := os.Getenv("POSTGRES_DB")
	if db == "" {
		db = "tearest"
	}
	return fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", user, pass, db)
}
