package mock

import (
	"github.com/pinub/pinub"
)

// Ensure Client implements pinub.Client
var _ pinub.Client = &Client{}

type Client struct {
	Us UserService
	Ls LinkService
}

func (c *Client) UserService() pinub.UserService {
	return &c.Us
}
func (c *Client) LinkService() pinub.LinkService {
	return &c.Ls
}
