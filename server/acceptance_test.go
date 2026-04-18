package server

import (
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/phona/ubox-crosser/models/config"
	"github.com/phona/ubox-crosser/models/errors"
	"github.com/phona/ubox-crosser/models/message"
	"github.com/phona/ubox-crosser/utils/connector"
)

// helpers

func setupProxyServer(t *testing.T) *ProxyServer {
	t.Helper()
	configs := map[string]config.ServerConfig{
		"test_service": {
			LoginPass: "correct-password",
			AuthPass:  "correct-auth-pass",
			Config:    config.Config{Key: "testkey", Method: ""},
		},
	}
	return &ProxyServer{
		controllers: make(map[string]*controller),
		errs:        make(chan error, 100),
		context:     configs,
	}
}

func newPipePair(t *testing.T) (*connector.Coordinator, net.Conn) {
	t.Helper()
	serverConn, clientConn := net.Pipe()
	t.Cleanup(func() {
		serverConn.Close()
		clientConn.Close()
	})
	return connector.AsCoordinator(serverConn), clientConn
}

func sendMessage(t *testing.T, conn net.Conn, msg message.Message) {
	t.Helper()
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	coord := connector.AsCoordinator(conn)
	if err := coord.SendMsg(string(data)); err != nil {
		t.Fatalf("send message: %v", err)
	}
}

func readResponse(t *testing.T, conn net.Conn) message.ResultMessage {
	t.Helper()
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	coord := connector.AsCoordinator(conn)
	content, err := coord.ReadMsg()
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	var resp message.ResultMessage
	if err := json.Unmarshal([]byte(content), &resp); err != nil {
		t.Fatalf("unmarshal response %q: %v", content, err)
	}
	return resp
}

// ---------------------------------------------------------------------------
// Scenario: Client login with valid credentials
// ---------------------------------------------------------------------------

func TestAcceptance_LoginSuccess(t *testing.T) {
	// Given a proxy server configured with "test_service" and password "correct-password"
	proxy := setupProxyServer(t)
	serverCoord, clientConn := newPipePair(t)

	// When the client sends a LOGIN message with correct serve_name and password
	go proxy.handleLoginRequest("test_service", "correct-password", serverCoord)

	// Then the server responds with SUCCESS and reason OK
	resp := readResponse(t, clientConn)
	if resp.Result != message.SUCCESS {
		t.Errorf("Given valid credentials, expected SUCCESS (%d), got %d", message.SUCCESS, resp.Result)
	}
	if resp.Reason != errors.OK {
		t.Errorf("Given valid credentials, expected reason OK (%d), got %d", errors.OK, resp.Reason)
	}

	// And a controller is registered for the service
	if _, ok := proxy.controllers["test_service"]; !ok {
		t.Error("Given successful login, expected controller to be registered for test_service")
	}
}

// ---------------------------------------------------------------------------
// Scenario: Client login with wrong password
// ---------------------------------------------------------------------------

func TestAcceptance_LoginWrongPassword(t *testing.T) {
	// Given a proxy server configured with "test_service"
	proxy := setupProxyServer(t)
	serverCoord, clientConn := newPipePair(t)

	// When the client sends a LOGIN message with wrong password
	go proxy.handleLoginRequest("test_service", "wrong-password", serverCoord)

	// Then the server responds with FAILED and reason INVALID_PASSWORD
	resp := readResponse(t, clientConn)
	if resp.Result != message.FAILED {
		t.Errorf("Given wrong password, expected FAILED (%d), got %d", message.FAILED, resp.Result)
	}
	if uint8(resp.Reason) != errors.INVALID_PASSWORD {
		t.Errorf("Given wrong password, expected reason INVALID_PASSWORD (%d), got %d", errors.INVALID_PASSWORD, resp.Reason)
	}

	// And no controller is registered
	if _, ok := proxy.controllers["test_service"]; ok {
		t.Error("Given failed login, expected no controller to be registered")
	}
}

