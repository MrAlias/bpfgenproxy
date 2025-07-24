package main

import (
	"context"
	"log"

	"go.opentelemetry.io/auto"
)

func main() {
	_, err := auto.NewInstrumentation(context.Background())
	if err != nil {
		log.Fatalf("failed to create instrumentation: %v", err)
	}
}
