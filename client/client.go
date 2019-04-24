package client

type Client struct {
	controller *Controller
}

func NewClient() *Client {
	return &Client{}
}

func (cli *Client) Connect(address string) {
	cli.controller = NewController(address)
	cli.controller.Run()
}
