package http

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pinub/mux"
	"github.com/pinub/pinub"
	"github.com/pinub/pinub/http/gzip"
)

const (
	cookieName     = "keks"
	cookieDays     = 14
	deleteMeCookie = "deleteMe"

	homeURL     = "/"
	signinURL   = "/signin"
	signoutURL  = "/signout"
	registerURL = "/register"
	profileURL  = "/profile"
)

const layoutTmpl = "_layout.html"
const signinTmpl = "signin.html"
const registerTmpl = "register.html"
const profileTmpl = "profile.html"
const listLinksTmpl = "list_links.html"
const homeTmpl = "home.html"

var ignoredFiles = map[string]bool{
	"/apple-touch-icon-152x152-precomposed.png": true,
	"/apple-touch-icon-152x152.png":             true,
	"/apple-touch-icon-120x120-precomposed.png": true,
	"/apple-touch-icon-120x120.png":             true,
	"/apple-touch-icon-precomposed.png":         true,
	"/apple-touch-icon.png":                     true,
	"/favicon.ico":                              true,
}

type context struct {
	client pinub.Client
	user   *pinub.User
	tmpl   map[string]*template.Template
}

// New creates a new context and a handler that gets returned.
func New(c pinub.Client, tmplPath string) http.Handler {
	ctx := &context{
		client: c,
		tmpl:   prepareTemplates(tmplPath),
	}

	m := mux.New()
	// Public
	m.Get(signinURL, ctx.public(ctx.getSignin))
	m.Post(signinURL, ctx.public(ctx.postSignin))
	m.Get(registerURL, ctx.public(ctx.getRegister))
	m.Post(registerURL, ctx.public(ctx.postRegister))
	// Private
	m.Get(signoutURL, ctx.private(ctx.getSignout))
	m.Get(profileURL, ctx.private(ctx.getProfile))
	m.Put(profileURL, ctx.private(ctx.putProfile))
	m.Get(homeURL, ctx.getLinks)
	m.NotFound = ctx.getNotFound

	var h http.Handler
	h = ctx.deleteLinks(m)
	//	h = http2.New([]string{
	//		cssFilePath(),
	//		jsFilePath(),
	//	}, h)
	h = ctx.auth(h)
	h = gzip.New(h)

	return h //ctx.auth(gzip.New(ctx.deleteLinks(pusher(m))))
}

func (ctx *context) render(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	if err := ctx.tmpl[name].Execute(w, data); err != nil {
		log.Fatal(err)
	}
}

func (ctx *context) private(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ctx.user == nil {
			http.Redirect(w, r, homeURL, http.StatusSeeOther)
		} else {
			f(w, r)
		}
	})
}

func (ctx *context) public(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ctx.user != nil {
			http.Redirect(w, r, homeURL, http.StatusSeeOther)
		} else {
			f(w, r)
		}
	})
}

func (ctx *context) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(cookieName); err == nil {
			if u, err := ctx.client.UserService().UserByToken(c.Value); err == nil {
				ctx.user = u
				// do not send cookie when we want to sign out
				if !strings.HasPrefix(r.URL.String(), signoutURL) {
					refreshCookie(w, c)
				}
				go ctx.client.UserService().RefreshToken(u) // nolint: errcheck
			} else {
				deleteCookie(w, c)
			}
		}

		next.ServeHTTP(w, r)

		// set to nil for next request
		ctx.user = nil
	})
}

// deleteLinks middleware - makes sure, that links get deleted on every
// request, not only when viewing the linking list.
func (ctx *context) deleteLinks(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(deleteMeCookie); err == nil && ctx.user != nil {
			delete := strings.Split(c.Value, ",")
			// this is not working due a bug in aws api gateway:
			// https://forums.aws.amazon.com/thread.jspa?messageID=701434
			// the fix is done on the client side
			//deleteCookie(w, c)

			for _, id := range delete {
				ctx.client.LinkService().DeleteLink(&pinub.Link{ID: id}, ctx.user) // nolint: errcheck
			}

		}

		next.ServeHTTP(w, r)
	})
}

func pusher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if pusher, ok := w.(http.Pusher); ok {
			if err := pusher.Push(cssFilePath(), nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
			if err := pusher.Push(jsFilePath(), nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Helper methods

func createCookie(w http.ResponseWriter, value string) {
	c := &http.Cookie{
		Name:    cookieName,
		Value:   value,
		Path:    "/",
		Expires: time.Now().Add(cookieDays * time.Hour * 24),
	}

	http.SetCookie(w, c)
}

func refreshCookie(w http.ResponseWriter, c *http.Cookie) {
	c.Path = "/"
	c.Expires = time.Now().Add(cookieDays * time.Hour * 24)
	http.SetCookie(w, c)
}

func deleteCookie(w http.ResponseWriter, c *http.Cookie) {
	c.Path = "/"
	c.MaxAge = -1
	http.SetCookie(w, c)
}

func prepareTemplates(tmplPath string) map[string]*template.Template {
	tmpl := make(map[string]*template.Template)
	layout := template.Must(template.New(layoutTmpl).Funcs(template.FuncMap{
		"time":     formatTime,
		"datetime": datetime,
		"url":      formatURL,
		"css":      cssFilePath,
		"js":       jsFilePath,
	}).ParseFiles(tmplPath + "/" + layoutTmpl))

	files, _ := filepath.Glob(tmplPath + "/*.html")
	for _, t := range files {
		if path.Base(t) == layoutTmpl {
			continue
		}
		tmpl[path.Base(t)] = template.Must(template.Must(layout.Clone()).ParseFiles(t))
	}

	return tmpl
}

func datetime(t time.Time) string {
	return t.Format("2006-01-02T15:04Z")
}

func formatTime(t time.Time) string {
	diff := time.Since(t)
	switch {
	case diff < time.Minute:
		return fmt.Sprintf("%.0fs ago", math.Abs(diff.Seconds()))
	case diff < time.Hour:
		// contains seconds which we do not want here
		// return diff.Round(time.Minute).String()
		return fmt.Sprintf("%.0fm ago", diff.Minutes())
	case diff < time.Hour*8:
		return fmt.Sprintf("%.0fh ago", diff.Hours())
	}

	return t.Format("02.01.06 15:04")
}

func formatURL(url string) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	return url
}

func cssFilePath() string {
	return os.Getenv("CSS")
}

func jsFilePath() string {
	return os.Getenv("JS")
}
