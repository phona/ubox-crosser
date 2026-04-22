package contract

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ApiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginData struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type Service struct {
	Name          string `json:"name"`
	Key           string `json:"key"`
	Method        string `json:"method"`
	Address       string `json:"address"`
	ExposeAddress string `json:"expose_address,omitempty"`
	LoginPassword string `json:"login_password,omitempty"`
	AuthPassword  string `json:"auth_password,omitempty"`
	LogFile       string `json:"log_file,omitempty"`
	LogLevel      string `json:"log_level,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

type ServiceListData struct {
	Services []Service `json:"services"`
	Total    int       `json:"total"`
}

type ServiceDetailData struct {
	Service        Service         `json:"service"`
	ProxyInstances []ProxyInstance `json:"proxy_instances"`
}

type ProxyInstance struct {
	InstanceID      string `json:"instance_id"`
	ServiceName     string `json:"service_name"`
	Address         string `json:"address"`
	Status          string `json:"status"`
	LastHeartbeat   string `json:"last_heartbeat,omitempty"`
	ConnectionCount int    `json:"connection_count,omitempty"`
}

type ProxyRegisterRequest struct {
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
	Address     string `json:"address"`
}

type ProxyRegisterData struct {
	Token string `json:"token"`
}

type ProxyHeartbeatRequest struct {
	InstanceID      string `json:"instance_id"`
	ConnectionCount int    `json:"connection_count"`
}

func jsonBody(t *testing.T, v interface{}) io.Reader {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}
	return bytes.NewReader(b)
}

func doRequest(t *testing.T, srv *httptest.Server, method, path string, body io.Reader, headers map[string]string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, srv.URL+path, body)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request %s %s failed: %v", method, path, err)
	}
	return resp
}

func parseResponse(t *testing.T, resp *http.Response) ApiResponse {
	t.Helper()
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	var ar ApiResponse
	if err := json.Unmarshal(body, &ar); err != nil {
		t.Fatalf("unmarshal response: %v (body: %s)", err, string(body))
	}
	return ar
}

func parseData[T any](t *testing.T, ar ApiResponse) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(ar.Data, &v); err != nil {
		t.Fatalf("unmarshal data: %v (raw: %s)", err, string(ar.Data))
	}
	return v
}
