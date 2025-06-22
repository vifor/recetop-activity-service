package main

import (
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vifor/recetop-activity-service/api" // Import our generated api package
)

// Server is our custom struct which will hold any dependencies and implement our API.
type Server struct {
	// In the future, you could add dependencies like a database connection here.
}

// Ensure our Server struct satisfies the generated ServerInterface.
// The compiler will error here if we don't implement all the interface's methods.
var _ api.ServerInterface = (*Server)(nil)

// HealthCheck implements the `healthCheck` operation defined in our openapi.yaml
func (s *Server) HealthCheck(ctx echo.Context) error {
	// This is our business logic from before.
	message := "Go serverless with OpenAPI spec!" // 1. Store the string in a variable.

	response := api.HealthStatus{
		Status:  "UP",
		Message: &message, // 2. Take the address of the variable. This is now valid.
	}

	return ctx.JSON(http.StatusOK, response)
}

func main() {
	// Create an instance of our server
	s := &Server{}

	// Create a new Echo instance
	e := echo.New()

	// Add standard middleware
	e.Use(middleware.Logger()) // This is what generated the helpful log!
	e.Use(middleware.Recover())

	// Define the base path that API Gateway sends to our Lambda.

	basePath := "/default/recetop-activity-service"

	// Create a group for our API handlers. All routes registered on this
	// group will be nested under the base path.
	apiGroup := e.Group(basePath)

	// Register our server's handlers (like GET /) onto the group.
	// oapi-codegen will now correctly register our HealthCheck on the path
	// "/default/recetop-activity-service/".
	api.RegisterHandlers(apiGroup, s)

	// The algnhsa library starts our Echo server and translates Lambda
	// events into standard HTTP requests that Echo can understand.
	algnhsa.ListenAndServe(e, nil)
}
