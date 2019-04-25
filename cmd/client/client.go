package main

import (
	"flag"
	"fmt"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"os"
	"ubox-crosser"
	"ubox-crosser/client"
	"ubox-crosser/log"
)

func main() {
	var cmdConfig crosser.ClientConfig

	flag.StringVar(&cmdConfig.Password, "p", "", "password")
	flag.StringVar(&cmdConfig.TargetAddress, "t", "", "target server address")
	flag.StringVar(&cmdConfig.Method, "m", "", "encryption method, default: aes-256-cfb")
	flag.Parse()

	if cmdConfig.TargetAddress == "" {
		fmt.Println("target address can't be empty")
	} else if cmdConfig.Method != "" && cmdConfig.Password == "" {
		fmt.Println("password can't be empty")
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
		cli := client.NewClient(cipher)
		if err := cli.Connect(cmdConfig.TargetAddress); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}
