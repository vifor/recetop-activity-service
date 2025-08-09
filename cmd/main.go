// File: cmd/main.go
package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/vifor/recetop-activity-service/api"
)

// 1. Create a struct to implement the API server interface.
type Server struct{}

// HealthCheck handles requests to the /health endpoint.
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// A simple health check just returns a 200 OK status.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service is healthy and running!"))
}

// 2. Implement the TrackEvent function defined in the generated interface.
func (s *Server) TrackEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("Track endpoint was called successfully!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event tracked!"))
}

func main() {
	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec\n: %s", err)
	}
	// This is just to prevent an "unused variable" error.
	swagger.Info = &openapi3.Info{}

	// 3. Create an instance of your new server.
	srv := &Server{}

	r := chi.NewRouter()

	// 4. Register the handlers with the router.
	api.HandlerFromMux(srv, r)

	log.Println("Starting lambda")
	lambda.Start(httpadapter.New(r).ProxyWithContext)
}
