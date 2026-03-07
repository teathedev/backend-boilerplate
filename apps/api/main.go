package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	_ "github.com/teathedev/backend-boilerplate/apps/api/controller"
	"github.com/teathedev/backend-boilerplate/apps/api/rest"
	"github.com/teathedev/backend-boilerplate/internal/actions"
	"github.com/teathedev/backend-boilerplate/pkg/env"
)

func main() {
	ctx := context.Background()
	if err := actions.RefreshJWKs(ctx); err != nil {
		fmt.Printf("Warning: failed to load JWKs: %v\n", err)
	}

	port := env.GetNumber("PORT", 8888)
	addr := fmt.Sprintf(":%d", port)

	log.Printf("Swagger UI: http://localhost%s/docs (select IAM / Todo / X from dropdown)", addr)
	log.Printf("OpenAPI:    %s/auth/openapi.json | %s/todos/openapi.json | %s/x/openapi.json", addr, addr, addr)
	log.Fatal(http.ListenAndServe(addr, rest.Router))
}
