// Package rest provides a rest generic things to use
package rest

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/teathedev/fullstack-boilerplate/backend/pkg/env"
)

var (
	Router  *chi.Mux
	API     huma.API
	Version = "1.0.0"
)

func init() {
	appName := env.GetString("APP_NAME", "TEARest")

	Router = chi.NewMux()
	API = humachi.New(
		Router,
		huma.DefaultConfig(appName, Version),
	)
}
