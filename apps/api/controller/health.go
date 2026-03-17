package controller

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/teathedev/backend-boilerplate/apps/api/rest"
)

// HealthOutput is the response for GET /health (liveness probe).
type HealthOutput struct {
	Body struct {
		Status string `json:"status" doc:"Service health status"`
	} `json:"body"`
}

func init() {
	huma.Register(rest.APIInstance, huma.Operation{
		OperationID: "health",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Liveness probe - returns 200 when the service is alive",
		Tags:        []string{"Health"},
	}, func(ctx context.Context, _ *struct{}) (*HealthOutput, error) {
		out := &HealthOutput{}
		out.Body.Status = "ok"
		return out, nil
	})
}