// ---------------------------------------------------------------------------
// Scenario: Client login with unknown service name
// ---------------------------------------------------------------------------

func TestAcceptance_LoginUnknownService(t *testing.T) {
	// Given a proxy server that does NOT have "unknown_svc" configured
	proxy := setupProxyServer(t)
	serverCoord, clientConn := newPipePair(t)

	// When the client sends a LOGIN message for an unknown service
	go proxy.handleLoginRequest("unknown_svc", "any-password", serverCoord)

	// Then the server responds with FAILED and reason INVALID_SERVE_NAME
	resp := readResponse(t, clientConn)
	if resp.Result != message.FAILED {
		t.Errorf("Given unknown service, expected FAILED (%d), got %d", message.FAILED, resp.Result)
	}
	if uint8(resp.Reason) != errors.INVALID_SERVE_NAME {
		t.Errorf("Given unknown service, expected reason INVALID_SERVE_NAME (%d), got %d", errors.INVALID_SERVE_NAME, resp.Reason)
	}
}

// ---------------------------------------------------------------------------
// Scenario: Auth request with valid credentials and active controller
// ---------------------------------------------------------------------------

func TestAcceptance_AuthSuccess(t *testing.T) {
	// Given a proxy server with a registered controller for "test_service"
	proxy := setupProxyServer(t)

	ctrlCoord, _ := newPipePair(t)
	ctrl := newController(ctrlCoord)
	proxy.controllers["test_service"] = ctrl

	// Pre-supply a worker connection so getConn doesn't block
	workerServer, workerClient := net.Pipe()
	t.Cleanup(func() {
		workerServer.Close()
		workerClient.Close()
	})
	ctrl.workConn <- workerClient

	authCoord, authClient := newPipePair(t)

	// When auth-server sends AUTHENTICATION with correct auth password
	go proxy.handleAuthRequest("test_service", "correct-auth-pass", authCoord)

	// Then the server responds with SUCCESS
	resp := readResponse(t, authClient)
	if resp.Result != message.SUCCESS {
		t.Errorf("Given valid auth credentials and active controller, expected SUCCESS (%d), got %d", message.SUCCESS, resp.Result)
	}
}

// ---------------------------------------------------------------------------
// Scenario: Auth request with wrong password
// ---------------------------------------------------------------------------

func TestAcceptance_AuthWrongPassword(t *testing.T) {
	// Given a proxy server with a registered controller for "test_service"
	proxy := setupProxyServer(t)

	ctrlCoord, _ := newPipePair(t)
	proxy.controllers["test_service"] = newController(ctrlCoord)

	authCoord, authClient := newPipePair(t)

	// When auth-server sends AUTHENTICATION with wrong password
	go proxy.handleAuthRequest("test_service", "wrong-auth-pass", authCoord)

	// Then the server responds with FAILED and reason INVALID_PASSWORD
	resp := readResponse(t, authClient)
	if resp.Result != message.FAILED {
		t.Errorf("Given wrong auth password, expected FAILED (%d), got %d", message.FAILED, resp.Result)
	}
	if uint8(resp.Reason) != errors.INVALID_PASSWORD {
		t.Errorf("Given wrong auth password, expected reason INVALID_PASSWORD (%d), got %d", errors.INVALID_PASSWORD, resp.Reason)
	}
}

// ---------------------------------------------------------------------------
// Scenario: Auth request for unknown service
// ---------------------------------------------------------------------------

func TestAcceptance_AuthUnknownService(t *testing.T) {
	// Given a proxy server that does NOT have "unknown_svc" configured
	proxy := setupProxyServer(t)
	authCoord, authClient := newPipePair(t)

	// When auth-server sends AUTHENTICATION for unknown service
	go proxy.handleAuthRequest("unknown_svc", "any-pass", authCoord)

	// Then the server responds with FAILED and reason INVALID_SERVE_NAME
	resp := readResponse(t, authClient)
	if resp.Result != message.FAILED {
		t.Errorf("Given unknown service, expected FAILED (%d), got %d", message.FAILED, resp.Result)
	}
	if uint8(resp.Reason) != errors.INVALID_SERVE_NAME {
		t.Errorf("Given unknown service, expected reason INVALID_SERVE_NAME (%d), got %d", errors.INVALID_SERVE_NAME, resp.Reason)
	}
}

