package main

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"github.com/spf13/cobra"
	"os"
	"ubox-crosser/log"
	"ubox-crosser/models/config"
	"ubox-crosser/server"
)

func main() {
	var cmdConfig config.ServerConfig
	cmd := &cobra.Command{
		Use: "UBox-crosser server",
		Run: func(cmd *cobra.Command, args []string) {
			if cmdConfig.ControllerAddress == "" {
				fmt.Println("control channel address can't be empty")
			} else if cmdConfig.ExposeAddress == "" {
				fmt.Println("external address can't be empty")
			}

			var cipher *shadowsocks.Cipher
			if cmdConfig.Method != "" {
				if err := shadowsocks.CheckCipherMethod(cmdConfig.Method); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				} else if cipher, err = shadowsocks.NewCipher(cmdConfig.Method, cmdConfig.Key); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}

			log.InitLog(cmdConfig.LogFile, cmdConfig.LogLevel)
			log.SetLogLevel("debug")
			content, _ := json.Marshal(cmdConfig)
			logrus.Infof("Config init: %s", content)
			proxy := server.NewProxyServer(cmdConfig.ExposeAddress, cmdConfig.ExposerPass,
				cmdConfig.ControllerAddress, cmdConfig.ControllerPass, cipher)
			proxy.Run()
		},
	}
	cmd.Flags().StringVarP(&cmdConfig.Key, "key", "k", "", "encrypt key")
	cmd.Flags().StringVarP(&cmdConfig.ControllerAddress, "controller-address", "c", "", "specify a address for communicate with ubox-client, example: 127.0.0.1:7000")
	cmd.Flags().StringVarP(&cmdConfig.ExposeAddress, "expose-address", "e", "", "specify a address for for accept request from internet, example: 127.0.0.1:7000")
	cmd.Flags().StringVarP(&cmdConfig.Method, "method", "m", "", "encryption method")
	cmd.Flags().StringVarP(&cmdConfig.ControllerPass, "south-password", "C", "", "expose password")
	cmd.Flags().StringVarP(&cmdConfig.ExposerPass, "north-password", "E", "", "controller password")
	cmd.Flags().StringVar(&cmdConfig.LogFile, "log-file", "", "log file path")
	cmd.Flags().StringVar(&cmdConfig.LogLevel, "log-level", "debug", "log file path")
	cmd.Flags().StringVar(&cmdConfig.ConfigFile, "config-file", "", "config file path")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
