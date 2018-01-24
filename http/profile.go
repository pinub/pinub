package http

import (
	"net/http"
	"strings"

	"github.com/pinub/pinub"
)

func (ctx *context) getProfile(w http.ResponseWriter, r *http.Request) {
	ctx.renderProfile(w, ctx.user)
}

// What should be checked:
//   1. email changed, but is in database already
//   2. entered password is wrong
// Password change cases
//   3. Password is not valid
//   4. Passwords are not equal
func (ctx *context) putProfile(w http.ResponseWriter, r *http.Request) {
	u := &pinub.User{
		Email:    strings.TrimSpace(r.FormValue("email")),
		Password: strings.TrimSpace(r.FormValue("password")),
	}
	if !u.IsValid() {
		ctx.renderProfile(w, u)
		return
	}

	// 1. email changed, but is in database already
	if u.Email != ctx.user.Email {
		if old, _ := ctx.client.UserService().User(u.Email); old != nil {
			u.Errors["Email"] = errUserEmailPresentInDatabase
			ctx.renderProfile(w, u)
			return
		}
	}

	// 2. entered password is wrong
	if !ctx.user.IsValidPassword(u.Password) {
		ctx.renderProfile(w, ctx.user)
		return
	}

	ctx.renderProfile(w, u)
}

func (ctx *context) renderProfile(w http.ResponseWriter, u *pinub.User) {
	ctx.render(w, profileTmpl, struct {
		User *pinub.User
	}{
		User: u,
	})
}