// ---------------------------------------------------------------------------
// Scenario: Auth request when controller is not alive
// ---------------------------------------------------------------------------

func TestAcceptance_AuthNoController(t *testing.T) {
	// Given a proxy server with "test_service" configured but NO active controller
	proxy := setupProxyServer(t)
	authCoord, authClient := newPipePair(t)

	// When auth-server sends AUTHENTICATION for a service with no controller
	go proxy.handleAuthRequest("test_service", "correct-auth-pass", authCoord)

	// Then the server responds with FAILED and reason INVALID_SERVE_NAME
	resp := readResponse(t, authClient)
	if resp.Result != message.FAILED {
		t.Errorf("Given no active controller, expected FAILED (%d), got %d", message.FAILED, resp.Result)
	}
	if uint8(resp.Reason) != errors.INVALID_SERVE_NAME {
		t.Errorf("Given no active controller, expected reason INVALID_SERVE_NAME (%d), got %d", errors.INVALID_SERVE_NAME, resp.Reason)
	}
}

// ---------------------------------------------------------------------------
// Scenario: Unknown message type dispatching
// ---------------------------------------------------------------------------

func TestAcceptance_UnknownMessageType(t *testing.T) {
	// Given a proxy server
	proxy := setupProxyServer(t)
	serverCoord, clientConn := newPipePair(t)

	// When a client sends a message with an unknown type code
	reqMsg := message.Message{Type: 99, ServeName: "test_service", Password: "pw"}

	go func() {
		proxy.handleConnection(serverCoord.Conn)
	}()

	sendMessage(t, clientConn, reqMsg)

	// Then the server responds with FAILED and reason UNKNOWN_CODE
	resp := readResponse(t, clientConn)
	if resp.Result != message.FAILED {
		t.Errorf("Given unknown message type, expected FAILED (%d), got %d", message.FAILED, resp.Result)
	}
	if uint8(resp.Reason) != errors.UNKNOWN_CODE {
		t.Errorf("Given unknown message type, expected reason UNKNOWN_CODE (%d), got %d", errors.UNKNOWN_CODE, resp.Reason)
	}
}

// ---------------------------------------------------------------------------
// Scenario: GEN_WORKER for unregistered service
// ---------------------------------------------------------------------------

func TestAcceptance_GenWorkerUnregistered(t *testing.T) {
	// Given a proxy server with no active controller for "test_service"
	proxy := setupProxyServer(t)
	serverCoord, clientConn := newPipePair(t)

	// When a GEN_WORKER message arrives for an unregistered service
	reqMsg := message.Message{Type: message.GEN_WORKER, ServeName: "test_service"}

	go func() {
		proxy.handleConnection(serverCoord.Conn)
	}()

	sendMessage(t, clientConn, reqMsg)

	// Then the server responds with FAILED and reason INVALID_SERVE_NAME
	resp := readResponse(t, clientConn)
	if resp.Result != message.FAILED {
		t.Errorf("Given unregistered service, expected FAILED (%d), got %d", message.FAILED, resp.Result)
	}
	if uint8(resp.Reason) != errors.INVALID_SERVE_NAME {
		t.Errorf("Given unregistered service, expected reason INVALID_SERVE_NAME (%d), got %d", errors.INVALID_SERVE_NAME, resp.Reason)
	}
}

// ---------------------------------------------------------------------------
// Scenario: Heartbeat echo on control channel
// ---------------------------------------------------------------------------

