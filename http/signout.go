package http

import "net/http"

func (ctx *context) getSignout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(cookieName); err == nil {
		deleteCookie(w, c)
	}

	ctx.user = nil
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}
