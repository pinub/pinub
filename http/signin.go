package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/pinub/pinub"
)

var (
	errUserEmailNotCorrect = pinub.ErrUserPasswordNotCorrect
	errUserCannotSignin    = errors.New("Cannot sign in user")
)

func (ctx *context) getSignin(w http.ResponseWriter, r *http.Request) {
	ctx.renderSignin(w, &pinub.User{})
}

func (ctx *context) postSignin(w http.ResponseWriter, r *http.Request) {
	tmp := &pinub.User{
		Email:    strings.TrimSpace(r.FormValue("email")),
		Password: strings.TrimSpace(r.FormValue("password")),
	}

	if !tmp.IsValid() {
		ctx.renderSignin(w, tmp)
		return
	}

	u, err := ctx.client.UserService().User(tmp.Email)
	if err != nil {
		tmp.Errors["Email"] = errUserEmailNotCorrect
		tmp.Errors["Password"] = tmp.Errors["Email"]
		ctx.renderSignin(w, tmp)
		return
	}

	if !u.IsValidPassword(tmp.Password) {
		u.Errors["Email"] = u.Errors["Password"]
		ctx.renderSignin(w, u)
		return
	}

	// everything is fine
	err = ctx.client.UserService().AddToken(u)
	if err != nil {
		tmp.Errors["Email"] = errUserCannotSignin
		ctx.renderSignin(w, tmp)
		return
	}

	createCookie(w, u.Token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (ctx *context) renderSignin(w http.ResponseWriter, u *pinub.User) {
	ctx.render(w, signinTmpl, struct {
		User *pinub.User
	}{
		User: u,
	})
}