func TestAcceptance_HeartbeatEcho(t *testing.T) {
	// Given a logged-in client with an active control connection
	serverConn, clientConn := net.Pipe()
	t.Cleanup(func() {
		serverConn.Close()
		clientConn.Close()
	})

	coord := connector.AsCoordinator(serverConn)
	ctrl := newController(coord)

	go ctrl.daemonize()

	// When the client sends a HEART_BEAT message on the control channel
	heartbeat := message.Message{Type: message.HEART_BEAT}
	data, _ := json.Marshal(heartbeat)
	clientCoord := connector.AsCoordinator(clientConn)
	if err := clientCoord.SendMsg(string(data)); err != nil {
		t.Fatalf("send heartbeat: %v", err)
	}

	// Then the server echoes back a HEART_BEAT message
	_ = clientConn.SetReadDeadline(time.Now().Add(2 * time.Second))
	content, err := clientCoord.ReadMsg()
	if err != nil {
		t.Fatalf("read heartbeat response: %v", err)
	}

	var resp message.Message
	if err := json.Unmarshal([]byte(content), &resp); err != nil {
		t.Fatalf("unmarshal heartbeat: %v", err)
	}
	if resp.Type != message.HEART_BEAT {
		t.Errorf("Given heartbeat, expected echoed HEART_BEAT type (%d), got %d", message.HEART_BEAT, resp.Type)
	}
}

// ---------------------------------------------------------------------------
// Scenario: Concurrent logins for different services
// ---------------------------------------------------------------------------

func TestAcceptance_ConcurrentLoginDifferentServices(t *testing.T) {
	// Given a proxy server configured with two services
	configs := map[string]config.ServerConfig{
		"service_a": {
			LoginPass: "pass_a",
			AuthPass:  "auth_a",
		},
		"service_b": {
			LoginPass: "pass_b",
			AuthPass:  "auth_b",
		},
	}
	proxy := &ProxyServer{
		controllers: make(map[string]*controller),
		errs:        make(chan error, 100),
		context:     configs,
	}

	coordA, clientA := newPipePair(t)
	coordB, clientB := newPipePair(t)

	// When two clients login concurrently for different services
	go proxy.handleLoginRequest("service_a", "pass_a", coordA)
	go proxy.handleLoginRequest("service_b", "pass_b", coordB)

	// Then both should succeed independently
	respA := readResponse(t, clientA)
	respB := readResponse(t, clientB)

	if respA.Result != message.SUCCESS {
		t.Errorf("Given concurrent login for service_a, expected SUCCESS, got %d", respA.Result)
	}
	if respB.Result != message.SUCCESS {
		t.Errorf("Given concurrent login for service_b, expected SUCCESS, got %d", respB.Result)
	}

	// And both controllers should be registered
	if _, ok := proxy.controllers["service_a"]; !ok {
		t.Error("Given concurrent logins, expected controller for service_a")
	}
	if _, ok := proxy.controllers["service_b"]; !ok {
		t.Error("Given concurrent logins, expected controller for service_b")
	}
}

// ---------------------------------------------------------------------------
// Scenario: Login replaces existing controller for same service
// ---------------------------------------------------------------------------

func TestAcceptance_LoginReplacesExistingController(t *testing.T) {
	// Given a proxy server with an already-registered controller for "test_service"
	proxy := setupProxyServer(t)

	oldCoord, _ := newPipePair(t)
	proxy.controllers["test_service"] = newController(oldCoord)

	newCoord, clientConn := newPipePair(t)

	// When a new client logs in with the same serve_name
	go proxy.handleLoginRequest("test_service", "correct-password", newCoord)

	// Then the login succeeds
	resp := readResponse(t, clientConn)
	if resp.Result != message.SUCCESS {
		t.Errorf("Given re-login, expected SUCCESS, got %d", resp.Result)
	}

	// And the controller is replaced (new coordinator)
	ctrl := proxy.controllers["test_service"]
	if ctrl == nil {
		t.Fatal("Given re-login, expected controller to exist")
	}
	if ctrl.conn != newCoord {
		t.Error("Given re-login, expected controller to use the new coordinator")
	}
}
