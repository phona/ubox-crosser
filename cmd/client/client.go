package main

import (
	"flag"
	"ubox-crosser"
	"ubox-crosser/client"
	"ubox-crosser/log"
)

func main() {
	var cmdConfig crosser.ClientConfig

	flag.StringVar(&cmdConfig.Password, "k", "", "password")
	flag.Uint64Var(&cmdConfig.MaxConnection, "c", 10, "how much connection will be created")
	flag.StringVar(&cmdConfig.TargetAddress, "t", "", "target server address")
	flag.Uint64Var(&cmdConfig.Timeout, "timeout", 300, "target server address")
	flag.StringVar(&cmdConfig.Method, "m", "", "encryption method, default: aes-256-cfb")
	flag.Parse()

	log.SetLogLevel("debug")
	cli := client.NewClient()
	cli.Connect(cmdConfig.TargetAddress)

	/*
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
			connector := crosser.NewConnector(server, cmdConfig.TargetAddress, time.Duration(cmdConfig.Timeout)*time.Second)
			for i := 0; i < int(cmdConfig.MaxConnection); i++ {
				go connector.RunWithCipher(cmdConfig.Method, cmdConfig.Password)
				//go connector.Run()
				defer connector.Close()
			}
			var wg sync.WaitGroup
			wg.Add(1)
			wg.Wait()
			log.Println("All done")
		}
	*/
}
