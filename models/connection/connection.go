package connection

import "time"

type Status string
type Type string

const (
	StatusActive     Status = "active"
	StatusIdle       Status = "idle"
	StatusTerminated Status = "terminated"
)

const (
	TypeControl Type = "control"
	TypeWorker  Type = "worker"
)

type Connection struct {
	ID            string     `json:"id"`
	ServeName     string     `json:"serve_name"`
	RemoteAddr    string     `json:"remote_addr"`
	Status        Status     `json:"status"`
	Type          Type       `json:"type"`
	ConnectedAt   time.Time  `json:"connected_at"`
	LastHeartbeat *time.Time `json:"last_heartbeat,omitempty"`
}
