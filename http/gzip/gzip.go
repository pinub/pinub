package gzip

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
)

// New returns a gziped response handler.
func New(next http.Handler) http.Handler {
	return gziphandler.GzipHandler(next)
}
