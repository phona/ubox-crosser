package main

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"ubox-crosser/log"
	"ubox-crosser/models/config"
	"ubox-crosser/utils/conf"
)

func main() {
	var cmdConfig config.ServerConfig
	cmd := &cobra.Command{
		Use: "UBox-crosser server",
		Run: func(cmd *cobra.Command, args []string) {
			var fileConfig config.ServerConfig
			if err := conf.ParseConfigFile(cmdConfig.ConfigFile, &fileConfig); err != nil {
				conf.CmdErrHandle(cmd, err)
			} else if err := (&cmdConfig).Update(fileConfig); err != nil {
				conf.CmdErrHandle(cmd, err)
			}

			if cmdConfig.Address == "" {
				conf.CmdErrHandle(cmd, "Address can't be empty")
			} else if cmdConfig.Method != "" && cmdConfig.Key == "" {
				conf.CmdErrHandle(cmd, "Password can't be empty")
			}

			//var cipher *shadowsocks.Cipher
			//if cmdConfig.Method != "" {
			//	if err := shadowsocks.CheckCipherMethod(cmdConfig.Method); err != nil {
			//		conf.CmdErrHandle(cmd, err)
			//	} else if cipher, err = shadowsocks.NewCipher(cmdConfig.Method, cmdConfig.Key); err != nil {
			//		conf.CmdErrHandle(cmd, err)
			//	}
			//}

			log.InitLog(cmdConfig.LogFile, cmdConfig.LogLevel)
			content, _ := json.Marshal(cmdConfig)
			logrus.Infof("Config init: %s", content)
			//proxy := server.NewProxyServer(cmdConfig.ExposerAddress, cmdConfig.ExposerPass,
			//	cmdConfig.ControllerAddress, cmdConfig.ControllerPass, cipher)
			//proxy.Run()
		},
	}
	cmd.Flags().StringVarP(&cmdConfig.Key, "key", "k", "", "encrypt key")
	cmd.Flags().StringVarP(&cmdConfig.Address, "exposer-address", "e", "", "specify a address for for accept request from internet, example: 127.0.0.1:7000")
	cmd.Flags().StringVarP(&cmdConfig.Method, "method", "m", "", "encryption method")
	cmd.Flags().StringVarP(&cmdConfig.LoginPass, "login-password", "C", "", "login password")
	cmd.Flags().StringVarP(&cmdConfig.AuthPass, "auth-password", "E", "", "authenticating password")
	cmd.Flags().StringVar(&cmdConfig.LogFile, "log-file", "", "log file path")
	cmd.Flags().StringVar(&cmdConfig.LogLevel, "log-level", "debug", "log file path")
	cmd.Flags().StringVar(&cmdConfig.ConfigFile, "config-file", "", "config file path")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
