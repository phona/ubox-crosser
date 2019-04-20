package main

import (
	"flag"
	"fmt"
	"github.com/armon/go-socks5"
	"log"
	"os"
	"sync"
	"crosser"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

func main() {
	var cmdConfig crosser.ClientConfig
	var err error

	flag.StringVar(&cmdConfig.Password, "k", "", "password")
	flag.IntVar(&cmdConfig.MaxConnection, "c", 10, "how much connection will be created")
	flag.StringVar(&cmdConfig.TargetAddress, "t", "", "target server address")
	flag.StringVar(&cmdConfig.Method, "m", "", "encryption method, default: aes-256-cfb")
	flag.Parse()

	if cmdConfig.TargetAddress == "" {
		fmt.Println("target address can't be empty")
	} else if cmdConfig.Password == "" {
		fmt.Println("password can't be empty")
	}

	if cmdConfig.Method == "" {
		cmdConfig.Method = "aes-256-cfb"
	}
	if err = ss.CheckCipherMethod(cmdConfig.Method); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	conf := &socks5.Config{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	if server, err := socks5.New(conf); err != nil {
		fmt.Println(err)
		os.Exit(0)
	} else {
		for i := 0; i < cmdConfig.MaxConnection; i++ {
			connector := crosser.NewConnector(server, cmdConfig.TargetAddress)
			go connector.RunWithCipher(cmdConfig.Method, cmdConfig.Password)
			//go connector.Run()
			defer connector.Close()
		}
		var wg sync.WaitGroup
		wg.Add(1)
		wg.Wait()
		log.Println("All done")
	}
}