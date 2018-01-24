package http_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	my "github.com/pinub/pinub/http"
	"github.com/pinub/pinub/mock"
)

func setUp() (http.Handler, *mock.Client, *httptest.ResponseRecorder) {
	c := &mock.Client{}
	c.Ls.Client = c
	c.Us.Client = c

	return my.New(c, "../views"), c, httptest.NewRecorder()
}

func TestHTTP_Index(t *testing.T) {
	t.Parallel()
	h, _, r := setUp()

	req, err := http.NewRequest("GET", "/", nil)
	ok(t, err)

	h.ServeHTTP(r, req)
	equals(t, r.Code, http.StatusOK)
}

func TestHTTP_IgnoredFiles(t *testing.T) {
	h, _, r := setUp()

	requests := []string{
		"/apple-touch-icon-120x120-precomposed.png",
		"/apple-touch-icon-120x120.png",
		"/apple-touch-icon.png",
		"/favicon.ico",
	}

	for _, request := range requests {
		req, err := http.NewRequest("GET", request, nil)
		ok(t, err)

		h.ServeHTTP(r, req)
		equals(t, r.Code, http.StatusNotFound)
	}
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if want is not equal to got.
func equals(tb testing.TB, want, got interface{}) {
	if !reflect.DeepEqual(want, got) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\twant: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, want, got)
		tb.FailNow()
	}
}
