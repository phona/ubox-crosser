package config

type Config struct {
	Key        string `json:"key"`
	Method     string `json:"method"` // encryption method
	LogFile    string `json:"log_file"`
	LogLevel   string `json:"log_level"`
	ConfigFile string `json:"config_file"`
}

type ClientConfig struct {
	LoginPassword string `json:"login_password"`
	TargetAddress string `json:"target_address"`
	Config
}

type ServerConfig struct {
	ControllerPass    string `json:"controller_password"`
	ExposerPass       string `json:"exposer_password"`
	ControllerAddress string `json:"controller_address"`
	ExposeAddress     string `json:"expose_address"`
	Config
}

type AuthServerConfig struct {
	ExposeAddress string `json:"expose_address"`
	ClientConfig
}
