package postgres_test

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pinub/pinub"
	"github.com/pinub/pinub/claim"
	"github.com/pinub/pinub/postgres"
)

const url = "http://example.com"

func TestLinkService(t *testing.T) {
	t.Parallel()

	client, err := postgres.New(os.Getenv("DATABASE_URL"))
	claim.Ok(t, err)

	link := pinub.Link{
		URL: url,
	}
	user := pinub.User{
		Email:    "link@example.de",
		Password: "link",
	}
	s := client.LinkService()
	claim.Ok(t, client.UserService().CreateUser(&user))

	t.Run("create link", func(t *testing.T) {
		claim.Ok(t, s.CreateLink(&link, &user))
		claim.Equals(t, url, link.URL)

		claim.Ok(t, s.DeleteLink(&link, &user))
	})

	t.Run("create link twice", func(t *testing.T) {
		claim.Ok(t, s.CreateLink(&link, &user))
		claim.Ok(t, s.CreateLink(&link, &user))

		claim.Ok(t, s.DeleteLink(&link, &user))
	})

	t.Run("delete link", func(t *testing.T) {
		claim.Ok(t, s.CreateLink(&link, &user))

		claim.Ok(t, s.DeleteLink(&link, &user))
	})

	t.Run("multiple links", func(t *testing.T) {
		claim.Ok(t, s.CreateLink(&pinub.Link{URL: "http://example.com"}, &user))
		claim.Ok(t, s.CreateLink(&pinub.Link{URL: "http://example.de"}, &user))
		claim.Ok(t, s.CreateLink(&pinub.Link{URL: "http://example.io"}, &user))
		claim.Ok(t, s.CreateLink(&pinub.Link{URL: "http://example.org"}, &user))
		claim.Ok(t, s.CreateLink(&pinub.Link{URL: "http://example.net"}, &user))

		links, err := s.Links(&user)
		claim.Ok(t, err)

		claim.Equals(t, 5, len(links))
	})

	claim.Ok(t, client.UserService().DeleteUser(&user))
}
