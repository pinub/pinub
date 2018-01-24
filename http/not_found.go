package http

import "net/http"

func (ctx *context) getNotFound(w http.ResponseWriter, r *http.Request) {
	if _, ok := ignoredFiles[r.URL.String()]; ok {
		http.NotFound(w, r)
		return
	}
	if ctx.user == nil {
		http.NotFound(w, r)
		return
	}

	ctx.createLink(w, r)
}
