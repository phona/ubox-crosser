package config

import (
	"encoding/json"
	"fmt"
	"reflect"
)

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
	ServeName     string `json:"serve_name"`
	Config
}

type ServerConfig struct {
	LoginPass string `json:"login_password"`
	AuthPass  string `json:"auth_password"`
	Address   string `json:"address"`
	Config
}

type AuthServerConfig struct {
	ExposeAddress string `json:"expose_address"`
	ClientConfig
}

func (old *Config) Update(new interface{}) error {
	var inInterface map[string]interface{}
	if inrec, err := json.Marshal(new); err != nil {
		return err
	} else if err := json.Unmarshal(inrec, &inInterface); err != nil {
		return err
	} else {
		oldType := reflect.TypeOf(old).Elem()
		oldVal := reflect.ValueOf(old).Elem()
		for i := 0; i < oldType.NumField(); i++ {
			oldVField := oldVal.Field(i)
			oldTField := oldType.Field(i)

			val, ok := inInterface[oldTField.Tag.Get("json")]
			if !ok {
				continue
			}
			newVal := reflect.ValueOf(val)
			switch oldVField.Kind() {
			case reflect.Interface:
				if fmt.Sprintf("%v", val) != "" {
					oldVField.Set(newVal)
				}
			case reflect.String:
				s := newVal.String()
				if s != "" {
					oldVField.SetString(s)
				}
			case reflect.Int:
				i := newVal.Int()
				if i != 0 {
					oldVField.SetInt(i)
				}
			}
		}
		return nil
	}
}
