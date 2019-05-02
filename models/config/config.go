package config

type Controller struct {
	LoginPassword string `json:"login_password"`
}

type Config struct {
	Key        string `json:"key"`
	Method     string `json:"method"` // encryption method
	LogFile    string `json:"log_file"`
	LogLevel   string `json:"log_level"`
	ConfigFile string `json:"config_file"`
	Controller
}

type ClientConfig struct {
	Config
	TargetAddress string `json:"target_address"`
	Name          string `json:"target_address"`
}

type ServerConfig struct {
	Config
	ControllerAddress string `json:"controller_address"`
	ExposeAddress     string `json:"expose_address"`
}
