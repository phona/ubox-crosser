package main

import (
	"fmt"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"github.com/spf13/cobra"
	"os"
	"ubox-crosser/client"
	"ubox-crosser/log"
	"ubox-crosser/models/config"
)

func main() {
	var cmdConfig config.ClientConfig
	cmd := &cobra.Command{
		Use: "UBox-crosser server",
		Run: func(cmd *cobra.Command, args []string) {
			if cmdConfig.TargetAddress == "" {
				fmt.Println("target address can't be empty")
			} else if cmdConfig.Method != "" && cmdConfig.Key == "" {
				fmt.Println("password can't be empty")
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
				log.InitLog(cmdConfig.LogFile, cmdConfig.LogLevel)
				if method == "" {
					cipher = nil
				}
				cli := client.NewClient(cipher)
				if err := cli.Connect(cmdConfig.TargetAddress, cmdConfig.Name, cmdConfig.Name, cmdConfig.LoginPassword); err != nil {
					fmt.Println(err)
					os.Exit(0)
				}
			}
		},
	}
	cmd.Flags().StringVarP(&cmdConfig.Key, "key", "k", "", "encrypt key")
	cmd.Flags().StringVarP(&cmdConfig.TargetAddress, "target-address", "t", "", "target server address, example: 127.0.0.1:7000")
	cmd.Flags().StringVarP(&cmdConfig.Method, "method", "m", "", "encryption method, default: aes-256-cfb")
	cmd.Flags().StringVarP(&cmdConfig.Name, "name", "n", "", "login username")
	cmd.Flags().StringVarP(&cmdConfig.LoginPassword, "login-password", "l", "", "login password")
	cmd.Flags().StringVar(&cmdConfig.LogFile, "log-file", "", "log file path")
	cmd.Flags().StringVar(&cmdConfig.LogLevel, "log-level", "debug", "log level: debug, info, error, warn")
	cmd.Flags().StringVar(&cmdConfig.ConfigFile, "config-file", "", "config file path")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
