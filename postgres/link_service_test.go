package postgres_test

import (
	"testing"

	"github.com/pinub/pinub"
)

const linkURL = "http://example.com"

func TestLinkService_CreateLink(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.LinkService()

	l := pinub.Link{URL: linkURL}
	u := createUser(c.UserService())

	ok(t, s.CreateLink(&l, u))
	equals(t, l.URL, linkURL)

	ok(t, s.DeleteLink(&l, u))
}

func TestLinkService_CreateLinkTwice(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.LinkService()

	l := pinub.Link{URL: linkURL}
	u := createUser(c.UserService())

	ok(t, s.CreateLink(&l, u))
	ok(t, s.CreateLink(&l, u))

	ok(t, s.DeleteLink(&l, u))
}

func TestLinkService_DeleteLink(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.LinkService()

	u := createUser(c.UserService())
	l := pinub.Link{URL: linkURL}

	ok(t, s.CreateLink(&l, u))

	ok(t, s.DeleteLink(&l, u))
}

func TestLinkService_Links(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.LinkService()

	u := createUser(c.UserService())
	s.CreateLink(&pinub.Link{URL: "http://example.com"}, u)
	s.CreateLink(&pinub.Link{URL: "http://example.de"}, u)
	s.CreateLink(&pinub.Link{URL: "http://example.io"}, u)
	s.CreateLink(&pinub.Link{URL: "http://example.org"}, u)
	s.CreateLink(&pinub.Link{URL: "http://example.net"}, u)

	links, err := s.Links(u)
	ok(t, err)

	equals(t, 5, len(links))
}
