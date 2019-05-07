package conf

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
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

func ParseServerConfigFile(filePath string, config interface{}) error {
	rootConfig := map[string]interface{}{}
	if err := ParseConfigFile(filePath, &rootConfig); err != nil {
		return err
	}

}

func CmdErrHandle(cmd *cobra.Command, msg ...interface{}) {
	cmd.Help()
	fmt.Println(msg)
	os.Exit(0)
}
