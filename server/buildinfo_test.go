package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestBuildInfoHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		buildIDEnv     string
		expectedStatus int
		expectedGitSHA bool
		expectedBuildID string
		expectedGoVersion string
	}{
		{
			name:              "GET request returns buildinfo",
			method:            "GET",
			buildIDEnv:        "",
			expectedStatus:    200,
			expectedGitSHA:    true,
			expectedBuildID:   "dev",
			expectedGoVersion: "go1.23",
		},
		{
			name:              "GET request with BUILD_ID env set",
			method:            "GET",
			buildIDEnv:        "build-123",
			expectedStatus:    200,
			expectedGitSHA:    true,
			expectedBuildID:   "build-123",
			expectedGoVersion: "go1.23",
		},
		{
			name:           "POST request returns method not allowed",
			method:         "POST",
			expectedStatus: 405,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set BUILD_ID env var if specified
			oldBuildID := os.Getenv("BUILD_ID")
			if tt.buildIDEnv != "" {
				os.Setenv("BUILD_ID", tt.buildIDEnv)
			} else {
				os.Unsetenv("BUILD_ID")
			}
			defer os.Setenv("BUILD_ID", oldBuildID)

			// Create proxy server
			configs := make(map[string]interface{})
			proxy := &ProxyServer{
				gitSHA: "abc1234",
			}

			// Create request and recorder
			req, err := http.NewRequest(tt.method, "/buildinfo", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()

			// Call handler
			proxy.buildInfoHandler(recorder, req)

			// Check status
			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}

			// Check response body for successful requests
			if tt.expectedStatus == 200 {
				var body map[string]interface{}
				if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				// Verify git_sha
				if gitSHA, ok := body["git_sha"].(string); !ok || gitSHA != "abc1234" {
					t.Errorf("Expected git_sha 'abc1234', got %v", body["git_sha"])
				}

				// Verify build_id
				if buildID, ok := body["build_id"].(string); !ok || buildID != tt.expectedBuildID {
					t.Errorf("Expected build_id '%s', got %v", tt.expectedBuildID, body["build_id"])
				}

				// Verify go_version
				if goVersion, ok := body["go_version"].(string); !ok || goVersion != "go1.23" {
					t.Errorf("Expected go_version 'go1.23', got %v", body["go_version"])
				}

				// Verify Content-Type
				contentType := recorder.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
				}
			}
		})
	}
}

func TestBuildInfoGitSHAFormat(t *testing.T) {
	proxy := &ProxyServer{
		gitSHA: "abc1234",
	}

	req, _ := http.NewRequest("GET", "/buildinfo", nil)
	recorder := httptest.NewRecorder()
	proxy.buildInfoHandler(recorder, req)

	var body map[string]interface{}
	json.NewDecoder(recorder.Body).Decode(&body)

	gitSHA := body["git_sha"].(string)
	if len(gitSHA) != 7 {
		t.Errorf("Expected git_sha to be 7 characters, got %d: %s", len(gitSHA), gitSHA)
	}

	// Verify it's hexadecimal
	for _, ch := range gitSHA {
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f')) {
			t.Errorf("git_sha contains non-hex character: %c", ch)
		}
	}
}

func TestBuildInfoEmptyGitSHA(t *testing.T) {
	proxy := &ProxyServer{
		gitSHA: "",
	}

	req, _ := http.NewRequest("GET", "/buildinfo", nil)
	recorder := httptest.NewRecorder()
	proxy.buildInfoHandler(recorder, req)

	var body map[string]interface{}
	json.NewDecoder(recorder.Body).Decode(&body)

	gitSHA := body["git_sha"].(string)
	if gitSHA != "" {
		t.Errorf("Expected empty git_sha when not set, got: %s", gitSHA)
	}
}

