package message

const (
	LOGIN = iota
	HEART_BEAT
	GEN_WORKER

	SUCCESS
	FAILED
)

type Message struct {
	Type     int64  `json:"type"`
	Password string `json:"password"`
}

type ResultMessage struct {
	Result int64 `json:"type"`
}
