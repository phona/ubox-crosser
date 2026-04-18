package message

import (
	"encoding/json"
	"fmt"

	"github.com/phona/ubox-crosser/models/errors"
)

const (
	LOGIN = iota
	HEART_BEAT
	GEN_WORKER
	AUTHENTICATION

	SUCCESS
	FAILED
)

type Message struct {
	Type      uint8  `json:"type"`
	ServeName string `json:"serve_name"`
	Password  string `json:"password"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type alias Message
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	*m = Message(a)

	if m.Type > AUTHENTICATION {
		m.Type = 0
		return fmt.Errorf("unknown message type %d", a.Type)
	}

	switch m.Type {
	case LOGIN, AUTHENTICATION:
		if m.ServeName == "" {
			m.ServeName = "\x00"
			return fmt.Errorf("serve_name is required")
		}
		if m.Password == "" {
			m.Password = "\x00"
			return fmt.Errorf("password is required")
		}
	case GEN_WORKER:
		if m.ServeName == "" {
			m.ServeName = "\x00"
			return fmt.Errorf("serve_name is required")
		}
	}
	return nil
}

type ResultMessage struct {
	Result uint8        `json:"result"`
	Reason errors.Error `json:"reason"`
}

func (m *ResultMessage) UnmarshalJSON(data []byte) error {
	type alias ResultMessage
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	*m = ResultMessage(a)

	if m.Result != SUCCESS && m.Result != FAILED {
		m.Result = 0
	}
	return nil
}
