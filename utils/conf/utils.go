package conf

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"ubox-crosser/models/config"
)

func ParseConfigFile(filePath string, config interface{}) error {
	if filePath == "" {
		return nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err := json.Unmarshal(data, config); err != nil {
		return err
	} else {
		return nil
	}
}

var CommonConfigName = "common"

func ParseServerConfigFile(filePath string) (map[string]config.ServerConfig, error) {
	rootConfig := make(map[string]config.ServerConfig, 10)
	if err := ParseConfigFile(filePath, &rootConfig); err != nil {
		return rootConfig, err
	}

	if common, ok := rootConfig[CommonConfigName]; ok {
		for k, v := range rootConfig {
			if k != CommonConfigName {
				newConfig := common
				if err := newConfig.Update(v); err != nil {
					return rootConfig, err
				} else {
					rootConfig[k] = newConfig
				}
			}
		}
	}
	return rootConfig, nil
}

func CmdErrHandle(cmd *cobra.Command, msg ...interface{}) {
	cmd.Help()
	fmt.Println(msg)
	os.Exit(0)
}
