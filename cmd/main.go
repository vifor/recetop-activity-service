package main

import (
	"log"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vifor/recetop-activity-service/api"
)

// Server struct is unchanged
type Server struct{}

// Ensure we satisfy the interface
var _ api.ServerInterface = (*Server)(nil)

// HealthCheck method is unchanged
func (s *Server) HealthCheck(ctx echo.Context) error {
	message := "Go serverless with OpenAPI spec!"
	response := api.HealthStatus{
		Status:  "UP",
		Message: &message,
	}
	return ctx.JSON(http.StatusOK, response)
}

// --- NEW DIAGNOSTIC MIDDLEWARE ---
// This middleware will run for EVERY request, before any other logic.
func diagnosticRequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		log.Println("-------------------------------------------------------")
		log.Printf("[DIAGNOSTIC] Received a request. Details below:")
		log.Printf("[DIAGNOSTIC]   - Method: %s", req.Method)
		log.Printf("[DIAGNOSTIC]   - Path from URL object: %s", req.URL.Path)
		log.Println("-------------------------------------------------------")
		return next(c)
	}
}

func main() {
	// Create an instance of our server
	s := &Server{}

	// Create a new Echo instance
	e := echo.New()

	// You can now remove the diagnosticRequestLogger middleware, as it has served its purpose!
	// The standard Logger is still useful for seeing the final status code.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// --- THE FINAL FIX ---
	// We manually register the exact path we know we are receiving from the logs
	// and point it directly to our type-safe HealthCheck handler method.
	// This removes all ambiguity.
	e.GET("/default/recetop-activity-service", s.HealthCheck)

	// We no longer need the apiGroup or the RegisterHandlers call.

	// Start the adapter. It will now find the explicit route and execute our handler.
	algnhsa.ListenAndServe(e, nil)
}
