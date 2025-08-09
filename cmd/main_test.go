package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert" // popular library for clean test assertions
)

// TestHealthCheck is our in-memory integration test for the HealthCheck endpoint.
// setupTestServer creates a new Echo server instance and registers all our handlers.
// It returns the Echo instance so our tests can use it.
func setupTestServer() *echo.Echo {
	s := &Server{}
	e := echo.New()

	basePath := "/default/recetop-activity-service"
	e.GET(basePath, s.HealthCheck)
	e.POST(basePath+"/track", s.TrackEvent)

	return e
}

func TestHealthCheck(t *testing.T) {
	// 1. SETUP: Get a configured server from our helper.
	e := setupTestServer()

	// 2. EXECUTION
	req := httptest.NewRequest(http.MethodGet, "/default/recetop-activity-service", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// 3. ASSERTION
	assert.Equal(t, http.StatusOK, rec.Code)
	// ... (rest of assertions)
}

func TestTrackEvent(t *testing.T) {
	// 1. SETUP: Get a configured server from our helper.
	e := setupTestServer()

	// --- THIS IS THE FIX ---
	// Provide a complete and valid JSON string for the request body.
	eventJSON := `{
		"type": "track",
		"event": "Test Event From Test Suite",
		"userId": "user-test-123",
		"anonymousId": "anon-test-456",
		"timestamp": "2025-07-01T10:00:00Z",
		"properties": {
			"screen_name": "TestScreen"
		}
	}`

	// 2. EXECUTION
	req := httptest.NewRequest(http.MethodPost, "/default/recetop-activity-service/track", strings.NewReader(eventJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// 3. ASSERTION
	assert.Equal(t, http.StatusAccepted, rec.Code)
}
