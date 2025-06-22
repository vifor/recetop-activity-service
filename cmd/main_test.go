package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert" // popular library for clean test assertions
	"github.com/vifor/recetop-activity-service/api"
)

// TestHealthCheck is our in-memory integration test for the HealthCheck endpoint.
func TestHealthCheck(t *testing.T) {
	// 1. SETUP: Create everything needed to simulate a real HTTP request.
	// --------------------------------------------------------------------
	e := echo.New() // A new Echo instance

	// We create a "dummy" HTTP request object.
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// We create a "recorder" which will capture the HTTP response that our handler writes.
	rec := httptest.NewRecorder()

	// We create an Echo Context, which is what our handler needs to receive.
	c := e.NewContext(req, rec)

	// Create an instance of our server, which contains the handler method.
	s := &Server{}

	// 2. EXECUTION: Call the handler method directly.
	// --------------------------------------------------------------------
	err := s.HealthCheck(c)

	// 3. ASSERTION: Check if the results are what we expect.
	// --------------------------------------------------------------------

	// First, check that our handler didn't return a Go error.
	assert.NoError(t, err)

	// Check that the HTTP status code written to the recorder is 200 OK.
	assert.Equal(t, http.StatusOK, rec.Code)

	// Now, let's check the JSON body of the response.
	var responseBody api.HealthStatus
	err = json.Unmarshal(rec.Body.Bytes(), &responseBody)

	assert.NoError(t, err, "Response body should be valid JSON")
	assert.Equal(t, "UP", responseBody.Status)
	assert.Equal(t, "Go serverless with OpenAPI spec!", *responseBody.Message)
}
