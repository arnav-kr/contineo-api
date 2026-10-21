package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arnav-kr/contineo-api/services"
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

func TestPortalHelpers(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectedURL string
	}{
		{"even", "even", "https://nmitparents.contineo.in"},
		{"", "even", "https://nmitparents.contineo.in"},
		{"odd", "odd", "https://nmitparents.contineo.in/parentsodd"},
		{"parentsodd", "odd", "https://nmitparents.contineo.in/parentsodd"},
		{"/parentsodd/", "odd", "https://nmitparents.contineo.in/parentsodd"},
	}

	for _, tt := range tests {
		t.Run("Portal_"+tt.input, func(t *testing.T) {
			got := services.NormalizePortal(tt.input)
			if got != tt.expected {
				t.Errorf("NormalizePortal(%q) = %q, want %q", tt.input, got, tt.expected)
			}
			url := services.GetBaseURL(tt.input)
			if url != tt.expectedURL {
				t.Errorf("GetBaseURL(%q) = %q, want %q", tt.input, url, tt.expectedURL)
			}
		})
	}
}

func TestMockStudentEndpoints(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Mock Dashboard",
			path:           "/api/dashboard?usn=1NT24CS006&dob=2000-01-01",
			expectedStatus: http.StatusOK,
			expectedBody:   `"usn":"1NT24CS006"`,
		},
		{
			name:           "Mock Attendance",
			path:           "/api/attendance?usn=1NT24CS006&dob=any-date&course=21CS41",
			expectedStatus: http.StatusOK,
			expectedBody:   `"totalClasses":40`,
		},
		{
			name:           "Mock CIE",
			path:           "/api/cie?usn=1NT24CS006&dob=random&course=21CS41",
			expectedStatus: http.StatusOK,
			expectedBody:   `"raw":"44.0 / 50"`,
		},
		{
			name:           "Mock Exam History",
			path:           "/api/exam-history?usn=1NT24CS006&dob=any",
			expectedStatus: http.StatusOK,
			expectedBody:   `"creditsEarned":64`,
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
