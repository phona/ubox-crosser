package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"

	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"

	"github.com/phona/ubox-crosser/models/message"
	"github.com/phona/ubox-crosser/utils/connector"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

var (
	proxyAddr     = getEnv("PROXY_SERVER_ADDR", "localhost:7000")
	authAddr      = getEnv("AUTH_SERVER_ADDR", "localhost:7002")
	cipherKey     = getEnv("TEST_CIPHER_KEY", "test-secret-key-123")
	cipherMethod  = getEnv("TEST_CIPHER_METHOD", "chacha20")
	loginPassword = getEnv("TEST_LOGIN_PASSWORD", "test-login-pass")
	_ = getEnv("TEST_AUTH_PASSWORD", "test-auth-pass") // reserved for future auth tests
	serveName     = getEnv("TEST_SERVE_NAME", "test_service")
)

func newCipher() *ss.Cipher {
	cipher, err := ss.NewCipher(cipherMethod, cipherKey)
	if err != nil {
		panic(fmt.Sprintf("failed to create cipher: %v", err))
	}
	return cipher
}

func dialEncrypted(addr string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	cipher := newCipher()
	return ss.NewConn(conn, cipher.Copy()), nil
}

// TestProxyServerLogin verifies client login to proxy server
func TestProxyServerLogin(t *testing.T) {
	conn, err := dialEncrypted(proxyAddr)
	if err != nil {
		t.Fatalf("failed to connect to proxy server: %v", err)
	}
	defer conn.Close()

	coordinator := connector.AsCoordinator(conn)

	// Send login message
	reqMsg := message.Message{
		Type:      message.LOGIN,
		ServeName: serveName,
		Password:  loginPassword,
	}
	buf, _ := json.Marshal(reqMsg)
	if err := coordinator.SendMsg(string(buf)); err != nil {
		t.Fatalf("failed to send login: %v", err)
	}

	// Read response
	content, err := coordinator.ReadMsg()
	if err != nil {
		t.Fatalf("failed to read login response: %v", err)
	}

	var respMsg message.ResultMessage
	if err := json.Unmarshal([]byte(content), &respMsg); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if respMsg.Result != message.SUCCESS {
		t.Fatalf("login failed: result=%d reason=%v", respMsg.Result, respMsg.Reason)
	}

	t.Log("Login successful")
}

// TestProxyServerLoginWrongPassword verifies wrong password is rejected
func TestProxyServerLoginWrongPassword(t *testing.T) {
	conn, err := dialEncrypted(proxyAddr)
	if err != nil {
		t.Fatalf("failed to connect to proxy server: %v", err)
	}
	defer conn.Close()

	coordinator := connector.AsCoordinator(conn)

	reqMsg := message.Message{
		Type:      message.LOGIN,
		ServeName: serveName,
		Password:  "wrong-password",
	}
	buf, _ := json.Marshal(reqMsg)
	if err := coordinator.SendMsg(string(buf)); err != nil {
		t.Fatalf("failed to send login: %v", err)
	}

	content, err := coordinator.ReadMsg()
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	var respMsg message.ResultMessage
	if err := json.Unmarshal([]byte(content), &respMsg); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if respMsg.Result != message.FAILED {
		t.Fatalf("expected FAILED result, got %d", respMsg.Result)
	}

	t.Log("Wrong password correctly rejected")
}

// TestProxyServerInvalidServeName verifies unknown service is rejected
func TestProxyServerInvalidServeName(t *testing.T) {
	conn, err := dialEncrypted(proxyAddr)
	if err != nil {
		t.Fatalf("failed to connect to proxy server: %v", err)
	}
	defer conn.Close()

	coordinator := connector.AsCoordinator(conn)

	reqMsg := message.Message{
		Type:      message.LOGIN,
		ServeName: "nonexistent_service",
		Password:  loginPassword,
	}
	buf, _ := json.Marshal(reqMsg)
	if err := coordinator.SendMsg(string(buf)); err != nil {
		t.Fatalf("failed to send login: %v", err)
	}

	content, err := coordinator.ReadMsg()
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	var respMsg message.ResultMessage
	if err := json.Unmarshal([]byte(content), &respMsg); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if respMsg.Result != message.FAILED {
		t.Fatalf("expected FAILED result, got %d", respMsg.Result)
	}

	t.Log("Invalid serve name correctly rejected")
}

