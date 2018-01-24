package http

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/pinub/pinub"
)

var (
	errUserEmailPresentInDatabase = errors.New("User exists already")
	errUserCreatePasswordHash     = errors.New("Error creating the password hash")
	errUserCannotPersist          = errors.New("Error persisting user to database")
	errUserPasswordNotMatching    = errors.New("Passwords do not match")
)

func (ctx *context) getRegister(w http.ResponseWriter, r *http.Request) {
	ctx.renderRegister(w, &pinub.User{})
}

func (ctx *context) postRegister(w http.ResponseWriter, r *http.Request) {
	new := &pinub.User{
		Email:    strings.TrimSpace(r.FormValue("email")),
		Password: strings.TrimSpace(r.FormValue("password")),
	}

	if !new.IsValid() {
		ctx.renderRegister(w, new)
		return
	}

	if new.Password != strings.TrimSpace(r.FormValue("password_confirm")) {
		new.Errors["PasswordConfirm"] = errUserPasswordNotMatching
		ctx.renderRegister(w, new)
		return
	}

	// check if in database already
	if u, _ := ctx.client.UserService().User(new.Email); u != nil {
		new.Errors["Password"] = errUserEmailPresentInDatabase
		ctx.renderRegister(w, new)
		return
	}

	hash, err := pinub.HashPassword(new.Password)
	if err != nil {
		new.Errors["Password"] = errUserCreatePasswordHash
		ctx.renderRegister(w, new)
		return
	}

	new.Password = hash
	if err = ctx.client.UserService().CreateUser(new); err != nil {
		new.Errors["Email"] = errUserCannotPersist
		ctx.renderRegister(w, new)
		return
	}

	// everything is fine
	err = ctx.client.UserService().AddToken(new)
	if err != nil {
		log.Printf("Error saving token %v", err)
	}

	createCookie(w, new.Token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (ctx *context) renderRegister(w http.ResponseWriter, u *pinub.User) {
	ctx.render(w, registerTmpl, struct {
		User *pinub.User
	}{
		User: u,
	})
}
