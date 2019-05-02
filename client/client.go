package client

import (
	"github.com/armon/go-socks5"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

type Client struct {
	controller *Controller
	cipher     *shadowsocks.Cipher
}

func NewClient(cipher *shadowsocks.Cipher) *Client {
	return &Client{cipher: cipher}
}

func (cli *Client) Connect(address, name, username, password string) error {
	conf := &socks5.Config{}
	if server, err := socks5.New(conf); err != nil {
		return err
	} else {
		cli.controller = NewController(address, server, cli.cipher, name, username, password)
		cli.controller.Run()
		return nil
	}
}
