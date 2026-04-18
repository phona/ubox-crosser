package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"github.com/spf13/cobra"
	"os"
	"github.com/phona/ubox-crosser/httpapi"
	"github.com/phona/ubox-crosser/log"
	"github.com/phona/ubox-crosser/models/config"
	"github.com/phona/ubox-crosser/server"
	"github.com/phona/ubox-crosser/utils/conf"
)

func main() {
	var cmdConfig config.AuthServerConfig
	cmd := &cobra.Command{
		Use: "UBox-crosser authenticated server",
		Run: func(cmd *cobra.Command, args []string) {
			var fileConfig config.AuthServerConfig
			if err := conf.ParseConfigFile(cmdConfig.ConfigFile, &fileConfig); err != nil {
				conf.CmdErrHandle(cmd, err)
			} else if err := fileConfig.Update(cmdConfig); err != nil {
				conf.CmdErrHandle(cmd, err)
			}

			if fileConfig.TargetAddress == "" {
				conf.CmdErrHandle(cmd, "target address can't be empty")
			} else if fileConfig.Method != "" && fileConfig.Key == "" {
				conf.CmdErrHandle(cmd, "password can't be empty")
			} else if fileConfig.ExposeAddress == "" {
				conf.CmdErrHandle(cmd, "external address can't be empty")
			}

			var cipher *shadowsocks.Cipher
			if fileConfig.Method != "" {
				if err := shadowsocks.CheckCipherMethod(fileConfig.Method); err != nil {
					conf.CmdErrHandle(cmd, err)
				} else if cipher, err = shadowsocks.NewCipher(fileConfig.Method, fileConfig.Key); err != nil {
					conf.CmdErrHandle(cmd, err)
				}
			}

			log.InitLog(fileConfig.LogFile, fileConfig.LogLevel)
			content, _ := json.Marshal(fileConfig)
			logrus.Infof("Config init: %s", content)

			httpapi.Init()
			if fileConfig.ApiAddr != "" {
				go httpapi.StartServer(fileConfig.ApiAddr)
			}
			proxy := server.NewAuthServer(fileConfig.TargetAddress, fileConfig.ServeName, fileConfig.LoginPassword, cipher)
			proxy.Listen(fileConfig.ExposeAddress)
		},
	}
	cmd.Flags().StringVarP(&cmdConfig.Key, "key", "k", "", "encrypt key")
	cmd.Flags().StringVarP(&cmdConfig.LoginPassword, "login-password", "p", "", "login password")
	cmd.Flags().StringVarP(&cmdConfig.TargetAddress, "target-address", "t", "", "specify a address to tunnel server, example: 127.0.0.1:7000")
	cmd.Flags().StringVarP(&cmdConfig.ExposeAddress, "expose-address", "e", "", "specify a address for proxy server, example: 127.0.0.1:7000")
	cmd.Flags().StringVarP(&cmdConfig.Method, "method", "m", "", "encryption method")
	cmd.Flags().StringVarP(&cmdConfig.ServeName, "serve-name", "n", "", "serve name")
	cmd.Flags().StringVar(&cmdConfig.LogFile, "log-file", "", "log file path")
	cmd.Flags().StringVar(&cmdConfig.LogLevel, "log-level", "", "log file path")
	cmd.Flags().StringVar(&cmdConfig.ConfigFile, "config-file", "", "config file path")
	cmd.Flags().StringVar(&cmdConfig.ApiAddr, "api-addr", "", "HTTP API listen address (e.g. :8080)")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
