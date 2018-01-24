package postgres_test

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/pinub/pinub/postgres"
)

type Client struct {
	*postgres.Client
}

func MustOpenClient() *Client {
	client, err := postgres.NewClient(os.Getenv("DATABASE"))
	if err != nil {
		panic(err)
	}

	c := &Client{Client: client}

	return c
}
