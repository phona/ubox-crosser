package client

import "github.com/armon/go-socks5"

type Client struct {
	controller *Controller
}

func NewClient() *Client {
	return &Client{}
}

func (cli *Client) Connect(address string) error {
	conf := &socks5.Config{}
	if server, err := socks5.New(conf); err != nil {
		return err
	} else {
		cli.controller = NewController(address, server)
		cli.controller.Run()
		return nil
	}
}
