package main

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"ubox-crosser/log"
	"ubox-crosser/models/config"
	"ubox-crosser/server"
	"ubox-crosser/utils/conf"
)

func main() {
	var cmdConfig config.ServerConfig
	cmd := &cobra.Command{
		Use: "UBox-crosser server",
		Run: func(cmd *cobra.Command, args []string) {
			var configs map[string]config.ServerConfig
			var err error
			if configs, err = conf.ParseServerConfigFile(cmdConfig.ConfigFile); err != nil {
				conf.CmdErrHandle(cmd, err)
			} else if len(configs) == 0 {
				log.InitLog(cmdConfig.LogFile, cmdConfig.LogLevel)
				configs["default"] = cmdConfig
				content, _ := json.Marshal(cmdConfig)
				logrus.Infoln("Using configuration from command line")
				logrus.Infof("Config init: %s", content)
			} else if commonConfig, ok := configs[conf.CommonConfigName]; ok {
				log.InitLog(commonConfig.LogFile, commonConfig.LogLevel)
				content, _ := json.Marshal(configs)
				logrus.Infoln("Using configuration from configure file")
				logrus.Infof("Config init: %s", content)
			} else {
				log.InitLog("", "")
				content, _ := json.Marshal(configs)
				logrus.Infoln("Log file and log level no defined, use default mode")
				logrus.Infoln("Using configuration from configure file")
				logrus.Infof("Config init: %s", content)
			}
			proxy := server.NewProxyServer(configs)
			go proxy.Process()
			func() {
				for {
					logrus.Errorln(proxy.Err())
				}
			}()

			//if err := (&cmdConfig).Update(fileConfig); err != nil {
			//	conf.CmdErrHandle(cmd, err)
			//}

			//if cmdConfig.Address == "" {
			//	conf.CmdErrHandle(cmd, "Address can't be empty")
			//} else if cmdConfig.Method != "" && cmdConfig.Key == "" {
			//	conf.CmdErrHandle(cmd, "Password can't be empty")
			//}

			//var cipher *shadowsocks.Cipher
			//if cmdConfig.Method != "" {
			//	if err := shadowsocks.CheckCipherMethod(cmdConfig.Method); err != nil {
			//		conf.CmdErrHandle(cmd, err)
			//	} else if cipher, err = shadowsocks.NewCipher(cmdConfig.Method, cmdConfig.Key); err != nil {
			//		conf.CmdErrHandle(cmd, err)
			//	}
			//}

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
