package postgres

import (
	"database/sql"

	"github.com/pinub/pinub"
)

const driver = "postgres"

// Ensure Client implements pinub.Client
var _ pinub.Client = &Client{}

// Client holds a link to the database and current implemented services.
type Client struct {
	*sql.DB

	linkService LinkService
	userService UserService
}

// New creates and returns a new Client struct.
func New(dataSource string) (*Client, error) {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	c := &Client{DB: db}
	c.linkService.Client = c
	c.userService.Client = c

	return c, nil
}

// LinkService returns the linked service.
func (c *Client) LinkService() pinub.LinkService {
	return &c.linkService
}

// UserService returns the linked services.
func (c *Client) UserService() pinub.UserService {
	return &c.userService
}
