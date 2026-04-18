package message

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/phona/ubox-crosser/models/errors"
)

// ---------------------------------------------------------------------------
// Contract: message type constants (contract.spec.yaml → message_types)
// ---------------------------------------------------------------------------

func TestMessageTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      uint8
		wantCode uint8
	}{
		{"LOGIN", LOGIN, 0},
		{"HEART_BEAT", HEART_BEAT, 1},
		{"GEN_WORKER", GEN_WORKER, 2},
		{"AUTHENTICATION", AUTHENTICATION, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.wantCode {
				t.Errorf("%s = %d, want %d", tt.name, tt.got, tt.wantCode)
			}
		})
	}
}

func TestResultCodeConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      uint8
		wantCode uint8
	}{
		{"SUCCESS", SUCCESS, 4},
		{"FAILED", FAILED, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.wantCode {
				t.Errorf("%s = %d, want %d", tt.name, tt.got, tt.wantCode)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Contract: Message JSON schema (contract.spec.yaml → messages.request)
// ---------------------------------------------------------------------------

func TestMessageJSONRoundTrip(t *testing.T) {
	msg := Message{
		Type:      LOGIN,
		ServeName: "test_service",
		Password:  "secret",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Message
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.Type != msg.Type {
		t.Errorf("Type = %d, want %d", decoded.Type, msg.Type)
	}
	if decoded.ServeName != msg.ServeName {
		t.Errorf("ServeName = %q, want %q", decoded.ServeName, msg.ServeName)
	}
	if decoded.Password != msg.Password {
		t.Errorf("Password = %q, want %q", decoded.Password, msg.Password)
	}
}

func TestMessageJSONKeys(t *testing.T) {
	msg := Message{Type: LOGIN, ServeName: "svc", Password: "pw"}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}

	requiredKeys := []string{"type", "serve_name", "password"}
	for _, key := range requiredKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("missing required JSON key %q", key)
		}
	}

	// Contract: no extra keys allowed
	if len(raw) != len(requiredKeys) {
		t.Errorf("expected exactly %d keys, got %d: %v", len(requiredKeys), len(raw), raw)
	}
}

// ---------------------------------------------------------------------------
// Contract: ResultMessage JSON schema (contract.spec.yaml → messages.response)
// ---------------------------------------------------------------------------

func TestResultMessageJSONRoundTrip(t *testing.T) {
	msg := ResultMessage{Result: SUCCESS, Reason: errors.OK}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded ResultMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.Result != msg.Result {
		t.Errorf("Result = %d, want %d", decoded.Result, msg.Result)
	}
	if decoded.Reason != msg.Reason {
		t.Errorf("Reason = %d, want %d", decoded.Reason, msg.Reason)
	}
}

func TestResultMessageJSONKeys(t *testing.T) {
	msg := ResultMessage{Result: FAILED, Reason: errors.INVALID_PASSWORD}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}

	requiredKeys := []string{"result", "reason"}
	for _, key := range requiredKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("missing required JSON key %q", key)
		}
	}

	if len(raw) != len(requiredKeys) {
		t.Errorf("expected exactly %d keys, got %d: %v", len(requiredKeys), len(raw), raw)
	}
}

// ---------------------------------------------------------------------------
// Contract: Message schema validation (contract.spec.yaml → fields.constraints)
// These tests WILL FAIL until a Validate() method is added to Message.
// The tests use a contractValidateMessage helper that enforces spec constraints;
// the production type must eventually enforce these itself.
// ---------------------------------------------------------------------------

