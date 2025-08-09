// File: cmd/main.go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/go-chi/chi/v5"
	"github.com/vifor/recetop-activity-service/api"
)

// Server struct remains the same.
type Server struct{}

// HealthCheck function remains the same.
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service is healthy and running!"))
}

// TrackEvent function remains the same.
func (s *Server) TrackEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("TrackEvent handler was called successfully!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event tracked!"))
}

func main() {
	// This log remains. It tells us the Go runtime has started.
	log.Println("--- LAMBDA INIT ---")

	// Router setup is the same.
	srv := &Server{}
	r := chi.NewRouter()
	api.HandlerFromMux(srv, r)

	// We create the adapter that connects our Chi router to Lambda.
	adapter := httpadapter.New(r)

	// --- NEW LOGGING HANDLER ---
	// Instead of starting the adapter directly, we start our own custom handler.
	// This lets us see the raw request from API Gateway.
	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		// LOG 1: The most important log. It prints the HTTP Method and Path.
		log.Printf("EVENT: Received %s request for %s", req.HTTPMethod, req.Path)

		// LOG 2: This prints the request body. For a GET it will be empty.
		// For your POST /track, it will show the full JSON payload.
		log.Printf("BODY: %s", req.Body)

		// After logging, we pass the request to the router adapter to handle it.
		return adapter.ProxyWithContext(ctx, req)
	})
}

// NOTE: We removed the api.GetSwagger() check from main() as it wasn't essential
// for the core logic and can sometimes clutter the startup process.
