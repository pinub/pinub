package pinub

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

var (
	// ErrLinkEmptyURL when no URL given.
	ErrLinkEmptyURL = errors.New("URL is empty")
	// ErrLinkInvalidURL when URL is invalid.
	ErrLinkInvalidURL = errors.New("Invalid URL")
	// ErrLinkInvalidHost when Host is invalid.
	ErrLinkInvalidHost = errors.New("Invalid Host")
)

// Link represents a stored url in database.
type Link struct {
	ID        string
	URL       string
	CreatedAt *time.Time
	Errors    map[string]error
}

// IsValid checks the current set values for validity.
func (l *Link) IsValid() bool {
	l.Errors = make(map[string]error)

	if strings.TrimSpace(l.URL) == "" {
		l.Errors["URL"] = ErrLinkEmptyURL
		return false
	}

	u, err := url.Parse(l.URL)
	if err != nil {
		l.Errors["URL"] = ErrLinkInvalidURL
		return false
	}

	if !strings.Contains(u.Host, ".") {
		l.Errors["URL"] = ErrLinkInvalidHost
		return false
	}

	return len(l.Errors) == 0
}
