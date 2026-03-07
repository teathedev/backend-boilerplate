package main

import (
	"context"
	"fmt"

	"github.com/teathedev/backend-boilerplate/internal/actions"
)

func main() {
	ctx := context.Background()
	if err := actions.RefreshJWKs(ctx); err != nil {
		fmt.Printf("Warning: failed to load JWKs: %v\n", err)
	}
	fmt.Println("Hello World")
}
