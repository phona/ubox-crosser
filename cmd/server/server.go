package main

import (
	"flag"
	"fmt"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"os"
	"ubox-crosser"
	"ubox-crosser/log"
	"ubox-crosser/server"
)

func main() {
	var cmdConfig crosser.ServerConfig

	flag.StringVar(&cmdConfig.Password, "p", "", "password")
	flag.StringVar(&cmdConfig.NorthAddress, "cP", "", "which port will be listened for control channel")
	flag.StringVar(&cmdConfig.SouthAddress, "eP", "", "which port will be listened for external address")
	flag.StringVar(&cmdConfig.Method, "m", "", "encryption method")
	flag.Parse()

	if cmdConfig.NorthAddress == "" {
		fmt.Println("control channel address can't be empty")
	} else if cmdConfig.Method != "" && cmdConfig.Password == "" {
		fmt.Println("password can't be empty")
	} else if cmdConfig.SouthAddress == "" {
		fmt.Println("external address can't be empty")
	}

	method := cmdConfig.Method
	if cmdConfig.Method == "" {
		cmdConfig.Method = "chacha20"
	}
	if err := shadowsocks.CheckCipherMethod(cmdConfig.Method); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else if cipher, err := shadowsocks.NewCipher(cmdConfig.Method, cmdConfig.Password); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		log.SetLogLevel("debug")
		if method == "" {
			cipher = nil
		}
		proxy := server.NewProxyServer(cipher)
		proxy.Listen(cmdConfig.SouthAddress, cmdConfig.NorthAddress)
	}
}
