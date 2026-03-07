// Package rest provides the API application's REST/OpenAPI setup: router,
// multiple servers (IAM, Todo, X), and Swagger UI. Generic REST utilities
// live in pkg/rest.
package rest

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/teathedev/backend-boilerplate/pkg/env"
	"github.com/teathedev/backend-boilerplate/pkg/errors"
	pkgrest "github.com/teathedev/backend-boilerplate/pkg/rest"
)

var (
	Router  *chi.Mux
	Version = "1.0.0"
)

// Server holds one API surface: its own router, Huma API, and OpenAPI spec.
type Server struct {
	Name     string
	BasePath string
	Router   chi.Router
	API      huma.API
}

var (
	IAMServer  *Server
	TodoServer *Server
	XServer    *Server
)

// API returns the default API (IAM) for backward compatibility.
func API() huma.API {
	if IAMServer != nil {
		return IAMServer.API
	}
	return nil
}

// IAMBadInputSchema is the 400 (validation) error schema: message + error array. Use only for 400.
var IAMBadInputSchema *huma.Schema

// IAMErrorSchema is the generic error schema for 401/403/404/500: message only.
var IAMErrorSchema *huma.Schema

// ErrorResponses returns a Responses map: 400 uses bad-input schema (message + error array),
// 401/403/404/500 use simple message-only schema.
func ErrorResponses(badInputSchema, simpleSchema *huma.Schema) map[string]*huma.Response {
	badContent := map[string]*huma.MediaType{"application/json": {Schema: badInputSchema}}
	simpleContent := map[string]*huma.MediaType{"application/json": {Schema: simpleSchema}}
	return map[string]*huma.Response{
		"400": {Description: "Bad request / validation failed", Content: badContent},
		"401": {Description: "Unauthorized", Content: simpleContent},
		"403": {Description: "Forbidden", Content: simpleContent},
		"404": {Description: "Not found", Content: simpleContent},
		"500": {Description: "Internal server error", Content: simpleContent},
	}
}

func init() {
	pkgrest.SetupErrorFactory()
	huma.DefaultArrayNullable = false

	Router = chi.NewMux()
	appName := env.GetString("APP_NAME", "TEARest")

	// --- IAM Server ---
	iamCfg := huma.DefaultConfig("IAM Server", Version)
	iamCfg.DocsPath = ""
	iamCfg.Servers = []*huma.Server{{URL: "/auth", Description: "IAM (authentication)"}}
	iamRouter := chi.NewRouter()
	IAMServer = &Server{
		Name:     "IAM Server",
		BasePath: "/auth",
		Router:   iamRouter,
		API:      humachi.New(iamRouter, iamCfg),
	}
	IAMBadInputSchema = huma.SchemaFromType(IAMServer.API.OpenAPI().Components.Schemas, reflect.TypeOf(errors.APIErrorResponse{}))
	IAMErrorSchema = huma.SchemaFromType(IAMServer.API.OpenAPI().Components.Schemas, reflect.TypeOf(errors.APIErrorSimple{}))
	serveOpenAPI(iamRouter, IAMServer)
	Router.Mount("/auth", iamRouter)

	// --- Todo Server ---
	todoCfg := huma.DefaultConfig("Todo Server", Version)
	todoCfg.DocsPath = ""
	todoCfg.Servers = []*huma.Server{{URL: "/todos", Description: "Todo tasks"}}
	todoRouter := chi.NewRouter()
	TodoServer = &Server{
		Name:     "Todo Server",
		BasePath: "/todos",
		Router:   todoRouter,
		API:      humachi.New(todoRouter, todoCfg),
	}
	_ = huma.SchemaFromType(TodoServer.API.OpenAPI().Components.Schemas, reflect.TypeOf(errors.APIErrorResponse{}))
	_ = huma.SchemaFromType(TodoServer.API.OpenAPI().Components.Schemas, reflect.TypeOf(errors.APIErrorSimple{}))
	serveOpenAPI(todoRouter, TodoServer)
	Router.Mount("/todos", todoRouter)

	// --- X Server ---
	xCfg := huma.DefaultConfig("X Server", Version)
	xCfg.DocsPath = ""
	xCfg.Servers = []*huma.Server{{URL: "/x", Description: "X tasks"}}
	xRouter := chi.NewRouter()
	XServer = &Server{
		Name:     "X Server",
		BasePath: "/x",
		Router:   xRouter,
		API:      humachi.New(xRouter, xCfg),
	}
	_ = huma.SchemaFromType(XServer.API.OpenAPI().Components.Schemas, reflect.TypeOf(errors.APIErrorResponse{}))
	_ = huma.SchemaFromType(XServer.API.OpenAPI().Components.Schemas, reflect.TypeOf(errors.APIErrorSimple{}))
	serveOpenAPI(xRouter, XServer)
	Router.Mount("/x", xRouter)

	Router.Get("/docs", serveSwaggerUI(appName))
}

func serveOpenAPI(sub chi.Router, s *Server) {
	spec := s.API.OpenAPI()
	sub.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(spec)
	})
}

func serveSwaggerUI(title string) http.HandlerFunc {
	html := `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>` + title + ` — API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        urls: [
          { url: "/auth/openapi.json", name: "IAM Server" },
          { url: "/todos/openapi.json", name: "Todo Server" },
          { url: "/x/openapi.json", name: "X Server" }
        ],
        dom_id: "#swagger-ui",
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>`
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}
