package crosser

type Config struct {
	Password      string `json:"password"`
	Method        string `json:"method"` // encryption method
	MaxConnection int    `json:"max_connection"`
	Timeout       int64  `json:"timeout"`
}

type ClientConfig struct {
	Config
	TargetAddress string `json:"target_address"`
}

type ServerConfig struct {
	Config
	NorthAddress string `json:"north_address"`
	SouthAddress string `json:"south_address"`
}
