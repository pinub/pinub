package http

import (
	"log"
	"net/http"
	"net/url"

	"github.com/pinub/pinub"
)

func (ctx *context) getLinks(w http.ResponseWriter, r *http.Request) {
	if ctx.user == nil {
		ctx.render(w, homeTmpl, nil)
		return
	}

	links, err := ctx.client.LinkService().Links(ctx.user)
	ctx.render(w, listLinksTmpl, struct {
		Links  []pinub.Link
		Errors error
	}{
		Links:  links,
		Errors: err,
	})
}

func (ctx *context) createLink(w http.ResponseWriter, r *http.Request) {
	link := &pinub.Link{URL: r.URL.String()[1:]}

	var isValid bool
	for _, v := range []string{"", "http://"} {
		link.URL = v + link.URL
		if isValid = link.IsValid(); isValid {
			break
		}
	}
	if !isValid {
		log.Print(link.Errors["URL"])
		http.NotFound(w, r)
		return
	}

	cleanURL, _ := url.ParseRequestURI(link.URL)
	link.URL = cleanURL.String()
	if err := ctx.client.LinkService().CreateLink(link, ctx.user); err != nil {
		log.Printf("cannot create link %v", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
