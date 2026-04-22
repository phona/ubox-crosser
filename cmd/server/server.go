package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/phona/ubox-crosser/log"
	"github.com/phona/ubox-crosser/models/config"
	"github.com/phona/ubox-crosser/server"
	"github.com/phona/ubox-crosser/utils/conf"
)

func main() {
	var cmdConfig config.ServerConfig
	var adminAddr string
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

			adminMux := server.NewAdminMux()
			go func() {
				logrus.Infof("Admin HTTP server listening on %s", adminAddr)
				if err := http.ListenAndServe(adminAddr, adminMux); err != nil {
					logrus.Errorf("Admin HTTP server error: %v", err)
				}
			}()

			proxy := server.NewProxyServer(configs)
			go proxy.Process()
			func() {
				for {
					logrus.Errorln(proxy.Err())
				}
			}()
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
	cmd.Flags().StringVar(&adminAddr, "admin-addr", ":8080", "admin HTTP server listen address")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
