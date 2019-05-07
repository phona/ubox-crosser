package main

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"github.com/spf13/cobra"
	"os"
	"ubox-crosser/client"
	"ubox-crosser/log"
	"ubox-crosser/models/config"
	"ubox-crosser/utils/conf"
)

func main() {
	var cmdConfig config.ClientConfig
	cmd := &cobra.Command{
		Use: "UBox-crosser server",
		Run: func(cmd *cobra.Command, args []string) {
			var fileConfig config.ClientConfig
			if err := conf.ParseConfigFile(cmdConfig.ConfigFile, &fileConfig); err != nil {
				conf.CmdErrHandle(cmd, err)
			} else if err := (&cmdConfig).Update(fileConfig); err != nil {
				conf.CmdErrHandle(cmd, err)
			}

			if cmdConfig.TargetAddress == "" {
				conf.CmdErrHandle(cmd, "target address can't be empty")
			} else if cmdConfig.Method != "" && cmdConfig.Key == "" {
				conf.CmdErrHandle(cmd, "password can't be empty")
			}

			var cipher *shadowsocks.Cipher
			if cmdConfig.Method != "" {
				if err := shadowsocks.CheckCipherMethod(cmdConfig.Method); err != nil {
					conf.CmdErrHandle(cmd, err)
				} else if cipher, err = shadowsocks.NewCipher(cmdConfig.Method, cmdConfig.Key); err != nil {
					conf.CmdErrHandle(cmd, err)
				}
			}

			log.InitLog(cmdConfig.LogFile, cmdConfig.LogLevel)
			content, _ := json.Marshal(cmdConfig)
			logrus.Infof("Config init: %s", content)
			cli := client.NewClient(cipher)
			if err := cli.Connect(cmdConfig.TargetAddress, cmdConfig.LoginPassword); err != nil {
				conf.CmdErrHandle(cmd, err)
			}
		},
	}
	cmd.Flags().StringVarP(&cmdConfig.Key, "key", "k", "", "encrypt key")
	cmd.Flags().StringVarP(&cmdConfig.TargetAddress, "target-address", "t", "", "target server address, example: 127.0.0.1:7000")
	cmd.Flags().StringVarP(&cmdConfig.Method, "method", "m", "", "encryption method, default: aes-256-cfb")
	cmd.Flags().StringVarP(&cmdConfig.LoginPassword, "login-password", "p", "", "login password")
	cmd.Flags().StringVar(&cmdConfig.LogFile, "log-file", "", "log file path")
	cmd.Flags().StringVar(&cmdConfig.LogLevel, "log-level", "debug", "log level: debug, info, error, warn")
	cmd.Flags().StringVar(&cmdConfig.ConfigFile, "config-file", "", "config file path")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
