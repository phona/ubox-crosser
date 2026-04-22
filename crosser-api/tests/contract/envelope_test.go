package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// REQ-997-S15: All responses follow unified envelope {code, message, data}
func TestSuccessEnvelopeStructure(t *testing.T) {
	srv := httptest.NewServer(stubServicesMux())
	defer srv.Close()

	doRequest(t, srv, http.MethodPost, "/api/v1/services",
		jsonBody(t, map[string]string{
			"name": "env-svc", "key": "k1", "method": "aes-256-cfb", "address": ":8388",
		}), nil)

	resp := doRequest(t, srv, http.MethodGet, "/api/v1/services", nil, nil)
	ar := parseResponse(t, resp)

	if ar.Code != 0 {
		t.Errorf("success response code should be 0, got %d", ar.Code)
	}
	if ar.Message == "" {
		t.Error("success response should have non-empty message")
	}
	if ar.Data == nil {
		t.Error("success response should have data field")
	}
}

func TestErrorEnvelopeStructure(t *testing.T) {
	srv := httptest.NewServer(stubServicesMux())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/v1/services/nonexistent", nil, nil)

	defer func() { _ = resp.Body.Close() }()
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if _, ok := raw["code"]; !ok {
		t.Error("error response must have \"code\" field")
	}
	if _, ok := raw["message"]; !ok {
		t.Error("error response must have \"message\" field")
	}

	var code int
	if err := json.Unmarshal(raw["code"], &code); err != nil {
		t.Fatalf("code is not integer: %v", err)
	}
	if code == 0 {
		t.Error("error response code must be non-zero")
	}

	if dataRaw, ok := raw["data"]; ok {
		if string(dataRaw) != "null" {
			t.Error("error response should not have non-null data field")
		}
	}
}
