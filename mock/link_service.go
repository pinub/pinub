package mock

import "github.com/pinub/pinub"

// Ensure LinkService implements pinub.LinkService
var _ pinub.LinkService = &LinkService{}

type LinkService struct {
	Client *Client

	LinksFn      func(*pinub.User) ([]pinub.Link, error)
	CreateLinkFn func(*pinub.Link, *pinub.User) error
	DeleteLinkFn func(*pinub.Link, *pinub.User) error

	LinksInvoked      bool
	CreateLinkInvoked bool
	DeleteLinkInvoked bool
}

func (s *LinkService) Links(user *pinub.User) ([]pinub.Link, error) {
	s.LinksInvoked = true
	return s.LinksFn(user)
}

func (s *LinkService) CreateLink(link *pinub.Link, user *pinub.User) error {
	s.CreateLinkInvoked = true
	return s.CreateLinkFn(link, user)
}

func (s *LinkService) DeleteLink(link *pinub.Link, user *pinub.User) error {
	s.DeleteLinkInvoked = true
	return s.DeleteLinkFn(link, user)
}
