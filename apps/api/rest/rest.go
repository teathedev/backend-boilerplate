// Package rest provides the API application's REST/OpenAPI setup: router,
// multiple servers (IAM, Todo, X), and Swagger UI. Generic REST utilities
// live in pkg/rest.
package rest

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/teathedev/pkg/env"
	"github.com/teathedev/pkg/errors"
)

var (
	Router      *chi.Mux
	Version     = "1.0.0"
	APIInstance huma.API
)

var (
	IAMGroup       *huma.Group
	ProtectedGroup *huma.Group
)

// API returns the default API (IAM) for backward compatibility.
func API() huma.API {
	return APIInstance
}

// BadInputSchema is the 400 (validation) error schema: message + error array. Use only for 400.
var BadInputSchema *huma.Schema

// ErrorSchema is the generic error schema for 401/403/404/500: message only.
var ErrorSchema *huma.Schema

// ErrorResponses returns a Responses map: 400 uses bad-input schema (message + error array),
// 401/403/404/500 use simple message-only schema.
func ErrorResponses() map[string]*huma.Response {
	badContent := map[string]*huma.MediaType{"application/json": {Schema: BadInputSchema}}
	simpleContent := map[string]*huma.MediaType{"application/json": {Schema: ErrorSchema}}
	return map[string]*huma.Response{
		"400": {Description: "Bad request / validation failed", Content: badContent},
		"401": {Description: "Unauthorized", Content: simpleContent},
		"403": {Description: "Forbidden", Content: simpleContent},
		"404": {Description: "Not found", Content: simpleContent},
		"500": {Description: "Internal server error", Content: simpleContent},
	}
}

func init() {
	SetupErrorFactory()
	huma.DefaultArrayNullable = false

	Router = chi.NewMux()
	appName := env.GetString("APP_NAME", "TEARest")

	// Single API with groups for IAM and Protected. Huma validation is disabled for
	// auth payloads; we only use it for mapping & docs there.
	cfg := huma.DefaultConfig(appName, Version)
	cfg.DocsPath = "/docs"
	cfg.DocsRenderer = huma.DocsRendererScalar
	APIInstance = humachi.New(Router, cfg)

	IAMGroup = huma.NewGroup(APIInstance, "/auth")
	IAMGroup.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = append(op.Tags, "IAM")
	})

	ProtectedGroup = huma.NewGroup(APIInstance, "/protected")
	ProtectedGroup.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = append(op.Tags, "Protected")
	})

	BadInputSchema = huma.SchemaFromType(APIInstance.OpenAPI().Components.Schemas, reflect.TypeOf(errors.APIErrorResponse{}))
	ErrorSchema = huma.SchemaFromType(APIInstance.OpenAPI().Components.Schemas, reflect.TypeOf(errors.APIErrorSimple{}))
}
