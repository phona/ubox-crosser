package main

import (
	"flag"
	"fmt"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"os"
	"ubox-crosser"
)

func main() {
	var cmdConfig crosser.ServerConfig
	var err error

	flag.StringVar(&cmdConfig.Password, "k", "", "password")
	flag.Uint64Var(&cmdConfig.MaxConnection, "c", 10, "how much connection will be created")
	flag.StringVar(&cmdConfig.NorthAddress, "n", "", "which port will be listened for north serve")
	flag.StringVar(&cmdConfig.SouthAddress, "s", "", "which port will be listened for south serve")
	flag.StringVar(&cmdConfig.Method, "m", "", "encryption method, default: aes-256-cfb")
	flag.Parse()

	if cmdConfig.NorthAddress == "" {
		fmt.Println("north address can't be empty")
	} else if cmdConfig.Password == "" {
		fmt.Println("password can't be empty")
	} else if cmdConfig.SouthAddress == "" {
		fmt.Println("south address can't be empty")
	}

	if cmdConfig.Method == "" {
		cmdConfig.Method = "aes-256-cfb"
	}
	if err = ss.CheckCipherMethod(cmdConfig.Method); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	tunnel := crosser.NewTunnel(cmdConfig.MaxConnection)
	ss.Debug = true

	go tunnel.OpenSouthWithCipher(cmdConfig.SouthAddress, cmdConfig.Method, cmdConfig.Password)
	tunnel.OpenNorth(cmdConfig.NorthAddress)
	//go tunnel.OpenSouth("127.0.0.1:7000")
}
