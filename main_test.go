package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEndpointsViaHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Root status",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Contineo API is running",
		},
		{
			name:           "Health check",
			path:           "/health",
			expectedStatus: http.StatusOK,
			expectedBody:   `"status":"ok"`,
		},
		{
			name:           "API Health check",
			path:           "/api/health",
			expectedStatus: http.StatusOK,
			expectedBody:   `"status":"ok"`,
		},
		{
			name:           "OpenAPI JSON",
			path:           "/openapi.json",
			expectedStatus: http.StatusOK,
			expectedBody:   `"openapi"`,
		},
		{
			name:           "Swagger Docs",
			path:           "/docs",
			expectedStatus: http.StatusOK,
			expectedBody:   "Contineo API Docs",
		},
		{
			name:           "Dashboard missing params",
			path:           "/api/dashboard",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "usn and dob query parameters are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			Handler(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			body := rec.Body.String()
			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, body)
			}
		})
	}
}
