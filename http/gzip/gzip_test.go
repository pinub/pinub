package gzip

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pinub/pinub/claim"
)

var (
	body  = strings.Repeat("hello world", 1000)
	hello = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})
)

func TestGzip(t *testing.T) {
	t.Parallel()

	h := New(hello)

	t.Run("accepts gzip", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept-Encoding", "gzip")

		h.ServeHTTP(res, req)

		header := make(http.Header)
		header.Add("Content-Type", "text/plain; charset=utf-8")
		header.Add("Content-Encoding", "gzip")
		header.Add("Vary", "Accept-Encoding")

		claim.Equals(t, 200, res.Code)
		claim.Equals(t, header, res.HeaderMap)

		gz, err := gzip.NewReader(res.Body)
		claim.Ok(t, err)

		b, err := ioutil.ReadAll(gz)
		claim.Ok(t, err)
		claim.Ok(t, gz.Close())

		claim.Equals(t, body, string(b))
	})

	t.Run("accepts identity", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		h.ServeHTTP(res, req)

		header := make(http.Header)
		header.Add("Content-Type", "text/plain; charset=utf-8")
		header.Add("Vary", "Accept-Encoding")

		claim.Equals(t, 200, res.Code)
		claim.Equals(t, header, res.HeaderMap)

		claim.Equals(t, body, res.Body.String())
	})
}
