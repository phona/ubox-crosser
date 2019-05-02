package main

import (
	"fmt"
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
			} else if cmdConfig.Method != "" && cmdConfig.Key == "" {
				fmt.Println("password can't be empty")
			} else if cmdConfig.ExposeAddress == "" {
				fmt.Println("external address can't be empty")
			}

			method := cmdConfig.Method
			if cmdConfig.Method == "" {
				cmdConfig.Method = "aes-256-cfb"
			}
			if err := shadowsocks.CheckCipherMethod(cmdConfig.Method); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			} else if cipher, err := shadowsocks.NewCipher(cmdConfig.Method, cmdConfig.Key); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			} else {
				log.SetLogLevel("debug")
				if method == "" {
					cipher = nil
				}
				proxy := server.NewProxyServer(cipher)
				proxy.Listen(cmdConfig.ExposeAddress, cmdConfig.ControllerAddress, cmdConfig.LoginPassword)
			}
		},
	}
	cmd.Flags().StringVarP(&cmdConfig.Key, "key", "k", "", "encrypt key")
	cmd.Flags().StringVarP(&cmdConfig.ControllerAddress, "controller-address", "c", "", "specify a address for communicate with ubox-client, example: 127.0.0.1:7000")
	cmd.Flags().StringVarP(&cmdConfig.ExposeAddress, "expose-address", "e", "", "specify a address for for accept request from internet, example: 127.0.0.1:7000")
	cmd.Flags().StringVarP(&cmdConfig.Method, "method", "m", "", "encryption method")
	cmd.Flags().StringVarP(&cmdConfig.LoginPassword, "login-password", "p", "", "root password for client login")
	cmd.Flags().StringVar(&cmdConfig.LogFile, "log-file", "", "log file path")
	cmd.Flags().StringVar(&cmdConfig.ConfigFile, "config-file", "", "config file path")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
