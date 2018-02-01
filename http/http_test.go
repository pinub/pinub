package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pinub/pinub/claim"
	my "github.com/pinub/pinub/http"
	"github.com/pinub/pinub/mock"
)

const indexURL = "/"

func setUp() (http.Handler, *mock.Client, *httptest.ResponseRecorder) {
	c := &mock.Client{}
	c.Ls.Client = c
	c.Us.Client = c

	return my.New(c, "../views"), c, httptest.NewRecorder()
}

func TestHTTP(t *testing.T) {
	t.Parallel()

	t.Run("index", func(t *testing.T) {
		handler, _, rec := setUp()

		req := httptest.NewRequest("GET", indexURL, nil)
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
	})

	t.Run("ignored files", func(t *testing.T) {
		handler, _, rec := setUp()

		requests := []string{
			"/apple-touch-icon-120x120-precomposed.png",
			"/apple-touch-icon-120x120.png",
			"/apple-touch-icon.png",
			"/favicon.ico",
		}

		for _, url := range requests {
			req := httptest.NewRequest("GET", url, nil)
			handler.ServeHTTP(rec, req)

			claim.Equals(t, 404, rec.Code)
		}
	})
}
