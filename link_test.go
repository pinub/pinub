package pinub_test

import (
	"testing"

	"github.com/pinub/pinub"
	"github.com/pinub/pinub/claim"
)

func TestLink_IsValid(t *testing.T) {
	t.Parallel()
	t.Run("valid", func(t *testing.T) {
		urls := []string{
			"http://www.google.com/",
			"https://www.google.com/",
			"https://www.google.com:443/",
		}

		for _, url := range urls {
			link := &pinub.Link{URL: url}
			claim.Assert(t, link.IsValid(), "link should be valid")
		}
	})

	t.Run("invalid", func(t *testing.T) {
		urls := []string{
			"",
			"unknown",
			"www.google.com",
			"google.com",
		}

		for _, url := range urls {
			link := &pinub.Link{URL: url}
			claim.Assert(t, !link.IsValid(), "link should be invalid")
		}
	})
}