// TestAuthServerTunnel verifies the full tunnel: external client → auth-server → proxy-server → client → echo-server
func TestAuthServerTunnel(t *testing.T) {
	// The auth-server exposes a plain TCP port.
	// When we connect, the auth-server gets a worker connection from proxy-server,
	// which in turn asks the client to create a SOCKS5 worker.
	// The data is then bridged bidirectionally.

	// Give services a moment to establish control channels
	time.Sleep(3 * time.Second)

	conn, err := net.DialTimeout("tcp", authAddr, 10*time.Second)
	if err != nil {
		t.Fatalf("failed to connect to auth server: %v", err)
	}
	defer conn.Close()

	// The auth-server bridges us to a SOCKS5 worker on the client side.
	// The client wraps the worker connection in a SOCKS5 server.
	// So we need to speak SOCKS5 protocol here to reach the echo-server.

	// SOCKS5 handshake: version=5, 1 method, no auth
	_, err = conn.Write([]byte{0x05, 0x01, 0x00})
	if err != nil {
		t.Fatalf("failed to send SOCKS5 handshake: %v", err)
	}

	// Read server choice
	buf := make([]byte, 2)
	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		t.Fatalf("failed to read SOCKS5 handshake response: %v", err)
	}
	if buf[0] != 0x05 || buf[1] != 0x00 {
		t.Fatalf("unexpected SOCKS5 handshake response: %x", buf)
	}

	// SOCKS5 connect request to echo-server (echo-server:9000)
	echoHost := "echo-server"
	echoPort := uint16(9000)
	connectReq := []byte{
		0x05,                    // version
		0x01,                    // connect
		0x00,                    // reserved
		0x03,                    // domain name
		byte(len(echoHost)),     // domain length
	}
	connectReq = append(connectReq, []byte(echoHost)...)
	connectReq = append(connectReq, byte(echoPort>>8), byte(echoPort&0xff))

	_, err = conn.Write(connectReq)
	if err != nil {
		t.Fatalf("failed to send SOCKS5 connect: %v", err)
	}

	// Read SOCKS5 connect response (at least 10 bytes for IPv4)
	respBuf := make([]byte, 256)
	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	n, err := conn.Read(respBuf)
	if err != nil {
		t.Fatalf("failed to read SOCKS5 connect response: %v", err)
	}

	if n < 4 {
		t.Fatalf("SOCKS5 connect response too short: %d bytes", n)
	}
	if respBuf[0] != 0x05 {
		t.Fatalf("unexpected SOCKS5 version in response: %x", respBuf[0])
	}
	if respBuf[1] != 0x00 {
		t.Fatalf("SOCKS5 connect failed with status: %x", respBuf[1])
	}

	// Now we have a tunnel to echo-server. Send an HTTP GET request.
	httpReq := "GET / HTTP/1.0\r\nHost: echo-server\r\n\r\n"
	_, err = conn.Write([]byte(httpReq))
	if err != nil {
		t.Fatalf("failed to send HTTP request through tunnel: %v", err)
	}

	// Read HTTP response
	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	response := make([]byte, 4096)
	n, err = conn.Read(response)
	if err != nil && err != io.EOF {
		t.Fatalf("failed to read HTTP response: %v", err)
	}

	responseStr := string(response[:n])
	t.Logf("Received response through tunnel:\n%s", responseStr)

	if n == 0 {
		t.Fatal("empty response from echo server through tunnel")
	}

	t.Log("Full tunnel test passed: external → auth-server → proxy-server → client(SOCKS5) → echo-server")
}