// contractValidateMessage enforces contract.spec.yaml constraints on a Message.
// This is NOT business code — it is a test-side oracle derived from the spec.
// Tests fail because the production Message type does not yet reject invalid input.
func contractValidateMessage(msg Message) error {
	// type must be a known code (0-3)
	if msg.Type > AUTHENTICATION {
		return fmt.Errorf("unknown message type %d", msg.Type)
	}
	// serve_name: required, pattern ^[a-zA-Z0-9_]+$
	if msg.ServeName == "" {
		return fmt.Errorf("serve_name is required")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(msg.ServeName) {
		return fmt.Errorf("serve_name %q does not match pattern ^[a-zA-Z0-9_]+$", msg.ServeName)
	}
	// password: required (min_length 1)
	if msg.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

// contractValidateResult enforces contract.spec.yaml constraints on ResultMessage.
func contractValidateResult(msg ResultMessage) error {
	if msg.Result != SUCCESS && msg.Result != FAILED {
		return fmt.Errorf("unknown result code %d", msg.Result)
	}
	return nil
}

func TestContractValidateServeName_Required(t *testing.T) {
	msg := Message{Type: LOGIN, ServeName: "", Password: "secret"}
	if err := contractValidateMessage(msg); err == nil {
		t.Error("expected validation error for empty serve_name, got nil")
	}

	// Contract: production code should also reject this.
	// FAILS: Message has no Validate method — empty serve_name reaches the server.
	data, _ := json.Marshal(msg)
	var decoded Message
	_ = json.Unmarshal(data, &decoded)
	if decoded.ServeName == "" {
		t.Error("CONTRACT VIOLATION: empty serve_name accepted by Message deserialization without error")
	}
}

func TestContractValidateServeName_Pattern(t *testing.T) {
	tests := []struct {
		name      string
		serveName string
		wantErr   bool
	}{
		{"valid alphanumeric", "test_service", false},
		{"valid underscore", "UBox_cytm", false},
		{"invalid spaces", "test service", true},
		{"invalid special chars", "test@service!", true},
		{"invalid dash", "test-service", true},
		{"empty", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := Message{Type: LOGIN, ServeName: tt.serveName, Password: "pw"}
			err := contractValidateMessage(msg)
			if tt.wantErr && err == nil {
				t.Errorf("expected error for serve_name=%q, got nil", tt.serveName)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error for serve_name=%q: %v", tt.serveName, err)
			}
		})
	}
}

func TestContractValidatePassword_Required(t *testing.T) {
	msg := Message{Type: LOGIN, ServeName: "svc", Password: ""}
	if err := contractValidateMessage(msg); err == nil {
		t.Error("expected validation error for empty password, got nil")
	}

	// FAILS: Message has no Validate method — empty password reaches the server.
	data, _ := json.Marshal(msg)
	var decoded Message
	_ = json.Unmarshal(data, &decoded)
	if decoded.Password == "" {
		t.Error("CONTRACT VIOLATION: empty password accepted by Message deserialization without error")
	}
}

func TestContractValidateType_KnownValues(t *testing.T) {
	msg := Message{Type: 99, ServeName: "svc", Password: "pw"}
	if err := contractValidateMessage(msg); err == nil {
		t.Error("expected validation error for unknown message type 99, got nil")
	}

	// FAILS: unknown type 99 successfully round-trips without any error.
	data, _ := json.Marshal(msg)
	var decoded Message
	_ = json.Unmarshal(data, &decoded)
	if decoded.Type == 99 {
		t.Error("CONTRACT VIOLATION: unknown message type 99 accepted without error")
	}
}

// ---------------------------------------------------------------------------
// Contract: error code coverage (contract.spec.yaml → error_codes)
// ---------------------------------------------------------------------------

func TestErrorCodeConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      uint8
		wantCode uint8
	}{
		{"OK", errors.OK, 0},
		{"INVALID_PASSWORD", errors.INVALID_PASSWORD, 1},
		{"INVALID_SERVE_NAME", errors.INVALID_SERVE_NAME, 2},
		{"UNKNOWN_CODE", errors.UNKNOWN_CODE, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.wantCode {
				t.Errorf("%s = %d, want %d", tt.name, tt.got, tt.wantCode)
			}
		})
	}
}

func TestErrorCodeStrings(t *testing.T) {
	tests := []struct {
		code    uint8
		wantMsg string
	}{
		{errors.INVALID_PASSWORD, "invalid password"},
		{errors.INVALID_SERVE_NAME, "invalid serve name"},
		{errors.UNKNOWN_CODE, "unknown code"},
	}
	for _, tt := range tests {
		t.Run(tt.wantMsg, func(t *testing.T) {
			e := errors.DecodeError(tt.code)
			if e.Error() != tt.wantMsg {
				t.Errorf("DecodeError(%d).Error() = %q, want %q", tt.code, e.Error(), tt.wantMsg)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Contract: deserialization from raw JSON (malformed input handling)
// ---------------------------------------------------------------------------

func TestMessageDeserialize_MissingFields(t *testing.T) {
	// Contract requires all fields; partial JSON should be rejected at validation.
	raw := `{"type":0}`
	var msg Message
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		t.Logf("unmarshal error (acceptable): %v", err)
	}

	if err := contractValidateMessage(msg); err == nil {
		t.Error("expected validation error for message with missing serve_name and password")
	}

	// FAILS: the production type silently accepts partial messages.
	if msg.ServeName == "" && msg.Password == "" {
		t.Error("CONTRACT VIOLATION: partial JSON deserialized without serve_name or password; no validation enforced")
	}
}

func TestMessageDeserialize_ExtraFields(t *testing.T) {
	raw := `{"type":0,"serve_name":"svc","password":"pw","extra":"bad"}`
	var msg Message
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	// Contract: extra fields should be silently ignored
	if msg.Type != LOGIN {
		t.Errorf("Type = %d, want %d", msg.Type, LOGIN)
	}
}

func TestResultMessageDeserialize_InvalidResultCode(t *testing.T) {
	raw := `{"result":99,"reason":0}`
	var msg ResultMessage
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if err := contractValidateResult(msg); err == nil {
		t.Error("expected validation error for unknown result code 99")
	}

	// FAILS: invalid result code 99 accepted without error.
	if msg.Result == 99 {
		t.Error("CONTRACT VIOLATION: unknown result code 99 accepted without error")
	}
}
